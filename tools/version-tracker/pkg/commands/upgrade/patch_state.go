package upgrade

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type repairFileState struct {
	content []byte
	mode    os.FileMode
}

// patchRepairState covers repair output only; version metadata belongs to the
// upgrade and survives a failed repair.
type patchRepairState struct {
	projectRoot      string
	repoPath         string
	patchesDir       string
	baseRevision     string
	patchesDirExists bool
	patchesDirMode   os.FileMode
	patches          map[string]repairFileState
	generatedDirs    []string
	generatedFiles   map[string]repairFileState
	untracked        map[string]struct{}
}

func captureRepairState(projectRoot, repoPath, patchesDir, releaseBranch string, appliedPatches int) (*patchRepairState, error) {
	if appliedPatches < 0 {
		return nil, fmt.Errorf("invalid applied patch count %d", appliedPatches)
	}
	root, repo, err := repairRepositoryPaths(context.Background(), projectRoot, repoPath)
	if err != nil {
		return nil, err
	}
	patches, err := repairPathInProject(root, patchesDir)
	if err != nil {
		return nil, err
	}
	if repairPathContains(repo, patches) || repairPathContains(patches, repo) {
		return nil, fmt.Errorf("patch directory %q overlaps upstream repository %q", patches, repo)
	}
	state := &patchRepairState{
		projectRoot:    root,
		repoPath:       repo,
		patchesDir:     patches,
		patches:        make(map[string]repairFileState),
		generatedDirs:  []string{root},
		generatedFiles: make(map[string]repairFileState),
		untracked:      make(map[string]struct{}),
	}
	if releaseBranch != "" {
		if filepath.Base(releaseBranch) != releaseBranch || releaseBranch == "." || releaseBranch == ".." {
			return nil, fmt.Errorf("invalid release branch directory %q", releaseBranch)
		}
		dir, err := repairPathInProject(root, filepath.Join(root, releaseBranch))
		if err != nil {
			return nil, err
		}
		if repairPathContains(repo, dir) || repairPathContains(dir, repo) {
			return nil, fmt.Errorf("release directory %q overlaps upstream repository %q", dir, repo)
		}
		state.generatedDirs = append(state.generatedDirs, dir)
	}

	// A failed git am has already committed the clean prefix.
	output, err := repairGitOutput(context.Background(), repo, "rev-parse", "--verify", fmt.Sprintf("HEAD~%d^{commit}", appliedPatches))
	if err != nil {
		return nil, fmt.Errorf("finding upstream baseline: %w", err)
	}
	state.baseRevision = strings.TrimSpace(string(output))
	if info, err := os.Stat(patches); err == nil {
		if !info.IsDir() {
			return nil, fmt.Errorf("patch path %q is not a directory", patches)
		}
		state.patchesDirExists = true
		state.patchesDirMode = info.Mode().Perm()
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading patch directory: %w", err)
	}
	if err := captureRepairFiles(patches, isRepairPatch, state.patches); err != nil {
		return nil, err
	}
	for _, dir := range state.generatedDirs {
		files := make(map[string]repairFileState)
		if err := captureRepairFiles(dir, isRepairGeneratedFile, files); err != nil {
			return nil, err
		}
		for name, file := range files {
			state.generatedFiles[filepath.Join(dir, name)] = file
		}
	}
	untracked, err := repairUntrackedPaths(repo)
	if err != nil {
		return nil, fmt.Errorf("listing original untracked source files: %w", err)
	}
	for _, name := range untracked {
		state.untracked[name] = struct{}{}
	}
	return state, nil
}

// restore removes generated repair output and returns the task's upstream
// checkout to its baseline. It never resets the surrounding build-tooling repo.
func (s *patchRepairState) restore() error {
	if err := s.checkPaths(context.Background()); err != nil {
		return err
	}
	var restoreErrors []error
	if err := restoreRepairFiles(s.patchesDir, isRepairPatch, s.patches); err != nil {
		restoreErrors = append(restoreErrors, err)
	} else if s.patchesDirExists {
		if err := os.MkdirAll(s.patchesDir, 0o755); err != nil {
			restoreErrors = append(restoreErrors, err)
		} else if err := os.Chmod(s.patchesDir, s.patchesDirMode); err != nil {
			restoreErrors = append(restoreErrors, err)
		}
	} else if err := removeEmptyRepairDirectory(s.patchesDir); err != nil {
		restoreErrors = append(restoreErrors, err)
	}
	for _, dir := range s.generatedDirs {
		files := make(map[string]repairFileState)
		for path, file := range s.generatedFiles {
			if filepath.Dir(path) == dir {
				files[filepath.Base(path)] = file
			}
		}
		if err := restoreRepairFiles(dir, isRepairGeneratedFile, files); err != nil {
			restoreErrors = append(restoreErrors, err)
		}
	}

	gitDir, err := repairGitOutput(context.Background(), s.repoPath, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return errors.Join(append(restoreErrors, err)...)
	}
	if _, err := os.Stat(filepath.Join(strings.TrimSpace(string(gitDir)), "rebase-apply")); err == nil {
		if _, err := repairGitOutput(context.Background(), s.repoPath, "am", "--abort"); err != nil {
			// Reset must not hide a failure to clear the active am session.
			return errors.Join(append(restoreErrors, err)...)
		}
	} else if !os.IsNotExist(err) {
		return errors.Join(append(restoreErrors, err)...)
	}
	if _, err := repairGitOutput(context.Background(), s.repoPath, "reset", "--hard", s.baseRevision); err != nil {
		return errors.Join(append(restoreErrors, err)...)
	}
	untracked, err := repairUntrackedPaths(s.repoPath)
	if err != nil {
		return errors.Join(append(restoreErrors, err)...)
	}
	for _, name := range untracked {
		if _, original := s.untracked[name]; original {
			continue
		}
		if strings.HasSuffix(name, "/") {
			hasOriginal := false
			for original := range s.untracked {
				if strings.HasPrefix(original, name) {
					hasOriginal = true
					break
				}
			}
			if hasOriginal {
				continue
			}
		}
		// Literal pathspecs keep unusual source filenames from selecting other
		// paths. Include ignored build output, preserving pre-existing files.
		if _, err := repairGitOutput(context.Background(), s.repoPath, "clean", "-fdx", "--", name); err != nil {
			restoreErrors = append(restoreErrors, err)
		}
	}
	return errors.Join(restoreErrors...)
}

// dockerfiles returns repository-relative paths changed by either committed
// patches or working tree edits since the upstream baseline.
func (s *patchRepairState) dockerfiles(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := s.checkPaths(ctx); err != nil {
		return nil, err
	}
	output, err := repairGitOutput(ctx, s.repoPath, "diff", "--name-only", "-z", s.baseRevision, "--")
	if err != nil {
		return nil, fmt.Errorf("listing repaired Dockerfiles: %w", err)
	}
	paths := make(map[string]struct{})
	for _, path := range strings.Split(string(output), "\x00") {
		name := strings.ToLower(filepath.Base(path))
		if name == "dockerfile" || strings.HasPrefix(name, "dockerfile.") || strings.HasSuffix(name, ".dockerfile") {
			paths[path] = struct{}{}
		}
	}
	files := make([]string, 0, len(paths))
	for path := range paths {
		files = append(files, path)
	}
	slices.Sort(files)
	return files, nil
}

func (s *patchRepairState) checkPaths(ctx context.Context) error {
	root, repo, err := repairRepositoryPaths(ctx, s.projectRoot, s.repoPath)
	if err != nil {
		return err
	}
	if root != s.projectRoot || repo != s.repoPath {
		return fmt.Errorf("repair repository paths changed after capture")
	}
	for _, path := range append([]string{s.patchesDir}, s.generatedDirs[1:]...) {
		resolved, err := repairPathInProject(root, path)
		if err != nil {
			return err
		}
		if resolved != path {
			return fmt.Errorf("repair output path %q changed after capture", path)
		}
	}
	return nil
}

func repairRepositoryPaths(ctx context.Context, projectRoot, repoPath string) (string, string, error) {
	root, err := repairResolvedPath(projectRoot)
	if err != nil {
		return "", "", err
	}
	info, err := os.Stat(root)
	if err != nil {
		return "", "", err
	}
	if !info.IsDir() {
		return "", "", fmt.Errorf("project root %q is not a directory", root)
	}
	repo, err := repairPathInProject(root, repoPath)
	if err != nil {
		return "", "", err
	}
	top, err := repairGitOutput(ctx, repo, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", err
	}
	topPath, err := repairResolvedPath(strings.TrimSpace(string(top)))
	if err != nil || topPath != repo {
		return "", "", fmt.Errorf("upstream path %q is not a separate repository root", repo)
	}
	gitDir, err := repairGitOutput(ctx, repo, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return "", "", err
	}
	if _, err := repairPathInProject(repo, strings.TrimSpace(string(gitDir))); err != nil {
		return "", "", fmt.Errorf("upstream Git directory is outside its checkout: %w", err)
	}
	commonDir, err := repairGitOutput(ctx, repo, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", "", err
	}
	commonPath := strings.TrimSpace(string(commonDir))
	if !filepath.IsAbs(commonPath) {
		commonPath = filepath.Join(repo, commonPath)
	}
	if _, err := repairPathInProject(repo, commonPath); err != nil {
		return "", "", fmt.Errorf("upstream common Git directory is outside its checkout: %w", err)
	}
	return root, repo, nil
}

func repairPathInProject(root, path string) (string, error) {
	resolved, err := repairResolvedPath(path)
	if err != nil {
		return "", err
	}
	if resolved == root || !repairPathContains(root, resolved) {
		return "", fmt.Errorf("repair path %q must be strictly inside project root %q", path, root)
	}
	return resolved, nil
}

func repairPathContains(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

// Resolve existing ancestors too, so a missing patch or release directory
// cannot escape the project through a symlink.
func repairResolvedPath(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err == nil {
		return resolved, nil
	}
	if !os.IsNotExist(err) {
		return "", err
	}
	if info, err := os.Lstat(absolute); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("repair path %q is a dangling symlink", absolute)
	}
	parent, err := repairResolvedPath(filepath.Dir(absolute))
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, filepath.Base(absolute)), nil
}

func isRepairPatch(name string) bool {
	return strings.HasSuffix(name, ".patch")
}

func isRepairGeneratedFile(name string) bool {
	return name == "CHECKSUMS" || strings.HasSuffix(name, "ATTRIBUTION.txt")
}

func captureRepairFiles(dir string, covered func(string) bool, files map[string]repairFileState) error {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading repair directory %q: %w", dir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !covered(entry.Name()) {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("repair file %q is not a regular file", path)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading repair file %q: %w", path, err)
		}
		files[entry.Name()] = repairFileState{content: content, mode: info.Mode().Perm()}
	}
	return nil
}

func restoreRepairFiles(dir string, covered func(string) bool, files map[string]repairFileState) error {
	entries, err := os.ReadDir(dir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading repair directory %q: %w", dir, err)
	}
	var restoreErrors []error
	blocked := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() || !covered(entry.Name()) {
			continue
		}
		if err := os.Remove(filepath.Join(dir, entry.Name())); err != nil {
			restoreErrors = append(restoreErrors, err)
			blocked[entry.Name()] = true
		}
	}
	if len(files) > 0 {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return errors.Join(append(restoreErrors, err)...)
		}
	}
	for name, file := range files {
		// A failed unlink may leave a replacement symlink. Never write or
		// chmod through that path while restoring the remaining snapshots.
		if blocked[name] {
			continue
		}
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, file.content, file.mode); err != nil {
			restoreErrors = append(restoreErrors, err)
		} else if err := os.Chmod(path, file.mode); err != nil {
			restoreErrors = append(restoreErrors, err)
		}
	}
	return errors.Join(restoreErrors...)
}

func removeEmptyRepairDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return os.Remove(dir)
	}
	return nil
}

func repairUntrackedPaths(repoPath string) ([]string, error) {
	names := make(map[string]struct{})
	// The second listing includes empty directories and directory roots, so
	// new generated trees can be cleaned without deleting original files.
	for _, args := range [][]string{
		{"ls-files", "--others", "-z"},
		{"ls-files", "--others", "--directory", "-z"},
	} {
		output, err := repairGitOutput(context.Background(), repoPath, args...)
		if err != nil {
			return nil, err
		}
		for _, name := range strings.Split(string(output), "\x00") {
			if name != "" {
				names[name] = struct{}{}
			}
		}
	}
	paths := make([]string, 0, len(names))
	for name := range names {
		paths = append(paths, name)
	}
	slices.Sort(paths)
	return paths, nil
}

func repairGitOutput(ctx context.Context, repoPath string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"--literal-pathspecs", "-C", repoPath}, args...)...)
	// Git location/config overrides inherited from a caller must never point a
	// repair command at the surrounding repository.
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "GIT_") {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return nil, fmt.Errorf("git %v in %q: %w\n%s", args, repoPath, err, output)
	}
	return output, nil
}

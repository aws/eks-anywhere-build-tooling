package agentpatchfixer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	maxReadLines     = 400
	maxSearchResults = 200
	maxToolOutput    = 50000
)

type patchMetadata struct {
	AuthorName  string
	AuthorEmail string
	AuthorDate  string
	Subject     string
	Body        string
}

type patchSeed struct {
	RejectFiles []string
	Output      string
}

type workspace struct {
	root          string
	allowedFiles  map[string]bool
	allowedRoots  []string
	originalPatch string
}

func newWorkspace(root string, allowedFiles map[string]bool, allowedRoots []string, originalPatch string) (*workspace, error) {
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, fmt.Errorf("resolving workspace root: %v", err)
	}
	return &workspace{
		root:          resolvedRoot,
		allowedFiles:  allowedFiles,
		allowedRoots:  allowedRoots,
		originalPatch: originalPatch,
	}, nil
}

func (w *workspace) resolve(relativePath string, requireAllowed bool, requireExisting bool) (string, error) {
	if filepath.IsAbs(relativePath) || pathEscapesRoot(relativePath) {
		return "", fmt.Errorf("path escapes workspace: %s", relativePath)
	}
	normalized := filepath.ToSlash(filepath.Clean(relativePath))
	if normalized == "." {
		normalized = ""
	}
	if requireAllowed && !w.editAllowed(normalized) {
		return "", fmt.Errorf("editing %s is not allowed", normalized)
	}

	candidate := filepath.Join(w.root, filepath.FromSlash(normalized))
	checkPath := candidate
	if !requireExisting {
		checkPath = filepath.Dir(candidate)
		for {
			if _, err := os.Lstat(checkPath); err == nil {
				break
			} else if !os.IsNotExist(err) {
				return "", fmt.Errorf("checking %s: %v", relativePath, err)
			}
			parent := filepath.Dir(checkPath)
			if parent == checkPath {
				return "", fmt.Errorf("no existing parent for %s", relativePath)
			}
			checkPath = parent
		}
	}
	resolved, err := filepath.EvalSymlinks(checkPath)
	if err != nil {
		return "", fmt.Errorf("resolving %s: %v", relativePath, err)
	}
	relative, err := filepath.Rel(w.root, resolved)
	if err != nil || pathEscapesRoot(relative) {
		return "", fmt.Errorf("path escapes workspace: %s", relativePath)
	}
	return candidate, nil
}

func (w *workspace) editAllowed(normalized string) bool {
	if w.allowedFiles[normalized] {
		return true
	}
	for _, root := range w.allowedRoots {
		if normalized == root || strings.HasPrefix(normalized, root+"/") {
			return true
		}
	}
	return false
}

func (w *workspace) listFiles(pattern string) (string, error) {
	if pattern == "" {
		pattern = "**/*"
	}
	if filepath.IsAbs(pattern) || pathEscapesRoot(pattern) {
		return "", fmt.Errorf("invalid file pattern: %s", pattern)
	}
	var files []string
	err := filepath.WalkDir(w.root, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(w.root, current)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if entry.IsDir() {
			if relative == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if globMatches(pattern, relative) {
			files = append(files, relative)
			if len(files) >= maxSearchResults {
				return fs.SkipAll
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(files)
	return strings.Join(files, "\n"), nil
}

func (w *workspace) readFile(path string, startLine, endLine int) (string, error) {
	target, err := w.resolve(path, false, true)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}
	return numberedLines(string(content), startLine, endLine), nil
}

func (w *workspace) searchRepository(query, searchPath string) (string, error) {
	pattern, err := regexp.Compile(query)
	if err != nil {
		return "", fmt.Errorf("invalid search expression: %v", err)
	}
	if searchPath == "" {
		searchPath = "."
	}
	target, err := w.resolve(searchPath, false, true)
	if err != nil {
		return "", err
	}

	var matches []string
	err = filepath.WalkDir(target, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Size() > 2*1024*1024 {
			return nil
		}
		content, err := os.ReadFile(current)
		if err != nil || bytes.IndexByte(content, 0) >= 0 {
			return nil
		}
		relative, err := filepath.Rel(w.root, current)
		if err != nil {
			return err
		}
		for index, line := range strings.Split(string(content), "\n") {
			if pattern.MatchString(line) {
				matches = append(matches, fmt.Sprintf("%s:%d:%s", filepath.ToSlash(relative), index+1, line))
				if len(matches) >= maxSearchResults {
					return fs.SkipAll
				}
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return strings.Join(matches, "\n"), nil
}

func (w *workspace) replaceText(path, oldText, newText string) (string, error) {
	target, err := w.resolve(path, true, true)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}
	if count := bytes.Count(content, []byte(oldText)); count != 1 {
		return "", fmt.Errorf("expected exactly one match in %s, found %d", path, count)
	}
	updated := bytes.Replace(content, []byte(oldText), []byte(newText), 1)
	if err := os.WriteFile(target, updated, 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf("updated %s", path), nil
}

func (w *workspace) replaceLines(path string, startLine, endLine int, newText string) (string, error) {
	target, err := w.resolve(path, true, true)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}
	hasTrailingNewline := bytes.HasSuffix(content, []byte("\n"))
	lines := strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	if startLine < 1 || endLine < startLine || endLine > len(lines) {
		return "", fmt.Errorf("invalid line range for %s: %d-%d, file has %d lines", path, startLine, endLine, len(lines))
	}
	replacement := strings.Split(newText, "\n")
	updated := append(append(append([]string{}, lines[:startLine-1]...), replacement...), lines[endLine:]...)
	output := strings.Join(updated, "\n")
	if hasTrailingNewline {
		output += "\n"
	}
	if err := os.WriteFile(target, []byte(output), 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf("replaced lines %d-%d in %s", startLine, endLine, path), nil
}

func (w *workspace) writeFile(path, content string) (string, error) {
	if len(content) > 2*1024*1024 {
		return "", fmt.Errorf("file content exceeds 2 MiB limit")
	}
	target, err := w.resolve(path, true, false)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf("wrote %s", path), nil
}

func (w *workspace) deleteFile(path string) (string, error) {
	target, err := w.resolve(path, true, true)
	if err != nil {
		return "", err
	}
	if err := os.Remove(target); err != nil {
		return "", err
	}
	return fmt.Sprintf("deleted %s", path), nil
}

func (w *workspace) readOriginalPatch(startLine, endLine int) string {
	return numberedLines(w.originalPatch, startLine, endLine)
}

func (w *workspace) showDiff(ctx context.Context) (string, error) {
	output, err := runCommand(ctx, w.root, nil, nil, "git", "diff", "--")
	if err != nil {
		return "", err
	}
	if len(output) > maxToolOutput {
		output = output[:maxToolOutput]
	}
	return output, nil
}

func parsePatchMetadata(ctx context.Context, patchText string) (patchMetadata, error) {
	tempDir, err := os.MkdirTemp("", "patch-fixer-mailinfo-")
	if err != nil {
		return patchMetadata{}, err
	}
	defer os.RemoveAll(tempDir)

	messagePath := filepath.Join(tempDir, "message")
	patchPath := filepath.Join(tempDir, "patch")
	output, err := runCommand(ctx, tempDir, strings.NewReader(patchText), nil, "git", "mailinfo", messagePath, patchPath)
	if err != nil {
		return patchMetadata{}, fmt.Errorf("unable to parse patch metadata: %v", err)
	}
	fields := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		key, value, found := strings.Cut(line, ": ")
		if found {
			fields[key] = value
		}
	}
	body, err := os.ReadFile(messagePath)
	if err != nil {
		return patchMetadata{}, err
	}
	metadata := patchMetadata{
		AuthorName:  fields["Author"],
		AuthorEmail: fields["Email"],
		AuthorDate:  fields["Date"],
		Subject:     fields["Subject"],
		Body:        strings.TrimSpace(string(body)),
	}
	if metadata.AuthorName == "" || metadata.AuthorEmail == "" || metadata.AuthorDate == "" || metadata.Subject == "" {
		return patchMetadata{}, fmt.Errorf("patch metadata is incomplete")
	}
	return metadata, nil
}

func patchFiles(patchText string) (map[string]bool, error) {
	lines := strings.Split(patchText, "\n")
	separator := -1
	for index, line := range lines {
		if line == "---" {
			separator = index
		}
	}
	if separator < 0 {
		return nil, fmt.Errorf("mail patch is missing the format-patch separator")
	}
	diffText := strings.Join(lines[separator+1:], "\n")
	for _, unsupported := range []string{"new file mode 120000", "old mode 120000", "GIT binary patch"} {
		if strings.Contains(diffText, unsupported) {
			return nil, fmt.Errorf("symlink and binary patches are not supported by the agent fixer")
		}
	}
	if strings.Contains(diffText, "rename from ") || strings.Contains(diffText, "rename to ") {
		return nil, fmt.Errorf("rename patches are not supported by the agent fixer")
	}

	files := make(map[string]bool)
	for _, line := range strings.Split(diffText, "\n") {
		if !strings.HasPrefix(line, "diff --git ") {
			continue
		}
		parts := splitGitHeader(line)
		if len(parts) >= 4 {
			files[strings.TrimPrefix(parts[3], "b/")] = true
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("patch contains no diff files")
	}
	return files, nil
}

func prepareWorkspace(ctx context.Context, sourceRepo, parentDir string) (string, error) {
	target := filepath.Join(parentDir, "workspace")
	if _, err := runCommand(ctx, parentDir, nil, nil, "git", "clone", "--shared", "--no-checkout", sourceRepo, target); err != nil {
		return "", err
	}
	for _, args := range [][]string{
		{"checkout", "--detach", "HEAD"},
		{"reset", "--hard", "HEAD"},
		{"clean", "-fdx"},
	} {
		if _, err := runCommand(ctx, target, nil, nil, "git", args...); err != nil {
			return "", err
		}
	}
	return target, nil
}

func seedWorkspace(ctx context.Context, workspacePath, patchPath string) (patchSeed, error) {
	output, err := runCommandAllowFailure(ctx, workspacePath, nil, nil, "git", "apply", "--reject", "--whitespace=nowarn", patchPath)
	var rejects []string
	walkErr := filepath.WalkDir(workspacePath, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".rej") {
			relative, relErr := filepath.Rel(workspacePath, current)
			if relErr != nil {
				return relErr
			}
			rejects = append(rejects, filepath.ToSlash(relative))
		}
		return nil
	})
	if walkErr != nil {
		return patchSeed{}, walkErr
	}
	sort.Strings(rejects)
	if err != nil && len(rejects) == 0 {
		return patchSeed{}, fmt.Errorf("unable to seed workspace from failed patch: %v", err)
	}
	return patchSeed{RejectFiles: rejects, Output: strings.TrimSpace(output)}, nil
}

func removeRejectFiles(workspacePath string, rejectFiles []string) {
	for _, relativePath := range rejectFiles {
		_ = os.Remove(filepath.Join(workspacePath, filepath.FromSlash(relativePath)))
	}
}

func generateCandidate(
	ctx context.Context,
	workspacePath string,
	metadata patchMetadata,
	outputPath string,
	allowedFiles map[string]bool,
	allowedRoots []string,
) ([]string, error) {
	if _, err := runCommand(ctx, workspacePath, nil, nil, "git", "diff", "--check"); err != nil {
		return nil, err
	}
	status, err := runCommand(ctx, workspacePath, nil, nil, "git", "status", "--porcelain")
	if err != nil {
		return nil, err
	}
	var editedFiles []string
	for _, line := range strings.Split(strings.TrimSuffix(status, "\n"), "\n") {
		if len(line) > 3 {
			editedFiles = append(editedFiles, filepath.ToSlash(line[3:]))
		}
	}
	if len(editedFiles) == 0 {
		return nil, nil
	}
	sort.Strings(editedFiles)
	for _, editedFile := range editedFiles {
		allowed := allowedFiles[editedFile]
		for _, root := range allowedRoots {
			allowed = allowed || editedFile == root || strings.HasPrefix(editedFile, root+"/")
		}
		if !allowed {
			return nil, fmt.Errorf("candidate changed file outside allowlist: %s", editedFile)
		}
	}

	if _, err := runCommand(ctx, workspacePath, nil, nil, "git", "add", "--all"); err != nil {
		return nil, err
	}
	message := metadata.Subject
	if metadata.Body != "" {
		message += "\n\n" + metadata.Body
	}
	messagePath := filepath.Join(workspacePath, ".git", "patch-fixer-message")
	if err := os.WriteFile(messagePath, []byte(message+"\n"), 0o600); err != nil {
		return nil, err
	}
	env := map[string]string{
		"GIT_AUTHOR_NAME":     metadata.AuthorName,
		"GIT_AUTHOR_EMAIL":    metadata.AuthorEmail,
		"GIT_AUTHOR_DATE":     metadata.AuthorDate,
		"GIT_COMMITTER_NAME":  "EKS Distro PR Bot",
		"GIT_COMMITTER_EMAIL": "aws-model-rocket-bots+eksdistroprbot@amazon.com",
	}
	if _, err := runCommand(ctx, workspacePath, nil, env, "git", "commit", "--file", messagePath); err != nil {
		return nil, err
	}
	patch, err := runCommand(ctx, workspacePath, nil, nil, "git", "format-patch", "-1", "--stdout", "--no-signature")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(outputPath, []byte(patch), 0o600); err != nil {
		return nil, err
	}
	return editedFiles, nil
}

func runCommand(ctx context.Context, dir string, stdin io.Reader, extraEnv map[string]string, name string, args ...string) (string, error) {
	output, err := runCommandAllowFailure(ctx, dir, stdin, extraEnv, name, args...)
	if err != nil {
		return "", err
	}
	return output, nil
}

func runCommandAllowFailure(ctx context.Context, dir string, stdin io.Reader, extraEnv map[string]string, name string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = dir
	command.Stdin = stdin
	command.Env = os.Environ()
	for key, value := range extraEnv {
		command.Env = append(command.Env, key+"="+value)
	}
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	err := command.Run()
	if err != nil {
		return output.String(), fmt.Errorf("command failed (%s %s): %s", name, strings.Join(args, " "), strings.TrimSpace(output.String()))
	}
	return output.String(), nil
}

func numberedLines(content string, startLine, endLine int) string {
	lines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")
	if startLine < 1 {
		startLine = 1
	}
	if endLine < startLine {
		endLine = startLine
	}
	if endLine > startLine+maxReadLines-1 {
		endLine = startLine + maxReadLines - 1
	}
	if endLine > len(lines) {
		endLine = len(lines)
	}
	if startLine > len(lines) {
		return ""
	}
	var output strings.Builder
	for index := startLine - 1; index < endLine; index++ {
		if output.Len() > 0 {
			output.WriteByte('\n')
		}
		fmt.Fprintf(&output, "%d: %s", index+1, lines[index])
	}
	return output.String()
}

func pathEscapesRoot(path string) bool {
	cleaned := filepath.Clean(path)
	return cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator))
}

func globMatches(pattern, relative string) bool {
	pattern = filepath.ToSlash(pattern)
	if pattern == "**/*" || pattern == "**" {
		return true
	}
	matched, _ := filepath.Match(filepath.FromSlash(pattern), filepath.FromSlash(relative))
	if matched {
		return true
	}
	if strings.HasPrefix(pattern, "**/") {
		matched, _ = filepath.Match(filepath.FromSlash(strings.TrimPrefix(pattern, "**/")), filepath.Base(relative))
	}
	return matched
}

func splitGitHeader(line string) []string {
	var parts []string
	for len(line) > 0 {
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if line[0] != '"' {
			if index := strings.IndexByte(line, ' '); index >= 0 {
				parts = append(parts, line[:index])
				line = line[index+1:]
			} else {
				parts = append(parts, line)
				break
			}
			continue
		}
		closed := false
		for index := 1; index < len(line); index++ {
			if line[index] == '"' && line[index-1] != '\\' {
				value, err := strconv.Unquote(line[:index+1])
				if err == nil {
					parts = append(parts, value)
				}
				line = line[index+1:]
				closed = true
				break
			}
		}
		if !closed {
			break
		}
	}
	return parts
}

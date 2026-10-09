package upgrade

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestPatchRepairStateRestoresRepairFilesAndRetainsVersionMetadata(t *testing.T) {
	projectRoot, repoPath, patchPath, _ := patchCandidateFixture(t)
	patchesDir := filepath.Dir(patchPath)
	releaseBranch := "1-99"
	releaseDir := filepath.Join(projectRoot, releaseBranch)
	if err := os.Mkdir(releaseDir, 0o755); err != nil {
		t.Fatal(err)
	}
	originals := map[string]repairFileState{
		patchPath: {content: []byte("first patch\x00\n"), mode: 0o640},
		filepath.Join(patchesDir, "0002-deleted.patch"):          {content: []byte("second patch\n"), mode: 0o600},
		filepath.Join(projectRoot, "CHECKSUMS"):                  {content: []byte("dummy root checksum\n"), mode: 0o600},
		filepath.Join(projectRoot, "ATTRIBUTION.txt"):            {content: []byte("dummy root attribution\n"), mode: 0o640},
		filepath.Join(releaseDir, "CHECKSUMS"):                   {content: []byte("dummy release checksum\n"), mode: 0o640},
		filepath.Join(releaseDir, "THIRD_PARTY_ATTRIBUTION.txt"): {content: []byte("dummy release attribution\n"), mode: 0o600},
	}
	for path, file := range originals {
		writeRepairTestFile(t, path, file.content, file.mode)
	}
	if err := os.Chmod(patchesDir, 0o750); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"GIT_TAG", "GOLANG_VERSION", "README.md"} {
		writeRepairTestFile(t, filepath.Join(projectRoot, name), []byte("original version\n"), 0o644)
		writeRepairTestFile(t, filepath.Join(releaseDir, name), []byte("original version\n"), 0o644)
	}
	runTestCommand(t, projectRoot, "git", "init", "-q")
	runTestCommand(t, projectRoot, "git", "config", "user.name", "Test")
	runTestCommand(t, projectRoot, "git", "config", "user.email", "test@example.com")
	runTestCommand(t, projectRoot, "git", "add", "GIT_TAG", "GOLANG_VERSION", "README.md")
	runTestCommand(t, projectRoot, "git", "commit", "-qm", "outer base")
	outerHead := runTestCommand(t, projectRoot, "git", "rev-parse", "HEAD")
	state, err := captureRepairState(projectRoot, repoPath, patchesDir, releaseBranch, 0)
	if err != nil {
		t.Fatalf("captureRepairState() error = %v", err)
	}

	// Replace one patch and delete another; generation can change the entire set.
	writeRepairTestFile(t, patchPath, []byte("repaired patch\n"), 0o644)
	if err := os.Remove(filepath.Join(patchesDir, "0002-deleted.patch")); err != nil {
		t.Fatal(err)
	}
	newPaths := []string{
		filepath.Join(patchesDir, "0003-generated.patch"),
		filepath.Join(projectRoot, "GO_ATTRIBUTION.txt"),
		filepath.Join(releaseDir, "NEW_ATTRIBUTION.txt"),
	}
	for _, path := range newPaths {
		writeRepairTestFile(t, path, []byte("generated\n"), 0o644)
	}
	for path := range originals {
		if strings.HasSuffix(path, ".patch") {
			continue
		}
		writeRepairTestFile(t, path, []byte("changed generated output\n"), 0o644)
	}
	if err := os.Remove(filepath.Join(projectRoot, "ATTRIBUTION.txt")); err != nil {
		t.Fatal(err)
	}
	retained := []string{
		filepath.Join(patchesDir, "README.md"),
		filepath.Join(projectRoot, "CHECKSUMS.backup"),
		filepath.Join(projectRoot, "notes.txt"),
		filepath.Join(projectRoot, "other-release", "CHECKSUMS"),
		filepath.Join(releaseDir, "nested", "ATTRIBUTION.txt"),
	}
	for _, dir := range []string{filepath.Join(projectRoot, "other-release"), filepath.Join(releaseDir, "nested")} {
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range retained {
		writeRepairTestFile(t, path, []byte("retain this\n"), 0o644)
	}
	for _, dir := range []string{projectRoot, releaseDir} {
		for _, name := range []string{"GIT_TAG", "GOLANG_VERSION", "README.md"} {
			path := filepath.Join(dir, name)
			writeRepairTestFile(t, path, []byte("new version\n"), 0o644)
		}
	}

	if err := state.restore(); err != nil {
		t.Fatalf("restore() error = %v", err)
	}
	for path, file := range originals {
		assertRepairTestFile(t, path, file.content, file.mode)
	}
	for _, path := range newPaths {
		assertRepairTestMissing(t, path)
	}
	for _, path := range retained {
		assertRepairTestFile(t, path, []byte("retain this\n"), 0o644)
	}
	for _, dir := range []string{projectRoot, releaseDir} {
		for _, name := range []string{"GIT_TAG", "GOLANG_VERSION", "README.md"} {
			assertRepairTestFile(t, filepath.Join(dir, name), []byte("new version\n"), 0o644)
		}
	}
	if head := runTestCommand(t, projectRoot, "git", "rev-parse", "HEAD"); head != outerHead {
		t.Fatalf("outer HEAD = %q, want %q", head, outerHead)
	}
	info, err := os.Stat(patchesDir)
	if err != nil || info.Mode().Perm() != 0o750 {
		t.Fatalf("patch directory mode was not restored: info %v, error %v", info, err)
	}
	if err := state.restore(); err != nil {
		t.Fatalf("second restore() error = %v", err)
	}
}

func TestPatchRepairStateRestoresPatchDirectoryExistence(t *testing.T) {
	for _, originallyExists := range []bool{false, true} {
		name := "missing"
		if originallyExists {
			name = "deleted"
		}
		t.Run(name, func(t *testing.T) {
			projectRoot, repoPath, patchPath, _ := patchCandidateFixture(t)
			patchesDir := filepath.Dir(patchPath)
			if !originallyExists {
				if err := os.RemoveAll(patchesDir); err != nil {
					t.Fatal(err)
				}
			}
			state, err := captureRepairState(projectRoot, repoPath, patchesDir, "1-99", 0)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.RemoveAll(patchesDir); err != nil {
				t.Fatal(err)
			}
			if !originallyExists {
				if err := os.Mkdir(patchesDir, 0o755); err != nil {
					t.Fatal(err)
				}
				writeRepairTestFile(t, filepath.Join(patchesDir, "new.patch"), []byte("dummy patch\n"), 0o644)
			}
			releaseDir := filepath.Join(projectRoot, "1-99")
			if err := os.Mkdir(releaseDir, 0o755); err != nil {
				t.Fatal(err)
			}
			newChecksum := filepath.Join(releaseDir, "CHECKSUMS")
			writeRepairTestFile(t, newChecksum, []byte("dummy generated checksum\n"), 0o644)
			if err := state.restore(); err != nil {
				t.Fatal(err)
			}
			assertRepairTestMissing(t, newChecksum)
			if originallyExists {
				assertRepairTestFile(t, patchPath, []byte("original-invalid-patch\n"), 0o644)
			} else {
				assertRepairTestMissing(t, patchesDir)
			}
		})
	}
}

func TestPatchRepairStateCapturesCleanPrefixAndAbortsGitAm(t *testing.T) {
	projectRoot, repoPath, patchesDir, candidates := multiPatchFixture(t)
	prefix, err := os.ReadFile(candidates["0001-stale.patch"])
	if err != nil {
		t.Fatal(err)
	}
	firstPatch := filepath.Join(patchesDir, "0001-stale.patch")
	writeRepairTestFile(t, firstPatch, prefix, 0o644)
	base := strings.TrimSpace(runTestCommand(t, repoPath, "git", "rev-parse", "HEAD"))
	cmd := exec.Command("git", "am", "--", firstPatch, filepath.Join(patchesDir, "0002-stale.patch"))
	cmd.Dir = repoPath
	if output, err := cmd.CombinedOutput(); err == nil || !strings.Contains(string(output), "Patch failed at 0002") {
		t.Fatalf("git am = error %v, output %q; want failure after clean prefix", err, output)
	}
	originalUntracked := filepath.Join(repoPath, "existing.txt")
	writeRepairTestFile(t, originalUntracked, []byte("retain original source\n"), 0o644)
	originalDir := filepath.Join(repoPath, "existing-output")
	if err := os.Mkdir(originalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeRepairTestFile(t, filepath.Join(originalDir, "keep.generated"), []byte("retain ignored source\n"), 0o644)
	writeRepairTestFile(t, filepath.Join(repoPath, ".git", "info", "exclude"), []byte("*.generated\n"), 0o644)
	state, err := captureRepairState(projectRoot, repoPath, patchesDir, "", 1)
	if err != nil {
		t.Fatal(err)
	}
	if state.baseRevision != base {
		t.Fatalf("baseRevision = %q, want %q", state.baseRevision, base)
	}
	writeRepairTestFile(t, filepath.Join(repoPath, "one.txt"), []byte("dirty prefix\n"), 0o644)
	runTestCommand(t, repoPath, "git", "add", "one.txt")
	writeRepairTestFile(t, filepath.Join(repoPath, "two.txt"), []byte("dirty source\n"), 0o644)
	generated := filepath.Join(repoPath, "output.generated")
	writeRepairTestFile(t, generated, []byte("dummy ignored output\n"), 0o644)
	untracked := filepath.Join(repoPath, "new source.txt")
	writeRepairTestFile(t, untracked, []byte("dummy untracked output\n"), 0o644)
	newDir := filepath.Join(repoPath, "new-output")
	if err := os.Mkdir(newDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeRepairTestFile(t, filepath.Join(newDir, "build.generated"), []byte("dummy generated source\n"), 0o644)
	emptyDir := filepath.Join(repoPath, "empty-output")
	if err := os.Mkdir(emptyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	newInOriginalDir := filepath.Join(originalDir, "new.generated")
	writeRepairTestFile(t, newInOriginalDir, []byte("dummy new output\n"), 0o644)
	literalPath := filepath.Join(repoPath, "[existing].txt")
	writeRepairTestFile(t, literalPath, []byte("dummy unusual filename\n"), 0o644)
	if err := state.restore(); err != nil {
		t.Fatal(err)
	}
	assertRepairGitFile(t, filepath.Join(repoPath, "one.txt"), []byte("current-one\n"), 0o644)
	assertRepairGitFile(t, filepath.Join(repoPath, "two.txt"), []byte("current-two\n"), 0o644)
	assertRepairTestFile(t, originalUntracked, []byte("retain original source\n"), 0o644)
	assertRepairTestFile(t, filepath.Join(originalDir, "keep.generated"), []byte("retain ignored source\n"), 0o644)
	assertRepairTestMissing(t, generated)
	assertRepairTestMissing(t, untracked)
	assertRepairTestMissing(t, newDir)
	assertRepairTestMissing(t, emptyDir)
	assertRepairTestMissing(t, newInOriginalDir)
	assertRepairTestMissing(t, literalPath)
	assertRepairTestMissing(t, filepath.Join(repoPath, ".git", "rebase-apply"))
	if head := strings.TrimSpace(runTestCommand(t, repoPath, "git", "rev-parse", "HEAD")); head != base {
		t.Fatalf("restored HEAD = %q, want %q", head, base)
	}
	if status := runTestCommand(t, repoPath, "git", "status", "--porcelain", "--untracked-files=no"); status != "" {
		t.Fatalf("upstream tracked files remain dirty: %s", status)
	}
}

func TestPatchRepairStateResetsCompletedRepair(t *testing.T) {
	projectRoot, repoPath, patchPath, candidatePath := patchCandidateFixture(t)
	state, err := captureRepairState(projectRoot, repoPath, filepath.Dir(patchPath), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, repoPath, "git", "am", candidatePath)
	if err := state.restore(); err != nil {
		t.Fatal(err)
	}
	assertRepairGitFile(t, filepath.Join(repoPath, "main.txt"), []byte("old\n"), 0o644)
	if head := strings.TrimSpace(runTestCommand(t, repoPath, "git", "rev-parse", "HEAD")); head != state.baseRevision {
		t.Fatalf("restored HEAD = %q, want %q", head, state.baseRevision)
	}
}

func TestPatchRepairStateRejectsUnsafeRepoPath(t *testing.T) {
	for _, name := range []string{"project root", "outside", "outer repository", "subdirectory", "symlink escape", "shared git directory", "shared common directory"} {
		t.Run(name, func(t *testing.T) {
			projectRoot, repoPath, patchPath, _ := patchCandidateFixture(t)
			runTestCommand(t, projectRoot, "git", "init", "-q")
			unsafePath := repoPath
			switch name {
			case "project root":
				unsafePath = projectRoot
			case "outside":
				unsafePath = t.TempDir()
			case "outer repository":
				unsafePath = filepath.Join(projectRoot, "inherits-outer-repo")
				if err := os.Mkdir(unsafePath, 0o755); err != nil {
					t.Fatal(err)
				}
			case "subdirectory":
				unsafePath = filepath.Join(repoPath, "subdir")
				if err := os.Mkdir(unsafePath, 0o755); err != nil {
					t.Fatal(err)
				}
			case "symlink escape":
				unsafePath = filepath.Join(projectRoot, "outside-link")
				if err := os.Symlink(t.TempDir(), unsafePath); err != nil {
					t.Fatal(err)
				}
			case "shared git directory":
				gitDir := filepath.Join(repoPath, ".git")
				if err := os.RemoveAll(gitDir); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(filepath.Join(projectRoot, ".git"), gitDir); err != nil {
					t.Fatal(err)
				}
			case "shared common directory":
				gitDir := filepath.Join(repoPath, ".git")
				commonDir := filepath.Join(projectRoot, "common-git")
				head, err := os.ReadFile(filepath.Join(gitDir, "HEAD"))
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(gitDir, commonDir); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(gitDir, 0o755); err != nil {
					t.Fatal(err)
				}
				writeRepairTestFile(t, filepath.Join(gitDir, "HEAD"), head, 0o644)
				writeRepairTestFile(t, filepath.Join(gitDir, "commondir"), []byte("../../common-git\n"), 0o644)
			}
			if _, err := captureRepairState(projectRoot, unsafePath, filepath.Dir(patchPath), "", 0); err == nil {
				t.Fatal("captureRepairState() error = nil, want unsafe path rejection")
			}
			assertRepairTestFile(t, patchPath, []byte("original-invalid-patch\n"), 0o644)
		})
	}
}

func TestPatchRepairStateRejectsUnsafeOutputPathsAndCounts(t *testing.T) {
	projectRoot, repoPath, patchPath, _ := patchCandidateFixture(t)
	for _, test := range []struct {
		name    string
		patches string
		release string
		count   int
	}{
		{name: "outside patches", patches: t.TempDir()},
		{name: "root patches", patches: projectRoot},
		{name: "upstream patches", patches: repoPath},
		{name: "outside release", patches: filepath.Dir(patchPath), release: "../outside"},
		{name: "upstream release", patches: filepath.Dir(patchPath), release: "repo"},
		{name: "negative count", patches: filepath.Dir(patchPath), count: -1},
		{name: "missing baseline", patches: filepath.Dir(patchPath), count: 99},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := captureRepairState(projectRoot, repoPath, test.patches, test.release, test.count); err == nil {
				t.Fatal("captureRepairState() error = nil, want rejection")
			}
		})
	}
}

func TestPatchRepairStateRejectsChangedRepoPathOnRestore(t *testing.T) {
	projectRoot, repoPath, patchPath, _ := patchCandidateFixture(t)
	state, err := captureRepairState(projectRoot, repoPath, filepath.Dir(patchPath), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(repoPath, filepath.Join(projectRoot, "saved-repo")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), repoPath); err != nil {
		t.Fatal(err)
	}
	writeRepairTestFile(t, patchPath, []byte("still pending repair\n"), 0o644)
	if err := state.restore(); err == nil {
		t.Fatal("restore() error = nil, want unsafe path rejection")
	}
	assertRepairTestFile(t, patchPath, []byte("still pending repair\n"), 0o644)
}

func TestPatchRepairStateIgnoresInheritedGitLocation(t *testing.T) {
	projectRoot, repoPath, patchPath, _ := patchCandidateFixture(t)
	base := strings.TrimSpace(runTestCommand(t, repoPath, "git", "rev-parse", "HEAD"))
	runTestCommand(t, projectRoot, "git", "init", "-q")
	t.Setenv("GIT_DIR", filepath.Join(projectRoot, ".git"))
	t.Setenv("GIT_WORK_TREE", projectRoot)
	state, err := captureRepairState(projectRoot, repoPath, filepath.Dir(patchPath), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if state.baseRevision != base {
		t.Fatalf("baseRevision = %q, want upstream revision %q", state.baseRevision, base)
	}
	writeRepairTestFile(t, filepath.Join(repoPath, "main.txt"), []byte("dirty source\n"), 0o644)
	if err := state.restore(); err != nil {
		t.Fatal(err)
	}
	assertRepairGitFile(t, filepath.Join(repoPath, "main.txt"), []byte("old\n"), 0o644)
}

func TestPatchRepairStateRestoresFilesWithoutFollowingReplacementSymlinks(t *testing.T) {
	projectRoot, repoPath, patchPath, _ := patchCandidateFixture(t)
	checksumPath := filepath.Join(projectRoot, "CHECKSUMS")
	writeRepairTestFile(t, checksumPath, []byte("original dummy checksum\n"), 0o600)
	state, err := captureRepairState(projectRoot, repoPath, filepath.Dir(patchPath), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "unrelated.txt")
	writeRepairTestFile(t, target, []byte("retain unrelated bytes\n"), 0o640)
	for _, path := range []string{patchPath, checksumPath} {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
	}
	if err := state.restore(); err != nil {
		t.Fatal(err)
	}
	assertRepairTestFile(t, patchPath, []byte("original-invalid-patch\n"), 0o644)
	assertRepairTestFile(t, checksumPath, []byte("original dummy checksum\n"), 0o600)
	assertRepairTestFile(t, target, []byte("retain unrelated bytes\n"), 0o640)
}

func TestRestoreRepairFilesDoesNotFollowSymlinkWhenRemovalFails(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(t.TempDir(), "unrelated.txt")
	writeRepairTestFile(t, target, []byte("retain unrelated bytes\n"), 0o640)
	checksumPath := filepath.Join(dir, "CHECKSUMS")
	if err := os.Symlink(target, checksumPath); err != nil {
		t.Fatal(err)
	}
	probe := filepath.Join(dir, "remove-probe")
	writeRepairTestFile(t, probe, []byte("dummy\n"), 0o600)
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(dir, 0o755); err != nil {
			t.Error(err)
		}
	})
	// Root can unlink in a read-only directory, so confirm this host can
	// exercise the failure before asserting the rollback behavior.
	if err := os.Remove(probe); err == nil {
		t.Skip("host permits unlink in a read-only directory")
	}
	err := restoreRepairFiles(dir, isRepairGeneratedFile, map[string]repairFileState{
		"CHECKSUMS": {content: []byte("original dummy checksum\n"), mode: 0o600},
	})
	if !errors.Is(err, os.ErrPermission) {
		t.Fatalf("restoreRepairFiles() error = %v, want removal permission error", err)
	}
	assertRepairTestFile(t, target, []byte("retain unrelated bytes\n"), 0o640)
	info, err := os.Lstat(checksumPath)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("blocked checksum path changed: info %v, error %v", info, err)
	}
}

func TestPatchRepairStateDockerfilesUsesBaselineAndUncommittedChanges(t *testing.T) {
	projectRoot, repoPath, patchPath, _ := patchCandidateFixture(t)
	for _, path := range []string{"Dockerfile", "Dockerfile.deleted", "Dockerfile.unchanged", "Dockerfileish", "README.Dockerfile.txt"} {
		writeRepairTestFile(t, filepath.Join(repoPath, path), []byte("original\n"), 0o644)
	}
	runTestCommand(t, repoPath, "git", "add", ".")
	runTestCommand(t, repoPath, "git", "commit", "-qm", "docker base")
	writeRepairTestFile(t, filepath.Join(repoPath, "Dockerfile"), []byte("prefix change\n"), 0o644)
	runTestCommand(t, repoPath, "git", "add", "Dockerfile")
	runTestCommand(t, repoPath, "git", "commit", "-qm", "clean prefix")
	state, err := captureRepairState(projectRoot, repoPath, filepath.Dir(patchPath), "", 1)
	if err != nil {
		t.Fatal(err)
	}
	committed := []string{"Dockerfile.dev", "nested/service.DoCkErFiLe", "nested/dOcKeRfIlE"}
	if err := os.Mkdir(filepath.Join(repoPath, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, path := range committed {
		writeRepairTestFile(t, filepath.Join(repoPath, path), []byte("committed change\n"), 0o644)
	}
	runTestCommand(t, repoPath, "git", "add", ".")
	runTestCommand(t, repoPath, "git", "commit", "-qm", "repair")
	writeRepairTestFile(t, filepath.Join(repoPath, "Dockerfile"), []byte("unstaged change\n"), 0o644)
	staged := "nested/Dockerfile.with space\nand newline"
	writeRepairTestFile(t, filepath.Join(repoPath, staged), []byte("staged change\n"), 0o644)
	runTestCommand(t, repoPath, "git", "add", "--", staged)
	if err := os.Remove(filepath.Join(repoPath, "Dockerfile.deleted")); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"Dockerfileish", "README.Dockerfile.txt", "Dockerfile.untracked"} {
		writeRepairTestFile(t, filepath.Join(repoPath, path), []byte("other change\n"), 0o644)
	}
	files, err := state.dockerfiles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := append(committed, "Dockerfile", "Dockerfile.deleted", staged)
	slices.Sort(want)
	if !slices.Equal(files, want) {
		t.Fatalf("dockerfiles() = %#v, want %#v", files, want)
	}
}

func TestPatchRepairStateDockerfilesHonorsContextAndEmptyDiff(t *testing.T) {
	projectRoot, repoPath, patchPath, _ := patchCandidateFixture(t)
	state, err := captureRepairState(projectRoot, repoPath, filepath.Dir(patchPath), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	files, err := state.dockerfiles(context.Background())
	if err != nil || len(files) != 0 {
		t.Fatalf("dockerfiles() = %#v, error %v; want empty diff", files, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := state.dockerfiles(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("dockerfiles() error = %v, want context canceled", err)
	}
}

func writeRepairTestFile(t *testing.T, path string, content []byte, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, content, mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}

func assertRepairTestFile(t *testing.T, path string, content []byte, mode os.FileMode) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Fatalf("%s content = %q, want %q", path, got, content)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != mode {
		t.Fatalf("%s mode = %v, want %v", path, got, mode)
	}
}

func assertRepairGitFile(t *testing.T, path string, content []byte, mode os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	// Git tracks executable bits; checkout read/write bits follow the host umask.
	expected := info.Mode().Perm()&0o666 | mode&0o111
	assertRepairTestFile(t, path, content, expected)
}

func assertRepairTestMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("%s remains after restore: %v", path, err)
	}
}

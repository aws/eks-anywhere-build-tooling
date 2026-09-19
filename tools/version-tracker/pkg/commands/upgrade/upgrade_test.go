package upgrade

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/commands/agentpatchfixer"
)

func TestFailedPatchPathUsesAppliedCountAsNextPatchIndex(t *testing.T) {
	patchesDir := t.TempDir()
	for _, name := range []string{"0002-second.patch", "README.md", "0001-first.patch", "0003-third.patch"} {
		if err := os.WriteFile(filepath.Join(patchesDir, name), []byte("test"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := failedPatchPath(patchesDir, 1)
	if err != nil {
		t.Fatalf("failedPatchPath() error = %v", err)
	}
	if want := filepath.Join(patchesDir, "0002-second.patch"); got != want {
		t.Fatalf("failedPatchPath() = %q, want %q", got, want)
	}
}

func TestFailedPatchPathRejectsOutOfRangeIndex(t *testing.T) {
	patchesDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(patchesDir, "0001-first.patch"), []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := failedPatchPath(patchesDir, 1); err == nil {
		t.Fatal("failedPatchPath() error = nil, want out-of-range error")
	}
}

func TestPatchFilesInDirectoryIgnoresNonPatchEntries(t *testing.T) {
	patchesDir := t.TempDir()
	for _, name := range []string{"0001-first.patch", "README.md"} {
		if err := os.WriteFile(filepath.Join(patchesDir, name), []byte("test"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(patchesDir, "archive.patch"), 0o755); err != nil {
		t.Fatal(err)
	}

	patches, err := patchFilesInDirectory(patchesDir)
	if err != nil {
		t.Fatalf("patchFilesInDirectory() error = %v", err)
	}
	if len(patches) != 1 || patches[0] != filepath.Join(patchesDir, "0001-first.patch") {
		t.Fatalf("patchFilesInDirectory() = %#v", patches)
	}
}

func TestIsPatchConflictOutput(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   bool
	}{
		{name: "does not apply", output: "error: main.go: patch does not apply", want: true},
		{name: "missing file", output: "error: main.go: does not exist in index", want: true},
		{name: "infrastructure failure", output: "fatal: unable to access repository", want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isPatchConflictOutput(test.output); got != test.want {
				t.Fatalf("isPatchConflictOutput() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestPatchesAppliedFromOutput(t *testing.T) {
	if got := patchesAppliedFromOutput("Patch failed at 0003 subject"); got != 2 {
		t.Fatalf("patchesAppliedFromOutput() = %d, want 2", got)
	}
	if got := patchesAppliedFromOutput("patch does not apply"); got != 0 {
		t.Fatalf("patchesAppliedFromOutput() = %d, want 0", got)
	}
}

func TestShouldRunUpgradeInBuildUsesCanonicalArchitecture(t *testing.T) {
	t.Setenv("CODEBUILD_BATCH_BUILD_IDENTIFIER", "linuxkit_linuxkit_linux_arm64")
	if shouldRunUpgradeInBuild() {
		t.Fatal("shouldRunUpgradeInBuild() = true for paired arm64 build")
	}
	t.Setenv("CODEBUILD_BATCH_BUILD_IDENTIFIER", "kubernetes_autoscaler_1_36")
	if !shouldRunUpgradeInBuild() {
		t.Fatal("shouldRunUpgradeInBuild() = false for canonical autoscaler build")
	}
}

func TestFormatPatchFailureDetails(t *testing.T) {
	got := formatPatchFailureDetails(
		2,
		4,
		"Patch failed at `0003-third.patch`",
		[]string{"one.go", "two.go"},
	)
	for _, want := range []string{
		"Only 2/4 patches were applied",
		"0003-third.patch",
		"`one.go`,`two.go`",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("formatPatchFailureDetails() = %q, want substring %q", got, want)
		}
	}
}

func TestFormatAgentPatchReviewComment(t *testing.T) {
	got := formatAgentPatchReviewComment([]string{
		"projects/example/patches/0002-second.patch",
		"projects/example/patches/0001-first.patch",
		"projects/example/patches/0002-second.patch",
	})
	want := "## AI-assisted patch repair\n\n" +
		"The patch fixer repaired the following patches using AI assistance:\n" +
		"- `0001-first.patch`\n" +
		"- `0002-second.patch`\n\n" +
		"Human review is required before merge to confirm the repaired patches preserve their original intent."
	if got != want {
		t.Fatalf("formatAgentPatchReviewComment() = %q, want %q", got, want)
	}
}

func TestRepairGenericPatchSeriesKeepsValidatedCandidate(t *testing.T) {
	projectRoot, repoPath, failedPatchPath, candidatePath := patchCandidateFixture(t)
	applied, _, failedFiles, failureOutput, applyErr := applyPatchesToRepo(projectRoot, "repo", 1)
	if applyErr == nil {
		t.Fatal("initial patch application unexpectedly succeeded")
	}
	configureFakeAgentRunner(t, candidatePath, "candidate_generated")

	fixed, changed, err := repairGenericPatchSeries(
		context.Background(),
		"example/project",
		projectRoot,
		"repo",
		filepath.Dir(failedPatchPath),
		"v1.0.0",
		"v1.1.0",
		"",
		1,
		applied,
		failedFiles,
		failureOutput,
	)
	if err != nil {
		t.Fatalf("repairGenericPatchSeries() error = %v", err)
	}
	if !fixed {
		t.Fatal("repairGenericPatchSeries() fixed = false, want true")
	}
	if len(changed) != 1 || changed[0] != failedPatchPath {
		t.Fatalf("changed paths = %#v", changed)
	}
	candidate, err := os.ReadFile(candidatePath)
	if err != nil {
		t.Fatal(err)
	}
	stored, err := os.ReadFile(failedPatchPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(stored) != string(candidate) {
		t.Fatal("validated candidate was not retained")
	}
	content, err := os.ReadFile(filepath.Join(repoPath, "main.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "new\n" {
		t.Fatalf("patched content = %q, want new", content)
	}
}

func TestRepairGenericPatchSeriesRestoresRejectedCandidate(t *testing.T) {
	projectRoot, repoPath, failedPatchPath, _ := patchCandidateFixture(t)
	original := []byte("original-invalid-patch\n")
	if err := os.WriteFile(failedPatchPath, original, 0o644); err != nil {
		t.Fatal(err)
	}
	applied, _, failedFiles, failureOutput, applyErr := applyPatchesToRepo(projectRoot, "repo", 1)
	if applyErr == nil {
		t.Fatal("initial patch application unexpectedly succeeded")
	}
	configureFakeAgentRunner(t, "", "agent_incomplete")
	fixed, _, err := repairGenericPatchSeries(
		context.Background(),
		"example/project",
		projectRoot,
		"repo",
		filepath.Dir(failedPatchPath),
		"v1.0.0",
		"v1.1.0",
		"",
		1,
		applied,
		failedFiles,
		failureOutput,
	)
	if err == nil {
		t.Fatal("repairGenericPatchSeries() error = nil, want rejection")
	}
	var repairFailure *patchRepairFailure
	if !errors.As(err, &repairFailure) {
		t.Fatalf("repairGenericPatchSeries() error = %T, want *patchRepairFailure", err)
	}
	if repairFailure.failedPatch != "Patch failed at `0001-change.patch`" {
		t.Fatalf("repair failure patch = %q", repairFailure.failedPatch)
	}
	if fixed {
		t.Fatal("applyAgentPatchCandidate() fixed = true, want false")
	}
	restored, readErr := os.ReadFile(failedPatchPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(restored) != string(original) {
		t.Fatalf("restored patch = %q, want %q", restored, original)
	}
	if status := runTestCommand(t, repoPath, "git", "status", "--porcelain"); status != "" {
		t.Fatalf("repository is dirty after rollback: %s", status)
	}
	if _, err := os.Stat(filepath.Join(repoPath, ".git", "rebase-apply")); !os.IsNotExist(err) {
		t.Fatalf("git am state remains after rollback: %v", err)
	}
}

func TestRepairGenericPatchSeriesRepairsMultipleConflicts(t *testing.T) {
	projectRoot, repoPath, patchesDir, candidates := multiPatchFixture(t)
	applied, _, failedFiles, failureOutput, applyErr := applyPatchesToRepo(projectRoot, "repo", 2)
	if applyErr == nil {
		t.Fatal("initial patch application unexpectedly succeeded")
	}
	configureMultiPatchAgentRunner(t, candidates)

	fixed, changed, err := repairGenericPatchSeries(
		context.Background(),
		"example/project",
		projectRoot,
		"repo",
		patchesDir,
		"v1.0.0",
		"v1.1.0",
		"",
		2,
		applied,
		failedFiles,
		failureOutput,
	)
	if err != nil {
		t.Fatalf("repairGenericPatchSeries() error = %v", err)
	}
	if !fixed || len(changed) != 2 {
		t.Fatalf("repairGenericPatchSeries() = fixed %t, changed %#v", fixed, changed)
	}
	if content, err := os.ReadFile(filepath.Join(repoPath, "one.txt")); err != nil || string(content) != "fixed-one\n" {
		t.Fatalf("one.txt = %q, error %v", content, err)
	}
	if content, err := os.ReadFile(filepath.Join(repoPath, "two.txt")); err != nil || string(content) != "fixed-two\n" {
		t.Fatalf("two.txt = %q, error %v", content, err)
	}
}

func TestRepairGenericPatchSeriesReportsTerminalConflict(t *testing.T) {
	projectRoot, _, patchesDir, candidates := multiPatchFixture(t)
	applied, _, failedFiles, failureOutput, applyErr := applyPatchesToRepo(projectRoot, "repo", 2)
	if applyErr == nil {
		t.Fatal("initial patch application unexpectedly succeeded")
	}
	configurePartialPatchAgentRunner(t, candidates["0001-stale.patch"])

	fixed, _, err := repairGenericPatchSeries(
		context.Background(),
		"example/project",
		projectRoot,
		"repo",
		patchesDir,
		"v1.0.0",
		"v1.1.0",
		"",
		2,
		applied,
		failedFiles,
		failureOutput,
	)
	if err == nil {
		t.Fatal("repairGenericPatchSeries() error = nil, want terminal conflict")
	}
	if fixed {
		t.Fatal("repairGenericPatchSeries() fixed = true, want false")
	}
	var repairFailure *patchRepairFailure
	if !errors.As(err, &repairFailure) {
		t.Fatalf("repairGenericPatchSeries() error = %T, want *patchRepairFailure", err)
	}
	if repairFailure.appliedPatches != 1 || repairFailure.totalPatches != 2 {
		t.Fatalf(
			"repair failure progress = %d/%d, want 1/2",
			repairFailure.appliedPatches,
			repairFailure.totalPatches,
		)
	}
	if repairFailure.failedPatch != "Patch failed at `0002-stale.patch`" {
		t.Fatalf("repair failure patch = %q", repairFailure.failedPatch)
	}
}

func TestRepairGenericPatchSeriesRemovesObsoletePatch(t *testing.T) {
	projectRoot, repoPath, failedPatchPath, _ := patchCandidateFixture(t)
	patchesDir := filepath.Dir(failedPatchPath)
	applied, _, failedFiles, failureOutput, applyErr := applyPatchesToRepo(projectRoot, "repo", 1)
	if applyErr == nil {
		t.Fatal("initial patch application unexpectedly succeeded")
	}
	configureFakeAgentRunner(t, "", "no_changes")

	fixed, changed, err := repairGenericPatchSeries(
		context.Background(),
		"example/project",
		projectRoot,
		"repo",
		patchesDir,
		"v1.0.0",
		"v1.1.0",
		"",
		1,
		applied,
		failedFiles,
		failureOutput,
	)
	if err != nil {
		t.Fatalf("repairGenericPatchSeries() error = %v", err)
	}
	if !fixed || len(changed) != 1 || changed[0] != failedPatchPath {
		t.Fatalf("repairGenericPatchSeries() = fixed %t, changed %#v", fixed, changed)
	}
	if _, err := os.Stat(patchesDir); !os.IsNotExist(err) {
		t.Fatalf("obsolete patch directory still exists: %v", err)
	}
	if status := runTestCommand(t, repoPath, "git", "status", "--porcelain"); status != "" {
		t.Fatalf("repository is dirty after obsolete patch removal: %s", status)
	}
}

func configureFakeAgentRunner(t *testing.T, candidatePath, status string) {
	t.Helper()
	t.Setenv(agentpatchfixer.EnabledEnv, "true")
	original := runAgentPatchFixer
	runAgentPatchFixer = func(context.Context, agentpatchfixer.Request) (*agentpatchfixer.Result, error) {
		return &agentpatchfixer.Result{
			Status:             status,
			StopReason:         "end_turn",
			CandidatePatchPath: candidatePath,
		}, nil
	}
	t.Cleanup(func() {
		runAgentPatchFixer = original
	})
}

func configureMultiPatchAgentRunner(t *testing.T, candidates map[string]string) {
	t.Helper()
	t.Setenv(agentpatchfixer.EnabledEnv, "true")
	original := runAgentPatchFixer
	runAgentPatchFixer = func(_ context.Context, request agentpatchfixer.Request) (*agentpatchfixer.Result, error) {
		candidate, found := candidates[filepath.Base(request.FailedPatchPath)]
		if !found {
			return nil, fmt.Errorf("unexpected patch %s", request.FailedPatchPath)
		}
		return &agentpatchfixer.Result{
			Status:             "candidate_generated",
			StopReason:         "end_turn",
			CandidatePatchPath: candidate,
		}, nil
	}
	t.Cleanup(func() {
		runAgentPatchFixer = original
	})
}

func configurePartialPatchAgentRunner(t *testing.T, firstCandidate string) {
	t.Helper()
	t.Setenv(agentpatchfixer.EnabledEnv, "true")
	original := runAgentPatchFixer
	runAgentPatchFixer = func(_ context.Context, request agentpatchfixer.Request) (*agentpatchfixer.Result, error) {
		if filepath.Base(request.FailedPatchPath) == "0001-stale.patch" {
			return &agentpatchfixer.Result{
				Status:             "candidate_generated",
				StopReason:         "end_turn",
				CandidatePatchPath: firstCandidate,
			}, nil
		}
		return &agentpatchfixer.Result{
			Status:     "agent_incomplete",
			StopReason: "limit_turns",
		}, nil
	}
	t.Cleanup(func() {
		runAgentPatchFixer = original
	})
}

func TestApplyGeneratedPatchDirectoryReplacesCompleteSet(t *testing.T) {
	projectRoot, _, failedPatchPath, candidatePath := patchCandidateFixture(t)
	generatedDir := t.TempDir()
	generatedPath := filepath.Join(generatedDir, "0001-generated.patch")
	candidate, err := os.ReadFile(candidatePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(generatedPath, candidate, 0o644); err != nil {
		t.Fatal(err)
	}

	fixed, count, changed, err := applyGeneratedPatchDirectory(projectRoot, "repo", filepath.Dir(failedPatchPath), generatedDir)
	if err != nil {
		t.Fatalf("applyGeneratedPatchDirectory() error = %v", err)
	}
	if !fixed || count != 1 {
		t.Fatalf("applyGeneratedPatchDirectory() = fixed %t, count %d", fixed, count)
	}
	if len(changed) != 2 {
		t.Fatalf("changed paths = %#v, want old and new patch paths", changed)
	}
	if _, err := os.Stat(failedPatchPath); !os.IsNotExist(err) {
		t.Fatalf("old patch still exists: %v", err)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(failedPatchPath), "0001-generated.patch")); err != nil {
		t.Fatalf("generated patch missing: %v", err)
	}
}

func TestApplyGeneratedPatchDirectoryRestoresRejectedSet(t *testing.T) {
	projectRoot, repoPath, failedPatchPath, _ := patchCandidateFixture(t)
	original, err := os.ReadFile(failedPatchPath)
	if err != nil {
		t.Fatal(err)
	}
	generatedDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(generatedDir, "0001-invalid.patch"), []byte("invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fixed, _, _, err := applyGeneratedPatchDirectory(projectRoot, "repo", filepath.Dir(failedPatchPath), generatedDir)
	if err == nil {
		t.Fatal("applyGeneratedPatchDirectory() error = nil, want rejection")
	}
	if fixed {
		t.Fatal("applyGeneratedPatchDirectory() fixed = true, want false")
	}
	restored, readErr := os.ReadFile(failedPatchPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(restored) != string(original) {
		t.Fatal("original patch set was not restored")
	}
	if status := runTestCommand(t, repoPath, "git", "status", "--porcelain"); status != "" {
		t.Fatalf("repository is dirty after rollback: %s", status)
	}
}

func patchCandidateFixture(t *testing.T) (projectRoot, repoPath, failedPatchPath, candidatePath string) {
	t.Helper()
	projectRoot = t.TempDir()
	repoPath = filepath.Join(projectRoot, "repo")
	if err := os.Mkdir(repoPath, 0o755); err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, repoPath, "git", "init", "-q")
	runTestCommand(t, repoPath, "git", "config", "user.name", "Test")
	runTestCommand(t, repoPath, "git", "config", "user.email", "test@example.com")
	if err := os.WriteFile(filepath.Join(repoPath, "main.txt"), []byte("old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, repoPath, "git", "add", "main.txt")
	runTestCommand(t, repoPath, "git", "commit", "-qm", "base")
	runTestCommand(t, repoPath, "git", "tag", "v1.0.0")
	if err := os.WriteFile(filepath.Join(repoPath, "main.txt"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, repoPath, "git", "add", "main.txt")
	runTestCommand(t, repoPath, "git", "commit", "-qm", "candidate")

	candidatePath = filepath.Join(projectRoot, "candidate.patch")
	candidate := runTestCommand(t, repoPath, "git", "format-patch", "-1", "--stdout", "--no-signature")
	if err := os.WriteFile(candidatePath, []byte(candidate), 0o644); err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, repoPath, "git", "reset", "--hard", "v1.0.0")

	patchesDir := filepath.Join(projectRoot, "patches")
	if err := os.Mkdir(patchesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	failedPatchPath = filepath.Join(patchesDir, "0001-change.patch")
	if err := os.WriteFile(failedPatchPath, []byte("original-invalid-patch\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	makefile := "patch-repo:\n\tgit -C repo am --committer-date-is-author-date $(CURDIR)/patches/*.patch\n"
	if err := os.WriteFile(filepath.Join(projectRoot, "Makefile"), []byte(makefile), 0o644); err != nil {
		t.Fatal(err)
	}
	return projectRoot, repoPath, failedPatchPath, candidatePath
}

func multiPatchFixture(t *testing.T) (projectRoot, repoPath, patchesDir string, candidates map[string]string) {
	t.Helper()
	projectRoot = t.TempDir()
	repoPath = filepath.Join(projectRoot, "repo")
	if err := os.Mkdir(repoPath, 0o755); err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, repoPath, "git", "init", "-q")
	runTestCommand(t, repoPath, "git", "config", "user.name", "Test")
	runTestCommand(t, repoPath, "git", "config", "user.email", "test@example.com")
	for name, content := range map[string]string{"one.txt": "current-one\n", "two.txt": "current-two\n"} {
		if err := os.WriteFile(filepath.Join(repoPath, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runTestCommand(t, repoPath, "git", "add", ".")
	runTestCommand(t, repoPath, "git", "commit", "-qm", "base")
	runTestCommand(t, repoPath, "git", "tag", "v1.0.0")

	candidates = make(map[string]string)
	for index, change := range []struct {
		name    string
		content string
		patch   string
	}{
		{name: "one.txt", content: "fixed-one\n", patch: "0001-stale.patch"},
		{name: "two.txt", content: "fixed-two\n", patch: "0002-stale.patch"},
	} {
		if err := os.WriteFile(filepath.Join(repoPath, change.name), []byte(change.content), 0o644); err != nil {
			t.Fatal(err)
		}
		runTestCommand(t, repoPath, "git", "add", change.name)
		runTestCommand(t, repoPath, "git", "commit", "-qm", fmt.Sprintf("candidate %d", index+1))
		candidatePath := filepath.Join(projectRoot, fmt.Sprintf("candidate-%d.patch", index+1))
		candidate := runTestCommand(t, repoPath, "git", "format-patch", "-1", "--stdout", "--no-signature")
		if err := os.WriteFile(candidatePath, []byte(candidate), 0o644); err != nil {
			t.Fatal(err)
		}
		candidates[change.patch] = candidatePath
	}
	runTestCommand(t, repoPath, "git", "reset", "--hard", "v1.0.0")

	staleRepo := filepath.Join(projectRoot, "stale")
	if err := os.Mkdir(staleRepo, 0o755); err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, staleRepo, "git", "init", "-q")
	runTestCommand(t, staleRepo, "git", "config", "user.name", "Test")
	runTestCommand(t, staleRepo, "git", "config", "user.email", "test@example.com")
	for name, content := range map[string]string{"one.txt": "old-one\n", "two.txt": "old-two\n"} {
		if err := os.WriteFile(filepath.Join(staleRepo, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runTestCommand(t, staleRepo, "git", "add", ".")
	runTestCommand(t, staleRepo, "git", "commit", "-qm", "stale base")

	patchesDir = filepath.Join(projectRoot, "patches")
	if err := os.Mkdir(patchesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for index, change := range []struct {
		name    string
		content string
		patch   string
	}{
		{name: "one.txt", content: "stale-one\n", patch: "0001-stale.patch"},
		{name: "two.txt", content: "stale-two\n", patch: "0002-stale.patch"},
	} {
		if err := os.WriteFile(filepath.Join(staleRepo, change.name), []byte(change.content), 0o644); err != nil {
			t.Fatal(err)
		}
		runTestCommand(t, staleRepo, "git", "add", change.name)
		runTestCommand(t, staleRepo, "git", "commit", "-qm", fmt.Sprintf("stale %d", index+1))
		patch := runTestCommand(t, staleRepo, "git", "format-patch", "-1", "--stdout", "--no-signature")
		if err := os.WriteFile(filepath.Join(patchesDir, change.patch), []byte(patch), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	makefile := "patch-repo:\n\tgit -C repo am --committer-date-is-author-date $(CURDIR)/patches/*.patch\n"
	if err := os.WriteFile(filepath.Join(projectRoot, "Makefile"), []byte(makefile), 0o644); err != nil {
		t.Fatal(err)
	}
	return projectRoot, repoPath, patchesDir, candidates
}

func runTestCommand(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, output)
	}
	return string(output)
}

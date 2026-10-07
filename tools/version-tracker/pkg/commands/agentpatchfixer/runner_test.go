package agentpatchfixer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testPatch = `From 0123456789012345678901234567890123456789 Mon Sep 17 00:00:00 2001
From: Example Author <author@example.com>
Date: Tue, 16 Sep 2025 10:30:00 +0000
Subject: [PATCH 1/2] Preserve useful behavior

Keep this explanatory body.

---
 main.go | 2 +-
 1 file changed, 1 insertion(+), 1 deletion(-)

diff --git a/main.go b/main.go
index 1111111..2222222 100644
--- a/main.go
+++ b/main.go
@@ -1 +1 @@
-old
+new
`

func TestEnabledDefaultsToFalse(t *testing.T) {
	t.Setenv(EnabledEnv, "")
	if Enabled() {
		t.Fatal("Enabled() = true, want false")
	}
}

func TestEnabledRejectsInvalidValue(t *testing.T) {
	t.Setenv(EnabledEnv, "sometimes")
	if Enabled() {
		t.Fatal("Enabled() = true, want false")
	}
}

func TestRunDisabled(t *testing.T) {
	t.Setenv(EnabledEnv, "false")
	result, err := Run(context.Background(), Request{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result != nil {
		t.Fatalf("Run() result = %#v, want nil", result)
	}
}

func TestParsePatchMetadata(t *testing.T) {
	metadata, err := parsePatchMetadata(context.Background(), testPatch)
	if err != nil {
		t.Fatal(err)
	}
	if metadata.AuthorName != "Example Author" || metadata.AuthorEmail != "author@example.com" {
		t.Fatalf("author = %s <%s>", metadata.AuthorName, metadata.AuthorEmail)
	}
	if metadata.Subject != "Preserve useful behavior" {
		t.Fatalf("subject = %q", metadata.Subject)
	}
	if metadata.Body != "Keep this explanatory body." {
		t.Fatalf("body = %q", metadata.Body)
	}
}

func TestPatchFilesIgnoresDiffTextInCommitMessage(t *testing.T) {
	patch := strings.Replace(testPatch, "Keep this explanatory body.", "diff --git a/Makefile b/Makefile\nKeep this explanatory body.", 1)
	files, err := patchFiles(context.Background(), patch)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || !files["main.go"] {
		t.Fatalf("files = %#v, want main.go", files)
	}
}

func TestPatchFilesHandlesDeletedDoubleDashLine(t *testing.T) {
	repo := initializeTestRepository(t, map[string]string{
		"first.txt":  "old-first\n",
		"second.txt": "--\nkeep\n",
	})
	writeTestFile(t, repo, "first.txt", "new-first\n")
	writeTestFile(t, repo, "second.txt", "keep\n")
	runTestCommand(t, repo, "git", "add", ".")
	runTestCommand(t, repo, "git", "commit", "-qm", "change both")

	patch := runTestCommand(t, repo, "git", "format-patch", "-1", "--stdout", "--no-signature")
	files, err := patchFiles(context.Background(), patch)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 || !files["first.txt"] || !files["second.txt"] {
		t.Fatalf("files = %#v, want first.txt and second.txt", files)
	}
}

func TestPatchFilesHandlesSpacesInPath(t *testing.T) {
	repo := initializeTestRepository(t, map[string]string{"file with spaces.txt": "old\n"})
	writeTestFile(t, repo, "file with spaces.txt", "new\n")
	runTestCommand(t, repo, "git", "add", ".")
	runTestCommand(t, repo, "git", "commit", "-qm", "change spaced path")

	patch := runTestCommand(t, repo, "git", "format-patch", "-1", "--stdout", "--no-signature")
	files, err := patchFiles(context.Background(), patch)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || !files["file with spaces.txt"] {
		t.Fatalf("files = %#v, want spaced path", files)
	}
}

func TestGenerateCandidatePreservesAuthorAndSubject(t *testing.T) {
	repo := initializeRepository(t, "old\n")
	if err := os.WriteFile(filepath.Join(repo, "main.go"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	metadata, err := parsePatchMetadata(context.Background(), testPatch)
	if err != nil {
		t.Fatal(err)
	}
	outputPath := filepath.Join(t.TempDir(), "candidate.patch")
	edited, err := generateCandidate(
		context.Background(),
		repo,
		metadata,
		outputPath,
		map[string]bool{"main.go": true},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(edited) != 1 || edited[0] != "main.go" {
		t.Fatalf("edited = %#v", edited)
	}
	candidate, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(candidate), "From: Example Author <author@example.com>") {
		t.Fatal("candidate did not preserve author")
	}
	if !strings.Contains(string(candidate), "Subject: [PATCH] Preserve useful behavior") {
		t.Fatal("candidate did not preserve subject")
	}
	runTestCommand(t, repo, "git", "apply", "--check", "--reverse", outputPath)
}

func TestGenerateCandidateAllowsNewFileInNewDirectory(t *testing.T) {
	repo := initializeRepository(t, "old\n")
	if err := os.Mkdir(filepath.Join(repo, "newdir"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "newdir", "new.go"), []byte("package newdir\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	metadata, err := parsePatchMetadata(context.Background(), testPatch)
	if err != nil {
		t.Fatal(err)
	}
	outputPath := filepath.Join(t.TempDir(), "candidate.patch")

	edited, err := generateCandidate(
		context.Background(),
		repo,
		metadata,
		outputPath,
		map[string]bool{"newdir/new.go": true},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(edited) != 1 || edited[0] != "newdir/new.go" {
		t.Fatalf("edited = %#v, want newdir/new.go", edited)
	}
}

func TestGenerateCandidateRejectsUnexpectedFile(t *testing.T) {
	repo := initializeRepository(t, "old\n")
	if err := os.Mkdir(filepath.Join(repo, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "pkg", "allowed.go"), []byte("old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "pkg", "sibling.go"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	metadata, err := parsePatchMetadata(context.Background(), testPatch)
	if err != nil {
		t.Fatal(err)
	}
	_, err = generateCandidate(
		context.Background(),
		repo,
		metadata,
		filepath.Join(t.TempDir(), "candidate.patch"),
		map[string]bool{"pkg/allowed.go": true},
	)
	if err == nil {
		t.Fatal("generateCandidate() accepted unexpected file")
	}
}

func TestChangedFilesHandlesRenameRecords(t *testing.T) {
	repo := initializeTestRepository(t, map[string]string{"old name.txt": "content\n"})
	runTestCommand(t, repo, "git", "mv", "old name.txt", "new name.txt")

	files, err := changedFiles(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 || files[0] != "new name.txt" || files[1] != "old name.txt" {
		t.Fatalf("files = %#v, want old and new paths", files)
	}
}

func TestSeedWorkspaceAppliesCleanHunksAndRecordsRejects(t *testing.T) {
	repo := t.TempDir()
	runTestCommand(t, repo, "git", "init", "-q")
	runTestCommand(t, repo, "git", "config", "user.name", "Test")
	runTestCommand(t, repo, "git", "config", "user.email", "test@example.com")
	writeTestFile(t, repo, "one.txt", "old-one\n")
	writeTestFile(t, repo, "two.txt", "old-two\n")
	runTestCommand(t, repo, "git", "add", ".")
	runTestCommand(t, repo, "git", "commit", "-qm", "base")
	base := strings.TrimSpace(runTestCommand(t, repo, "git", "rev-parse", "HEAD"))
	writeTestFile(t, repo, "one.txt", "patched-one\n")
	writeTestFile(t, repo, "two.txt", "patched-two\n")
	patchPath := filepath.Join(t.TempDir(), "change.patch")
	if err := os.WriteFile(patchPath, []byte(runTestCommand(t, repo, "git", "diff", "--binary")), 0o644); err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, repo, "git", "reset", "--hard", base)
	writeTestFile(t, repo, "two.txt", "current-two\n")

	seed, err := seedWorkspace(context.Background(), repo, patchPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, repo, "one.txt"); got != "patched-one\n" {
		t.Fatalf("one.txt = %q", got)
	}
	if got := readTestFile(t, repo, "two.txt"); got != "current-two\n" {
		t.Fatalf("two.txt = %q", got)
	}
	if len(seed.RejectFiles) != 1 || seed.RejectFiles[0] != "two.txt.rej" {
		t.Fatalf("rejects = %#v", seed.RejectFiles)
	}
	removeRejectFiles(repo, seed.RejectFiles)
	if _, err := os.Stat(filepath.Join(repo, "two.txt.rej")); !os.IsNotExist(err) {
		t.Fatalf("reject file still exists: %v", err)
	}
}

func TestSeedWorkspaceAllowsMissingFileWithoutReject(t *testing.T) {
	repo := initializeTestRepository(t, map[string]string{
		"keep.txt":    "old-keep\n",
		"removed.txt": "old-removed\n",
	})
	writeTestFile(t, repo, "keep.txt", "patched-keep\n")
	writeTestFile(t, repo, "removed.txt", "patched-removed\n")
	patchPath := filepath.Join(t.TempDir(), "change.patch")
	if err := os.WriteFile(patchPath, []byte(runTestCommand(t, repo, "git", "diff", "--binary")), 0o644); err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, repo, "git", "reset", "--hard", "HEAD")
	runTestCommand(t, repo, "git", "rm", "removed.txt")
	runTestCommand(t, repo, "git", "commit", "-qm", "remove file upstream")

	seed, err := seedWorkspace(context.Background(), repo, patchPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(seed.RejectFiles) != 0 {
		t.Fatalf("reject files = %#v, want none", seed.RejectFiles)
	}
	if len(seed.MissingFiles) != 1 || seed.MissingFiles[0] != "removed.txt" {
		t.Fatalf("missing files = %#v, want removed.txt", seed.MissingFiles)
	}
	if got := readTestFile(t, repo, "keep.txt"); got != "patched-keep\n" {
		t.Fatalf("keep.txt = %q, want patched content", got)
	}
}

func TestExecuteUsesStrandsProcess(t *testing.T) {
	source := initializeRepository(t, "current\n")
	patchPath := filepath.Join(t.TempDir(), "failed.patch")
	if err := os.WriteFile(patchPath, []byte(testPatch), 0o644); err != nil {
		t.Fatal(err)
	}
	runDir := t.TempDir()
	t.Setenv(AgentCommandEnv, fakeAgentRunner(t, "end_turn", true))
	result, err := execute(context.Background(), Request{
		SchemaVersion:    requestSchemaVersion,
		ProjectName:      "example/project",
		ProjectRoot:      filepath.Dir(source),
		UpstreamRepoPath: source,
		FailedPatchPath:  patchPath,
		FailureOutput:    "patch does not apply",
		CurrentRevision:  "old",
		TargetRevision:   "new",
		MaxTurns:         4,
		MaxTotalTokens:   1000,
	}, runDir)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "candidate_generated" || result.CandidatePatchPath != "candidate.patch" {
		t.Fatalf("result = %#v", result)
	}
	if result.Turns != 2 || result.TotalTokens != 38 {
		t.Fatalf("usage = %d turns, %d tokens", result.Turns, result.TotalTokens)
	}
	var invocation agentInvocation
	invocationBytes, err := os.ReadFile(filepath.Join(runDir, "agent-request.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(invocationBytes, &invocation); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(invocation.Prompt, "main.go.rej") {
		t.Fatal("prompt does not identify reject file")
	}
	if len(invocation.AllowedFiles) != 1 || invocation.AllowedFiles[0] != "main.go" {
		t.Fatalf("allowed files = %#v", invocation.AllowedFiles)
	}
	if _, err := os.Stat(filepath.Join(runDir, "candidate.patch")); err != nil {
		t.Fatalf("candidate patch missing: %v", err)
	}
}

func TestExecuteDoesNotPublishIncompleteEdits(t *testing.T) {
	source := initializeRepository(t, "current\n")
	patchPath := filepath.Join(t.TempDir(), "failed.patch")
	if err := os.WriteFile(patchPath, []byte(testPatch), 0o644); err != nil {
		t.Fatal(err)
	}
	runDir := t.TempDir()
	t.Setenv(AgentCommandEnv, fakeAgentRunner(t, "max_tokens", true))
	result, err := execute(context.Background(), Request{
		SchemaVersion:    requestSchemaVersion,
		ProjectName:      "example/project",
		ProjectRoot:      filepath.Dir(source),
		UpstreamRepoPath: source,
		FailedPatchPath:  patchPath,
		MaxTurns:         4,
		MaxTotalTokens:   1000,
	}, runDir)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "agent_incomplete" || result.CandidatePatchPath != "" {
		t.Fatalf("result = %#v", result)
	}
	if _, err := os.Stat(filepath.Join(runDir, "incomplete.diff")); err != nil {
		t.Fatalf("incomplete diff missing: %v", err)
	}
}

func TestRunTimesOut(t *testing.T) {
	source := initializeRepository(t, "current\n")
	patchPath := filepath.Join(t.TempDir(), "failed.patch")
	if err := os.WriteFile(patchPath, []byte(testPatch), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnabledEnv, "true")
	t.Setenv(ResultDirEnv, t.TempDir())
	t.Setenv(TimeoutEnv, "20ms")
	runnerPath := filepath.Join(t.TempDir(), "runner.sh")
	if err := os.WriteFile(runnerPath, []byte("#!/bin/sh\nsleep 2\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(AgentCommandEnv, runnerPath)

	start := time.Now()
	result, err := Run(context.Background(), Request{
		ProjectName:      "example/project",
		ProjectRoot:      filepath.Dir(source),
		UpstreamRepoPath: source,
		FailedPatchPath:  patchPath,
	})
	if err == nil {
		t.Fatal("Run() error = nil, want timeout")
	}
	if result == nil || result.Status != "runner_timeout" {
		t.Fatalf("result = %#v", result)
	}
	if time.Since(start) > time.Second {
		t.Fatalf("timeout took too long: %s", time.Since(start))
	}
}

func TestTotalTimeout(t *testing.T) {
	t.Setenv(TotalTimeoutEnv, "90s")
	timeout, err := TotalTimeout()
	if err != nil {
		t.Fatal(err)
	}
	if timeout != 90*time.Second {
		t.Fatalf("TotalTimeout() = %s, want 90s", timeout)
	}
}

func TestAgentEnvironmentExcludesGitHubToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "secret")
	t.Setenv("AWS_REGION", "us-west-2")
	environment := agentEnvironment(t.TempDir())
	for _, entry := range environment {
		if strings.HasPrefix(entry, "GITHUB_TOKEN=") {
			t.Fatal("agent environment contains GITHUB_TOKEN")
		}
	}
}

func fakeAgentRunner(t *testing.T, stopReason string, edit bool) string {
	t.Helper()
	runnerPath := filepath.Join(t.TempDir(), "runner.py")
	editStatement := ""
	if edit {
		editStatement = `(workspace / "main.go").write_text("new\n", encoding="utf-8")`
	}
	script := fmt.Sprintf(`#!/usr/bin/env python3
import argparse
import json
from pathlib import Path

parser = argparse.ArgumentParser()
parser.add_argument("--request", required=True)
parser.add_argument("--result", required=True)
args = parser.parse_args()
request = json.loads(Path(args.request).read_text(encoding="utf-8"))
workspace = Path(request["workspace_path"])
%s
Path(args.result).write_text(json.dumps({
    "schema_version": 1,
    "status": "completed",
    "stop_reason": %q,
    "summary": "done",
    "turns": 2,
    "total_tokens": 38,
}), encoding="utf-8")
`, editStatement, stopReason)
	if err := os.WriteFile(runnerPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return "python3 " + runnerPath
}

func initializeRepository(t *testing.T, content string) string {
	t.Helper()
	return initializeTestRepository(t, map[string]string{"main.go": content})
}

func initializeTestRepository(t *testing.T, files map[string]string) string {
	t.Helper()
	repo := t.TempDir()
	runTestCommand(t, repo, "git", "init", "-q")
	runTestCommand(t, repo, "git", "config", "user.name", "Test")
	runTestCommand(t, repo, "git", "config", "user.email", "test@example.com")
	for path, content := range files {
		writeTestFile(t, repo, path, content)
	}
	runTestCommand(t, repo, "git", "add", ".")
	runTestCommand(t, repo, "git", "commit", "-qm", "base")
	return repo
}

func runTestCommand(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	output, err := runCommand(context.Background(), dir, nil, nil, name, args...)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

func writeTestFile(t *testing.T, root, relativePath, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, relativePath)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, relativePath), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, root, relativePath string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(root, relativePath))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

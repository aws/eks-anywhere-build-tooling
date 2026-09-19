package agentpatchfixer

import (
	"context"
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

type scriptedModel struct {
	responses []modelResponse
	requests  []modelRequest
}

func (m *scriptedModel) Converse(_ context.Context, request modelRequest) (modelResponse, error) {
	m.requests = append(m.requests, request)
	response := m.responses[0]
	m.responses = m.responses[1:]
	return response, nil
}

type blockingModel struct{}

func (blockingModel) Converse(ctx context.Context, _ modelRequest) (modelResponse, error) {
	<-ctx.Done()
	return modelResponse{}, ctx.Err()
}

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
	files, err := patchFiles(patch)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || !files["main.go"] {
		t.Fatalf("files = %#v, want main.go", files)
	}
}

func TestWorkspaceBlocksPathEscapeAndUnlistedEdits(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	workspace, err := newWorkspace(root, map[string]bool{"main.go": true}, nil, testPatch)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.readFile("../secret", 1, 10); err == nil {
		t.Fatal("readFile() accepted path escape")
	}
	if _, err := workspace.listFiles("../*"); err == nil {
		t.Fatal("listFiles() accepted path escape")
	}
	if _, err := workspace.replaceText("other.go", "a", "b"); err == nil {
		t.Fatal("replaceText() accepted unlisted edit")
	}
}

func TestWorkspaceBlocksSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	workspace, err := newWorkspace(root, map[string]bool{"linked/secret": true}, []string{"linked"}, testPatch)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.readFile("linked/secret", 1, 10); err == nil {
		t.Fatal("readFile() followed symlink outside workspace")
	}
	if _, err := workspace.writeFile("linked/new", "data"); err == nil {
		t.Fatal("writeFile() followed symlink outside workspace")
	}
}

func TestWorkspaceReplacesLineRange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "main.go")
	if err := os.WriteFile(path, []byte("one\ntwo\nthree\nfour\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	workspace, err := newWorkspace(root, map[string]bool{"main.go": true}, nil, testPatch)
	if err != nil {
		t.Fatal(err)
	}
	result, err := workspace.replaceLines("main.go", 2, 3, "new-two\nnew-three")
	if err != nil {
		t.Fatal(err)
	}
	if result != "replaced lines 2-3 in main.go" {
		t.Fatalf("result = %q", result)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "one\nnew-two\nnew-three\nfour\n" {
		t.Fatalf("content = %q", content)
	}
}

func TestWorkspaceWritesNewFileBelowAllowedRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	workspace, err := newWorkspace(root, nil, []string{"pkg"}, testPatch)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.writeFile("pkg/new/nested.go", "package nested\n"); err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, root, "pkg/new/nested.go"); got != "package nested\n" {
		t.Fatalf("content = %q", got)
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
		nil,
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

func TestGenerateCandidateRejectsUnexpectedFile(t *testing.T) {
	repo := initializeRepository(t, "old\n")
	if err := os.WriteFile(filepath.Join(repo, "unexpected.go"), []byte("new\n"), 0o644); err != nil {
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
		map[string]bool{"main.go": true},
		nil,
	)
	if err == nil {
		t.Fatal("generateCandidate() accepted unexpected file")
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

func TestExecuteUsesNativeToolLoop(t *testing.T) {
	source := initializeRepository(t, "current\n")
	patchPath := filepath.Join(t.TempDir(), "failed.patch")
	if err := os.WriteFile(patchPath, []byte(testPatch), 0o644); err != nil {
		t.Fatal(err)
	}
	runDir := t.TempDir()
	client := &scriptedModel{responses: []modelResponse{
		{
			Message: modelMessage{
				Role: "assistant",
				Content: []contentBlock{{ToolUse: &toolUse{
					ID:   "tool-1",
					Name: "replace_text",
					Input: map[string]any{
						"path":     "main.go",
						"old_text": "current\n",
						"new_text": "new\n",
					},
				}}},
			},
			StopReason:   "tool_use",
			InputTokens:  10,
			OutputTokens: 5,
		},
		{
			Message:      modelMessage{Role: "assistant", Content: []contentBlock{{Text: "done"}}},
			StopReason:   "end_turn",
			InputTokens:  20,
			OutputTokens: 3,
		},
	}}
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
	}, runDir, func(context.Context, string, string) (model, error) {
		return client, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "candidate_generated" || result.CandidatePatchPath != "candidate.patch" {
		t.Fatalf("result = %#v", result)
	}
	if result.Turns != 2 || result.TotalTokens != 38 {
		t.Fatalf("usage = %d turns, %d tokens", result.Turns, result.TotalTokens)
	}
	if len(client.requests) != 2 {
		t.Fatalf("requests = %d, want 2", len(client.requests))
	}
	if got := client.requests[1].Messages[len(client.requests[1].Messages)-1].Content[0].ToolResult; got == nil || got.IsError {
		t.Fatalf("tool result = %#v", got)
	}
	if !strings.Contains(client.requests[0].Messages[0].Content[0].Text, "main.go.rej") {
		t.Fatal("prompt does not identify reject file")
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
	client := &scriptedModel{responses: []modelResponse{
		{
			Message: modelMessage{
				Role: "assistant",
				Content: []contentBlock{{ToolUse: &toolUse{
					ID:   "tool-1",
					Name: "replace_text",
					Input: map[string]any{
						"path":     "main.go",
						"old_text": "current\n",
						"new_text": "new\n",
					},
				}}},
			},
			StopReason: "tool_use",
		},
		{
			Message:    modelMessage{Role: "assistant", Content: []contentBlock{{Text: "out of room"}}},
			StopReason: "max_tokens",
		},
	}}
	result, err := execute(context.Background(), Request{
		SchemaVersion:    requestSchemaVersion,
		ProjectName:      "example/project",
		ProjectRoot:      filepath.Dir(source),
		UpstreamRepoPath: source,
		FailedPatchPath:  patchPath,
		MaxTurns:         4,
		MaxTotalTokens:   1000,
	}, runDir, func(context.Context, string, string) (model, error) {
		return client, nil
	})
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
	originalFactory := defaultModelFactory
	defaultModelFactory = func(context.Context, string, string) (model, error) {
		return blockingModel{}, nil
	}
	t.Cleanup(func() {
		defaultModelFactory = originalFactory
	})

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

func initializeRepository(t *testing.T, content string) string {
	t.Helper()
	repo := t.TempDir()
	runTestCommand(t, repo, "git", "init", "-q")
	runTestCommand(t, repo, "git", "config", "user.name", "Test")
	runTestCommand(t, repo, "git", "config", "user.email", "test@example.com")
	writeTestFile(t, repo, "main.go", content)
	runTestCommand(t, repo, "git", "add", "main.go")
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

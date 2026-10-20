package upgrade

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/commands/agentpatchfixer"
)

func TestInitialPatchApplyPreservesCodeCommitCredentialsAndReplayScrubsThem(t *testing.T) {
	projectRoot := t.TempDir()
	bin := t.TempDir()
	// Exercise the real Make subprocess and credential-helper boundary without
	// contacting AWS or putting credentials in command output.
	helper := `#!/bin/sh
set -eu
test "$1" = codecommit
test "$2" = credential-helper
test "$CLONE_URL" = https://codecommit.example.invalid/v1/repos/dummy
test "$CODEBUILD_CI" = true
test "$AWS_CONTAINER_CREDENTIALS_RELATIVE_URI" = /dummy/credentials
test "$AWS_ACCESS_KEY_ID" = dummy-key
test "$AWS_SECRET_ACCESS_KEY" = dummy-secret
test "$AWS_CONFIG_FILE" = /dummy/config
printf 'username=dummy\npassword=dummy\n'
`
	writeRepairTestFile(t, filepath.Join(bin, "aws"), []byte(helper), 0o755)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("CLONE_URL", "https://codecommit.example.invalid/v1/repos/dummy")
	t.Setenv("CODEBUILD_CI", "true")
	t.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "/dummy/credentials")
	t.Setenv("AWS_ACCESS_KEY_ID", "dummy-key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "dummy-secret")
	t.Setenv("AWS_CONFIG_FILE", "/dummy/config")
	t.Setenv("GITHUB_TOKEN", "dummy-token")
	t.Setenv("MAKEFLAGS", "")
	t.Setenv("MFLAGS", "")
	t.Setenv("MAKEOVERRIDES", "")
	makefile := `.PHONY: patch-repo
patch-repo:
	@set -eu; if [ ! -d repo ]; then aws codecommit credential-helper get >/dev/null; mkdir repo; echo checkout >> calls; fi
	@set -eu; if [ "$(CODEBUILD_CI)" = false ]; then \
		test -z "$${AWS_CONTAINER_CREDENTIALS_RELATIVE_URI:-}"; test -z "$${AWS_ACCESS_KEY_ID:-}"; test -z "$${AWS_SECRET_ACCESS_KEY:-}"; test -z "$${GITHUB_TOKEN:-}"; \
		test "$$AWS_EC2_METADATA_DISABLED" = true; test "$$AWS_CONFIG_FILE" = /dev/null; test "$$AWS_SHARED_CREDENTIALS_FILE" = /dev/null; \
		echo replay >> calls; \
	fi
`
	writeRepairTestFile(t, filepath.Join(projectRoot, "Makefile"), []byte(makefile), 0o644)
	if _, _, _, _, err := applyPatchesToRepoContext(context.Background(), projectRoot, "repo", 1); err == nil {
		t.Fatal("fresh CodeCommit checkout unexpectedly succeeded without credentials")
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "repo")); !os.IsNotExist(err) {
		t.Fatalf("failed credential helper still created checkout: %v", err)
	}
	if applied, _, _, _, err := applyPatchesToRepo(projectRoot, "repo", 1); err != nil || applied != 1 {
		t.Fatalf("initial CodeCommit apply = %d patches, error %v", applied, err)
	}
	if applied, _, _, _, err := applyPatchesToRepoContext(context.Background(), projectRoot, "repo", 1); err != nil || applied != 1 {
		t.Fatalf("scrubbed replay = %d patches, error %v", applied, err)
	}
	calls, err := os.ReadFile(filepath.Join(projectRoot, "calls"))
	if err != nil || string(calls) != "checkout\nreplay\n" {
		t.Fatalf("checkout/replay calls = %q, error %v", calls, err)
	}
}

func TestRepairGenericPatchSeriesRejectsFailedProjectValidation(t *testing.T) {
	projectRoot, _, failedPath, candidatePath := patchCandidateFixture(t)
	original, err := os.ReadFile(failedPath)
	if err != nil {
		t.Fatal(err)
	}
	applied, _, files, output, _ := applyPatchesToRepo(projectRoot, "repo", 1)
	configureFakeAgentRunner(t, candidatePath, "candidate_generated")
	validationFailure := errors.New("attribution requires a compatible Go toolchain")
	calls := 0
	fixed, changed, err := repairGenericPatchSeries(
		context.Background(), "example/project", projectRoot, "repo",
		filepath.Dir(failedPath), "v1.0.0", "v1.1.0", "", 1, applied, files, output,
		patchRepairOptions{Validate: func(context.Context) error {
			calls++
			return validationFailure
		}},
	)
	if fixed || len(changed) != 0 || !errors.Is(err, validationFailure) {
		t.Fatalf("rejected validation = fixed %t, changed %v, error %v", fixed, changed, err)
	}
	if calls != 1 {
		t.Fatalf("environmental validation was retried %d times", calls)
	}
	restored, err := os.ReadFile(failedPath)
	if err != nil || string(restored) != string(original) {
		t.Fatalf("rejected patch was retained: %q, error %v", restored, err)
	}
}

func TestRepairGenericPatchSeriesRetriesSourceFailureWithOriginalIntent(t *testing.T) {
	projectRoot, repoPath, failedPath, candidatePath := patchCandidateFixture(t)
	candidate, err := os.ReadFile(candidatePath)
	if err != nil {
		t.Fatal(err)
	}
	originalPatch := strings.Replace(string(candidate), "\n-old\n", "\n-stale\n", 1)
	if err := os.WriteFile(failedPath, []byte(originalPatch), 0o644); err != nil {
		t.Fatal(err)
	}
	// Match Common.mk: the patch marker depends on a checkout marker, rather
	// than patch timestamps. A retry must invalidate both cached targets.
	makefile := `.PHONY: patch-repo
patch-repo: repo/eks-anywhere-patched
repo/eks-anywhere-checkout-v1.0.0:
	git -C repo checkout --quiet -f v1.0.0
	touch $@
repo/eks-anywhere-patched: repo/eks-anywhere-checkout-v1.0.0
	git -C repo am $(CURDIR)/patches/*.patch
	touch $@
`
	if err := os.WriteFile(filepath.Join(projectRoot, "Makefile"), []byte(makefile), 0o644); err != nil {
		t.Fatal(err)
	}
	secondCandidate := filepath.Join(projectRoot, "second-candidate.patch")
	if err := os.WriteFile(secondCandidate, []byte(strings.Replace(string(candidate), "\n+new\n", "\n+validated\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	applied, _, files, output, applyErr := applyPatchesToRepo(projectRoot, "repo", 1)
	if applyErr == nil {
		t.Fatal("stale patch unexpectedly applied")
	}
	originalRunner := runAgentPatchFixer
	t.Cleanup(func() { runAgentPatchFixer = originalRunner })
	attempts, validations := 0, 0
	runAgentPatchFixer = func(_ context.Context, request agentpatchfixer.Request) (*agentpatchfixer.Result, error) {
		attempts++
		if attempts == 2 {
			content, err := os.ReadFile(request.FailedPatchPath)
			if err != nil || string(content) != originalPatch {
				t.Fatalf("retry lost original patch intent: %q, error %v", content, err)
			}
			if !strings.Contains(request.FailureOutput, "invalid generated source") {
				t.Fatalf("retry has no validation feedback: %q", request.FailureOutput)
			}
		}
		path := candidatePath
		if attempts == 2 {
			path = secondCandidate
		}
		return &agentpatchfixer.Result{Status: "candidate_generated", CandidatePatchPath: path}, nil
	}
	fixed, changed, err := repairGenericPatchSeries(
		context.Background(), "example/project", projectRoot, "repo",
		filepath.Dir(failedPath), "v1.0.0", "v1.1.0", "", 1, applied, files, output,
		patchRepairOptions{
			BaseRevision: "v1.0.0",
			Validate: func(context.Context) error {
				validations++
				if validations == 1 {
					return &repairSourceValidationError{Path: "main.txt", Detail: "invalid generated source"}
				}
				content, err := os.ReadFile(filepath.Join(repoPath, "main.txt"))
				if err != nil || string(content) != "validated\n" {
					t.Fatalf("validation checked an unapplied retry: %q, error %v", content, err)
				}
				return nil
			},
		},
	)
	if err != nil || !fixed || len(changed) != 1 {
		t.Fatalf("source retry = fixed %t, changed %v, error %v", fixed, changed, err)
	}
	if attempts != 2 || validations != 2 {
		t.Fatalf("attempts %d, validations %d; want 2 each", attempts, validations)
	}
}

func TestRepairGenericPatchSeriesBoundsSourceValidationRetries(t *testing.T) {
	projectRoot, repoPath, failedPath, candidatePath := patchCandidateFixture(t)
	candidate, err := os.ReadFile(candidatePath)
	if err != nil {
		t.Fatal(err)
	}
	originalPatch := strings.Replace(string(candidate), "\n-old\n", "\n-stale\n", 1)
	if err := os.WriteFile(failedPath, []byte(originalPatch), 0o644); err != nil {
		t.Fatal(err)
	}
	applied, _, files, output, _ := applyPatchesToRepo(projectRoot, "repo", 1)
	configureFakeAgentRunner(t, candidatePath, "candidate_generated")
	validations := 0
	fixed, _, err := repairGenericPatchSeries(
		context.Background(), "example/project", projectRoot, "repo",
		filepath.Dir(failedPath), "v1.0.0", "v1.1.0", "", 1, applied, files, output,
		patchRepairOptions{
			BaseRevision: "v1.0.0",
			Validate: func(context.Context) error {
				validations++
				return &repairSourceValidationError{Path: "main.txt", Detail: "invalid generated source"}
			},
		},
	)
	if fixed || err == nil || !strings.Contains(err.Error(), "exceeded 3 attempts") || validations != 3 {
		t.Fatalf("retry limit = fixed %t, validations %d, error %v", fixed, validations, err)
	}
	content, err := os.ReadFile(failedPath)
	if err != nil || string(content) != originalPatch {
		t.Fatalf("original patch not restored: %q, error %v", content, err)
	}
	if content, err := os.ReadFile(filepath.Join(repoPath, "main.txt")); err != nil || string(content) != "old\n" {
		t.Fatalf("rejected source remains: %q, error %v", content, err)
	}
}

func TestRepairGenericPatchSeriesRestoresRemovedPatchAfterValidationFailure(t *testing.T) {
	projectRoot, _, failedPath, _ := patchCandidateFixture(t)
	original, err := os.ReadFile(failedPath)
	if err != nil {
		t.Fatal(err)
	}
	applied, _, files, output, _ := applyPatchesToRepo(projectRoot, "repo", 1)
	configureFakeAgentRunner(t, "", "no_changes")
	fixed, _, err := repairGenericPatchSeries(
		context.Background(), "example/project", projectRoot, "repo",
		filepath.Dir(failedPath), "v1.0.0", "v1.1.0", "", 1, applied, files, output,
		patchRepairOptions{Validate: func(context.Context) error {
			return errors.New("license validation failed")
		}},
	)
	if fixed || err == nil {
		t.Fatalf("obsolete candidate = fixed %t, error %v", fixed, err)
	}
	content, err := os.ReadFile(failedPath)
	if err != nil || string(content) != string(original) {
		t.Fatalf("removed patch not restored: %q, error %v", content, err)
	}
}

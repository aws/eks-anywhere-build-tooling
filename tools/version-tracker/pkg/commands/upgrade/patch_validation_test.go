package upgrade

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestValidateRepairedProjectRejectsInvalidDockerfileBeforeMake(t *testing.T) {
	for _, test := range []struct {
		name    string
		content string
		detail  string
	}{
		{name: "unknown instruction", content: "FROM scratch\nilding runc\n", detail: "dockerfile parse error on line 2: unknown instruction: ilding"},
		{name: "parser failure", content: "FROM scratch\nRUN <<EOF\nunfinished\n", detail: "unterminated heredoc"},
		{name: "invalid arguments", content: "FROM\n", detail: "FROM requires"},
		{name: "no build stage", content: "ARG VERSION=example\n", detail: "no build stage"},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := repairValidationFixture(t, test.content)
			result, err := validateRepairedProject(context.Background(), request)
			var sourceError *repairSourceValidationError
			if !errors.As(err, &sourceError) {
				t.Fatalf("validateRepairedProject() error = %T %v, want *repairSourceValidationError", err, err)
			}
			if sourceError.Path != "Dockerfile" || !strings.Contains(sourceError.Detail, test.detail) {
				t.Fatalf("source error = %#v, want Dockerfile and %q", sourceError, test.detail)
			}
			if !strings.Contains(err.Error(), sourceError.Path) || !strings.Contains(err.Error(), sourceError.Detail) {
				t.Fatalf("source error message = %q", err.Error())
			}
			if len(result.Checks) != 0 || result.NativeImages {
				t.Fatalf("validation result = %#v, want no completed checks", result)
			}
			assertRepairValidationCalls(t, request)
		})
	}
}

func TestValidateRepairedProjectAcceptsValidDockerfileAndBuildsNativeOCI(t *testing.T) {
	validDockerfile := "# syntax=docker/dockerfile:1\n" +
		"ARG BASE=scratch\nFROM ${BASE} AS source\n" +
		"COPY --chown=1000:1000 [\"dummy.txt\", \"/dummy.txt\"]\n" +
		"FROM scratch\nCOPY --from=source /dummy.txt /dummy.txt\n"
	request := repairValidationFixture(t, validDockerfile)
	for name, value := range map[string]string{
		"IMAGE_OUTPUT_TYPE":  "image",
		"IMAGE_OUTPUT":       "push=true",
		"IMAGE_OUTPUT_DIR":   "/dummy/output",
		"IMAGE_REPO":         "example.invalid/dummy",
		"IMAGE_EXPORT_CACHE": "dummy-export-cache",
		"IMAGE_IMPORT_CACHE": "dummy-import-cache",
		"IMAGE_PLATFORMS":    "linux/amd64,linux/arm64",
		"AWS_ACCESS_KEY_ID":  "dummy-access-key",
		"GITHUB_TOKEN":       "dummy-token",
	} {
		t.Setenv(name, value)
	}
	result, err := validateRepairedProject(context.Background(), request)
	if err != nil {
		t.Fatalf("validateRepairedProject() error = %v", err)
	}
	want := repairValidationResult{
		Checks:       []string{"dockerfile-syntax", "attribution-checksums", "validate-checksums", "local-images"},
		NativeImages: true,
		Architecture: "arm64",
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("validateRepairedProject() = %#v, want %#v", result, want)
	}
	assertRepairValidationCalls(t, request, "attribution-checksums", "validate-checksums", "var-value-IMAGE_NAMES", "var-value-PATCH_VALIDATION_IMAGE_TARGET", "var-value-BUILDER_PLATFORM_ARCH", "local-images")
	output, err := os.ReadFile(filepath.Join(request.ProjectRoot, "image-settings"))
	if err != nil {
		t.Fatal(err)
	}
	wantSettings := "oci\n" + filepath.Join(request.RunDirectory, "images", "dummy.tar") +
		"\n\n\nlocal/patch-validation\nlinux/arm64\narm64\n"
	if string(output) != wantSettings {
		t.Fatalf("image settings = %q, want %q", output, wantSettings)
	}
}

func TestValidateRepairedProjectWithoutDockerfilesRunsOnlySourceTargets(t *testing.T) {
	request := repairValidationFixture(t, "")
	request.Dockerfiles = nil
	request.RunDirectory = ""
	result, err := validateRepairedProject(context.Background(), request)
	if err != nil {
		t.Fatalf("validateRepairedProject() error = %v", err)
	}
	want := repairValidationResult{Checks: []string{"attribution-checksums", "validate-checksums"}}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("validateRepairedProject() = %#v, want %#v", result, want)
	}
	assertRepairValidationCalls(t, request, "attribution-checksums", "validate-checksums")
}

func TestValidateRepairedProjectWithoutContainerImagesReportsSkippedBuild(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	t.Setenv("VALIDATION_IMAGE_NAMES", "")
	result, err := validateRepairedProject(context.Background(), request)
	if err != nil {
		t.Fatalf("validateRepairedProject() error = %v", err)
	}
	want := repairValidationResult{
		Checks:          []string{"dockerfile-syntax", "attribution-checksums", "validate-checksums"},
		ImageSkipReason: "project declares no container images",
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("validateRepairedProject() = %#v, want %#v", result, want)
	}
	assertRepairValidationCalls(t, request, "attribution-checksums", "validate-checksums", "var-value-IMAGE_NAMES")
	summary := patchValidationSummary(result)
	if !strings.Contains(summary, "image build skipped: project declares no container images") ||
		strings.Contains(summary, "native") {
		t.Fatalf("validation summary = %q, want explicit skip without a native build claim", summary)
	}
}

func TestValidateRepairedProjectUsesProjectImageValidationTarget(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	t.Setenv("VALIDATION_IMAGE_TARGET", "patch-validation-images")
	result, err := validateRepairedProject(context.Background(), request)
	if err != nil {
		t.Fatalf("validateRepairedProject() error = %v", err)
	}
	if !result.NativeImages || result.Architecture != "arm64" {
		t.Fatalf("validateRepairedProject() = %#v, want native arm64 image validation", result)
	}
	assertRepairValidationCalls(t, request, "attribution-checksums", "validate-checksums", "var-value-IMAGE_NAMES", "var-value-PATCH_VALIDATION_IMAGE_TARGET", "var-value-BUILDER_PLATFORM_ARCH", "patch-validation-images")
}

func TestValidateRepairedProjectDoesNotReloadCredentialsThroughMakeSettings(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	injected := filepath.Join(t.TempDir(), "injected.mk")
	content := "export GITHUB_TOKEN=dummy-token\nexport AWS_ACCESS_KEY_ID=dummy-key\nexport AWS_CONFIG_FILE=/dummy/config\n"
	if err := os.WriteFile(injected, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MAKEFLAGS", "-- GITHUB_TOKEN=dummy-token AWS_ACCESS_KEY_ID=dummy-key AWS_CONFIG_FILE=/dummy/config")
	t.Setenv("MFLAGS", "-- GITHUB_TOKEN=dummy-token")
	t.Setenv("MAKEOVERRIDES", "AWS_ACCESS_KEY_ID=dummy-key")
	t.Setenv("MAKEFILES", injected)
	if _, err := validateRepairedProject(context.Background(), request); err != nil {
		t.Fatalf("Make reintroduced credentials: %v", err)
	}
}

func TestValidateRepairedProjectRejectsSuccessfulTargetWithoutImageOutput(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	t.Setenv("VALIDATION_EMPTY_IMAGES", "true")
	result, err := validateRepairedProject(context.Background(), request)
	if err == nil || !strings.Contains(err.Error(), "without producing a local OCI archive") || result.NativeImages {
		t.Fatalf("empty image validation = %#v, error %v", result, err)
	}
}

func TestValidateRepairedProjectDeletedDockerfileStillRequiresImages(t *testing.T) {
	request := repairValidationFixture(t, "")
	request.Dockerfiles = []string{"deleted/Dockerfile"}
	result, err := validateRepairedProject(context.Background(), request)
	if err != nil {
		t.Fatalf("validateRepairedProject() error = %v", err)
	}
	if !result.NativeImages {
		t.Fatalf("validateRepairedProject() = %#v, want native image validation", result)
	}
	assertRepairValidationCalls(t, request, "attribution-checksums", "validate-checksums", "var-value-IMAGE_NAMES", "var-value-PATCH_VALIDATION_IMAGE_TARGET", "var-value-BUILDER_PLATFORM_ARCH", "local-images")
}

func TestValidateRepairedProjectParsesEveryDockerfileBeforeMake(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	if err := os.WriteFile(filepath.Join(request.UpstreamRepo, "Dockerfile.second"), []byte("FROM scratch\nilding runc\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	request.Dockerfiles = append(request.Dockerfiles, "Dockerfile.second")
	_, err := validateRepairedProject(context.Background(), request)
	var sourceError *repairSourceValidationError
	if !errors.As(err, &sourceError) || sourceError.Path != "Dockerfile.second" {
		t.Fatalf("validateRepairedProject() error = %v, want Dockerfile.second source error", err)
	}
	assertRepairValidationCalls(t, request)
}

func TestValidateRepairedProjectRejectsUnsafeDockerfilePaths(t *testing.T) {
	for _, name := range []string{"relative escape", "absolute escape", "file symlink", "directory symlink", "dangling symlink", "directory", "empty path"} {
		t.Run(name, func(t *testing.T) {
			request := repairValidationFixture(t, "FROM scratch\n")
			outside := filepath.Join(t.TempDir(), "Dockerfile")
			if err := os.WriteFile(outside, []byte("FROM scratch\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			switch name {
			case "relative escape":
				request.Dockerfiles = []string{"../Dockerfile"}
			case "absolute escape":
				request.Dockerfiles = []string{outside}
			case "file symlink", "dangling symlink":
				if name == "dangling symlink" {
					outside += ".deleted"
				}
				if err := os.Symlink(outside, filepath.Join(request.UpstreamRepo, "linked.Dockerfile")); err != nil {
					t.Fatal(err)
				}
				request.Dockerfiles = []string{"linked.Dockerfile"}
			case "directory symlink":
				if err := os.Symlink(filepath.Dir(outside), filepath.Join(request.UpstreamRepo, "linked")); err != nil {
					t.Fatal(err)
				}
				request.Dockerfiles = []string{"linked/Dockerfile"}
			case "directory":
				request.Dockerfiles = []string{"."}
			case "empty path":
				request.Dockerfiles = []string{""}
			}
			_, err := validateRepairedProject(context.Background(), request)
			var sourceError *repairSourceValidationError
			if !errors.As(err, &sourceError) {
				t.Fatalf("validateRepairedProject() error = %T %v, want *repairSourceValidationError", err, err)
			}
			assertRepairValidationCalls(t, request)
		})
	}
}

func TestValidateRepairedProjectAcceptsAbsoluteDockerfileWithinUpstream(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	request.Dockerfiles = []string{filepath.Join(request.UpstreamRepo, "Dockerfile")}
	if _, err := validateRepairedProject(context.Background(), request); err != nil {
		t.Fatalf("validateRepairedProject() error = %v", err)
	}
}

func TestValidateRepairedProjectAcceptsRelativeUpstreamDirectory(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	request.UpstreamRepo = "repo"
	if _, err := validateRepairedProject(context.Background(), request); err != nil {
		t.Fatalf("validateRepairedProject() error = %v", err)
	}
}

func TestValidateRepairedProjectRequiresImageRunDirectoryBeforeMake(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	request.RunDirectory = ""
	if _, err := validateRepairedProject(context.Background(), request); err == nil || !strings.Contains(err.Error(), "run directory") {
		t.Fatalf("validateRepairedProject() error = %v, want missing run directory", err)
	}
	assertRepairValidationCalls(t, request)
}

func TestValidateRepairedProjectStopsAtFailedMakeStage(t *testing.T) {
	stages := []string{"attribution-checksums", "validate-checksums", "var-value-IMAGE_NAMES", "var-value-PATCH_VALIDATION_IMAGE_TARGET", "var-value-BUILDER_PLATFORM_ARCH", "local-images"}
	for index, stage := range stages {
		t.Run(stage, func(t *testing.T) {
			request := repairValidationFixture(t, "FROM scratch\n")
			t.Setenv("VALIDATION_FAIL_STAGE", stage)
			result, err := validateRepairedProject(context.Background(), request)
			var exitError *exec.ExitError
			if !errors.As(err, &exitError) {
				t.Fatalf("validateRepairedProject() error = %T %v, want command failure", err, err)
			}
			var projectError *repairProjectValidationError
			if !errors.As(err, &projectError) || projectError.Stage != stage {
				t.Fatalf("validateRepairedProject() error = %v, want project validation stage %q", err, stage)
			}
			if !strings.Contains(projectError.Detail, "dummy failure output") {
				t.Fatalf("validation detail = %q, want command output", projectError.Detail)
			}
			if !strings.Contains(err.Error(), stage) || !strings.Contains(err.Error(), "dummy failure output") {
				t.Fatalf("validateRepairedProject() error = %v, want stage and output", err)
			}
			if result.NativeImages {
				t.Fatalf("validateRepairedProject() = %#v, want no successful image validation", result)
			}
			assertRepairValidationCalls(t, request, stages[:index+1]...)
		})
	}
}

func TestRepairValidationMakeStagesHonorDeadline(t *testing.T) {
	stages := []string{"attribution-checksums", "validate-checksums", "var-value-IMAGE_NAMES", "var-value-PATCH_VALIDATION_IMAGE_TARGET", "var-value-BUILDER_PLATFORM_ARCH", "local-images"}
	for _, stage := range stages {
		t.Run(stage, func(t *testing.T) {
			request := repairValidationFixture(t, "FROM scratch\n")
			t.Setenv("VALIDATION_WAIT_STAGE", stage)
			// Exercise cancellation at this stage without spending the deadline
			// on preceding Make queries or unrelated host startup work.
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			start := time.Now()
			_, err := runRepairValidationMake(ctx, request.ProjectRoot, stage)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("runRepairValidationMake() error = %v, want deadline exceeded", err)
			}
			var projectError *repairProjectValidationError
			if !errors.As(err, &projectError) || projectError.Stage != stage {
				t.Fatalf("runRepairValidationMake() error = %v, want project validation stage %q", err, stage)
			}
			if !strings.Contains(err.Error(), stage) {
				t.Fatalf("runRepairValidationMake() error = %v, want stage %q", err, stage)
			}
			if elapsed := time.Since(start); elapsed > 5*time.Second {
				t.Fatalf("validation deadline took too long: %s", elapsed)
			}
			assertRepairValidationCalls(t, request, stage)
		})
	}
}

func TestValidateRepairedProjectHonorsCanceledContext(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := validateRepairedProject(ctx, request); !errors.Is(err, context.Canceled) {
		t.Fatalf("validateRepairedProject() error = %v, want context canceled", err)
	}
	assertRepairValidationCalls(t, request)
}

func TestValidateRepairedProjectRejectsUnknownNativeArchitecture(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	t.Setenv("VALIDATION_ARCH", "dummy-arch")
	result, err := validateRepairedProject(context.Background(), request)
	if err == nil || !strings.Contains(err.Error(), "unsupported native architecture") {
		t.Fatalf("validateRepairedProject() error = %v, want unsupported architecture", err)
	}
	var projectError *repairProjectValidationError
	if !errors.As(err, &projectError) || projectError.Stage != "var-value-BUILDER_PLATFORM_ARCH" {
		t.Fatalf("validateRepairedProject() error = %v, want architecture validation failure", err)
	}
	if result.NativeImages || result.Architecture != "" {
		t.Fatalf("validateRepairedProject() = %#v, want no image validation", result)
	}
	assertRepairValidationCalls(t, request, "attribution-checksums", "validate-checksums", "var-value-IMAGE_NAMES", "var-value-PATCH_VALIDATION_IMAGE_TARGET", "var-value-BUILDER_PLATFORM_ARCH")
}

func TestValidateRepairedProjectReportsGoToolchainPreflightFailure(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	if err := os.WriteFile(filepath.Join(request.UpstreamRepo, "go.mod"), []byte("module example.invalid/dummy\n\ngo 1.25.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	script := "#!/bin/sh\nprintf 'dummy Go toolchain unavailable\\n' >&2\nexit 1\n"
	if err := os.WriteFile(filepath.Join(bin, "go"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	_, err := validateRepairedProject(context.Background(), request)
	var projectError *repairProjectValidationError
	if !errors.As(err, &projectError) || projectError.Stage != "Go toolchain preflight" {
		t.Fatalf("validateRepairedProject() error = %v, want Go toolchain validation failure", err)
	}
	for _, want := range []string{"Go toolchain preflight", "dummy Go toolchain unavailable", "builder-base"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("validateRepairedProject() error = %v, want %q", err, want)
		}
	}
	assertRepairValidationCalls(t, request)
}

func TestValidateRepairedProjectDiagnosesIncompatibleGoToolchain(t *testing.T) {
	request := repairValidationFixture(t, "FROM scratch\n")
	t.Setenv("VALIDATION_FAIL_STAGE", "attribution-checksums")
	t.Setenv("VALIDATION_FAILURE_OUTPUT", "go: go.mod requires go >= 1.25.0 (running go 1.21.0; GOTOOLCHAIN=local)")
	_, err := validateRepairedProject(context.Background(), request)
	var projectError *repairProjectValidationError
	if !errors.As(err, &projectError) || projectError.Stage != "attribution-checksums" {
		t.Fatalf("validateRepairedProject() error = %v, want attribution validation failure", err)
	}
	for _, want := range []string{"incompatible Go toolchain", "go.mod requires go >= 1.25.0", "builder-base"} {
		if !strings.Contains(projectError.Detail, want) {
			t.Fatalf("validation detail = %q, want %q", projectError.Detail, want)
		}
	}
	assertRepairValidationCalls(t, request, "attribution-checksums")
}

func repairValidationFixture(t *testing.T, dockerfile string) repairValidationRequest {
	t.Helper()
	projectRoot := filepath.Join(t.TempDir(), "dummy project")
	if err := os.Mkdir(projectRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	upstream := filepath.Join(projectRoot, "repo")
	if err := os.Mkdir(upstream, 0o755); err != nil {
		t.Fatal(err)
	}
	if dockerfile != "" {
		if err := os.WriteFile(filepath.Join(upstream, "Dockerfile"), []byte(dockerfile), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("MAKEFLAGS", "")
	t.Setenv("MFLAGS", "")
	t.Setenv("MAKEOVERRIDES", "")
	t.Setenv("VALIDATION_FAIL_STAGE", "")
	t.Setenv("VALIDATION_FAILURE_OUTPUT", "dummy failure output")
	t.Setenv("VALIDATION_WAIT_STAGE", "")
	t.Setenv("VALIDATION_ARCH", "arm64")
	t.Setenv("VALIDATION_EMPTY_IMAGES", "")
	t.Setenv("VALIDATION_IMAGE_NAMES", "dummy")
	t.Setenv("VALIDATION_IMAGE_TARGET", "local-images")
	makefile := `BUILDER_PLATFORM_ARCH=$(VALIDATION_ARCH)
IMAGE_OUTPUT_NAME=dummy
IMAGE_NAMES=$(VALIDATION_IMAGE_NAMES)
PATCH_VALIDATION_IMAGE_TARGET=$(VALIDATION_IMAGE_TARGET)
.PHONY: attribution-checksums validate-checksums var-value-IMAGE_NAMES var-value-PATCH_VALIDATION_IMAGE_TARGET var-value-BUILDER_PLATFORM_ARCH local-images patch-validation-images build images
attribution-checksums validate-checksums var-value-IMAGE_NAMES var-value-PATCH_VALIDATION_IMAGE_TARGET var-value-BUILDER_PLATFORM_ARCH local-images patch-validation-images:
	@echo '$@' >> calls
	@test -z "$$AWS_ACCESS_KEY_ID" && test -z "$$GITHUB_TOKEN" && test "$$AWS_EC2_METADATA_DISABLED" = true && test "$$CODEBUILD_CI" = false && test "$$AWS_CONFIG_FILE" = /dev/null && test "$$AWS_SHARED_CREDENTIALS_FILE" = /dev/null
	@if [ "$(VALIDATION_FAIL_STAGE)" = "$@" ]; then echo "$$VALIDATION_FAILURE_OUTPUT" >&2; exit 1; fi
	@if [ "$(VALIDATION_WAIT_STAGE)" = "$@" ]; then sleep 5; fi
	@if [ "$@" = var-value-IMAGE_NAMES ]; then echo "$(IMAGE_NAMES)"; fi
	@if [ "$@" = var-value-PATCH_VALIDATION_IMAGE_TARGET ]; then echo "$(PATCH_VALIDATION_IMAGE_TARGET)"; fi
	@if [ "$@" = var-value-BUILDER_PLATFORM_ARCH ]; then echo "dummy make warning" >&2; echo "$(BUILDER_PLATFORM_ARCH)"; fi
	@if [ "$@" = local-images ] || [ "$@" = patch-validation-images ]; then printf '%s\n' '$(IMAGE_OUTPUT_TYPE)' '$(patsubst dest=%,%,$(IMAGE_OUTPUT))' '$(IMAGE_EXPORT_CACHE)' '$(IMAGE_IMPORT_CACHE)' '$(IMAGE_REPO)' '$(IMAGE_PLATFORMS)' '$(BUILDER_PLATFORM_ARCH)' > image-settings; fi
	@if { [ "$@" = local-images ] || [ "$@" = patch-validation-images ]; } && [ "$(VALIDATION_EMPTY_IMAGES)" != true ]; then mkdir -p "$(IMAGE_OUTPUT_DIR)"; printf 'dummy OCI archive\n' > "$(IMAGE_OUTPUT_DIR)/dummy.tar"; fi
build images:
	@echo 'push target must not run' >> calls
	@exit 1
`
	if err := os.WriteFile(filepath.Join(projectRoot, "Makefile"), []byte(makefile), 0o644); err != nil {
		t.Fatal(err)
	}
	return repairValidationRequest{
		ProjectRoot:  projectRoot,
		UpstreamRepo: upstream,
		RunDirectory: t.TempDir(),
		Dockerfiles:  []string{"Dockerfile"},
	}
}

func assertRepairValidationCalls(t *testing.T, request repairValidationRequest, want ...string) {
	t.Helper()
	output, err := os.ReadFile(filepath.Join(request.ProjectRoot, "calls"))
	if os.IsNotExist(err) && len(want) == 0 {
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(output)); got != strings.Join(want, "\n") {
		t.Fatalf("make calls = %q, want %q", got, strings.Join(want, "\n"))
	}
}

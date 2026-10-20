package upgrade

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type repairValidationRequest struct {
	ProjectRoot  string
	UpstreamRepo string
	RunDirectory string
	// Dockerfiles contains upstream-relative paths or absolute paths within UpstreamRepo.
	Dockerfiles []string
}

type repairValidationResult struct {
	Checks          []string `json:"checks"`
	NativeImages    bool     `json:"native_images"`
	Architecture    string   `json:"architecture,omitempty"`
	ImageSkipReason string   `json:"image_skip_reason,omitempty"`
}

// repairSourceValidationError identifies invalid repaired source, rather than a failed build.
type repairSourceValidationError struct {
	Path   string
	Detail string
}

func (e *repairSourceValidationError) Error() string {
	return fmt.Sprintf("validating repaired source %q: %s", e.Path, e.Detail)
}

// repairProjectValidationError keeps project gates distinct from patch application failures.
type repairProjectValidationError struct {
	Stage  string
	Detail string
	Cause  error
}

func (e *repairProjectValidationError) Error() string {
	detail := e.Detail
	if detail == "" && e.Cause != nil {
		detail = e.Cause.Error()
	}
	return fmt.Sprintf("%s validation: %s", e.Stage, detail)
}

func (e *repairProjectValidationError) Unwrap() error {
	return e.Cause
}

func validateRepairedProject(ctx context.Context, request repairValidationRequest) (repairValidationResult, error) {
	result := repairValidationResult{}
	if err := ctx.Err(); err != nil {
		return result, &repairProjectValidationError{Stage: "preflight", Detail: err.Error(), Cause: err}
	}
	if request.ProjectRoot == "" || request.UpstreamRepo == "" {
		return result, &repairProjectValidationError{Stage: "preflight", Detail: "repaired-project validation requires project and upstream directories"}
	}
	upstream := request.UpstreamRepo
	if !filepath.IsAbs(upstream) {
		upstream = filepath.Join(request.ProjectRoot, upstream)
	}
	info, err := os.Lstat(upstream)
	if err != nil {
		return result, &repairSourceValidationError{Path: request.UpstreamRepo, Detail: err.Error()}
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return result, &repairSourceValidationError{Path: request.UpstreamRepo, Detail: "upstream root must be a directory, without a symlink"}
	}
	upstream, err = filepath.Abs(upstream)
	if err != nil {
		return result, &repairSourceValidationError{Path: request.UpstreamRepo, Detail: err.Error()}
	}
	for _, path := range request.Dockerfiles {
		if err := ctx.Err(); err != nil {
			return result, &repairProjectValidationError{Stage: "dockerfile-syntax", Detail: err.Error(), Cause: err}
		}
		if err := validateRepairDockerfile(upstream, path); err != nil {
			return result, err
		}
	}
	var runDirectory string
	if len(request.Dockerfiles) > 0 {
		result.Checks = append(result.Checks, "dockerfile-syntax")
		if request.RunDirectory == "" {
			return result, &repairProjectValidationError{Stage: "local-images", Detail: "image validation requires a run directory"}
		}
		runDirectory, err = filepath.Abs(request.RunDirectory)
		if err != nil {
			return result, &repairProjectValidationError{Stage: "local-images", Detail: "resolving image output directory: " + err.Error(), Cause: err}
		}
	}

	// Querying the module detects missing or incompatible Go toolchains before
	// source targets run. Projects without a root Go module still use make.
	if info, err := os.Lstat(filepath.Join(upstream, "go.mod")); err == nil {
		if !info.Mode().IsRegular() {
			return result, &repairSourceValidationError{Path: "go.mod", Detail: "Go module must be a regular file, without a symlink"}
		}
		if _, err := runRepairValidationCommand(ctx, "Go toolchain preflight", upstream, "go", "list", "-m"); err != nil {
			return result, err
		}
	} else if !os.IsNotExist(err) {
		return result, &repairSourceValidationError{Path: filepath.Join(upstream, "go.mod"), Detail: err.Error()}
	}

	for _, target := range []string{"attribution-checksums", "validate-checksums"} {
		if _, err := runRepairValidationMake(ctx, request.ProjectRoot, target); err != nil {
			return result, err
		}
		result.Checks = append(result.Checks, target)
	}
	if len(request.Dockerfiles) == 0 {
		return result, nil
	}
	output, err := runRepairValidationMake(ctx, request.ProjectRoot, "var-value-IMAGE_NAMES")
	if err != nil {
		return result, err
	}
	if strings.TrimSpace(output) == "" {
		result.ImageSkipReason = "project declares no container images"
		logger.Info("Skipping repaired-project image build", "reason", result.ImageSkipReason)
		return result, nil
	}
	output, err = runRepairValidationMake(ctx, request.ProjectRoot, "var-value-PATCH_VALIDATION_IMAGE_TARGET")
	if err != nil {
		return result, err
	}
	imageTarget := strings.TrimSpace(output)
	if imageTarget == "" {
		return result, &repairProjectValidationError{Stage: "local-images", Detail: "project has no patch validation image target"}
	}
	output, err = runRepairValidationMake(ctx, request.ProjectRoot, "var-value-BUILDER_PLATFORM_ARCH")
	if err != nil {
		return result, err
	}
	architecture := strings.TrimSpace(output)
	if architecture != "amd64" && architecture != "arm64" {
		return result, &repairProjectValidationError{
			Stage:  "var-value-BUILDER_PLATFORM_ARCH",
			Detail: fmt.Sprintf("unsupported native architecture %q", architecture),
		}
	}
	result.Architecture = architecture
	if _, err := runRepairValidationMake(ctx, request.ProjectRoot, imageTarget,
		"BUILDER_PLATFORM_ARCH="+architecture,
		"IMAGE_PLATFORMS=linux/"+architecture,
		"IMAGE_OUTPUT_TYPE=oci",
		"IMAGE_OUTPUT=dest=$(IMAGE_OUTPUT_DIR)/$(IMAGE_OUTPUT_NAME).tar",
		"IMAGE_OUTPUT_DIR="+filepath.Join(runDirectory, "images"),
		"IMAGE_EXPORT_CACHE=",
		"IMAGE_IMPORT_CACHE=",
		"IMAGE_REPO=local/patch-validation",
	); err != nil {
		return result, err
	}
	archives, err := filepath.Glob(filepath.Join(runDirectory, "images", "*.tar"))
	if err != nil {
		return result, &repairProjectValidationError{Stage: "local-images", Detail: err.Error(), Cause: err}
	}
	hasArchive := false
	for _, archive := range archives {
		info, err := os.Lstat(archive)
		if err != nil {
			return result, &repairProjectValidationError{Stage: "local-images", Detail: err.Error(), Cause: err}
		}
		if info.Mode().IsRegular() && info.Size() > 0 {
			hasArchive = true
		}
	}
	if !hasArchive {
		return result, &repairProjectValidationError{Stage: "local-images", Detail: "target succeeded without producing a local OCI archive"}
	}
	result.Checks = append(result.Checks, "local-images")
	result.NativeImages = true
	return result, nil
}

func validateRepairDockerfile(root, path string) error {
	fail := func(detail string) error {
		return &repairSourceValidationError{Path: path, Detail: detail}
	}
	if path == "" {
		return fail("empty Dockerfile path")
	}
	fullPath := path
	if !filepath.IsAbs(fullPath) {
		fullPath = filepath.Join(root, fullPath)
	}
	relative, err := filepath.Rel(root, fullPath)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fail("Dockerfile path escapes upstream root")
	}
	// Check every component before opening, including dangling symlinks. Missing
	// paths may have been intentionally deleted by the final patch series.
	current := root
	for _, component := range strings.Split(relative, string(filepath.Separator)) {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			logger.Info("Skipping deleted repaired Dockerfile", "path", path)
			return nil
		}
		if err != nil {
			return fail(err.Error())
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fail("Dockerfile path contains a symlink")
		}
		if current == filepath.Clean(fullPath) && !info.Mode().IsRegular() {
			return fail("Dockerfile is not a regular file")
		}
	}
	file, err := os.OpenInRoot(root, relative)
	if err != nil {
		return fail(err.Error())
	}
	defer file.Close()
	parsed, err := parser.Parse(file)
	if err != nil {
		return fail(err.Error())
	}
	stages, _, err := instructions.Parse(parsed.AST, nil)
	if err != nil {
		return fail(err.Error())
	}
	if len(stages) == 0 {
		return fail("Dockerfile has no build stage")
	}
	return nil
}

func runRepairValidationMake(ctx context.Context, projectRoot, target string, variables ...string) (string, error) {
	return runProjectMake(ctx, projectRoot, target, patchValidationEnvironment(), variables...)
}

func runProjectMake(ctx context.Context, projectRoot, target string, environment []string, variables ...string) (string, error) {
	args := append([]string{"--no-print-directory", "-C", projectRoot, target}, variables...)
	return runProjectCommand(ctx, target, "", "make", environment, args...)
}

func runRepairValidationCommand(ctx context.Context, stage, directory, name string, args ...string) (string, error) {
	return runProjectCommand(ctx, stage, directory, name, patchValidationEnvironment(), args...)
}

func runProjectCommand(ctx context.Context, stage, directory, name string, environment []string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = directory
	cmd.Env = environment
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return os.ErrProcessDone
		}
		return err
	}
	// make recipes can retain output pipes briefly after make is canceled.
	cmd.WaitDelay = 250 * time.Millisecond
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	text := strings.TrimSpace(stdout.String())
	errorText := strings.TrimSpace(stderr.String())
	logger.Info("Project command", "stage", stage, "stdout", text, "stderr", errorText, "error", err)
	if err != nil {
		output := strings.TrimSpace(text + "\n" + errorText)
		failure := &repairProjectValidationError{
			Stage:  stage,
			Detail: fmt.Sprintf("%v\nOutput: %s", err, output),
			Cause:  err,
		}
		if ctx.Err() != nil {
			failure.Detail = fmt.Sprintf("%v (%v)\nOutput: %s", ctx.Err(), err, output)
			failure.Cause = ctx.Err()
		} else if strings.Contains(output, "go.mod requires go >=") {
			failure.Detail = "incompatible Go toolchain: " + failure.Detail +
				"\nUse the project's required Go version in builder-base or the lab environment."
		} else if stage == "Go toolchain preflight" {
			failure.Detail += "\nUse the project Go toolchain in builder-base or the lab environment."
		}
		return output, failure
	}
	return text, nil
}

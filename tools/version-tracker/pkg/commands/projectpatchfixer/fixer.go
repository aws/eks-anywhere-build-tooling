package projectpatchfixer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const autoscalerCommandEnv = "PATCH_FIXER_AUTOSCALER_COMMAND"

// Request contains the exact project state available when patch application fails.
type Request struct {
	ProjectName    string
	ProjectRoot    string
	UpstreamRepo   string
	PatchesDir     string
	TargetRevision string
	ReleaseBranch  string
}

// Result describes a deterministic project-specific repair.
type Result struct {
	Handled      bool
	SkipAgent    bool
	Reason       string
	PatchesDir   string
	PatchCount   int
	RunDirectory string
}

type fixer func(context.Context, Request) (Result, error)

var fixers = map[string]fixer{
	"kubernetes/autoscaler": fixAutoscaler,
}

var skippedAgentProjects = map[string]string{
	"goharbor/harbor":             "swagger patches contain generated output and require deterministic regeneration",
	"kubernetes-sigs/cluster-api": "patch series is intentionally excluded pending a dedicated validation run",
}

// Run executes a project-specific deterministic fixer when one is registered.
func Run(ctx context.Context, request Request) (Result, error) {
	if reason, ok := skippedAgentProjects[request.ProjectName]; ok {
		return Result{Handled: true, SkipAgent: true, Reason: reason}, nil
	}
	fix, ok := fixers[request.ProjectName]
	if !ok {
		return Result{Handled: false}, nil
	}
	return fix(ctx, request)
}

func fixAutoscaler(ctx context.Context, request Request) (Result, error) {
	runDir, err := newRunDirectory(request.ProjectName)
	if err != nil {
		return Result{Handled: true}, err
	}
	outputDir := filepath.Join(runDir, "patches")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return Result{Handled: true}, fmt.Errorf("creating autoscaler patch output directory: %v", err)
	}
	manifestPath := filepath.Join(runDir, "manifest.json")

	command := autoscalerCommand()
	args := append(command[1:],
		"--source", request.UpstreamRepo,
		"--patches", request.PatchesDir,
		"--output", outputDir,
		"--manifest", manifestPath,
		"--revision", request.TargetRevision,
		"--release-branch", request.ReleaseBranch,
	)
	cmd := exec.CommandContext(ctx, command[0], args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf(
			"regenerating autoscaler patches: %v\nOutput: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("reading autoscaler patch manifest: %v", err)
	}
	var manifest struct {
		Patches []string `json:"patches"`
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("parsing autoscaler patch manifest: %v", err)
	}
	if len(manifest.Patches) == 0 {
		return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("autoscaler regenerator produced no patches")
	}
	for _, patch := range manifest.Patches {
		if filepath.Dir(patch) != outputDir {
			return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("autoscaler patch escaped output directory: %s", patch)
		}
		if _, err := os.Stat(patch); err != nil {
			return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("autoscaler patch %q is not accessible: %v", patch, err)
		}
	}

	return Result{
		Handled:      true,
		PatchesDir:   outputDir,
		PatchCount:   len(manifest.Patches),
		RunDirectory: runDir,
	}, nil
}

func autoscalerCommand() []string {
	if configured := strings.TrimSpace(os.Getenv(autoscalerCommandEnv)); configured != "" {
		return strings.Fields(configured)
	}
	candidates := []string{
		filepath.Join("..", "..", "projects", "kubernetes", "autoscaler", "regenerate_patches.py"),
		filepath.Join("projects", "kubernetes", "autoscaler", "regenerate_patches.py"),
	}
	for _, candidate := range candidates {
		if absolute, err := filepath.Abs(candidate); err == nil {
			if _, err := os.Stat(absolute); err == nil {
				return []string{"python3", absolute}
			}
		}
	}
	return []string{"python3", candidates[0]}
}

func newRunDirectory(projectName string) (string, error) {
	root := strings.TrimSpace(os.Getenv("PATCH_FIXER_RESULT_DIR"))
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("getting working directory for patch fixer results: %v", err)
		}
		root = filepath.Join(cwd, "patch-fixer-runs")
	}
	projectDir := strings.NewReplacer("/", "_", "\\", "_").Replace(projectName)
	runDir := filepath.Join(root, projectDir, time.Now().UTC().Format("20060102T150405.000000000Z"))
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return "", fmt.Errorf("creating patch fixer run directory: %v", err)
	}
	return runDir, nil
}

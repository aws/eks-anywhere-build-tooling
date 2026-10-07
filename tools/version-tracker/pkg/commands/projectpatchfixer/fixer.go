package projectpatchfixer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Request contains the exact project state available when patch application fails.
type Request struct {
	ProjectName    string
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

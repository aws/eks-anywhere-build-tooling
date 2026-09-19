package agentpatchfixer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

const (
	EnabledEnv     = "PATCH_FIXER_ENABLED"
	ResultDirEnv   = "PATCH_FIXER_RESULT_DIR"
	TimeoutEnv     = "PATCH_FIXER_AGENT_TIMEOUT"
	ModelIDEnv     = "PATCH_FIXER_MODEL_ID"
	DiagnosticsEnv = "PATCH_FIXER_DIAGNOSTICS"

	requestSchemaVersion  = 1
	defaultTimeout        = 10 * time.Minute
	defaultMaxTurns       = 20
	defaultMaxTotalTokens = 1000000
	defaultModelID        = "global.anthropic.claude-sonnet-4-6"
)

var defaultModelFactory modelFactory = newBedrockModel

// Request is the immutable input passed from the upgrade workflow to the agent.
type Request struct {
	SchemaVersion    int      `json:"schema_version"`
	ProjectName      string   `json:"project_name"`
	ProjectRoot      string   `json:"project_root"`
	UpstreamRepoPath string   `json:"upstream_repo_path"`
	FailedPatchPath  string   `json:"failed_patch_path"`
	FailedFiles      []string `json:"failed_files"`
	EditableFiles    []string `json:"editable_files,omitempty"`
	FailureOutput    string   `json:"failure_output"`
	CurrentRevision  string   `json:"current_revision"`
	TargetRevision   string   `json:"target_revision"`
	ReleaseBranch    string   `json:"release_branch,omitempty"`
	MaxTurns         int      `json:"max_turns"`
	MaxTotalTokens   int      `json:"max_total_tokens"`
}

// Result is emitted after producing a candidate patch.
type Result struct {
	SchemaVersion      int      `json:"schema_version"`
	Status             string   `json:"status"`
	StopReason         string   `json:"stop_reason,omitempty"`
	Summary            string   `json:"summary,omitempty"`
	CandidatePatchPath string   `json:"candidate_patch_path,omitempty"`
	EditedFiles        []string `json:"edited_files,omitempty"`
	Turns              int      `json:"turns,omitempty"`
	TotalTokens        int      `json:"total_tokens,omitempty"`
	RunDirectory       string   `json:"-"`
}

// Enabled returns whether agent patch fixing is enabled for this upgrade run.
func Enabled() bool {
	enabled, err := strconv.ParseBool(strings.TrimSpace(os.Getenv(EnabledEnv)))
	return err == nil && enabled
}

// Run executes the native Go patch repair agent.
func Run(ctx context.Context, request Request) (*Result, error) {
	if !Enabled() {
		return nil, nil
	}

	request.SchemaVersion = requestSchemaVersion
	request.FailureOutput = boundedTail(request.FailureOutput, 64*1024)
	if request.MaxTurns <= 0 {
		request.MaxTurns = defaultMaxTurns
	}
	if request.MaxTotalTokens <= 0 {
		request.MaxTotalTokens = defaultMaxTotalTokens
	}
	if err := validateRequest(request); err != nil {
		return nil, err
	}

	runDir, err := resultDirectory(request.ProjectName)
	if err != nil {
		return nil, err
	}
	if err := writeJSON(filepath.Join(runDir, "request.json"), request); err != nil {
		return nil, fmt.Errorf("writing patch fixer request: %v", err)
	}

	duration, err := timeout()
	if err != nil {
		return nil, err
	}
	runCtx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	result, err := execute(runCtx, request, runDir, defaultModelFactory)
	if err != nil {
		failed := &Result{
			SchemaVersion: requestSchemaVersion,
			Status:        "runner_failed",
			Summary:       boundedTail(err.Error(), 8*1024),
			RunDirectory:  runDir,
		}
		if runCtx.Err() == context.DeadlineExceeded {
			failed.Status = "runner_timeout"
			failed.Summary = fmt.Sprintf("agent patch fixer timed out after %s", duration)
		}
		_ = writeJSON(filepath.Join(runDir, "result.json"), failed)
		return failed, err
	}
	result.RunDirectory = runDir
	if err := writeJSON(filepath.Join(runDir, "result.json"), result); err != nil {
		return nil, fmt.Errorf("writing patch fixer result: %v", err)
	}

	logger.Info("Agent patch fixer attempt complete",
		"project", request.ProjectName,
		"status", result.Status,
		"stop_reason", result.StopReason,
		"candidate_patch", result.CandidatePatchPath,
		"result_dir", runDir)
	return result, nil
}

func validateRequest(request Request) error {
	requiredPaths := map[string]string{
		"project root":      request.ProjectRoot,
		"upstream repo":     request.UpstreamRepoPath,
		"failed patch path": request.FailedPatchPath,
	}
	for name, path := range requiredPaths {
		if strings.TrimSpace(path) == "" {
			return fmt.Errorf("%s is required", name)
		}
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("%s %q is not accessible: %v", name, path, err)
		}
	}
	if strings.Count(request.ProjectName, "/") != 1 {
		return fmt.Errorf("invalid project name %q", request.ProjectName)
	}
	for _, editableFile := range request.EditableFiles {
		if filepath.IsAbs(editableFile) || pathEscapesRoot(editableFile) {
			return fmt.Errorf("invalid editable file path %q", editableFile)
		}
	}
	return nil
}

func resultDirectory(projectName string) (string, error) {
	root := strings.TrimSpace(os.Getenv(ResultDirEnv))
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
		return "", fmt.Errorf("creating patch fixer result directory: %v", err)
	}
	return runDir, nil
}

func timeout() (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(TimeoutEnv))
	if value == "" {
		return defaultTimeout, nil
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parsing %s: %v", TimeoutEnv, err)
	}
	return duration, nil
}

func modelID() string {
	if value := strings.TrimSpace(os.Getenv(ModelIDEnv)); value != "" {
		return value
	}
	return defaultModelID
}

func diagnosticsFull() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv(DiagnosticsEnv)), "full")
}

func writeJSON(path string, value any) error {
	content, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')
	return os.WriteFile(path, content, 0o600)
}

func boundedTail(value string, maxBytes int) string {
	if len(value) <= maxBytes {
		return value
	}
	return value[len(value)-maxBytes:]
}

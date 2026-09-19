package agentpatchfixer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

const (
	EnabledEnv      = "PATCH_FIXER_ENABLED"
	AgentCommandEnv = "PATCH_FIXER_AGENT_COMMAND"
	ResultDirEnv    = "PATCH_FIXER_RESULT_DIR"
	TimeoutEnv      = "PATCH_FIXER_AGENT_TIMEOUT"
	ModelIDEnv      = "PATCH_FIXER_MODEL_ID"
	DiagnosticsEnv  = "PATCH_FIXER_DIAGNOSTICS"

	requestSchemaVersion  = 1
	defaultTimeout        = 10 * time.Minute
	defaultMaxTurns       = 20
	defaultMaxTotalTokens = 1000000
	defaultModelID        = "global.anthropic.claude-opus-5-5"
)

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

// Run executes the constrained Strands patch repair agent.
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

	duration, err := Timeout()
	if err != nil {
		return nil, err
	}
	runCtx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	result, err := execute(runCtx, request, runDir)
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

func runAgentProcess(ctx context.Context, invocation agentInvocation, runDir string) (agentOutcome, error) {
	requestPath := filepath.Join(runDir, "agent-request.json")
	resultPath := filepath.Join(runDir, "agent-result.json")
	if err := writeJSON(requestPath, invocation); err != nil {
		return agentOutcome{}, fmt.Errorf("writing Strands agent request: %v", err)
	}

	command := agentCommand()
	args := append(command[1:], "--request", requestPath, "--result", resultPath)
	cmd := exec.CommandContext(ctx, command[0], args...)
	cmd.Env = agentEnvironment(runDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	cmd.WaitDelay = 5 * time.Second
	output, runErr := cmd.CombinedOutput()

	var outcome agentOutcome
	resultBytes, readErr := os.ReadFile(resultPath)
	if readErr == nil {
		if err := json.Unmarshal(resultBytes, &outcome); err != nil {
			return agentOutcome{}, fmt.Errorf("parsing Strands agent result: %v", err)
		}
	}
	if ctx.Err() != nil {
		return agentOutcome{}, ctx.Err()
	}
	if runErr != nil {
		detail := strings.TrimSpace(string(output))
		if outcome.Summary != "" {
			detail = outcome.Summary
		}
		return agentOutcome{}, fmt.Errorf("running Strands agent: %v\nOutput: %s", runErr, boundedTail(detail, 8*1024))
	}
	if readErr != nil {
		return agentOutcome{}, fmt.Errorf("reading Strands agent result: %v\nOutput: %s", readErr, boundedTail(string(output), 8*1024))
	}
	if outcome.SchemaVersion != requestSchemaVersion {
		return agentOutcome{}, fmt.Errorf("unsupported Strands agent result schema version %d", outcome.SchemaVersion)
	}
	if outcome.Status != "completed" {
		return agentOutcome{}, fmt.Errorf("Strands agent failed: %s", outcome.Summary)
	}
	return outcome, nil
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

func agentCommand() []string {
	if configured := strings.TrimSpace(os.Getenv(AgentCommandEnv)); configured != "" {
		return strings.Fields(configured)
	}
	candidates := []string{
		filepath.Join("agent", "run.sh"),
		filepath.Join("tools", "version-tracker", "agent", "run.sh"),
	}
	for _, candidate := range candidates {
		if absolute, err := filepath.Abs(candidate); err == nil {
			if _, err := os.Stat(absolute); err == nil {
				return []string{absolute}
			}
		}
	}
	return []string{candidates[0]}
}

// Timeout returns the configured upper bound for one repair attempt.
func Timeout() (time.Duration, error) {
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

func agentEnvironment(runDir string) []string {
	allowed := map[string]bool{
		"PATH": true, "HOME": true, "TMPDIR": true, "LANG": true, "LC_ALL": true,
		"HTTP_PROXY": true, "HTTPS_PROXY": true, "NO_PROXY": true,
		"AWS_REGION": true, "AWS_DEFAULT_REGION": true, "AWS_CA_BUNDLE": true,
		"AWS_CONTAINER_CREDENTIALS_RELATIVE_URI": true,
		"AWS_CONTAINER_CREDENTIALS_FULL_URI":     true,
		"AWS_CONTAINER_AUTHORIZATION_TOKEN_FILE": true,
		"PATCH_FIXER_MODEL_ID":                   true,
		"PATCH_FIXER_BOOTSTRAP_DIR":              true,
		"PATCH_FIXER_BOOTSTRAP_TIMEOUT_SECONDS":  true,
		"PATCH_FIXER_DEPENDENCY_DIR":             true,
		"PATCH_FIXER_PYTHON":                     true,
		"PATCH_FIXER_PYTHON_INSTALL_DIR":         true,
		"PATCH_FIXER_PYTHON_TIMEOUT_SECONDS":     true,
		"PATCH_FIXER_PYTHON_VERSION":             true,
		"PATCH_FIXER_DIAGNOSTICS":                true,
		"PATCH_FIXER_UV_CACHE_DIR":               true,
		"PIP_INDEX_URL":                          true,
		"PIP_EXTRA_INDEX_URL":                    true,
		"PIP_TRUSTED_HOST":                       true,
		"PYTHONPATH":                             true,
		"SSL_CERT_FILE":                          true,
		"REQUESTS_CA_BUNDLE":                     true,
	}
	environment := make([]string, 0, len(allowed)+1)
	for _, entry := range os.Environ() {
		name, _, found := strings.Cut(entry, "=")
		if found && allowed[name] {
			environment = append(environment, entry)
		}
	}
	return append(environment, fmt.Sprintf("PATCH_FIXER_RUN_DIR=%s", runDir))
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

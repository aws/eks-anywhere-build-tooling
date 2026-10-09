package upgrade

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/commands/agentpatchfixer"
)

func patchValidationRejectionComment(err error) (string, bool) {
	var projectError *repairProjectValidationError
	var sourceError *repairSourceValidationError
	var stage, reason string
	switch {
	case errors.As(err, &projectError):
		stage, reason = projectError.Stage, projectError.Detail
	case errors.As(err, &sourceError):
		stage, reason = "source validation", sourceError.Error()
		if strings.Contains(strings.ToLower(filepath.Base(sourceError.Path)), "dockerfile") {
			stage = "Dockerfile syntax"
		}
	default:
		return "", false
	}
	for _, line := range strings.Split(reason, "\n") {
		if strings.Contains(line, "requires go >=") || strings.Contains(line, "No such file or directory") ||
			strings.Contains(line, "unknown instruction") || strings.Contains(line, "failed to solve") {
			reason = line
			if index := strings.Index(reason, "vendor/"); index >= 0 {
				reason = reason[index:]
			}
			break
		}
	}
	return strings.Join([]string{
		boundedPatchReviewLine("Patch fixer rejected candidate: " + stage + " failed."),
		boundedPatchReviewLine("Reason: " + patchReviewPlainText(reason)),
		"Original patches, checksums and attribution restored; upgrade remains incomplete.",
		"Human review required.",
	}, "\n"), true
}

func patchValidationSummary(result repairValidationResult) string {
	checks := append([]string{"patch replay"}, result.Checks...)
	for index, check := range checks {
		switch check {
		case "dockerfile-syntax":
			checks[index] = "Dockerfile syntax"
		case "attribution-checksums":
			checks[index] = "attribution/checksum generation"
		case "validate-checksums":
			checks[index] = "checksum verification"
		case "local-images":
			if result.NativeImages {
				checks[index] = fmt.Sprintf("native %s image build", result.Architecture)
			}
		}
	}
	if result.ImageSkipReason != "" {
		checks = append(checks, "image build skipped: "+result.ImageSkipReason)
	}
	return strings.Join(checks, ", ")
}

func newRepairValidationDirectory(root, project string) (string, error) {
	resultsRoot := strings.TrimSpace(os.Getenv(agentpatchfixer.ResultDirEnv))
	if resultsRoot == "" {
		resultsRoot = filepath.Join(root, "patch-fixer-runs")
	}
	parent := filepath.Join(resultsRoot, strings.ReplaceAll(project, "/", "_"))
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return "", fmt.Errorf("creating patch validation directory: %w", err)
	}
	directory, err := os.MkdirTemp(parent, "validation-")
	if err != nil {
		return "", fmt.Errorf("creating patch validation run: %w", err)
	}
	return directory, nil
}

func writeRepairValidationEvidence(directory string, result repairValidationResult, validationErr error) error {
	status, detail := "validated", ""
	if validationErr != nil {
		status, detail = "rejected", validationErr.Error()
	}
	payload := struct {
		Status     string                 `json:"status"`
		Validation repairValidationResult `json:"validation"`
		Error      string                 `json:"error,omitempty"`
	}{Status: status, Validation: result, Error: detail}
	content, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding patch validation result: %w", err)
	}
	if err := os.WriteFile(filepath.Join(directory, "result.json"), append(content, '\n'), 0o600); err != nil {
		return fmt.Errorf("saving patch validation result: %w", err)
	}
	return nil
}

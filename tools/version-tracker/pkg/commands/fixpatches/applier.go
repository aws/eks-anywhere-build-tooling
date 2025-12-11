package fixpatches

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// ApplyPatchFix applies the LLM-generated patch to files.
// This function applies the patch completely or fails - use ApplyPatchFixWithReject for partial application.
func ApplyPatchFix(fix *types.PatchFix, projectPath string) error {
	logger.Info("Applying LLM-generated patch", "path", projectPath)

	// Get the repo directory
	repoPath, err := GetRepoPath(projectPath)
	if err != nil {
		return err
	}

	// Save patch to temporary file
	tmpPatchFile := filepath.Join(projectPath, ".llm-patch.tmp")
	if err := os.WriteFile(tmpPatchFile, []byte(fix.Patch), 0644); err != nil {
		return fmt.Errorf("writing temporary patch file: %v", err)
	}
	defer os.Remove(tmpPatchFile) // Clean up temp file

	logger.Info("Saved patch to temporary file", "file", tmpPatchFile)

	// Also save to a debug file that persists (for debugging)
	debugPatchFile := filepath.Join(projectPath, ".llm-patch-debug.txt")
	if err := os.WriteFile(debugPatchFile, []byte(fix.Patch), 0644); err != nil {
		logger.Info("Warning: failed to write debug patch file", "error", err)
	} else {
		logger.Info("Saved debug patch file", "file", debugPatchFile)
	}

	// Apply patch using git apply
	// Note: We use git apply instead of git am because we're applying to an already-cloned repo
	output, err := GitApplyWithC(repoPath, tmpPatchFile, false, true)
	if err != nil {
		logger.Info("git apply failed", "error", err, "output", output)
		return err
	}

	logger.Info("Patch applied successfully")

	// Stage the changes
	if _, err := GitCommandWithC(repoPath, "add", "-A"); err != nil {
		return err
	}

	logger.Info("Changes staged successfully")

	return nil
}

// ApplyPatchFixWithReject applies the LLM-generated patch with --reject to allow partial success.
// Returns .rej files and patch application result for extracting context from failures.
func ApplyPatchFixWithReject(patchContent string, projectPath string) ([]string, *types.PatchApplicationResult, error) {
	logger.Info("Applying LLM-generated patch with --reject", "path", projectPath)

	// Get the repo directory
	repoPath, err := GetRepoPath(projectPath)
	if err != nil {
		return nil, nil, err
	}

	// Save patch to temporary file
	tmpPatchFile := filepath.Join(projectPath, ".llm-patch-attempt.tmp")
	if err := os.WriteFile(tmpPatchFile, []byte(patchContent), 0644); err != nil {
		return nil, nil, fmt.Errorf("writing temporary patch file: %v", err)
	}
	defer os.Remove(tmpPatchFile)

	logger.Info("Saved LLM patch to temporary file", "file", tmpPatchFile)

	// Apply patch with --reject
	outputStr, err := GitApplyWithC(repoPath, tmpPatchFile, true, true)

	// Parse output for offset information
	result := &types.PatchApplicationResult{
		OffsetFiles: make(map[string]int),
		GitOutput:   outputStr,
	}

	// Parse output line by line to detect offsets (same logic as applySinglePatchWithReject)
	var currentFile string
	for _, line := range strings.Split(outputStr, "\n") {
		if strings.HasPrefix(line, "Checking patch ") {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				currentFile = strings.TrimSuffix(parts[2], "...")
			}
		}
		if currentFile != "" && strings.Contains(line, "succeeded at") && strings.Contains(line, "offset") {
			// Extract offset amount
			if strings.Contains(line, "offset ") {
				offsetStr := strings.Split(strings.Split(line, "offset ")[1], " ")[0]
				var offset int
				fmt.Sscanf(offsetStr, "%d", &offset)
				result.OffsetFiles[currentFile] = offset
				logger.Info("Detected offset hunk in LLM patch", "file", currentFile, "offset", offset)
			}
		}
	}

	// Find .rej files
	rejFiles, findErr := findRejectionFiles(repoPath)
	if findErr != nil {
		logger.Info("Warning: failed to find rejection files", "error", findErr)
	}

	if err != nil {
		// Check if it's a patch conflict (expected) vs other error
		if strings.Contains(outputStr, "patch does not apply") ||
			strings.Contains(outputStr, "Rejected hunk") ||
			strings.Contains(outputStr, "does not exist in index") {
			logger.Info("LLM patch application had conflicts", "rej_files", len(rejFiles))
			return rejFiles, result, fmt.Errorf("patch conflicts: %s", outputStr)
		}
		return rejFiles, result, fmt.Errorf("git apply --reject failed: %v\nOutput: %s", err, outputStr)
	}

	logger.Info("LLM patch applied successfully without conflicts")
	return rejFiles, result, nil
}

// RevertPatchFix reverts a failed patch application.
func RevertPatchFix(projectPath string) error {
	logger.Info("Reverting patch changes", "path", projectPath)

	// Get the repo directory
	repoPath, err := GetRepoPath(projectPath)
	if err != nil {
		return err
	}

	// Reset to clean state
	return ResetToCleanStateWithC(repoPath)
}

// CommitPatchFix commits the successfully applied patch.
func CommitPatchFix(projectPath string, commitMessage string) error {
	logger.Info("Committing patch fix", "path", projectPath, "message", commitMessage)

	// Get the repo directory
	repoPath, err := GetRepoPath(projectPath)
	if err != nil {
		return err
	}

	// Commit the changes (GitCommit handles "nothing to commit" case)
	if err := GitCommit(repoPath, commitMessage); err != nil {
		return err
	}

	logger.Info("Patch fix committed successfully")
	return nil
}

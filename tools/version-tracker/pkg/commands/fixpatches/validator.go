package fixpatches

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// ValidateBuild runs make build and handles checksum updates.
// During version bumps, checksums are expected to change, so we:
// 1. Run make build with SKIP_CHECKSUM_VALIDATION=true to skip validation
// 2. Run make checksums to update the CHECKSUMS file
// 3. Stage the updated CHECKSUMS file for commit
func ValidateBuild(projectPath string, releaseBranch string) error {
	// Check if SKIP_VALIDATION env var is set (for testing)
	if os.Getenv("SKIP_VALIDATION") == "true" {
		logger.Info("Skipping build validation (SKIP_VALIDATION=true)")
		return nil
	}

	logger.Info("Running build validation", "path", projectPath, "release_branch", releaseBranch)

	// Run make build with SKIP_CHECKSUM_VALIDATION=true and IMAGE_NAMES=""
	// - SKIP_CHECKSUM_VALIDATION: skips validate-checksums which would fail during version bumps
	// - IMAGE_NAMES="": skips building Docker images (only builds binaries for validation)
	// - BUILD_TARGETS_OVERRIDE: overrides BUILD_TARGETS to exclude attribution-pr (we don't need separate attribution PRs)
	//   Common.mk uses: build: $(or $(BUILD_TARGETS_OVERRIDE),$(BUILD_TARGETS))
	buildCmd := exec.Command("make", "-C", projectPath, "build")
	buildEnv := os.Environ()
	buildEnv = append(buildEnv, "SKIP_CHECKSUM_VALIDATION=true")
	buildEnv = append(buildEnv, "IMAGE_NAMES=")
	buildEnv = append(buildEnv, "BUILD_TARGETS_OVERRIDE=github-rate-limit-pre validate-checksums attribution github-rate-limit-post")
	if releaseBranch != "" {
		buildEnv = append(buildEnv, fmt.Sprintf("RELEASE_BRANCH=%s", releaseBranch))
	}
	buildCmd.Env = buildEnv

	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %v\nOutput: %s", err, string(buildOutput))
	}

	logger.Info("Build succeeded")

	// Now run make checksums to update the CHECKSUMS file with the new binary checksums
	logger.Info("Updating checksums for new binaries...")
	checksumCmd := exec.Command("make", "-C", projectPath, "checksums")
	checksumEnv := os.Environ()
	if releaseBranch != "" {
		checksumEnv = append(checksumEnv, fmt.Sprintf("RELEASE_BRANCH=%s", releaseBranch))
	}
	checksumCmd.Env = checksumEnv

	checksumOutput, err := checksumCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update checksums: %v\nOutput: %s", err, string(checksumOutput))
	}

	logger.Info("Checksums updated successfully")

	// Stage the updated CHECKSUMS or expected-artifacts file
	// For release-branched projects, the CHECKSUMS file is in projectPath/releaseBranch/
	if err := stageChecksumFiles(projectPath, releaseBranch); err != nil {
		logger.Info("Warning: failed to stage checksum files", "error", err)
	}

	logger.Info("Build validation completed successfully")
	return nil
}

// extractChecksumsFromOutput parses the validate_checksums.sh output to extract new checksums
// The output format is:
// *************** CHECKSUMS ***************
// <checksum>  <filepath>
// <checksum>  <filepath>
// *****************************************
func extractChecksumsFromOutput(output string) ([]string, error) {
	lines := strings.Split(output, "\n")
	var checksums []string
	inChecksumSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "*************** CHECKSUMS ***************") {
			inChecksumSection = true
			continue
		}

		if strings.Contains(line, "*****************************************") {
			inChecksumSection = false
			break
		}

		if inChecksumSection && line != "" {
			// Each line should be: "<checksum>  <filepath>"
			// We want to preserve the exact format
			checksums = append(checksums, line)
		}
	}

	if len(checksums) == 0 {
		return nil, fmt.Errorf("no checksums found in output")
	}

	return checksums, nil
}

// writeChecksumsFile writes the checksums to the CHECKSUMS file
func writeChecksumsFile(projectPath string, checksums []string) error {
	checksumsPath := filepath.Join(projectPath, "CHECKSUMS")

	// Write checksums to file
	content := strings.Join(checksums, "\n") + "\n"
	if err := os.WriteFile(checksumsPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing CHECKSUMS file: %v", err)
	}

	return nil
}

// runMakeChecksums runs make checksums to update the CHECKSUMS file (fallback method)
func runMakeChecksums(projectPath string, releaseBranch string) error {
	logger.Info("Running make checksums to update CHECKSUMS file")

	checksumCmd := exec.Command("make", "-C", projectPath, "checksums")
	if releaseBranch != "" {
		checksumCmd.Env = append(os.Environ(), fmt.Sprintf("RELEASE_BRANCH=%s", releaseBranch))
	}
	checksumOutput, err := checksumCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("make checksums failed: %v\nOutput: %s", err, string(checksumOutput))
	}

	logger.Info("Checksums updated successfully via make checksums")
	return nil
}

// stageChecksumFiles stages CHECKSUMS or expected-artifacts files for commit
// For release-branched projects, the files are in projectPath/releaseBranch/
func stageChecksumFiles(projectPath string, releaseBranch string) error {
	// Determine the directory where CHECKSUMS file is located
	// For release-branched projects: projectPath/releaseBranch/CHECKSUMS
	// For non-release-branched projects: projectPath/CHECKSUMS
	checksumsDir := projectPath
	if releaseBranch != "" {
		checksumsDir = filepath.Join(projectPath, releaseBranch)
	}

	// Stage CHECKSUMS if it exists
	checksumsPath := filepath.Join(checksumsDir, "CHECKSUMS")
	if _, err := os.Stat(checksumsPath); err == nil {
		// Use relative path for git add (relative to projectPath)
		relPath := "CHECKSUMS"
		if releaseBranch != "" {
			relPath = filepath.Join(releaseBranch, "CHECKSUMS")
		}

		stageCmd := exec.Command("git", "add", relPath)
		stageCmd.Dir = projectPath
		if output, err := stageCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("staging CHECKSUMS: %v\nOutput: %s", err, string(output))
		}
		logger.Info("Staged updated CHECKSUMS file", "path", relPath)
	}

	// Stage expected-artifacts if it exists
	expectedArtifactsPath := filepath.Join(checksumsDir, "expected-artifacts")
	if _, err := os.Stat(expectedArtifactsPath); err == nil {
		// Use relative path for git add
		relPath := "expected-artifacts"
		if releaseBranch != "" {
			relPath = filepath.Join(releaseBranch, "expected-artifacts")
		}

		stageCmd := exec.Command("git", "add", relPath)
		stageCmd.Dir = projectPath
		if output, err := stageCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("staging expected-artifacts: %v\nOutput: %s", err, string(output))
		}
		logger.Info("Staged updated expected-artifacts file", "path", relPath)
	}

	return nil
}

// ValidateSemantics checks if fix preserves original intent.
func ValidateSemantics(fix *types.PatchFix, ctx *types.PatchContext) error {
	logger.Info("Running semantic validation")

	// Validate patch metadata is preserved
	if ctx.PatchAuthor != "" && !strings.Contains(fix.Patch, ctx.PatchAuthor) {
		logger.Info("Warning: patch author not preserved", "expected", ctx.PatchAuthor)
		// Don't fail - this is a warning
	}

	if ctx.PatchDate != "" && !strings.Contains(fix.Patch, ctx.PatchDate) {
		logger.Info("Warning: patch date not preserved", "expected", ctx.PatchDate)
	}

	if ctx.PatchSubject != "" {
		subjectCore := strings.TrimPrefix(ctx.PatchSubject, "[PATCH]")
		subjectCore = strings.TrimSpace(subjectCore)
		if !strings.Contains(fix.Patch, subjectCore) {
			logger.Info("Warning: patch subject not preserved", "expected", subjectCore)
		}
	}

	// Count lines changed in original patch
	originalLines := countChangedLines(ctx.OriginalPatch)
	fixLines := countChangedLines(fix.Patch)

	// Check for excessive drift (>50% more changes)
	if fixLines > originalLines*3/2 {
		return fmt.Errorf("semantic drift: fix changes %d lines vs %d in original (>50%% increase)",
			fixLines, originalLines)
	}

	logger.Info("Semantic validation passed", "original_lines", originalLines, "fix_lines", fixLines)

	return nil
}

// countChangedLines counts the number of changed lines in a patch (+ and - lines).
func countChangedLines(patch string) int {
	lines := strings.Split(patch, "\n")
	count := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			count++
		}
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			count++
		}
	}
	return count
}

package fixpatches

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// Run executes the patch fixing workflow, processing each patch file sequentially.
func Run(opts *types.FixPatchesOptions) error {
	logger.Info("Starting patch fixing workflow", "project", opts.ProjectName, "pr", opts.PRNumber)

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %v", err)
	}

	// Extract org and repo from project name
	projectRepo := strings.Split(opts.ProjectName, "/")[1]

	// Construct project path: projects/<org>/<repo>
	projectPath := filepath.Join(cwd, "projects", opts.ProjectName)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	logger.Info("Project directory located", "path", projectPath)

	// SPECIAL CASE: kubernetes/autoscaler requires hardcoded fix
	// See: projects/kubernetes/autoscaler/README.md lines 25-60
	if strings.Contains(opts.ProjectName, "kubernetes/autoscaler") {
		logger.Info("Detected kubernetes/autoscaler - using hardcoded fix")

		// Get release branch from the commit message
		// The commit message format is: "Bump kubernetes/autoscaler 1-XX release branch to latest release"
		getReleaseBranchCmd := exec.Command("git", "log", "-1", "--format=%s")
		getReleaseBranchCmd.Dir = projectPath
		commitMsgBytes, err := getReleaseBranchCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("getting commit message: %v", err)
		}
		commitMsg := string(commitMsgBytes)

		// Extract release branch from commit message (e.g., "1-29" from "Bump kubernetes/autoscaler 1-29 release branch...")
		var releaseBranch string
		if strings.Contains(commitMsg, "autoscaler") {
			// Look for pattern like "1-29" or "1-30"
			parts := strings.Fields(commitMsg)
			for _, part := range parts {
				if strings.HasPrefix(part, "1-") && len(part) <= 5 {
					releaseBranch = part
					break
				}
			}
		}

		if releaseBranch == "" {
			// Fallback: use the latest supported branch
			supportedBranchesFile := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(projectPath))), "release", "SUPPORTED_RELEASE_BRANCHES")
			branchesContent, err := os.ReadFile(supportedBranchesFile)
			if err != nil {
				return fmt.Errorf("reading SUPPORTED_RELEASE_BRANCHES: %v", err)
			}
			branches := strings.Split(strings.TrimSpace(string(branchesContent)), "\n")
			if len(branches) == 0 {
				return fmt.Errorf("no release branches found")
			}
			releaseBranch = strings.TrimSpace(branches[len(branches)-1])
		}

		logger.Info("Detected release branch", "branch", releaseBranch)

		// For autoscaler, we need to checkout the source WITHOUT applying patches
		// (since the patches are what we're trying to fix)
		logger.Info("Checking out autoscaler source", "release_branch", releaseBranch)

		// Read GIT_TAG to know which version to checkout
		gitTagFile := filepath.Join(projectPath, releaseBranch, "GIT_TAG")
		gitTagBytes, err := os.ReadFile(gitTagFile)
		if err != nil {
			return fmt.Errorf("reading GIT_TAG: %v", err)
		}
		gitTag := strings.TrimSpace(string(gitTagBytes))
		logger.Info("Found git tag", "tag", gitTag)

		// Clone the repo if it doesn't exist
		repoPath := filepath.Join(projectPath, "autoscaler")
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			logger.Info("Cloning autoscaler repository")
			cloneCmd := exec.Command("git", "clone", "https://github.com/kubernetes/autoscaler.git", repoPath)
			cloneCmd.Dir = projectPath
			if output, err := cloneCmd.CombinedOutput(); err != nil {
				return fmt.Errorf("cloning repo: %v\nOutput: %s", err, output)
			}
		}

		// Checkout the specific tag
		logger.Info("Checking out tag", "tag", gitTag)
		checkoutCmd := exec.Command("git", "checkout", "-f", gitTag)
		checkoutCmd.Dir = repoPath
		if output, err := checkoutCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("checking out tag: %v\nOutput: %s", err, output)
		}

		// Clean any existing changes
		cleanCmd := exec.Command("git", "clean", "-fd")
		cleanCmd.Dir = repoPath
		if output, err := cleanCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("cleaning repo: %v\nOutput: %s", err, output)
		}

		logger.Info("Source checked out successfully")

		// Abort any existing git am session from failed patch application
		abortCmd := exec.Command("git", "am", "--abort")
		abortCmd.Dir = repoPath
		abortCmd.Run() // Ignore errors - there might not be an am session

		return fixAutoscalerPatches(projectPath, releaseBranch)
	}

	// Check if project has binaries that are release-branched
	// For such projects, patches are in <project>/<release-branch>/patches/
	// instead of <project>/patches/
	binariesReleaseBranchedCmd := exec.Command("make", "-C", projectPath, "var-value-BINARIES_ARE_RELEASE_BRANCHED")
	binariesReleaseBranchedCmd.Env = append(os.Environ(), "RELEASE_BRANCH=dummy")
	binariesReleaseBranchedOutput, _ := binariesReleaseBranchedCmd.CombinedOutput()
	// Get the last line of output (Makefile may output errors to stderr which get captured)
	outputLines := strings.Split(strings.TrimSpace(string(binariesReleaseBranchedOutput)), "\n")
	lastLine := strings.TrimSpace(outputLines[len(outputLines)-1])
	binariesReleaseBranched := lastLine == "true"
	logger.Info("Checked BINARIES_ARE_RELEASE_BRANCHED", "value", binariesReleaseBranched, "last_line", lastLine)

	// Determine patches directory based on where CHECKSUMS file is located
	// If CHECKSUMS is in project root, patches are in <project>/patches/
	// If CHECKSUMS is in a release branch dir, patches are in <project>/<release-branch>/patches/
	var patchesDir string
	var releaseBranchForMake string // Release branch to pass to make commands

	// Get all supported release branches
	supportedBranchesFile := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(projectPath))), "release", "SUPPORTED_RELEASE_BRANCHES")
	var latestBranch string
	if branchesContent, err := os.ReadFile(supportedBranchesFile); err == nil {
		branches := strings.Split(strings.TrimSpace(string(branchesContent)), "\n")
		if len(branches) > 0 {
			latestBranch = strings.TrimSpace(branches[len(branches)-1])
		}
	}

	// Check if project has CHECKSUMS in root (indicates root structure)
	// Some projects (like image-builder) use expected-artifacts instead
	checksumsInRoot := filepath.Join(projectPath, "CHECKSUMS")
	hasRootStructure := false
	if _, err := os.Stat(checksumsInRoot); err == nil {
		hasRootStructure = true
		logger.Info("Found CHECKSUMS in root - using root structure")
	} else {
		// Fallback: check for expected-artifacts (used by image-builder and other oddballs)
		expectedArtifactsInRoot := filepath.Join(projectPath, "expected-artifacts")
		if _, err := os.Stat(expectedArtifactsInRoot); err == nil {
			hasRootStructure = true
			logger.Info("Found expected-artifacts in root - using root structure")
		}
	}

	// Try to extract release branch from commit message (for release-branched projects)
	// Commit message format: "Bump <project> 1-XX release branch to latest release"
	var branchFromCommit string
	getReleaseBranchCmd := exec.Command("git", "log", "-1", "--format=%s")
	getReleaseBranchCmd.Dir = projectPath
	if commitMsgBytes, err := getReleaseBranchCmd.CombinedOutput(); err == nil {
		commitMsg := string(commitMsgBytes)
		// Look for pattern like "1-29" or "1-30" in commit message
		parts := strings.Fields(commitMsg)
		for _, part := range parts {
			if strings.HasPrefix(part, "1-") && len(part) <= 5 {
				branchFromCommit = part
				logger.Info("Extracted release branch from commit message", "branch", branchFromCommit, "commit", commitMsg)
				break
			}
		}
	}

	// Determine patches directory based on structure indicators and commit message
	if hasRootStructure {
		// Project has root structure (CHECKSUMS or expected-artifacts in root)
		patchesDir = filepath.Join(projectPath, constants.PatchesDirectory)
		releaseBranchForMake = latestBranch // Use latest release branch for make commands
		logger.Info("Project structure: patches in root directory", "patches_dir", patchesDir, "release_branch_for_make", releaseBranchForMake)
	} else if branchFromCommit != "" {
		// Check if patches exist in the branch from commit message
		branchPatchesDir := filepath.Join(projectPath, branchFromCommit, constants.PatchesDirectory)
		if _, err := os.Stat(branchPatchesDir); err == nil {
			// Patches are in the specific release branch from commit
			patchesDir = branchPatchesDir
			releaseBranchForMake = branchFromCommit
			logger.Info("Project structure: patches in release branch directory (from commit)", "release_branch", branchFromCommit, "patches_dir", patchesDir)
		} else {
			// Commit mentioned a branch but patches aren't there - fall through to root check
			logger.Info("Commit mentioned branch but patches not found there", "branch", branchFromCommit)
		}
	}

	// If we haven't found patches yet, check root directory
	if patchesDir == "" {
		rootPatchesDir := filepath.Join(projectPath, constants.PatchesDirectory)
		if _, err := os.Stat(rootPatchesDir); err == nil {
			// Patches are in root directory (e.g., image-builder, most projects)
			patchesDir = rootPatchesDir
			releaseBranchForMake = latestBranch // Use latest release branch for make commands
			logger.Info("Project structure: patches in root directory", "patches_dir", patchesDir, "release_branch_for_make", releaseBranchForMake)
		} else {
			// Last resort: check latest release branch directory
			latestBranchPatchesDir := filepath.Join(projectPath, latestBranch, constants.PatchesDirectory)
			if _, err := os.Stat(latestBranchPatchesDir); err == nil {
				patchesDir = latestBranchPatchesDir
				releaseBranchForMake = latestBranch
				logger.Info("Project structure: patches in latest release branch directory", "release_branch", latestBranch, "patches_dir", patchesDir)
			}
		}
	}

	// Get sorted list of patch files
	patchFiles, err := filepath.Glob(filepath.Join(patchesDir, "*.patch"))
	if err != nil {
		return fmt.Errorf("finding patch files: %v", err)
	}

	if len(patchFiles) == 0 {
		logger.Info("No patch files found - nothing to fix")
		return nil
	}

	// Sort patch files to ensure sequential processing (0001, 0002, 0003...)
	sort.Strings(patchFiles)

	logger.Info("Found patch files", "count", len(patchFiles), "files", patchFiles)

	// Track fixed patches for reporting
	var fixedPatches []string
	var failedPatches []string
	var lastError error

	// Process each patch file sequentially
	for patchIndex, patchFile := range patchFiles {
		logger.Info("Processing patch", "index", patchIndex+1, "total", len(patchFiles), "file", filepath.Base(patchFile))

		// Try to fix this specific patch
		if err := fixSinglePatch(patchFile, projectPath, projectRepo, releaseBranchForMake, opts); err != nil {
			logger.Info("Failed to fix patch", "patch", filepath.Base(patchFile), "error", err)
			failedPatches = append(failedPatches, filepath.Base(patchFile))
			lastError = err
			// Continue processing other patches instead of failing immediately
			continue
		}

		logger.Info("Patch processed successfully", "file", filepath.Base(patchFile))
		fixedPatches = append(fixedPatches, filepath.Base(patchFile))
	}

	// After all patches are processed, run build validation once
	if len(fixedPatches) > 0 {
		logger.Info("Running final build validation for all fixed patches", "fixed_count", len(fixedPatches))
		if err := ValidateBuild(projectPath, releaseBranchForMake); err != nil {
			logger.Info("Final build validation failed", "error", err)
			return fmt.Errorf("build validation failed after fixing patches: %v", err)
		}
		logger.Info("Final build validation passed")
	}

	// Commit and push if requested
	if opts.Push && len(fixedPatches) > 0 {
		if opts.Branch == "" {
			return fmt.Errorf("--branch is required when using --push")
		}

		commitMsg := fmt.Sprintf("Auto-fix patches for %s\n\nFixed patches:\n", opts.ProjectName)
		for _, patch := range fixedPatches {
			commitMsg += fmt.Sprintf("- %s\n", patch)
		}

		if err := CommitAndPush(projectPath, opts.Branch, commitMsg); err != nil {
			logger.Info("Failed to commit and push", "error", err)
			// Don't fail the whole operation if push fails
		}
	}

	// Comment on PR if requested
	if opts.CommentOnPR && opts.PRNumber > 0 {
		ghClient, err := NewGitHubClient()
		if err != nil {
			logger.Info("Failed to create GitHub client", "error", err)
		} else {
			var comment string
			if len(failedPatches) > 0 {
				comment = FormatFailureComment(opts.ProjectName, failedPatches)
			} else {
				comment = FormatSuccessComment(opts.ProjectName, fixedPatches)
			}

			if err := ghClient.CommentOnPR(opts.PRNumber, comment); err != nil {
				logger.Info("Failed to comment on PR", "error", err)
			}
		}
	}

	// Return error if any patches failed
	if len(failedPatches) > 0 {
		return fmt.Errorf("failed to fix %d patch(es): %v", len(failedPatches), lastError)
	}

	logger.Info("All patches processed successfully")
	return nil
}

// fixSinglePatch processes a single patch file through the fix-validate cycle.
func fixSinglePatch(patchFile string, projectPath string, projectRepo string, releaseBranch string, opts *types.FixPatchesOptions) error {
	logger.Info("Fixing single patch", "patch", filepath.Base(patchFile))

	// Apply this specific patch with git apply --reject
	rejFiles, patchResult, err := applySinglePatchWithReject(patchFile, projectPath, projectRepo, releaseBranch)
	if err != nil {
		return fmt.Errorf("applying patch with reject: %v", err)
	}

	// If no .rej files, patch applied successfully
	if len(rejFiles) == 0 {
		logger.Info("Patch applied successfully without conflicts", "patch", filepath.Base(patchFile))
		return nil
	}

	logger.Info("Patch has conflicts", "patch", filepath.Base(patchFile), "rej_files", len(rejFiles), "offset_files", len(patchResult.OffsetFiles))

	// Calculate complexity for this patch
	// TODO(Phase 2): Consider PR-level complexity gating instead of per-patch
	// If any single patch exceeds threshold, skip entire PR for better UX
	// Rationale: Avoid mixed state where some patches fixed, others need manual work
	complexity, err := calculateComplexity(rejFiles)
	if err != nil {
		return fmt.Errorf("calculating complexity: %v", err)
	}

	logger.Info("Calculated patch complexity", "score", complexity, "threshold", opts.ComplexityThreshold)

	// Check if complexity exceeds threshold
	// TODO(Phase 2): Refine complexity calculation based on PoC metrics
	// Current: complexity = hunks + files
	// Consider: weighted scoring based on hunk type, file type, lines changed
	// Track success rates by complexity level to optimize threshold
	if complexity > opts.ComplexityThreshold {
		logger.Info("Complexity exceeds threshold - skipping this patch",
			"complexity", complexity,
			"threshold", opts.ComplexityThreshold)
		return &types.PatchFixError{
			Code:    types.ErrorComplexityTooHigh,
			Message: fmt.Sprintf("Patch %s complexity (%d) exceeds threshold (%d)", filepath.Base(patchFile), complexity, opts.ComplexityThreshold),
			Details: map[string]interface{}{
				"patch":      filepath.Base(patchFile),
				"complexity": complexity,
				"threshold":  opts.ComplexityThreshold,
				"rej_files":  rejFiles,
			},
		}
	}

	// Extract context ONCE from the original patch application
	// This context will be reused for all attempts to avoid state pollution
	baseContext, err := ExtractPatchContext(rejFiles, patchFile, projectPath, 1, patchResult)
	if err != nil {
		return fmt.Errorf("extracting patch context: %v", err)
	}

	logger.Info("Extracted base patch context", "token_count", baseContext.TokenCount, "hunks", len(baseContext.FailedHunks))

	// Build a test prompt to measure actual size with template overhead
	testPrompt := BuildPrompt(baseContext, 1)
	promptChars := len(testPrompt)
	estimatedPromptTokens := promptChars / 3 // Conservative: 3 chars per token for code

	logger.Info("Built test prompt to measure size",
		"prompt_chars", promptChars,
		"estimated_tokens", estimatedPromptTokens,
		"context_tokens", baseContext.TokenCount)

	// Claude Sonnet 4.5 has 200K input token limit
	// We target 60K tokens (180K chars) to stay safely under the limit
	maxPromptTokens := 60000
	maxPromptChars := maxPromptTokens * 3

	if promptChars > maxPromptChars {
		logger.Info("Prompt too large, need to prune context",
			"current_chars", promptChars,
			"max_chars", maxPromptChars,
			"reduction_needed", promptChars-maxPromptChars)

		// Calculate overhead multiplier: prompt_size / context_size
		// This tells us how much the template expands the context
		overheadMultiplier := float64(promptChars) / float64(baseContext.TokenCount*3)
		logger.Info("Calculated template overhead", "multiplier", overheadMultiplier)

		// Work backwards: if we want maxPromptChars and template adds overhead,
		// then context should be: maxPromptChars / overhead
		targetContextChars := int(float64(maxPromptChars) / overheadMultiplier)
		targetContextTokens := targetContextChars / 3

		logger.Info("Calculated target context size",
			"target_chars", targetContextChars,
			"target_tokens", targetContextTokens)

		// Prune context to target size
		if err := PruneContext(baseContext, targetContextTokens); err != nil {
			return fmt.Errorf("pruning context: %v", err)
		}

		// Rebuild prompt and verify size
		testPrompt = BuildPrompt(baseContext, 1)
		promptChars = len(testPrompt)
		estimatedPromptTokens = promptChars / 3

		logger.Info("After pruning",
			"prompt_chars", promptChars,
			"estimated_tokens", estimatedPromptTokens,
			"context_tokens", baseContext.TokenCount)

		if promptChars > maxPromptChars {
			logger.Info("Warning: prompt still too large after pruning, may hit API limit")
		}
	} else {
		logger.Info("Prompt size acceptable, no pruning needed")
	}

	// Iterative refinement loop for this patch
	// Start with base context, then extract NEW context from each LLM attempt's failures
	currentContext := baseContext
	var lastBuildError string

	for attempt := 1; attempt <= opts.MaxAttempts; attempt++ {
		logger.Info("Starting fix attempt for patch", "patch", filepath.Base(patchFile), "attempt", attempt, "max_attempts", opts.MaxAttempts)

		// Use context from previous attempt (or base context for attempt 1)
		ctx := *currentContext // Create a copy
		ctx.BuildError = lastBuildError

		logger.Info("Using context for attempt", "token_count", ctx.TokenCount, "hunks", len(ctx.FailedHunks))

		// Call LLM to generate fix
		fix, err := CallBedrockForPatchFix(&ctx, opts.Model, attempt)
		if err != nil {
			logger.Info("Bedrock API call failed", "error", err, "attempt", attempt)
			if attempt == opts.MaxAttempts {
				return &types.PatchFixError{
					Code:    types.ErrorBedrockAPI,
					Message: fmt.Sprintf("Bedrock API failed for patch %s after %d attempts: %v", filepath.Base(patchFile), opts.MaxAttempts, err),
					Details: map[string]interface{}{
						"patch":    filepath.Base(patchFile),
						"attempts": opts.MaxAttempts,
						"error":    err.Error(),
					},
				}
			}
			continue
		}

		logger.Info("LLM generated patch fix", "tokens_used", fix.TokensUsed, "patch_length", len(fix.Patch))

		// Log first 500 chars of generated patch for debugging
		patchPreview := fix.Patch
		if len(patchPreview) > 500 {
			patchPreview = patchPreview[:500] + "...(truncated)"
		}
		logger.Info("Generated patch preview", "preview", patchPreview)

		// CRITICAL: Revert to clean state BEFORE applying LLM's patch
		// This ensures we're not applying on top of the original patch's modifications
		logger.Info("Reverting to clean state before applying LLM patch")
		if revertErr := RevertPatchFix(projectPath); revertErr != nil {
			logger.Info("Warning: failed to revert to clean state", "error", revertErr)
		}

		// Apply the LLM's patch with --reject to see what fails
		// This allows partial success and lets us extract context from actual failures
		rejFiles, patchResult, applyErr := ApplyPatchFixWithReject(fix.Patch, projectPath)

		if len(rejFiles) == 0 && applyErr == nil {
			// Success! Patch applied completely
			logger.Info("LLM patch applied successfully without conflicts")

			// Now validate build and semantics
			// (validation code continues below)
		} else {
			// Patch had conflicts - extract NEW context from THIS attempt's failures
			logger.Info("LLM patch had conflicts", "rej_files", len(rejFiles), "error", applyErr)

			// Store error for next attempt
			if applyErr != nil {
				lastBuildError = applyErr.Error()
			}

			// Extract NEW context from the LLM's patch failures
			// This shows what ACTUALLY failed in this attempt, not the original patch
			if len(rejFiles) > 0 {
				logger.Info("Extracting NEW context from LLM patch failures")
				newContext, extractErr := ExtractPatchContext(rejFiles, patchFile, projectPath, attempt+1, patchResult)
				if extractErr != nil {
					logger.Info("Warning: failed to extract new context", "error", extractErr)
					// Fall back to reusing current context
				} else {
					logger.Info("Extracted new context from LLM patch failures", "hunks", len(newContext.FailedHunks))
					// Use this NEW context for the next attempt
					currentContext = newContext
				}
			}

			// Revert changes to clean state
			if revertErr := RevertPatchFix(projectPath); revertErr != nil {
				logger.Info("Failed to revert patch", "error", revertErr)
			}

			continue
		}

		logger.Info("Patch fix applied successfully")

		// Validate semantics (per-patch validation)
		if err := ValidateSemantics(fix, &ctx); err != nil {
			logger.Info("Semantic validation failed", "error", err, "attempt", attempt)

			// Store ONLY the current error for next attempt (simplified approach)
			lastBuildError = err.Error()

			// Revert changes to clean state
			if revertErr := RevertPatchFix(projectPath); revertErr != nil {
				logger.Info("Failed to revert patch", "error", revertErr)
			}

			// DON'T re-apply original patch - we reuse the base context instead

			if attempt == opts.MaxAttempts {
				return &types.PatchFixError{
					Code:    types.ErrorSemanticDrift,
					Message: fmt.Sprintf("Semantic validation failed for patch %s after %d attempts", filepath.Base(patchFile), opts.MaxAttempts),
					Details: map[string]interface{}{
						"patch":    filepath.Base(patchFile),
						"attempts": opts.MaxAttempts,
						"error":    err.Error(),
					},
				}
			}
			continue
		}

		logger.Info("Semantic validation passed")

		// Success! This patch is fixed
		logger.Info("Patch fix successful", "patch", filepath.Base(patchFile), "attempt", attempt, "tokens_used", fix.TokensUsed)

		// Write the fixed patch back to the original patch file
		logger.Info("Writing fixed patch to file", "file", patchFile, "patch_length", len(fix.Patch))
		if err := WritePatchToFile(fix.Patch, patchFile); err != nil {
			return fmt.Errorf("writing fixed patch to file: %v", err)
		}

		logger.Info("Fixed patch written to file successfully", "file", patchFile)

		// Stage the fixed patch file and any updated checksums/artifacts
		// make checksums updates CHECKSUMS or expected-artifacts, which should be committed
		stageCmd := exec.Command("git", "add", filepath.Base(patchFile))
		stageCmd.Dir = filepath.Dir(patchFile) // Run in patches directory
		if output, err := stageCmd.CombinedOutput(); err != nil {
			logger.Info("Warning: failed to stage fixed patch", "error", err, "output", string(output))
		}

		// Stage CHECKSUMS or expected-artifacts if they exist
		checksumsPath := filepath.Join(projectPath, "CHECKSUMS")
		if _, err := os.Stat(checksumsPath); err == nil {
			stageChecksumsCmd := exec.Command("git", "add", "CHECKSUMS")
			stageChecksumsCmd.Dir = projectPath
			if output, err := stageChecksumsCmd.CombinedOutput(); err != nil {
				logger.Info("Warning: failed to stage CHECKSUMS", "error", err, "output", string(output))
			} else {
				logger.Info("Staged updated CHECKSUMS file")
			}
		}

		expectedArtifactsPath := filepath.Join(projectPath, "expected-artifacts")
		if _, err := os.Stat(expectedArtifactsPath); err == nil {
			stageArtifactsCmd := exec.Command("git", "add", "expected-artifacts")
			stageArtifactsCmd.Dir = projectPath
			if output, err := stageArtifactsCmd.CombinedOutput(); err != nil {
				logger.Info("Warning: failed to stage expected-artifacts", "error", err, "output", string(output))
			} else {
				logger.Info("Staged updated expected-artifacts file")
			}
		}

		// Clean up .rej files for this patch
		for _, rejFile := range rejFiles {
			os.Remove(rejFile)
		}

		return nil
	}

	// All attempts exhausted for this patch
	return &types.PatchFixError{
		Code:    types.ErrorMaxAttemptsExceeded,
		Message: fmt.Sprintf("Failed to fix patch %s after %d attempts", filepath.Base(patchFile), opts.MaxAttempts),
		Details: map[string]interface{}{
			"patch":    filepath.Base(patchFile),
			"attempts": opts.MaxAttempts,
		},
	}
}

// applyPatches attempts to apply patches using git apply --reject to generate .rej files.
// This function:
// 1. Ensures the upstream repo is checked out (via GIT_CHECKOUT_TARGET, NOT checkout-repo)
// 2. Applies patches using git apply --reject to generate .rej files for conflicts
//
// Note: We use GIT_CHECKOUT_TARGET instead of checkout-repo because checkout-repo
// will try to apply patches (via GIT_PATCH_TARGET), which will fail if patches don't apply cleanly.
// We want to apply patches ourselves with --reject to generate .rej files.
func applyPatches(projectPath string, repoName string) error {
	logger.Info("Checking out upstream repository", "path", projectPath)

	// Get the GIT_TAG from the project's Makefile
	gitTagCmd := exec.Command("make", "-C", projectPath, "var-value-GIT_TAG")
	gitTagOutput, err := gitTagCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("getting GIT_TAG: %v\nOutput: %s", err, gitTagOutput)
	}
	gitTag := strings.TrimSpace(string(gitTagOutput))

	// Build the GIT_CHECKOUT_TARGET: $(REPO)/eks-anywhere-checkout-$(GIT_TAG)
	checkoutTarget := fmt.Sprintf("%s/eks-anywhere-checkout-%s", repoName, gitTag)

	// Ensure the repo is checked out (but don't apply patches)
	// This creates the marker file $(REPO)/eks-anywhere-checkout-$(GIT_TAG)
	checkoutCmd := exec.Command("make", "-C", projectPath, checkoutTarget)
	checkoutOutput, err := checkoutCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("make %s failed: %v\nOutput: %s", checkoutTarget, err, checkoutOutput)
	}

	logger.Info("Repository checked out successfully", "tag", gitTag)

	// Find the patches directory
	patchesDir := filepath.Join(projectPath, constants.PatchesDirectory)
	if _, err := os.Stat(patchesDir); os.IsNotExist(err) {
		logger.Info("No patches directory found - nothing to fix")
		return nil
	}

	// Get list of patch files
	patchFiles, err := filepath.Glob(filepath.Join(patchesDir, "*.patch"))
	if err != nil {
		return fmt.Errorf("finding patch files: %v", err)
	}

	if len(patchFiles) == 0 {
		logger.Info("No patch files found in patches directory")
		return nil
	}

	logger.Info("Found patch files", "count", len(patchFiles))

	// The cloned repo directory is named after the repository
	repoPath := filepath.Join(projectPath, repoName)

	// Check if repo was cloned
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("cloned repository not found at %s", repoPath)
	}

	// Configure git in the cloned repo (same as Common.mk does for patch application)
	configEmailCmd := exec.Command("git", "-C", repoPath, "config", "user.email", constants.PatchApplyGitUserEmail)
	if err := configEmailCmd.Run(); err != nil {
		return fmt.Errorf("configuring git user.email: %v", err)
	}

	configNameCmd := exec.Command("git", "-C", repoPath, "config", "user.name", constants.PatchApplyGitUserName)
	if err := configNameCmd.Run(); err != nil {
		return fmt.Errorf("configuring git user.name: %v", err)
	}

	// Apply patches using git apply --reject
	// This will:
	// - Apply successful hunks
	// - Create .rej files for failed hunks
	// - Return error if any hunks fail
	logger.Info("Applying patches with git apply --reject", "repo", repoPath)

	for _, patchFile := range patchFiles {
		logger.Info("Applying patch", "file", filepath.Base(patchFile))

		cmd := exec.Command("git", "-C", repoPath, "apply", "--reject", "--whitespace=fix", patchFile)
		output, err := cmd.CombinedOutput()

		if err != nil {
			// Check if it's a patch conflict (expected) vs other error
			outputStr := string(output)
			if strings.Contains(outputStr, "patch does not apply") ||
				strings.Contains(outputStr, "Rejected hunk") ||
				strings.Contains(outputStr, "does not exist in index") {
				logger.Info("Patch application failed with conflicts (expected)",
					"patch", filepath.Base(patchFile),
					"output", outputStr)
				// Continue to next patch - we want to apply as many as possible
				continue
			}
			return fmt.Errorf("git apply failed for %s: %v\nOutput: %s", patchFile, err, output)
		}

		logger.Info("Patch applied successfully", "file", filepath.Base(patchFile))
	}

	// If we got here, at least one patch had conflicts (which is what we expect)
	// Return an error to signal that .rej files were created
	return fmt.Errorf("patch conflicts detected - .rej files generated")
}

// WritePatchToFile writes the fixed patch content to the original patch file.
func WritePatchToFile(patchContent string, patchFile string) error {
	logger.Info("Writing fixed patch to file", "file", patchFile)

	// Ensure the patch content ends with a newline
	if !strings.HasSuffix(patchContent, "\n") {
		patchContent += "\n"
	}

	// Write the patch to the file
	if err := os.WriteFile(patchFile, []byte(patchContent), 0644); err != nil {
		return fmt.Errorf("writing patch file: %v", err)
	}

	logger.Info("Patch file updated successfully", "file", patchFile)
	return nil
}

// findRejectionFiles locates all .rej files in the cloned repository directory.
// .rej files are created by git am when patches fail to apply.
func findRejectionFiles(repoPath string) ([]string, error) {
	var rejFiles []string

	// Check if repo directory exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return rejFiles, nil // No repo directory means no .rej files
	}

	// Walk through the entire cloned repo to find .rej files
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".rej") {
			rejFiles = append(rejFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking repository directory: %v", err)
	}

	return rejFiles, nil
}

// calculateComplexity scores patch failure complexity.
// TODO(Phase 2): Refine complexity calculation based on PoC metrics
// Current formula: complexity = total_hunks + num_files
// Future considerations:
// - Weighted scoring: different weights for hunk types (context vs logic changes)
// - File type weights: go.mod (predictable) vs core logic (complex)
// - Lines changed: larger changes = higher complexity
// - Historical success rates: learn optimal weights from data
// Track metrics: success_rate_by_complexity, avg_attempts_by_complexity, cost_by_complexity
func calculateComplexity(rejFiles []string) (int, error) {
	// Complexity is based on number of failed hunks across all .rej files
	totalHunks := 0

	for _, rejFile := range rejFiles {
		content, err := os.ReadFile(rejFile)
		if err != nil {
			return 0, fmt.Errorf("reading rejection file %s: %v", rejFile, err)
		}

		// Count hunks by counting "@@" markers in the .rej file
		hunks := strings.Count(string(content), "@@")
		// Each hunk has 2 @@ markers (start and end), so divide by 2
		if hunks > 0 {
			totalHunks += hunks / 2
		}
	}

	// Complexity score = number of failed hunks + number of affected files
	complexity := totalHunks + len(rejFiles)

	return complexity, nil
}

// applySinglePatchWithReject applies a single patch file and returns any .rej files generated and application info.
func applySinglePatchWithReject(patchFile string, projectPath string, repoName string, releaseBranch string) ([]string, *types.PatchApplicationResult, error) {
	logger.Info("Applying single patch with reject", "patch", filepath.Base(patchFile), "release_branch", releaseBranch)

	// Determine GIT_TAG file location
	// Try root first, then release branch directory
	var gitTagPath string
	rootGitTag := filepath.Join(projectPath, "GIT_TAG")
	if _, err := os.Stat(rootGitTag); err == nil {
		// GIT_TAG in root
		gitTagPath = rootGitTag
		logger.Info("Project structure: GIT_TAG in root directory")
	} else if releaseBranch != "" {
		// Try GIT_TAG in release branch directory
		branchGitTag := filepath.Join(projectPath, releaseBranch, "GIT_TAG")
		if _, err := os.Stat(branchGitTag); err == nil {
			gitTagPath = branchGitTag
			logger.Info("Project structure: GIT_TAG in release branch directory", "release_branch", releaseBranch)
		} else {
			// Last resort: assume root (will error if not found)
			gitTagPath = rootGitTag
			logger.Info("Project structure: GIT_TAG in root directory (fallback)")
		}
	} else {
		// No release branch specified, assume root
		gitTagPath = rootGitTag
		logger.Info("Project structure: GIT_TAG in root directory (no release branch)")
	}

	// Read GIT_TAG directly from file to avoid Makefile chicken-and-egg problem
	// (Makefile sets GIT_TAG=non-existent when RELEASE_BRANCH not provided for release-branched projects)
	gitTagBytes, err := os.ReadFile(gitTagPath)
	if err != nil {
		return nil, nil, fmt.Errorf("reading GIT_TAG file at %s: %v", gitTagPath, err)
	}
	gitTag := strings.TrimSpace(string(gitTagBytes))

	// Build the GIT_CHECKOUT_TARGET: $(REPO)/eks-anywhere-checkout-$(GIT_TAG)
	checkoutTarget := fmt.Sprintf("%s/eks-anywhere-checkout-%s", repoName, gitTag)

	// Ensure the repo is checked out (but don't apply patches)
	// ALWAYS pass RELEASE_BRANCH to avoid Makefile errors (some projects require it)
	checkoutCmd := exec.Command("make", "-C", projectPath, checkoutTarget)
	if releaseBranch != "" {
		checkoutCmd.Env = append(os.Environ(), fmt.Sprintf("RELEASE_BRANCH=%s", releaseBranch))
		logger.Info("Running checkout with RELEASE_BRANCH", "branch", releaseBranch)
	}
	checkoutOutput, err := checkoutCmd.CombinedOutput()
	if err != nil {
		return nil, nil, fmt.Errorf("make %s failed: %v\nOutput: %s", checkoutTarget, err, checkoutOutput)
	}

	logger.Info("Repository checked out successfully", "tag", gitTag)

	// The cloned repo directory
	repoPath := filepath.Join(projectPath, repoName)

	// Check if repo was cloned
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("cloned repository not found at %s", repoPath)
	}

	// CRITICAL: Reset repository to clean state BEFORE extracting context
	// This ensures we're not reading from files modified by previous patch attempts
	logger.Info("Resetting repository to clean state")
	resetCmd := exec.Command("git", "-C", repoPath, "reset", "--hard", "HEAD")
	if err := resetCmd.Run(); err != nil {
		logger.Info("Warning: git reset failed", "error", err)
		// Continue anyway - might be first time
	}

	cleanCmd := exec.Command("git", "-C", repoPath, "clean", "-fd")
	if err := cleanCmd.Run(); err != nil {
		logger.Info("Warning: git clean failed", "error", err)
		// Continue anyway
	}

	logger.Info("Repository reset to clean state")

	// Configure git in the cloned repo (same as Common.mk does for patch application)
	configEmailCmd := exec.Command("git", "-C", repoPath, "config", "user.email", constants.PatchApplyGitUserEmail)
	if err := configEmailCmd.Run(); err != nil {
		return nil, nil, fmt.Errorf("configuring git user.email: %v", err)
	}

	configNameCmd := exec.Command("git", "-C", repoPath, "config", "user.name", constants.PatchApplyGitUserName)
	if err := configNameCmd.Run(); err != nil {
		return nil, nil, fmt.Errorf("configuring git user.name: %v", err)
	}

	// CRITICAL: Extract pristine content BEFORE applying patch
	// This ensures we capture the original state before git apply modifies files
	// Now that we've reset to clean state, this will be truly pristine
	logger.Info("Extracting pristine file content before applying patch")
	pristineContent, err := extractPristineContent(patchFile, repoPath)
	if err != nil {
		logger.Info("Warning: failed to extract pristine content", "error", err)
		// Continue anyway - we'll try to work with what we have
	} else {
		logger.Info("Extracted pristine content", "files", len(pristineContent))
	}

	// Apply this specific patch using git apply --reject
	// Need to use absolute path for patch file since we're running git from the repo directory
	absPatchFile, err := filepath.Abs(patchFile)
	if err != nil {
		return nil, nil, fmt.Errorf("getting absolute path for patch file: %v", err)
	}

	logger.Info("Applying patch with git apply --reject",
		"patch", filepath.Base(patchFile),
		"repo_path", repoPath,
		"patch_path", absPatchFile)

	cmd := exec.Command("git", "-C", repoPath, "apply", "--reject", "--whitespace=fix", absPatchFile)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// Parse git apply output to detect offset hunks
	result := &types.PatchApplicationResult{
		OffsetFiles:     make(map[string]int),
		GitOutput:       outputStr,
		PristineContent: pristineContent, // Store pristine content for LLM
	}

	// Parse output line by line to detect offsets
	// Git output format:
	//   Checking patch go.sum...
	//   Hunk #1 succeeded at 935 (offset 2 lines).
	var currentFile string
	scanner := bufio.NewScanner(strings.NewReader(outputStr))
	for scanner.Scan() {
		line := scanner.Text()

		// Track current file being checked: "Checking patch go.sum..."
		if strings.HasPrefix(line, "Checking patch ") {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				currentFile = strings.TrimSuffix(parts[2], "...")
			}
		}

		// Detect offset for current file: "Hunk #1 succeeded at 935 (offset 2 lines)."
		if currentFile != "" && strings.Contains(line, "succeeded at") && strings.Contains(line, "offset") {
			offsetRegex := regexp.MustCompile(`offset (\d+) lines?`)
			if match := offsetRegex.FindStringSubmatch(line); len(match) >= 2 {
				offset, _ := strconv.Atoi(match[1])
				result.OffsetFiles[currentFile] = offset
				logger.Info("Detected offset hunk", "file", currentFile, "offset", offset)
			}
		}
	}

	if err != nil {
		// Check if it's a patch conflict (expected) vs other error
		if strings.Contains(outputStr, "patch does not apply") ||
			strings.Contains(outputStr, "Rejected hunk") ||
			strings.Contains(outputStr, "does not exist in index") {
			logger.Info("Patch application failed with conflicts (expected)",
				"patch", filepath.Base(patchFile),
				"output", outputStr)
			// Continue - we'll find the .rej files
		} else {
			return nil, nil, fmt.Errorf("git apply failed for %s: %v\nOutput: %s", patchFile, err, output)
		}
	} else {
		logger.Info("Patch applied successfully without conflicts", "patch", filepath.Base(patchFile))
	}

	// Find .rej files generated for this patch
	rejFiles, err := findRejectionFiles(repoPath)
	if err != nil {
		return nil, nil, fmt.Errorf("finding rejection files: %v", err)
	}

	return rejFiles, result, nil
}

// extractPristineContent reads the original content of all files in the patch BEFORE git apply modifies them.
// This is critical because git apply --reject will modify files that apply successfully (even with offset),
// and we need the ORIGINAL content to show the LLM what needs to be changed.
func extractPristineContent(patchFile string, repoPath string) (map[string]string, error) {
	pristineContent := make(map[string]string)

	// Read the patch file to find all files being modified
	patchContent, err := os.ReadFile(patchFile)
	if err != nil {
		return nil, fmt.Errorf("reading patch file: %v", err)
	}

	// Parse patch to extract filenames
	// Look for: diff --git a/file b/file
	scanner := bufio.NewScanner(strings.NewReader(string(patchContent)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "diff --git") {
			// Extract filename from "diff --git a/file b/file"
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				filename := strings.TrimPrefix(parts[3], "b/")

				// Read the pristine content of this file
				filePath := filepath.Join(repoPath, filename)
				content, err := os.ReadFile(filePath)
				if err != nil {
					logger.Info("Warning: could not read pristine file", "file", filename, "error", err)
					continue
				}

				pristineContent[filename] = string(content)
				logger.Info("Captured pristine content", "file", filename, "size", len(content))
			}
		}
	}

	return pristineContent, nil
}

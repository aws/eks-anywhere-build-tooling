package fixpatches

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// restorePatches restores patches from backup directory to patches directory
func restorePatches(backupDir, patchesDir string) {
	// Remove any new patches that were copied
	newPatches, _ := filepath.Glob(filepath.Join(patchesDir, "*.patch"))
	for _, p := range newPatches {
		os.Remove(p)
	}

	// Restore from backup
	backupPatches, _ := filepath.Glob(filepath.Join(backupDir, "*.patch"))
	for _, backupPatch := range backupPatches {
		destPath := filepath.Join(patchesDir, filepath.Base(backupPatch))
		os.Rename(backupPatch, destPath)
	}
	logger.Info("Restored old patches from backup")
}

// fixAutoscalerPatches handles the special case for kubernetes/autoscaler project.
//
// Unlike other projects, autoscaler patches are regenerated rather than fixed:
// 1. Remove cloud provider code (keeping only Cluster-API)
// 2. Run go mod tidy to update dependencies
// 3. Generate new patches from these changes
//
// This hardcoded approach is necessary because autoscaler's patch workflow
// is fundamentally different from the standard LLM-based patch fixing.
//
// See: projects/kubernetes/autoscaler/README.md lines 25-60
func fixAutoscalerPatches(projectPath string, releaseBranch string, prNumber int) error {
	logger.Info("Detected kubernetes/autoscaler project - using hardcoded fix")

	// Path to the autoscaler source after checkout
	// The repo is cloned at projectPath/autoscaler (not projectPath/releaseBranch/autoscaler)
	autoscalerPath := filepath.Join(projectPath, "autoscaler")
	clusterAutoscalerPath := filepath.Join(autoscalerPath, "cluster-autoscaler")
	buildersPath := filepath.Join(clusterAutoscalerPath, "cloudprovider", "builder")

	// Step 1: Remove builder files (except CAPI-related ones)
	// NOTE: We do NOT remove the cloud provider directories here - that's handled by the
	// REMOVE_CLOUD_PROVIDERS_TARGET in the Makefile during build. The patch only removes builder files.
	// Per README: cd autoscaler/cluster-autoscaler/cloudprovider/builder
	//             ls . | grep -v -e _all.go -e clusterapi.go -e _builder.go | xargs rm
	cloudProviderPath := filepath.Dir(buildersPath) // cloudprovider directory

	// Step 1a: Remove cloud provider directories (except clusterapi, builder, mocks, test)
	// This matches the REMOVE_CLOUD_PROVIDERS_TARGET in Makefile
	logger.Info("Step 1a: Removing cloud provider directories", "path", cloudProviderPath)
	removeDirsCmd := exec.Command("bash", "-c", "ls . | grep -v -e builder -e clusterapi -e mocks -e test -e .go | xargs rm -rf")
	removeDirsCmd.Dir = cloudProviderPath
	if output, err := removeDirsCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("removing cloud provider directories: %v\nOutput: %s", err, output)
	}
	logger.Info("Cloud provider directories removed")

	// Step 1b: Remove non-CAPI builder files
	logger.Info("Step 1b: Removing non-CAPI builder files", "path", buildersPath)
	removeFilesCmd := exec.Command("bash", "-c", "ls . | grep -v -e _all.go -e clusterapi -e _builder.go | xargs rm -f")
	removeFilesCmd.Dir = buildersPath
	if output, err := removeFilesCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("removing builder files: %v\nOutput: %s", err, output)
	}
	logger.Info("Builder files removed successfully")

	// Step 2: Git add only the builder directory changes (not the entire cloudprovider directory)
	logger.Info("Step 2: Staging builder directory changes")
	if err := GitAdd(cloudProviderPath, "builder"); err != nil {
		return fmt.Errorf("git add failed: %v", err)
	}

	// Step 3: Clean references in builder files
	logger.Info("Step 3: Cleaning references in builder files")

	builderFiles := []string{
		filepath.Join(buildersPath, "builder_all.go"),
		filepath.Join(buildersPath, "builder_clusterapi.go"),
		filepath.Join(buildersPath, "cloud_provider_builder.go"),
	}

	for _, file := range builderFiles {
		if err := cleanBuilderFile(file); err != nil {
			logger.Info("Warning: failed to clean builder file", "file", filepath.Base(file), "error", err)
			// Continue anyway - the build might still work
		}
	}

	// Step 3b: Stage the cleaned builder files
	logger.Info("Step 3b: Staging cleaned builder files")
	if err := GitAdd(cloudProviderPath, "builder"); err != nil {
		return fmt.Errorf("git add cleaned files failed: %v", err)
	}

	// Step 4: Commit the cloud provider changes
	// NOTE: Commit message determines patch filename via git format-patch
	// Extract commit message from original patch filename to preserve naming convention
	// e.g., "0001-Remove-Cloud-Provider-Builders-Except-CAPI.patch" -> "Remove Cloud Provider Builders Except CAPI"
	patchesDir := filepath.Join(projectPath, releaseBranch, "patches")
	originalPatch0001, _ := filepath.Glob(filepath.Join(patchesDir, "0001-*.patch"))
	commitMsg := "Remove Cloud Provider Builders Except CAPI" // default
	if len(originalPatch0001) > 0 {
		// Extract commit message from filename: remove "0001-" prefix and ".patch" suffix, replace "-" with " "
		baseName := filepath.Base(originalPatch0001[0])
		baseName = strings.TrimPrefix(baseName, "0001-")
		baseName = strings.TrimSuffix(baseName, ".patch")
		commitMsg = strings.ReplaceAll(baseName, "-", " ")
		logger.Info("Extracted commit message from original patch", "filename", filepath.Base(originalPatch0001[0]), "commitMsg", commitMsg)
	}

	logger.Info("Step 4: Committing cloud provider removal")
	if err := GitCommit(buildersPath, commitMsg); err != nil {
		return fmt.Errorf("git commit failed: %v", err)
	}

	// Step 5: Generate patches to a TEMP directory first
	// We'll validate them before replacing the old patches
	logger.Info("Step 5: Generating patches to temp directory for validation")
	tempPatchesDir := filepath.Join(projectPath, releaseBranch, "patches-temp")

	// Create temp directory
	if err := os.MkdirAll(tempPatchesDir, 0755); err != nil {
		return fmt.Errorf("creating temp patches directory: %v", err)
	}
	// Ensure cleanup on any exit path
	defer os.RemoveAll(tempPatchesDir)

	// Copy the 0002 patch (GCE dependencies) to temp dir - it's static and doesn't change
	gcePatchSrc := filepath.Join(patchesDir, "0002-Remove-additional-GCE-Dependencies.patch")
	gcePatchDst := filepath.Join(tempPatchesDir, "0002-Remove-additional-GCE-Dependencies.patch")
	if _, err := os.Stat(gcePatchSrc); err == nil {
		gcePatchContent, err := os.ReadFile(gcePatchSrc)
		if err != nil {
			return fmt.Errorf("reading GCE patch: %v", err)
		}
		if err := os.WriteFile(gcePatchDst, gcePatchContent, 0644); err != nil {
			return fmt.Errorf("copying GCE patch to temp: %v", err)
		}
	}

	// Generate patch 0001 to temp directory
	logger.Info("Step 5a: Generating cloud provider removal patch")
	output, err := GitFormatPatch(buildersPath, tempPatchesDir, 1)
	if err != nil {
		return fmt.Errorf("git format-patch failed: %v", err)
	}

	// Extract the generated patch filename from output
	patchFile1 := strings.TrimSpace(output)
	logger.Info("Generated cloud provider patch", "file", filepath.Base(patchFile1))

	// Step 5b: Apply any additional patches (like 0002-Remove-additional-GCE-Dependencies.patch)
	// This must be done AFTER generating patch 0001, so it doesn't get included in that patch
	// Note: The GCE patch already exists and shouldn't be regenerated
	patchFiles, err2 := filepath.Glob(filepath.Join(patchesDir, "0002-*.patch"))
	if err2 == nil && len(patchFiles) > 0 {
		for _, patchFile := range patchFiles {
			patchName := filepath.Base(patchFile)
			// Skip the go.mod patch (we'll regenerate that)
			if strings.Contains(patchName, "go.mod") || strings.Contains(patchName, "go-mod") {
				continue
			}

			logger.Info("Applying additional patch", "patch", patchName)
			if output, err := GitApply(autoscalerPath, patchFile, false, false); err != nil {
				logger.Info("Warning: failed to apply patch", "patch", patchName, "error", err, "output", output)
			} else {
				// Don't stage the GCE patch changes - they're already in patch 0002
				// We just apply them so go mod tidy can see the removed imports
				logger.Info("Applied additional patch (not staged)", "patch", patchName)
			}
		}
	}

	// Step 6: Run go mod tidy
	// This should detect removed imports from both the cloud provider removal and GCE patch
	logger.Info("Step 6: Running go mod tidy")

	// Read project's Go version from GOLANG_VERSION file
	golangVersionFile := filepath.Join(projectPath, releaseBranch, "GOLANG_VERSION")
	golangVersion := ""
	if versionBytes, err := os.ReadFile(golangVersionFile); err == nil {
		golangVersion = strings.TrimSpace(string(versionBytes))
		logger.Info("Using project Go version", "version", golangVersion)
	}

	// Construct go binary path - builder image has versions at /go/go${version}/bin
	goBinary := "go"
	if golangVersion != "" {
		goPath := fmt.Sprintf("/go/go%s/bin/go", golangVersion)
		if _, err := os.Stat(goPath); err == nil {
			goBinary = goPath
		} else {
			logger.Info("Go version not found at expected path, using default", "path", goPath)
		}
	}

	goModTidyCmd := exec.Command(goBinary, "mod", "tidy")
	goModTidyCmd.Dir = clusterAutoscalerPath
	if output, err := goModTidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %v\nOutput: %s", err, output)
	}

	// Step 7: Check if go.mod or go.sum changed
	logger.Info("Step 7: Checking for go.mod/go.sum changes")
	statusOutput, _ := GitStatus(clusterAutoscalerPath, true, "go.mod", "go.sum")
	hasGoModChanges := len(strings.TrimSpace(statusOutput)) > 0

	if hasGoModChanges {
		logger.Info("go.mod/go.sum have changes", "status", statusOutput)

		// Step 7a: Git add go.mod and go.sum
		logger.Info("Step 7a: Staging go.mod and go.sum")
		if err := GitAdd(clusterAutoscalerPath, "go.mod", "go.sum"); err != nil {
			return fmt.Errorf("git add go.mod/go.sum failed: %v", err)
		}

		// Step 8: Commit go.mod changes
		logger.Info("Step 8: Committing go.mod changes")
		commitMsg2 := "Update go.mod Dependencies"
		if err := GitCommit(clusterAutoscalerPath, commitMsg2); err != nil {
			return fmt.Errorf("git commit go.mod failed: %v", err)
		}
	} else {
		logger.Info("No go.mod/go.sum changes detected - skipping commit")
	}

	// Step 9: Generate patch for go.mod changes to temp directory
	logger.Info("Step 9: Generating patch for go.mod changes")
	output2, err := GitFormatPatch(clusterAutoscalerPath, tempPatchesDir, 1)
	if err != nil {
		return fmt.Errorf("git format-patch go.mod failed: %v", err)
	}

	// Extract the generated patch filename from output
	patchFile2 := strings.TrimSpace(output2)
	logger.Info("Generated go.mod patch", "file", filepath.Base(patchFile2))

	// Step 10: Rename patches to match expected names
	// Check if we have a 0002 patch (GCE dependencies), if so, go.mod should be 0003
	gcePatches, _ := filepath.Glob(filepath.Join(tempPatchesDir, "0002-Remove-additional-GCE-Dependencies.patch"))
	hasGCEPatch := len(gcePatches) > 0

	var expectedGoModPatchName string
	if hasGCEPatch {
		expectedGoModPatchName = "0003-Update-go.mod-Dependencies.patch"
	} else {
		expectedGoModPatchName = "0002-Update-go.mod-Dependencies.patch"
	}

	// Rename the go.mod patch to the expected name
	expectedGoModPatchPath := filepath.Join(tempPatchesDir, expectedGoModPatchName)
	if patchFile2 != expectedGoModPatchPath {
		logger.Info("Renaming go.mod patch", "from", filepath.Base(patchFile2), "to", expectedGoModPatchName)
		if err := os.Rename(patchFile2, expectedGoModPatchPath); err != nil {
			logger.Info("Warning: failed to rename patch", "error", err)
		}
		patchFile2 = expectedGoModPatchPath
	}

	logger.Info("✅ Patches generated in temp directory",
		"patches_generated", 2,
		"patch1", filepath.Base(patchFile1),
		"patch2", filepath.Base(patchFile2))

	// Step 10b: VALIDATE the new patches by running make attribution-checksums with temp patches
	// We need to temporarily swap the patches to validate
	logger.Info("Step 10b: Validating new patches apply cleanly")

	// Backup old patches
	backupDir := filepath.Join(projectPath, releaseBranch, "patches-backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("creating backup directory: %v", err)
	}
	defer os.RemoveAll(backupDir)

	// Move old patches to backup
	oldPatches, _ := filepath.Glob(filepath.Join(patchesDir, "*.patch"))
	for _, oldPatch := range oldPatches {
		backupPath := filepath.Join(backupDir, filepath.Base(oldPatch))
		if err := os.Rename(oldPatch, backupPath); err != nil {
			logger.Info("Warning: failed to backup patch", "patch", filepath.Base(oldPatch), "error", err)
		}
	}

	// Copy new patches from temp to patches dir for validation
	newPatches, _ := filepath.Glob(filepath.Join(tempPatchesDir, "*.patch"))
	for _, newPatch := range newPatches {
		destPath := filepath.Join(patchesDir, filepath.Base(newPatch))
		content, err := os.ReadFile(newPatch)
		if err != nil {
			// Restore old patches and fail
			restorePatches(backupDir, patchesDir)
			return fmt.Errorf("reading new patch %s: %v", filepath.Base(newPatch), err)
		}
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			restorePatches(backupDir, patchesDir)
			return fmt.Errorf("writing new patch %s: %v", filepath.Base(newPatch), err)
		}
	}

	// Now validate by running make attribution-checksums
	logger.Info("Step 10c: Running make attribution-checksums to validate patches")
	makeCmd := exec.Command("make", "-C", projectPath, "attribution-checksums")
	makeCmd.Env = append(os.Environ(), fmt.Sprintf("RELEASE_BRANCH=%s", releaseBranch))
	makeOutput, makeErr := makeCmd.CombinedOutput()

	if makeErr != nil {
		// Validation failed! Restore old patches
		logger.Info("❌ Patch validation FAILED - restoring old patches", "error", makeErr, "output", string(makeOutput))
		restorePatches(backupDir, patchesDir)
		return fmt.Errorf("new patches failed validation (make attribution-checksums): %v\nOutput: %s", makeErr, makeOutput)
	}

	logger.Info("✅ Patch validation PASSED - new patches apply cleanly")

	// Update patchFile1 and patchFile2 to point to the final locations
	patchFile1 = filepath.Join(patchesDir, filepath.Base(patchFile1))
	patchFile2 = filepath.Join(patchesDir, filepath.Base(patchFile2))

	// Step 11: Stage, commit, and push the generated patches to the PR branch
	// The patches are in projectPath/<releaseBranch>/patches/
	// We need to commit from the repo root (cwd)
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return fmt.Errorf("getting current working directory: %v", cwdErr)
	}

	logger.Info("Step 11: Staging generated patches in build-tooling repo")
	patchesRelPath := filepath.Join("projects", "kubernetes", "autoscaler", releaseBranch, "patches")
	if err := GitAdd(cwd, patchesRelPath); err != nil {
		return fmt.Errorf("staging patches in build-tooling repo: %v", err)
	}

	// Also stage CHECKSUMS and ATTRIBUTION.txt if they exist
	checksumsRelPath := filepath.Join("projects", "kubernetes", "autoscaler", releaseBranch, "CHECKSUMS")
	if _, err := os.Stat(filepath.Join(cwd, checksumsRelPath)); err == nil {
		if err := GitAdd(cwd, checksumsRelPath); err != nil {
			logger.Info("Warning: failed to stage CHECKSUMS", "error", err)
		}
	}

	attributionRelPath := filepath.Join("projects", "kubernetes", "autoscaler", releaseBranch, "ATTRIBUTION.txt")
	if _, err := os.Stat(filepath.Join(cwd, attributionRelPath)); err == nil {
		if err := GitAdd(cwd, attributionRelPath); err != nil {
			logger.Info("Warning: failed to stage ATTRIBUTION.txt", "error", err)
		}
	}

	// Commit the changes
	logger.Info("Step 12: Committing patch changes")
	finalCommitMsg := fmt.Sprintf("Auto-fix patches for kubernetes/autoscaler %s\n\nRegenerated patches:\n- %s\n- %s",
		releaseBranch, filepath.Base(patchFile1), filepath.Base(patchFile2))
	if err := GitCommit(cwd, finalCommitMsg); err != nil {
		return fmt.Errorf("committing patches: %v", err)
	}

	// Push to origin (the PR branch)
	// Use --force since this is a bot-owned PR branch and we want idempotent behavior
	// (re-running the tool should succeed even if previous run pushed something)
	logger.Info("Step 13: Pushing to PR branch")
	if _, err := GitCommand(cwd, "push", "--force", "origin", "HEAD"); err != nil {
		return fmt.Errorf("pushing to PR branch: %v", err)
	}

	logger.Info("✅ Patches committed and pushed to PR branch")

	// Step 14: Add comment on PR explaining the automated fix
	if prNumber > 0 {
		logger.Info("Step 14: Adding comment to PR", "pr", prNumber)
		ghClient, err := NewGitHubClient()
		if err != nil {
			logger.Info("Warning: failed to create GitHub client for PR comment", "error", err)
		} else {
			comment := fmt.Sprintf(`## 🤖 Automated Patch Fix

This PR has been automatically updated by the LLM-powered patch fixer.

### What was done:
- Regenerated patches for **kubernetes/autoscaler %s** release branch
- Generated patches:
  - %s (cloud provider removal)
  - %s (go.mod dependencies update)

### ⚠️ Manual Review Required
Please review the regenerated patches to ensure they are correct before merging.

The patches were generated using the standard autoscaler patch workflow:
1. Checkout upstream tag
2. Remove non-CAPI cloud providers
3. Run go mod tidy
4. Generate patches from changes

---
*This comment was automatically generated by the patch-fixer tool.*`,
				releaseBranch, filepath.Base(patchFile1), filepath.Base(patchFile2))

			if err := ghClient.CommentOnPR(prNumber, comment); err != nil {
				logger.Info("Warning: failed to comment on PR", "error", err)
			} else {
				logger.Info("Added comment to PR")
			}
		}
	}

	return nil
}

// cleanBuilderFile removes references to non-CAPI cloud providers from builder files.
// This removes import lines and AvailableCloudProviders entries for all providers except clusterapi.
func cleanBuilderFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %v", err)
	}

	original := string(content)
	lines := strings.Split(original, "\n")
	var cleaned []string
	inImportBlock := false
	inAvailableProvidersBlock := false
	inSwitchBlock := false
	switchBraceCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track if we're in the import block
		if strings.HasPrefix(trimmed, "import (") {
			inImportBlock = true
			cleaned = append(cleaned, line)
			continue
		}
		if inImportBlock && trimmed == ")" {
			inImportBlock = false
			cleaned = append(cleaned, line)
			continue
		}

		// Track if we're in the AvailableCloudProviders array
		if strings.Contains(line, "var AvailableCloudProviders") {
			inAvailableProvidersBlock = true
			cleaned = append(cleaned, line)
			continue
		}
		if inAvailableProvidersBlock && trimmed == "}" {
			inAvailableProvidersBlock = false
			cleaned = append(cleaned, line)
			continue
		}

		// Track if we're in the switch statement
		if strings.Contains(line, "switch opts.CloudProviderName") {
			inSwitchBlock = true
			cleaned = append(cleaned, line)
			continue
		}
		if inSwitchBlock {
			// Count braces to know when switch ends
			switchBraceCount += strings.Count(line, "{") - strings.Count(line, "}")
			if switchBraceCount < 0 {
				inSwitchBlock = false
				cleaned = append(cleaned, line)
				switchBraceCount = 0
				continue
			}
		}

		// In import block: keep only clusterapi and standard imports
		if inImportBlock {
			// Keep lines that contain "clusterapi" or don't look like cloud provider imports
			if strings.Contains(line, "clusterapi") ||
				strings.Contains(line, "\"k8s.io/autoscaler/cluster-autoscaler/cloudprovider\"") ||
				strings.Contains(line, "\"k8s.io/autoscaler/cluster-autoscaler/config\"") ||
				strings.Contains(line, "\"k8s.io/client-go/informers\"") ||
				!strings.Contains(line, "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/") {
				cleaned = append(cleaned, line)
			} else {
				logger.Info("Removing import line", "line", trimmed)
			}
			continue
		}

		// In AvailableCloudProviders block: keep only ClusterAPIProviderName
		if inAvailableProvidersBlock {
			if strings.Contains(line, "ClusterAPIProviderName") ||
				!strings.Contains(line, "ProviderName") {
				cleaned = append(cleaned, line)
			} else {
				logger.Info("Removing provider entry", "line", trimmed)
			}
			continue
		}

		// In switch block: keep only ClusterAPI case
		if inSwitchBlock {
			if strings.Contains(line, "ClusterAPIProviderName") ||
				strings.Contains(line, "clusterapi.Build") ||
				(!strings.Contains(line, "case cloudprovider.") && !strings.Contains(line, "return ")) {
				cleaned = append(cleaned, line)
			} else {
				logger.Info("Removing switch case", "line", trimmed)
			}
			continue
		}

		// Keep all other lines
		cleaned = append(cleaned, line)
	}

	cleanedContent := strings.Join(cleaned, "\n")

	// Only write if we made changes
	if cleanedContent != original {
		if err := os.WriteFile(filePath, []byte(cleanedContent), 0644); err != nil {
			return fmt.Errorf("writing file: %v", err)
		}
		logger.Info("Cleaned builder file", "file", filepath.Base(filePath), "removed_lines", len(lines)-len(cleaned))
	}

	return nil
}

package fixpatches

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

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
func fixAutoscalerPatches(projectPath string, releaseBranch string) error {
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
	logger.Info("Step 4: Committing cloud provider removal")
	commitMsg := "Remove Cloud Provider Builders Except Cluster-API"
	if err := GitCommit(buildersPath, commitMsg); err != nil {
		return fmt.Errorf("git commit failed: %v", err)
	}

	// Step 5: Generate patch for cloud provider removal
	logger.Info("Step 5: Generating patch for cloud provider removal")
	patchesDir := filepath.Join(projectPath, releaseBranch, "patches")
	output, err := GitFormatPatch(buildersPath, patchesDir, 1)
	if err != nil {
		return fmt.Errorf("git format-patch failed: %v", err)
	}

	// Extract the generated patch filename from output
	patchFile1 := strings.TrimSpace(output)
	logger.Info("Generated cloud provider patch", "file", filepath.Base(patchFile1))

	// Step 5b: Apply any additional patches (like 0002-Remove-additional-GCE-Dependencies.patch)
	// This must be done AFTER generating patch 0001, so it doesn't get included in that patch
	// Note: The GCE patch already exists and shouldn't be regenerated
	patchesDirTemp := filepath.Join(projectPath, releaseBranch, "patches")
	patchFiles, err2 := filepath.Glob(filepath.Join(patchesDirTemp, "0002-*.patch"))
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
	goModTidyCmd := exec.Command("go", "mod", "tidy")
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

	// Step 9: Generate patch for go.mod changes
	logger.Info("Step 9: Generating patch for go.mod changes")
	output2, err := GitFormatPatch(clusterAutoscalerPath, patchesDir, 1)
	if err != nil {
		return fmt.Errorf("git format-patch go.mod failed: %v", err)
	}

	// Extract the generated patch filename from output
	patchFile2 := strings.TrimSpace(output2)
	logger.Info("Generated go.mod patch", "file", filepath.Base(patchFile2))

	// Step 10: Rename patches to match expected names
	// Check if we have a 0002 patch (GCE dependencies), if so, go.mod should be 0003
	gcePatches, _ := filepath.Glob(filepath.Join(patchesDir, "0002-Remove-additional-GCE-Dependencies.patch"))
	hasGCEPatch := len(gcePatches) > 0

	var expectedGoModPatchName string
	if hasGCEPatch {
		expectedGoModPatchName = "0003-Update-go.mod-Dependencies.patch"
	} else {
		expectedGoModPatchName = "0002-Update-go.mod-Dependencies.patch"
	}

	// Rename the go.mod patch to the expected name
	expectedGoModPatchPath := filepath.Join(patchesDir, expectedGoModPatchName)
	if patchFile2 != expectedGoModPatchPath {
		logger.Info("Renaming go.mod patch", "from", filepath.Base(patchFile2), "to", expectedGoModPatchName)
		if err := os.Rename(patchFile2, expectedGoModPatchPath); err != nil {
			logger.Info("Warning: failed to rename patch", "error", err)
		}
		patchFile2 = expectedGoModPatchPath
	}

	logger.Info("✅ Autoscaler hardcoded fix complete",
		"patches_generated", 2,
		"patch1", filepath.Base(patchFile1),
		"patch2", filepath.Base(patchFile2))

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

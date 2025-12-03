package upgrade

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	eksdistrorelease "github.com/aws/eks-distro-build-tooling/release/api/v1alpha1"
	"github.com/ghodss/yaml"
	gogithub "github.com/google/go-github/v53/github"
	"github.com/pelletier/go-toml/v2"
	goyamlv3 "gopkg.in/yaml.v3"
	sigsyaml "sigs.k8s.io/yaml"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/commands/fixpatches"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/ecrpublic"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/git"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/github"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/command"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/file"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// Run contains the business logic to execute the `upgrade` subcommand.
func Run(upgradeOptions *types.UpgradeOptions) error {
	var currentRevision, latestRevision, patchesWarningComment string
	var isTrackedByCommitHash, patchApplySucceeded bool
	var totalPatchCount int
	var updatedFiles []string
	var pullRequest *gogithub.PullRequest
	failedSteps := map[string]error{}

	projectName := upgradeOptions.ProjectName

	// Get org and repository name from project name.
	projectOrg := strings.Split(projectName, "/")[0]
	projectRepo := strings.Split(projectName, "/")[1]

	// Check if branch name environment variable has been set.
	branchName, ok := os.LookupEnv(constants.BranchNameEnvVar)
	if !ok {
		branchName = constants.MainBranchName
	}

	// Check if base repository owner environment variable has been set.
	baseRepoOwner, ok := os.LookupEnv(constants.BaseRepoOwnerEnvvar)
	if !ok {
		return fmt.Errorf("BASE_REPO_OWNER environment variable is not set")
	}

	// Check if head repository owner environment variable has been set.
	headRepoOwner, ok := os.LookupEnv(constants.HeadRepoOwnerEnvvar)
	if !ok {
		return fmt.Errorf("HEAD_REPO_OWNER environment variable is not set")
	}

	// Check if GitHub token environment variable has been set.
	githubToken, ok := os.LookupEnv(constants.GitHubTokenEnvvar)
	if !ok {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}
	client := gogithub.NewTokenClient(context.Background(), githubToken)

	// Skip project upgrade if it is in the ProjectsUpgradedOnlyOnMainBranch list and branch is not main
	if branchName != constants.MainBranchName && slices.Contains(constants.ProjectsUpgradedOnlyOnMainBranch, projectName) {
		logger.Info(fmt.Sprintf("Skipping upgrade for project %s on %s branch", projectName, branchName))
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("retrieving current working directory: %v", err)
	}

	skippedProjectsFilepath := filepath.Join(cwd, constants.SkippedProjectsFile)
	contents, err := os.ReadFile(skippedProjectsFilepath)
	if err != nil {
		return fmt.Errorf("reading skipped projects file: %v", err)
	}
	skippedProjects := strings.Split(string(contents), "\n")
	if slices.Contains(skippedProjects, projectName) {
		logger.Info("Project is in SKIPPED_PROJECTS list. Skipping upgrade")
		return nil
	}

	// Clone the eks-anywhere-build-tooling repository.
	buildToolingRepoPath := filepath.Join(cwd, constants.BuildToolingRepoName)
	repo, headCommit, err := git.CloneRepo(fmt.Sprintf(constants.BuildToolingRepoURL, baseRepoOwner), buildToolingRepoPath, headRepoOwner, branchName)
	if err != nil {
		return fmt.Errorf("cloning build-tooling repo: %v", err)
	}

	// Get the worktree corresponding to the cloned repository.
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("getting repo's current worktree: %v", err)
	}

	// Checkout the eks-anywhere-build-tooling repo at the provided branch name.
	createBranch := (branchName != constants.MainBranchName)
	err = git.Checkout(worktree, branchName, createBranch)
	if err != nil {
		return fmt.Errorf("checking out worktree at branch %s: %v", branchName, err)
	}

	// Reset current worktree to get a clean index.
	err = git.ResetToHEAD(worktree, headCommit)
	if err != nil {
		return fmt.Errorf("resetting new branch to [origin/%s] HEAD: %v", branchName, err)
	}

	var headBranchName, baseBranchName, commitMessage, pullRequestBody string
	if isEKSDistroUpgrade(projectName) {
		headBranchName = fmt.Sprintf("update-eks-distro-latest-releases-%s", branchName)
		baseBranchName = branchName
		commitMessage = "Bump EKS Distro releases to latest"
		pullRequestBody = constants.EKSDistroUpgradePullRequestBody

		// Checkout a new branch to keep track of version upgrade chaneges.
		err = git.Checkout(worktree, headBranchName, true)
		if err != nil {
			return fmt.Errorf("checking out worktree at branch %s: %v", headBranchName, err)
		}

		// Reset current worktree to get a clean index.
		err = git.ResetToHEAD(worktree, headCommit)
		if err != nil {
			return fmt.Errorf("resetting new branch to [origin/%s] HEAD: %v", branchName, err)
		}

		isUpdated, err := updateEKSDistroReleasesFile(buildToolingRepoPath)
		if err != nil {
			return fmt.Errorf("updating EKS Distro releases file: %v", err)
		}
		if isUpdated {
			updatedFiles = append(updatedFiles, constants.EKSDistroLatestReleasesFile)
		}
	} else if isEKSDistroBuildToolingUpgrade(projectName) {
		headBranchName = fmt.Sprintf("update-eks-distro-base-image-tag-files-%s", branchName)
		baseBranchName = branchName
		commitMessage = "Bump EKS Distro base image tag files to latest"

		// Checkout a new branch to keep track of version upgrade chaneges.
		err = git.Checkout(worktree, headBranchName, true)
		if err != nil {
			return fmt.Errorf("checking out worktree at branch %s: %v", headBranchName, err)
		}

		// Reset current worktree to get a clean index.
		err = git.ResetToHEAD(worktree, headCommit)
		if err != nil {
			return fmt.Errorf("resetting new branch to [origin/%s] HEAD: %v", branchName, err)
		}

		eksDistroBaseTagFilesGlobPattern := filepath.Join(buildToolingRepoPath, constants.EKSDistroBaseTagFilesPattern)
		eksDistroBaseTagFilesGlob, err := filepath.Glob(eksDistroBaseTagFilesGlobPattern)
		if err != nil {
			return fmt.Errorf("finding filenames matching EKS Distro Base tag file pattern [%s]: %v", constants.EKSDistroBaseTagFilesPattern, err)
		}

		updatedPackages, isUpdated, err := updateEKSDistroBaseImageTagFiles(client, buildToolingRepoPath, eksDistroBaseTagFilesGlob)
		if err != nil {
			return fmt.Errorf("updating EKS Distro base tag files: %v", err)
		}
		if isUpdated {
			pullRequestBody = fmt.Sprintf(constants.EKSDistroBuildToolingUpgradePullRequestBody, updatedPackages)
			for _, tagFile := range eksDistroBaseTagFilesGlob {
				tagFileRelativePath, err := filepath.Rel(buildToolingRepoPath, tagFile)
				if err != nil {
					return fmt.Errorf("getting relative path for tag file: %v", err)
				}
				updatedFiles = append(updatedFiles, tagFileRelativePath)
			}
		}
	} else {
		// Validate if the project name provided exists in the repository.
		projectPath := filepath.Join("projects", projectName)
		projectRootFilepath := filepath.Join(buildToolingRepoPath, projectPath)
		if _, err := os.Stat(projectRootFilepath); os.IsNotExist(err) {
			return fmt.Errorf("invalid project name %s", projectName)
		}

		// Load upstream projects tracker file.
		upstreamProjectsTrackerFilePath := filepath.Join(buildToolingRepoPath, constants.UpstreamProjectsTrackerFile)
		_, targetRepo, err := loadUpstreamProjectsTrackerFile(upstreamProjectsTrackerFilePath, projectOrg, projectRepo)
		if err != nil {
			return fmt.Errorf("loading upstream projects tracker file: %v", err)
		}

		// Validate whether the given project is release-branched.
		var isReleaseBranched bool
		var currentVersion types.Version
		var versionIndex int
		if len(targetRepo.Versions) > 1 {
			isReleaseBranched = true
		}
		releaseBranch := os.Getenv(constants.ReleaseBranchEnvvar)
		if releaseBranch == "" {
			releaseBranch, err = getDefaultReleaseBranch(buildToolingRepoPath)
			if err != nil {
				return fmt.Errorf("getting default EKS Distro release branch: %v", err)
			}
			os.Setenv(constants.ReleaseBranchEnvvar, releaseBranch)
		}
		if isReleaseBranched {
			supportedReleaseBranches, err := getSupportedReleaseBranches(buildToolingRepoPath)
			if err != nil {
				return fmt.Errorf("getting supported EKS Distro release branches: %v", err)
			}

			versionIndex = slices.Index(supportedReleaseBranches, releaseBranch)
		} else {
			versionIndex = 0
		}
		currentVersion = targetRepo.Versions[versionIndex]

		if currentVersion.Tag != "" {
			currentRevision = currentVersion.Tag
		} else if currentVersion.Commit != "" {
			currentRevision = currentVersion.Commit
			isTrackedByCommitHash = true
		}

		// Check if project to be upgraded has patches
		projectHasPatches := false
		patchesDirectory := filepath.Join(projectRootFilepath, constants.PatchesDirectory)
		if isReleaseBranched {
			patchesDirectory = filepath.Join(projectRootFilepath, releaseBranch, constants.PatchesDirectory)
		}
		if _, err := os.Stat(patchesDirectory); err == nil {
			projectHasPatches = true
			patchFiles, err := os.ReadDir(patchesDirectory)
			if err != nil {
				return fmt.Errorf("reading patches directory: %v", err)
			}
			totalPatchCount = len(patchFiles)
		}

		headBranchName = fmt.Sprintf("update-%s-%s-%s", projectOrg, projectRepo, branchName)
		baseBranchName = branchName
		commitMessage = fmt.Sprintf("Bump %s to latest release", projectName)
		if isReleaseBranched {
			headBranchName = fmt.Sprintf("update-%s-%s-%s-%s", projectOrg, projectRepo, releaseBranch, branchName)
			commitMessage = fmt.Sprintf("Bump %s %s release branch to latest release", projectName, releaseBranch)
		}

		var latestRevision string
		var needsUpgrade bool
		if projectName == "cilium/cilium" {
			latestRevision, needsUpgrade, err = ecrpublic.GetLatestRevision(constants.CiliumImageRepository, currentRevision, branchName)
			if err != nil {
				return fmt.Errorf("getting latest revision from ECR Public: %v", err)
			}
		} else {
			// Get latest revision for the project from GitHub.
			latestRevision, needsUpgrade, err = github.GetLatestRevision(client, projectOrg, projectRepo, currentRevision, branchName, isTrackedByCommitHash, isReleaseBranched)
			if err != nil {
				return fmt.Errorf("getting latest revision from GitHub: %v", err)
			}
		}

		prLabels := constants.DefaultProjectUpgradePRLabels
		if slices.Contains(constants.CuratedPackagesProjects, projectName) {
			prLabels = constants.PackagesProjectUpgradePRLabels
		}
		pullRequestBody = fmt.Sprintf(constants.DefaultUpgradePullRequestBody, projectOrg, projectRepo, currentRevision, latestRevision, strings.Join(prLabels, "\n"))

		// Upgrade project if latest commit was made after current commit and the semver of the latest revision is
		// greater than the semver of the current version.
		if needsUpgrade || slices.Contains(constants.ProjectsWithUnconventionalUpgradeFlows, projectName) {
			// Checkout a new branch to keep track of version upgrade chaneges.
			err = git.Checkout(worktree, headBranchName, true)
			if err != nil {
				return fmt.Errorf("checking out worktree at branch %s: %v", headBranchName, err)
			}

			// Reset current worktree to get a clean index.
			err = git.ResetToHEAD(worktree, headCommit)
			if err != nil {
				return fmt.Errorf("resetting new branch to [origin/%s] HEAD: %v", branchName, err)
			}

			if needsUpgrade {
				logger.Info("Project is out of date.", "Current version", currentRevision, "Latest version", latestRevision)

				// Reload upstream projects tracker file to get its original value instead of
				// the updated one from another project's previous upgrade
				projectsList, targetRepo, err := loadUpstreamProjectsTrackerFile(upstreamProjectsTrackerFilePath, projectOrg, projectRepo)
				if err != nil {
					return fmt.Errorf("reloading upstream projects tracker file: %v", err)
				}
				if isTrackedByCommitHash {
					targetRepo.Versions[versionIndex].Commit = latestRevision
				} else {
					targetRepo.Versions[versionIndex].Tag = latestRevision
				}

				// Update the Git tag file corresponding to the project
				logger.Info("Updating Git tag file corresponding to the project")
				projectGitTagRelativePath, err := updateProjectVersionFile(buildToolingRepoPath, constants.GitTagFile, projectName, latestRevision, releaseBranch, isReleaseBranched)
				if err != nil {
					return fmt.Errorf("updating project GIT_TAG file: %v", err)
				}
				updatedFiles = append(updatedFiles, projectGitTagRelativePath)

				var latestGoVersion string
				if currentVersion.GoVersion != "N/A" {
					currentGoVersion := currentVersion.GoVersion
					// Get Go version corresponding to the latest revision of the project.
					latestGoVersion, err := github.GetGoVersionForLatestRevision(client, projectOrg, projectRepo, latestRevision)
					if err != nil {
						return fmt.Errorf("getting latest Go version for release %s: %v", latestRevision, err)
					}

					// Get the minor version for the current revision's Go version.
					currentGoMinorVersion, err := strconv.Atoi(strings.Split(currentGoVersion, ".")[1])
					if err != nil {
						return fmt.Errorf("getting current Go minor version: %v", err)
					}

					// Get the major version for the latest revision's Go version.
					latestGoMinorVersion, err := strconv.Atoi(strings.Split(latestGoVersion, ".")[1])
					if err != nil {
						return fmt.Errorf("getting latest Go minor version: %v", err)
					}

					// If the Go version has been updated in the latest revision, then update the Go version file corresponding to the project.
					if latestGoMinorVersion > currentGoMinorVersion {
						logger.Info("Project Go version needs to be updated.", "Current Go version", currentGoVersion, "Latest Go version", latestGoVersion)
						targetRepo.Versions[versionIndex].GoVersion = latestGoVersion

						logger.Info("Updating Go version file corresponding to the project")
						projectGoVersionRelativePath, err := updateProjectVersionFile(buildToolingRepoPath, constants.GoVersionFile, projectName, latestGoVersion, releaseBranch, isReleaseBranched)
						if err != nil {
							return fmt.Errorf("updating project GOLANG_VERSION file: %v", err)
						}
						updatedFiles = append(updatedFiles, projectGoVersionRelativePath)
					}
				} else {
					latestGoVersion = "N/A"
					targetRepo.Versions[versionIndex].GoVersion = latestGoVersion
				}

				// Update the tag and Go version in the section of the upstream projects tracker file corresponding to the given project.
				logger.Info("Updating Git tag and Go version in upstream projects tracker file")
				err = updateUpstreamProjectsTrackerFile(&projectsList, buildToolingRepoPath, upstreamProjectsTrackerFilePath)
				if err != nil {
					return fmt.Errorf("updating upstream projects tracker file: %v", err)
				}
				updatedFiles = append(updatedFiles, constants.UpstreamProjectsTrackerFile)

				// Update the version in the project's README file.
				logger.Info("Updating project README file")
				projectReadmePath := filepath.Join(projectPath, constants.ReadmeFile)
				err = updateProjectReadmeVersion(buildToolingRepoPath, projectOrg, projectRepo)
				if err != nil {
					return fmt.Errorf("updating version in project README: %v", err)
				}
				updatedFiles = append(updatedFiles, projectReadmePath)

				// If project has patches, attempt to apply them. Track failed patches and files that failed to apply, if any.
				if projectHasPatches {
					appliedPatchesCount, failedPatch, applyFailedFiles, err := applyPatchesToRepo(projectRootFilepath, projectRepo, totalPatchCount)
					if appliedPatchesCount == totalPatchCount {
						patchApplySucceeded = true
					}
					if !patchApplySucceeded {
						failedSteps["Patch application"] = err
						patchesWarningComment = fmt.Sprintf(constants.FailedPatchesCommentBody, appliedPatchesCount, totalPatchCount, failedPatch, applyFailedFiles)

						// Publish EventBridge event for automatic patch fixing
						if pullRequest != nil {
							event := fixpatches.PatchFailureEvent{
								Project:       projectName,
								PRNumber:      *pullRequest.Number,
								Branch:        branchName,
								FailedPatches: []string{failedPatch},
								Reason:        fmt.Sprintf("Patch failed to apply to files: %s", applyFailedFiles),
								RepoOwner:     baseRepoOwner,
								RepoName:      constants.BuildToolingRepoName,
							}

							if err := fixpatches.PublishPatchFailureEvent(event); err != nil {
								logger.Info("Failed to publish patch failure event", "error", err)
								// Don't fail the upgrade if event publishing fails
							}
						}
					}
				}

				// If project doesn't have patches, or it does and they were applied successfully, then update the checksums file
				// and attribution file(s) corresponding to the project.
				if !projectHasPatches || patchApplySucceeded {
					projectChecksumsFile := filepath.Join(projectRootFilepath, constants.ChecksumsFile)
					projectChecksumsFileRelativePath := filepath.Join(projectPath, constants.ChecksumsFile)
					projectAttributionFileGlob := filepath.Join(projectRootFilepath, constants.AttributionsFilePattern)
					if isReleaseBranched {
						projectChecksumsFile = filepath.Join(projectRootFilepath, releaseBranch, constants.ChecksumsFile)
						projectChecksumsFileRelativePath = filepath.Join(projectPath, releaseBranch, constants.ChecksumsFile)
						projectAttributionFileGlob = filepath.Join(projectRootFilepath, releaseBranch, constants.AttributionsFilePattern)
					}
					if _, err := os.Stat(projectChecksumsFile); err == nil {
						logger.Info("Updating project checksums and attribution files")
						err = updateChecksumsAttributionFiles(projectRootFilepath)
						if err != nil {
							failedSteps["Checksums and attribution generation"] = err
						} else {
							updatedFiles = append(updatedFiles, projectChecksumsFileRelativePath)

							// Attribution files can have a binary name prefix so we use a common prefix regular expression
							// and glob them to cover all possibilities.
							projectAttributionFileGlob, err := filepath.Glob(projectAttributionFileGlob)
							if err != nil {
								return fmt.Errorf("finding filenames matching attribution file pattern [%s]: %v", constants.AttributionsFilePattern, err)
							}
							for _, attributionFile := range projectAttributionFileGlob {
								attributionFileRelativePath, err := filepath.Rel(buildToolingRepoPath, attributionFile)
								if err != nil {
									return fmt.Errorf("getting relative path for attribution file: %v", err)
								}
								updatedFiles = append(updatedFiles, attributionFileRelativePath)
							}
						}
					}
				}

				op, message := getProjectSpecificUpdateOperation(projectName)
				if op != nil {
					updatedProjectFiles, err := op(projectRootFilepath, projectPath)
					if err != nil {
						failedSteps[message] = err
					} else {
						updatedFiles = append(updatedFiles, updatedProjectFiles...)
					}
				}
			}

			if projectName == "kubernetes-sigs/image-builder" {
				currentBottlerocketVersion, latestBottlerocketVersion, updatedBRFiles, err := updateBottlerocketVersionFiles(client, projectRootFilepath, projectPath, branchName)
				if err != nil {
					failedSteps["Bottlerocket version upgrade"] = err
				} else {
					if len(updatedBRFiles) > 0 {
						updatedFiles = append(updatedFiles, updatedBRFiles...)
						if len(updatedFiles) == len(updatedBRFiles) {
							headBranchName = fmt.Sprintf("update-bottlerocket-releases-%s", branchName)
							commitMessage = "Bump Bottlerocket versions to latest release"
							pullRequestBody = fmt.Sprintf(constants.BottlerocketUpgradePullRequestBody, currentBottlerocketVersion, latestBottlerocketVersion)
						} else {
							headBranchName = fmt.Sprintf("update-%s-%s-and-bottlerocket-%s", projectOrg, projectRepo, branchName)
							commitMessage = fmt.Sprintf("Bump %s and Bottlerocket versions to latest release", projectName)
							pullRequestBody = fmt.Sprintf(constants.CombinedImageBuilderBottlerocketUpgradePullRequestBody, currentRevision, latestRevision, currentBottlerocketVersion, latestBottlerocketVersion)
						}

						err = git.Checkout(worktree, headBranchName, true)
						if err != nil {
							return fmt.Errorf("checking out worktree at branch %s: %v", headBranchName, err)
						}
					}
				}
			}
		} else if latestRevision == currentRevision {
			logger.Info("Project is at the latest available version.", "Current version", currentRevision, "Latest version", latestRevision)
		}
	}

	if len(updatedFiles) > 0 {
		// Add all the updated files to the index.
		err = git.Add(worktree, updatedFiles)
		if err != nil {
			return fmt.Errorf("adding updated files to index: %v", err)
		}

		// Create a new commit including the updated files, with an appropriate commit message.
		err = git.Commit(worktree, commitMessage)
		if err != nil {
			return fmt.Errorf("committing updated project version files for [%s] project: %v", projectName, err)
		}

		if !upgradeOptions.DryRun {
			// Push the changes to the target branch in the head repository.
			err = git.Push(repo, headRepoOwner, headBranchName, githubToken)
			if err != nil {
				return fmt.Errorf("pushing updated project version files for [%s] project: %v", projectName, err)
			}

			// Update the title of the pull request depending on the base branch name.
			title := commitMessage
			if baseBranchName != constants.MainBranchName {
				title = fmt.Sprintf("[%s] %s", baseBranchName, title)
			}

			// Create a pull request from the branch in the head repository to the target branch in the aws/eks-anywhere-build-tooling repository.
			logger.Info("Creating pull request with updated files")
			pullRequest, err = github.CreatePullRequest(client, projectOrg, projectRepo, title, pullRequestBody, baseRepoOwner, baseBranchName, headRepoOwner, headBranchName, currentRevision, latestRevision)
			if err != nil {
				return fmt.Errorf("creating pull request to %s repository: %v", constants.BuildToolingRepoName, err)
			}
		} else {
			logger.Info(fmt.Sprintf("Completed dry run of upgrade for project %s", projectName))
		}
	}

	if len(failedSteps) > 0 {
		var failedStepsList []string
		var errorsList []string
		for step, err := range failedSteps {
			if step == "Patch application" {
				step = fmt.Sprintf("%s\n%s", step, patchesWarningComment)
			}
			failedStepsList = append(failedStepsList, fmt.Sprintf("* %s", step))
			errorsList = append(errorsList, fmt.Sprintf("Error occured in %s step: %v", step, err))
		}
		failedUpgradeComment := fmt.Sprintf(constants.FailedUpgradeCommentBody, strings.Join(failedStepsList, "\n"))

		if !upgradeOptions.DryRun {
			err = github.AddCommentOnPR(client, baseRepoOwner, failedUpgradeComment, pullRequest)
			if err != nil {
				return fmt.Errorf("commenting failed upgrade comment on pull request [%s]: %v", *pullRequest.HTMLURL, err)
			}
		}

		return errors.New(strings.Join(errorsList, "\n"))
	}

	return nil
}

func updateEKSDistroReleasesFile(buildToolingRepoPath string) (bool, error) {
	var isUpdated bool
	eksDistroReleasesFilepath := filepath.Join(buildToolingRepoPath, constants.EKSDistroLatestReleasesFile)

	eksDistroReleasesFileContents, err := os.ReadFile(eksDistroReleasesFilepath)
	if err != nil {
		return false, fmt.Errorf("reading EKS Distro latest releases file: %v", err)
	}

	var eksDistroLatestReleases types.EKSDistroLatestReleases
	err = yaml.Unmarshal(eksDistroReleasesFileContents, &eksDistroLatestReleases)
	if err != nil {
		return false, fmt.Errorf("unmarshalling EKS Distro latest releases file: %v", err)
	}

	supportedReleaseBranches, err := getSupportedReleaseBranches(buildToolingRepoPath)
	if err != nil {
		return false, fmt.Errorf("getting supported EKS Distro release branches: %v", err)
	}

	for i := range eksDistroLatestReleases.Releases {
		if slices.Contains(supportedReleaseBranches, eksDistroLatestReleases.Releases[i].Branch) {
			number, kubeVersion, err := getLatestEKSDistroRelease(eksDistroLatestReleases.Releases[i].Branch)
			if err != nil {
				return false, fmt.Errorf("getting latest EKS Distro release for %s branch: %v", eksDistroLatestReleases.Releases[i].Branch, err)
			}
			if eksDistroLatestReleases.Releases[i].Number != number || eksDistroLatestReleases.Releases[i].KubeVersion != kubeVersion {
				isUpdated = true
				eksDistroLatestReleases.Releases[i].Number = number
				eksDistroLatestReleases.Releases[i].KubeVersion = kubeVersion
			}
		}
	}

	if isUpdated {
		updatedEKSDistroReleasesFileContents, err := yaml.Marshal(eksDistroLatestReleases)
		if err != nil {
			return false, fmt.Errorf("marshalling EKS Distro latest releases: %v", err)
		}

		err = os.WriteFile(eksDistroReleasesFilepath, updatedEKSDistroReleasesFileContents, 0o644)
		if err != nil {
			return false, fmt.Errorf("writing EKS Distro latest releases file: %v", err)
		}
	}

	return isUpdated, nil
}

func updateEKSDistroBaseImageTagFiles(client *gogithub.Client, buildToolingRepoPath string, tagFileGlob []string) (string, bool, error) {
	var updatedPackages string
	var isUpdated bool

	eksDistroBaseTagYAMLContents, err := github.GetFileContents(client, constants.AWSOrgName, constants.EKSDistroBuildToolingRepoName, constants.EKSDistroBaseTagsYAMLFile, "main")
	if err != nil {
		return "", false, fmt.Errorf("getting contents of EKS Distro Base tag file: %v", err)
	}

	var eksDistroBaseTagYAMLMap map[string]interface{}
	err = yaml.Unmarshal(eksDistroBaseTagYAMLContents, &eksDistroBaseTagYAMLMap)
	if err != nil {
		return "", false, fmt.Errorf("unmarshalling EKS Distro Base tag file: %v", err)
	}

	for _, tagFile := range tagFileGlob {
		tagFileContents, err := file.ReadContentsTrimmed(tagFile)
		if err != nil {
			return "", false, fmt.Errorf("reading tag file: %v", err)
		}
		tagFileName := filepath.Base(tagFile)
		imageName := strings.TrimSuffix(tagFileName, constants.TagFileSuffix)

		tagFileKey := strings.ReplaceAll(strings.ToLower(imageName), "_", "-")
		osFolder := "2"
		osKey := "al2"
		if strings.HasSuffix(tagFileKey, "al2023") {
			tagFileKey = strings.TrimSuffix(tagFileKey, constants.AL2023Suffix)
			osFolder = "2023"
			osKey = "al2023"
		}
		tagReleaseDate := eksDistroBaseTagYAMLMap[osKey].(map[string]interface{})[tagFileKey].(string)

		if string(tagFileContents) != tagReleaseDate {
			isUpdated = true
			err = os.WriteFile(tagFile, []byte(fmt.Sprintf("%s\n", tagReleaseDate)), 0o644)
			if err != nil {
				return "", false, fmt.Errorf("writing tag file: %v", err)
			}

			updatedPackagesFilesContents, err := github.GetFileContents(client, constants.AWSOrgName, constants.EKSDistroBuildToolingRepoName, fmt.Sprintf(constants.EKSDistroBaseUpdatedPackagesFileFormat, osFolder, tagFileKey), "main")
			if err != nil {
				return "", false, fmt.Errorf("getting contents of EKS Distro Base image updated packages: %v", err)
			}
			updatedPackages = fmt.Sprintf("%s#### %s\nThe following yum packages were updated:\n```bash\n%s```\n\n", updatedPackages, imageName, string(updatedPackagesFilesContents))
		}
	}
	return updatedPackages, isUpdated, nil
}

func getSupportedReleaseBranches(buildToolingRepoPath string) ([]string, error) {
	supportedReleaseBranchesFilepath := filepath.Join(buildToolingRepoPath, constants.SupportedReleaseBranchesFile)

	supportedReleaseBranchesFileContents, err := os.ReadFile(supportedReleaseBranchesFilepath)
	if err != nil {
		return nil, fmt.Errorf("reading supported release branches file: %v", err)
	}
	supportedK8sVersions := strings.Split(strings.TrimRight(string(supportedReleaseBranchesFileContents), "\n"), "\n")

	return supportedK8sVersions, nil
}

func getLatestEKSDistroRelease(branch string) (int, string, error) {
	var eksDistroReleaseChannel eksdistrorelease.ReleaseChannel
	var eksDistroRelease eksdistrorelease.Release
	var kubeVersion string

	eksDistroReleaseChannelsFileURL := fmt.Sprintf(constants.EKSDistroReleaseChannelsFileURLFormat, branch)
	eksDistroReleaseChannelsFileContents, err := file.ReadURL(eksDistroReleaseChannelsFileURL)
	if err != nil {
		return 0, "", fmt.Errorf("reading EKS Distro ReleaseChannels file URL: %v", err)
	}

	err = sigsyaml.Unmarshal(eksDistroReleaseChannelsFileContents, &eksDistroReleaseChannel)
	if err != nil {
		return 0, "", fmt.Errorf("unmarshalling EKS Distro ReleaseChannels file: %v", err)
	}
	releaseNumber := eksDistroReleaseChannel.Status.LatestRelease

	eksDistroReleaseManifestURL := fmt.Sprintf(constants.EKSDistroReleaseManifestURLFormat, branch, releaseNumber)
	eksDistroReleaseManifestContents, err := file.ReadURL(eksDistroReleaseManifestURL)
	if err != nil {
		return 0, "", fmt.Errorf("reading EKS Distro release manifest URL: %v", err)
	}

	err = sigsyaml.Unmarshal(eksDistroReleaseManifestContents, &eksDistroRelease)
	if err != nil {
		return 0, "", fmt.Errorf("unmarshalling EKS Distro release manifest: %v", err)
	}
	for _, component := range eksDistroRelease.Status.Components {
		if component.Name == "kubernetes" {
			kubeVersion = component.GitTag
			break
		}
	}

	return releaseNumber, kubeVersion, nil
}

// updateProjectVersionFile updates the version information stored in a specific file.
func updateProjectVersionFile(buildToolingRepoPath, filename, projectName, value, releaseBranch string, isReleaseBranched bool) (string, error) {
	fileRelativepath := filepath.Join("projects", projectName, filename)
	if isReleaseBranched {
		fileRelativepath = filepath.Join("projects", projectName, releaseBranch, filename)
	}
	fileAbsolutepath := filepath.Join(buildToolingRepoPath, fileRelativepath)
	fileAbsolutePathStat, err := os.Stat(fileAbsolutepath)
	if err != nil {
		return "", fmt.Errorf("unable to stat project %s file [%s]: %v", filename, fileAbsolutepath, err)
	}
	if err := os.WriteFile(fileAbsolutepath, []byte(fmt.Sprintf("%s\n", value)), fileAbsolutePathStat.Mode()); err != nil {
		return "", fmt.Errorf("writing project %s file [%s]: %v", filename, fileAbsolutepath, err)
	}

	return fileRelativepath, nil
}

// loadUpstreamProjectsTrackerFile reads and unmarshals the contents of the upstream projects tracker file and
// returns the target repository object corresponding to the project being upgraded.
func loadUpstreamProjectsTrackerFile(upstreamProjectsTrackerFilePath, org, repository string) (types.ProjectsList, types.Repo, error) {
	contents, err := os.ReadFile(upstreamProjectsTrackerFilePath)
	if err != nil {
		return types.ProjectsList{}, types.Repo{}, fmt.Errorf("reading upstream projects tracker file: %v", err)
	}

	var projectsList types.ProjectsList
	err = goyamlv3.Unmarshal(contents, &projectsList)
	if err != nil {
		return types.ProjectsList{}, types.Repo{}, fmt.Errorf("unmarshalling upstream projects tracker file: %v", err)
	}

	var targetRepo types.Repo
	for _, project := range projectsList.Projects {
		if project.Org != org {
			continue
		}
		for _, repo := range project.Repos {
			if repo.Name != repository {
				continue
			}
			targetRepo = repo
		}
	}

	return projectsList, targetRepo, nil
}

// updateUpstreamProjectsTrackerFile updates the Git tag and Go version in the section of the upstream projects
// tracker file corresponding to the project being upgraded.
func updateUpstreamProjectsTrackerFile(projectsList *types.ProjectsList, buildToolingRepoPath, upstreamProjectsTrackerFilePath string) error {
	// Load the boilerplate license text that is used as comment header for the upstream projects tracker file.
	licenseBoilerplateFilepath := filepath.Join(buildToolingRepoPath, constants.LicenseBoilerplateFile)
	licenseBoilerplateContents, err := os.ReadFile(licenseBoilerplateFilepath)
	if err != nil {
		return fmt.Errorf("reading license boilerplate file: %v", err)
	}

	// Prefix the non-empty lines of the boilerplate text with a `#` to render it as a YAML comment.
	var b bytes.Buffer
	s := bufio.NewScanner(strings.NewReader(strings.ReplaceAll(string(licenseBoilerplateContents), "\\\"", "\"")))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Bytes()
		if len(line) == 0 {
			b.Write([]byte(fmt.Sprintf("%s\n", line)))
		} else {
			b.Write([]byte(fmt.Sprintf("# %s\n", line)))
		}
	}
	b.Write([]byte("\n"))

	// Create a new YAML encoder with an appropriate indentation value and encode the project list into a byte buffer
	yamlEncoder := goyamlv3.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	err = yamlEncoder.Encode(&projectsList)
	if err != nil {
		return fmt.Errorf("encoding the project list into a byte buffer: %v", err)
	}

	err = os.WriteFile(upstreamProjectsTrackerFilePath, b.Bytes(), 0o644)
	if err != nil {
		return fmt.Errorf("writing upstream projects tracker file: %v", err)
	}

	return nil
}

// applyPatchesToRepo runs a Make command to apply patches to the cloned repository of the project
// being upgraded.
func applyPatchesToRepo(projectRootFilepath, projectRepo string, totalPatchCount int) (int, string, string, error) {
	var patchesApplied int
	var failedPatch, failedFilesInPatch string
	patchApplySucceeded := true

	applyPatchesCommandSequence := fmt.Sprintf("make -C %s patch-repo", projectRootFilepath)
	applyPatchesCmd := exec.Command("bash", "-c", applyPatchesCommandSequence)
	applyPatchesOutput, err := command.ExecCommand(applyPatchesCmd)
	if err != nil {
		if strings.Contains(applyPatchesOutput, constants.FailedPatchApplyMarker) || strings.Contains(applyPatchesOutput, constants.DoesNotExistInIndexMarker) {
			patchApplySucceeded = false
		} else {
			return 0, "", "", fmt.Errorf("running patch-repo Make command: %v", err)
		}
	}

	if patchApplySucceeded {
		patchesApplied = totalPatchCount
	} else {
		failedFiles := []string{}
		gitDescribeRegex := regexp.MustCompile(constants.GitDescribeRegex)
		gitDescribeCmd := exec.Command("git", "-C", filepath.Join(projectRootFilepath, projectRepo), "describe", "--tag")
		gitDescribeOutput, err := command.ExecCommand(gitDescribeCmd)
		if err != nil {
			return 0, "", "", fmt.Errorf("running git describe command: %v", err)
		}
		gitDescribeMatches := gitDescribeRegex.FindStringSubmatch(gitDescribeOutput)
		if gitDescribeMatches[1] != "" {
			patchesApplied, err = strconv.Atoi(gitDescribeMatches[2])
			if err != nil {
				return 0, "", "", fmt.Errorf("converting patch count to integer %v", err)
			}
		}

		failedPatchRegex := regexp.MustCompile(constants.FailedPatchApplyRegex)
		failedPatch = failedPatchRegex.FindString(applyPatchesOutput)

		failedPatchFileRegex := regexp.MustCompile(fmt.Sprintf("%s|%s", constants.FailedPatchFilesRegex, constants.DoesNotExistInIndexFilesRegex))
		applyFailedFiles := failedPatchFileRegex.FindAllStringSubmatch(applyPatchesOutput, -1)
		for _, files := range applyFailedFiles {
			if files[1] != "" {
				failedFiles = append(failedFiles, fmt.Sprintf("`%s`", files[1]))
			} else if files[2] != "" {
				failedFiles = append(failedFiles, fmt.Sprintf("`%s`", files[2]))
			}
		}

		failedFilesInPatch = strings.Join(failedFiles, ",")
	}

	return patchesApplied, failedPatch, failedFilesInPatch, nil
}

// updateChecksumsAttributionFiles runs a Make command to update the checksums and attribution files
// corresponding to the project being upgraded.
func updateChecksumsAttributionFiles(projectRootFilepath string) error {
	updateChecksumsAttributionCommandSequence := fmt.Sprintf("make -C %s attribution-checksums", projectRootFilepath)
	updateChecksumsAttributionCmd := exec.Command("bash", "-c", updateChecksumsAttributionCommandSequence)
	_, err := command.ExecCommand(updateChecksumsAttributionCmd)
	if err != nil {
		return fmt.Errorf("running checksums-attribution Make command: %v", err)
	}

	return nil
}

// updateProjectReadmeVersion runs a script to update the version in the README file corresponding
// to the project being upgraded.
func updateProjectReadmeVersion(buildToolingRepoPath, projectOrg, projectRepo string) error {
	readmeUpdateScriptFilepath := filepath.Join(buildToolingRepoPath, constants.ReadmeUpdateScriptFile)
	readmeUpdateCommandSequence := fmt.Sprintf("%s %s %s/%s", readmeUpdateScriptFilepath, buildToolingRepoPath, projectOrg, projectRepo)
	readmeUpdateCmd := exec.Command("bash", "-c", readmeUpdateCommandSequence)
	_, err := command.ExecCommand(readmeUpdateCmd)
	if err != nil {
		return fmt.Errorf("running README update command: %v", err)
	}

	return nil
}

func updateCiliumImageDigestFiles(projectRootFilepath, projectPath string) ([]string, error) {
	updateCiliumFiles := []string{}
	updateDigestsCommandSequence := fmt.Sprintf("make -C %s update-digests", projectRootFilepath)
	updateDigestsCmd := exec.Command("bash", "-c", updateDigestsCommandSequence)
	_, err := command.ExecCommand(updateDigestsCmd)
	if err != nil {
		return nil, fmt.Errorf("running update-digests Make command: %v", err)
	}

	for _, directory := range constants.CiliumImageDirectories {
		updateCiliumFiles = append(updateCiliumFiles, filepath.Join(projectPath, "images", directory, "IMAGE_DIGEST"))
	}
	return updateCiliumFiles, nil
}

func updateADOTImageDigestFiles(projectRootFilepath, projectPath string) ([]string, error) {
	updateADOTFiles := []string{}
	updateDigestsCommandSequence := fmt.Sprintf("make -C %s update-digests", projectRootFilepath)
	updateDigestsCmd := exec.Command("bash", "-c", updateDigestsCommandSequence)
	_, err := command.ExecCommand(updateDigestsCmd)
	if err != nil {
		return nil, fmt.Errorf("running update-digests Make command: %v", err)
	}

	for _, directory := range constants.ADOTImageDirectories {
		updateADOTFiles = append(updateADOTFiles, filepath.Join(projectPath, "images", directory, "IMAGE_DIGEST"))
	}
	return updateADOTFiles, nil
}

func updateCertManagerManifestFile(projectRootFilepath, projectPath string) ([]string, error) {
	updateCertManagerManifestCommandSequence := fmt.Sprintf("make -C %s update-cert-manager-manifest", projectRootFilepath)
	updateCertManagerManifestCmd := exec.Command("bash", "-c", updateCertManagerManifestCommandSequence)
	_, err := command.ExecCommand(updateCertManagerManifestCmd)
	if err != nil {
		return nil, fmt.Errorf("running update-cert-manager-manifest Make command: %v", err)
	}

	return []string{filepath.Join(projectPath, constants.ManifestsDirectory, constants.CertManagerManifestYAMLFile)}, nil
}

func updateKindManifestImages(projectRootFilepath, projectPath string) ([]string, error) {
	updateKindManifestImagesCommandSequence := fmt.Sprintf("make -C %s update-manifest-images", projectRootFilepath)
	updateKindManifestImagesCmd := exec.Command("bash", "-c", updateKindManifestImagesCommandSequence)
	_, err := command.ExecCommand(updateKindManifestImagesCmd)
	if err != nil {
		return nil, fmt.Errorf("running update-manifest-images Make command: %v", err)
	}

	return []string{filepath.Join(projectPath, constants.BuildDirectory, constants.KindNodeImageBuildArgsScriptFile)}, nil
}

func updateBottlerocketVersionFiles(client *gogithub.Client, projectRootFilepath, projectPath, branchName string) (string, string, []string, error) {
	updatedBRFiles := []string{}
	var bottlerocketReleaseMap map[string]interface{}
	bottlerocketReleasesFilePath := filepath.Join(projectRootFilepath, constants.BottlerocketReleasesFile)
	bottlerocketReleasesRelativeFilePath := filepath.Join(projectPath, constants.BottlerocketReleasesFile)
	bottlerocketReleasesFileContents, err := os.ReadFile(bottlerocketReleasesFilePath)
	if err != nil {
		return "", "", nil, fmt.Errorf("reading Bottlerocket releases file: %v", err)
	}

	err = yaml.Unmarshal(bottlerocketReleasesFileContents, &bottlerocketReleaseMap)
	if err != nil {
		return "", "", nil, fmt.Errorf("unmarshalling Bottlerocket releases file: %v", err)
	}

	var currentBottlerocketVersion string
	for channel := range bottlerocketReleaseMap {
		for _, format := range constants.BottlerocketImageFormats {
			releaseVersionByFormat := bottlerocketReleaseMap[channel].(map[string]interface{})[fmt.Sprintf("%s-release-version", format)]
			if releaseVersionByFormat != nil {
				currentBottlerocketVersion = releaseVersionByFormat.(string)
				break
			}
		}
		if currentBottlerocketVersion != "" {
			break
		}
	}

	latestBottlerocketVersion, needsUpgrade, err := github.GetLatestRevision(client, "bottlerocket-os", "bottlerocket", currentBottlerocketVersion, branchName, false, false)
	if err != nil {
		return "", "", nil, fmt.Errorf("getting latest Bottlerocket version from GitHub: %v", err)
	}

	if needsUpgrade {
		logger.Info("Bottlerocket version is out of date.", "Current version", currentBottlerocketVersion, "Latest version", latestBottlerocketVersion)

		err = updateBottlerocketReleasesFile(bottlerocketReleaseMap, bottlerocketReleasesFilePath, latestBottlerocketVersion)
		if err != nil {
			return "", "", nil, fmt.Errorf("updating Bottlerocket releases file: %v", err)
		}
		updatedBRFiles = append(updatedBRFiles, bottlerocketReleasesRelativeFilePath)

		updatedHostContainerFiles, err := updateBottlerocketHostContainerMetadata(client, projectRootFilepath, projectPath, latestBottlerocketVersion)
		if err != nil {
			return "", "", nil, fmt.Errorf("updating Bottlerocket host containers metadata files: %v", err)
		}
		updatedBRFiles = append(updatedBRFiles, updatedHostContainerFiles...)
	}

	return currentBottlerocketVersion, latestBottlerocketVersion, updatedBRFiles, nil
}

func updateBottlerocketReleasesFile(bottlerocketReleaseMap map[string]interface{}, bottlerocketReleasesFilePath, latestBottlerocketVersion string) error {
	for channel := range bottlerocketReleaseMap {
		for _, format := range constants.BottlerocketImageFormats {
			releaseVersionByFormat := bottlerocketReleaseMap[channel].(map[string]interface{})[fmt.Sprintf("%s-release-version", format)]
			if releaseVersionByFormat != nil {
				imageExists, err := verifyBRImageExists(channel, format, latestBottlerocketVersion)
				if err != nil {
					return fmt.Errorf("checking if Bottlerocket %s image exists for %s release branch: %v", format, channel, err)
				}

				if imageExists {
					bottlerocketReleaseMap[channel].(map[string]interface{})[fmt.Sprintf("%s-release-version", format)] = latestBottlerocketVersion
				}
			}
		}
	}
	updatedBottlerocketReleases, err := yaml.Marshal(bottlerocketReleaseMap)
	if err != nil {
		return fmt.Errorf("marshalling Bottlerocket releases: %v", err)
	}

	err = os.WriteFile(bottlerocketReleasesFilePath, updatedBottlerocketReleases, 0o644)
	if err != nil {
		return fmt.Errorf("writing Bottlerocket releases file: %v", err)
	}

	return nil
}

func verifyBRImageExists(channel, format, bottlerocketVersion string) (bool, error) {
	kubeVersion := strings.ReplaceAll(channel, "-", ".")
	var variant, imageTarget string
	switch format {
	case "ami":
		variant = "aws"
		imageTarget = fmt.Sprintf(constants.BottlerocketAMIImageTargetFormat, variant, kubeVersion, bottlerocketVersion)
	case "ova":
		variant = "vmware"
		imageTarget = fmt.Sprintf(constants.BottlerocketOVAImageTargetFormat, variant, kubeVersion, bottlerocketVersion)
	case "raw":
		variant = "metal"
		imageTarget = fmt.Sprintf(constants.BottlerocketRawImageTargetFormat, variant, kubeVersion, bottlerocketVersion)
	}

	timestampURL := fmt.Sprintf(constants.BottlerocketTimestampJSONURLFormat, variant, kubeVersion)
	timestampManifest, err := file.ReadURL(timestampURL)
	if err != nil {
		return false, fmt.Errorf("reading Bottlerocket timestamp URL: %v", err)
	}

	var timestampData interface{}
	err = json.Unmarshal(timestampManifest, &timestampData)
	if err != nil {
		return false, fmt.Errorf("unmarshalling Bottlerocket timestamp manifest: %v", err)
	}

	version := timestampData.(map[string]interface{})["signed"].(map[string]interface{})["version"].(float64)
	versionString := fmt.Sprintf("%.0f", version)

	targetsURL := fmt.Sprintf(constants.BottlerocketTargetsJSONURLFormat, variant, kubeVersion, versionString)
	targetsManifest, err := file.ReadURL(targetsURL)
	if err != nil {
		return false, fmt.Errorf("reading Bottlerocket targets URL: %v", err)
	}

	var targetsData interface{}
	err = json.Unmarshal(targetsManifest, &targetsData)
	if err != nil {
		return false, fmt.Errorf("unmarshalling Bottlerocket targets manifest: %v", err)
	}

	targets := targetsData.(map[string]interface{})["signed"].(map[string]interface{})["targets"].(map[string]interface{})
	for target := range targets {
		if target == imageTarget {
			return true, nil
		}
	}

	return false, nil
}

func updateBottlerocketHostContainerMetadata(client *gogithub.Client, projectRootFilepath, projectPath, latestBottlerocketVersion string) ([]string, error) {
	updatedHostContainerFiles := []string{}
	hostContainersTOMLContents, err := github.GetFileContents(client, constants.BottlerocketOrgName, constants.BottlerocketRepoName, constants.BottlerocketHostContainersTOMLFile, latestBottlerocketVersion)
	if err != nil {
		return nil, fmt.Errorf("getting contents of Bottlerocket host containers file: %v", err)
	}

	var hostContainersTOMLMap interface{}
	err = toml.Unmarshal(hostContainersTOMLContents, &hostContainersTOMLMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling Bottlerocket host containers file: %v", err)
	}

	for _, container := range constants.BottlerocketHostContainers {
		var hostContainerImageMetadata types.ImageMetadata
		hostContainerMetadataFilePath := filepath.Join(projectRootFilepath, fmt.Sprintf(constants.BottlerocketContainerMetadataFileFormat, strings.ToUpper(container)))
		hostContainerMetadataRelativeFilePath := filepath.Join(projectPath, fmt.Sprintf(constants.BottlerocketContainerMetadataFileFormat, strings.ToUpper(container)))
		hostContainerMetadataFileContents, err := os.ReadFile(hostContainerMetadataFilePath)
		if err != nil {
			return nil, fmt.Errorf("reading Bottlerocket %s container metadata file: %v", container, err)
		}
		err = yaml.Unmarshal(hostContainerMetadataFileContents, &hostContainerImageMetadata)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling Bottlerocket %s container metadata file: %v", container, err)
		}

		hostContainerSourceCommand := hostContainersTOMLMap.(map[string]interface{})["metadata"].(map[string]interface{})["settings"].(map[string]interface{})["host-containers"].(map[string]interface{})[container].(map[string]interface{})["source"].(map[string]interface{})["setting-generator"].(map[string]interface{})["command"].(string)
		hostContainerSourceImageRegex := regexp.MustCompile(constants.BottlerocketHostContainerSourceImageRegex)
		hostContainerSourceImage := hostContainerSourceImageRegex.FindAllStringSubmatch(hostContainerSourceCommand, -1)[0][1]
		hostContainerSourceImageTag := strings.Split(hostContainerSourceImage, ":")[1]

		if hostContainerImageMetadata.Tag != hostContainerSourceImageTag {
			hostContainerImageMetadata.Tag = hostContainerSourceImageTag
			skopeoInspectCmd := exec.Command("skopeo", "inspect", fmt.Sprintf("docker://%s", hostContainerSourceImage), "--override-os", "linux", "--override-arch", "amd64", "--format", "{{.Digest}}")
			stdout, err := command.ExecCommand(skopeoInspectCmd)
			if err != nil {
				return nil, fmt.Errorf("running skopeo inspect command: %v", err)
			}
			hostContainerImageMetadata.ImageDigest = stdout

			updatedHostContainerMetadataFileContents, err := yaml.Marshal(hostContainerImageMetadata)
			if err != nil {
				return nil, fmt.Errorf("marshalling updated Bottlerocket %s container: %v", container, err)
			}

			err = os.WriteFile(hostContainerMetadataFilePath, updatedHostContainerMetadataFileContents, 0o644)
			if err != nil {
				return nil, fmt.Errorf("writing Bottlerocket releases file: %v", err)
			}

			updatedHostContainerFiles = append(updatedHostContainerFiles, hostContainerMetadataRelativeFilePath)
		}
	}

	return updatedHostContainerFiles, nil
}

func isEKSDistroUpgrade(projectName string) bool {
	return projectName == "aws/eks-distro"
}

func isEKSDistroBuildToolingUpgrade(projectName string) bool {
	return projectName == "aws/eks-distro-build-tooling"
}

func getDefaultReleaseBranch(buildToolingRepoPath string) (string, error) {
	defaultReleaseBranchCommandSequence := fmt.Sprintf("make --no-print-directory -C %s get-default-release-branch", buildToolingRepoPath)
	defaultReleaseBranchCmd := exec.Command("bash", "-c", defaultReleaseBranchCommandSequence)
	defaultReleaseBranch, err := command.ExecCommand(defaultReleaseBranchCmd)
	if err != nil {
		return "", fmt.Errorf("running get-default-release-branch Make command: %v", err)
	}

	return defaultReleaseBranch, nil
}

func getProjectSpecificUpdateOperation(projectName string) (func(projectRootFilepath, projectPath string) ([]string, error), string) {
	switch projectName {
	case "aws-observability/aws-otel-collector":
		return updateADOTImageDigestFiles, "ADOT image digest update"
	case "cert-manager/cert-manager":
		return updateCertManagerManifestFile, "Cert-manager manifest file update"
	case "cilium/cilium":
		return updateCiliumImageDigestFiles, "Cilium image digest update"
	case "kubernetes-sigs/kind":
		return updateKindManifestImages, "Kind manifest images update"
	default:
		return nil, ""
	}
}

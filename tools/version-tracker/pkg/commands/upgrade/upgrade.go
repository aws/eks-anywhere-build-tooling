package upgrade

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/eks-anywhere/pkg/semver"
	gogithub "github.com/google/go-github/v53/github"
	"gopkg.in/yaml.v3"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/git"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/github"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/command"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/slices"
)

// Run contains the business logic to execute the `upgrade` subcommand.
func Run(upgradeOptions *types.UpgradeOptions) error {
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
	if slices.Contains(skippedProjects, upgradeOptions.ProjectName) {
		logger.Info("Project is in SKIPPED_PROJECTS list. Skipping upgrade")
		return nil
	}

	// Clone the eks-anywhere-build-tooling repository.
	buildToolingRepoPath := filepath.Join(cwd, constants.BuildToolingRepoName)
	repo, headCommit, err := git.CloneRepo(constants.BuildToolingRepoURL, buildToolingRepoPath, headRepoOwner)
	if err != nil {
		return fmt.Errorf("cloning build-tooling repo: %v", err)
	}

	// Get the worktree corresponding to the cloned repository.
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("getting repo's current worktree: %v", err)
	}

	// Validate if the project name provided exists in the repository.
	projectPath := filepath.Join("projects", upgradeOptions.ProjectName)
	projectRootFilepath := filepath.Join(buildToolingRepoPath, projectPath)
	if _, err := os.Stat(projectRootFilepath); os.IsNotExist(err) {
		return fmt.Errorf("invalid project name %s", upgradeOptions.ProjectName)
	}

	// Check if project to be upgraded has patches
	projectHasPatches := false
	if _, err := os.Stat(filepath.Join(projectRootFilepath, constants.PatchesDirectory)); err == nil {
		projectHasPatches = true
	}

	// Get org and repository name from project name.
	projectOrg := strings.Split(upgradeOptions.ProjectName, "/")[0]
	projectRepo := strings.Split(upgradeOptions.ProjectName, "/")[1]

	// Load upstream projects tracker file.
	upstreamProjectsTrackerFilePath := filepath.Join(buildToolingRepoPath, constants.UpstreamProjectsTrackerFile)
	_, targetRepo, err := loadUpstreamProjectsTrackerFile(upstreamProjectsTrackerFilePath, projectOrg, projectRepo)
	if err != nil {
		return fmt.Errorf("loading upstream projects tracker file: %v", err)
	}

	// Validate whether the given project is release-branched.
	if len(targetRepo.Versions) > 1 {
		return fmt.Errorf("release-branched projects not supported at this time")
	}

	// Get latest revision for the project from GitHub.
	latestRevision, latestRevisionCommit, allTags, err := github.GetLatestRevision(client, projectOrg, projectRepo)
	if err != nil {
		return fmt.Errorf("getting latest revision from GitHub: %v", err)
	}

	currentVersion := targetRepo.Versions[0]
	// Validate whether the project builds off a commit hash instead of a tag.
	if currentVersion.Tag == "" {
		return fmt.Errorf("projects tracked with commit hashes not supported at this time")
	}
	currentRevision := currentVersion.Tag

	// Get commit hash corresponding to current revision tag.
	currentRevisionCommit := github.GetCommitForTag(allTags, currentRevision)

	// Get Unix timestamp for current revision's commit
	currentRevisionCommitEpoch, err := github.GetCommitDateEpoch(client, projectOrg, projectRepo, currentRevisionCommit)
	if err != nil {
		return fmt.Errorf("getting epoch time corresponding to current revision commit: %v", err)
	}

	// Get Unix timestamp for latest revision's commit
	latestRevisionCommitEpoch, err := github.GetCommitDateEpoch(client, projectOrg, projectRepo, latestRevisionCommit)
	if err != nil {
		return fmt.Errorf("getting epoch time corresponding to latest revision commit: %v", err)
	}

	currentRevisionSemver, err := semver.New(currentRevision)
	if err != nil {
		return fmt.Errorf("getting semver for current version: %v", err)
	}

	latestRevisionSemver, err := semver.New(latestRevision)
	if err != nil {
		return fmt.Errorf("getting semver for latest version: %v", err)
	}

	// Upgrade project if latest commit was made after current commit.
	if latestRevisionCommitEpoch > currentRevisionCommitEpoch || latestRevisionSemver.GreaterThan(currentRevisionSemver) {
		logger.Info("Project is out of date.", "Current version", currentRevision, "Latest version", latestRevision)

		var updatedFiles []string
		headBranchName := fmt.Sprintf("update-%s-%s", projectOrg, projectRepo)
		baseBranchName := constants.MainBranchName

		// Checkout a new branch to keep track of version upgrade chaneges.
		err = git.Checkout(worktree, headBranchName)
		if err != nil {
			return fmt.Errorf("getting repo's current worktree: %v", err)
		}

		// Reset current worktree to get a clean index.
		err = git.ResetToMain(worktree, headCommit)
		if err != nil {
			return fmt.Errorf("resetting new branch to [origin/main] HEAD: %v", err)
		}

		// Reload upstream projects tracker file to get its original value instead of
		// the updated one from another project's previous upgrade
		projectsList, targetRepo, err := loadUpstreamProjectsTrackerFile(upstreamProjectsTrackerFilePath, projectOrg, projectRepo)
		if err != nil {
			return fmt.Errorf("reloading upstream projects tracker file: %v", err)
		}
		targetRepo.Versions[0].Tag = latestRevision

		// Update the Git tag file corresponding to the project
		logger.Info("Updating Git tag file corresponding to the project")
		projectGitTagRelativePath, err := updateProjectVersionFile(buildToolingRepoPath, constants.GitTagFile, upgradeOptions.ProjectName, latestRevision)
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
				targetRepo.Versions[0].GoVersion = latestGoVersion

				logger.Info("Updating Go version file corresponding to the project")
				projectGoVersionRelativePath, err := updateProjectVersionFile(buildToolingRepoPath, constants.GoVersionFile, upgradeOptions.ProjectName, latestGoVersion)
				if err != nil {
					return fmt.Errorf("updating project GOLANG_VERSION file: %v", err)
				}
				updatedFiles = append(updatedFiles, projectGoVersionRelativePath)
			}
		} else {
			latestGoVersion = "N/A"
			targetRepo.Versions[0].GoVersion = latestGoVersion
		}

		// Update the tag and Go version in the section of the upstream projects tracker file corresponding to the given project.
		logger.Info("Updating Git tag and Go version in upstream projects tracker file")
		err = updateUpstreamProjectsTrackerFile(&projectsList, targetRepo, buildToolingRepoPath, upstreamProjectsTrackerFilePath, latestRevision, latestGoVersion)
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

		// Update the checksums file and attribution file(s) corresponding to the project.
		if !projectHasPatches {
			if _, err := os.Stat(filepath.Join(projectRootFilepath, constants.ChecksumsFile)); err == nil {
				logger.Info("Updating project checksums and attribution files")
				projectChecksumsFileRelativePath := filepath.Join(projectPath, constants.ChecksumsFile)
				err = updateChecksumsAttributionFiles(projectRootFilepath)
				if err != nil {
					return fmt.Errorf("updating project checksums and attribution files: %v", err)
				}
				updatedFiles = append(updatedFiles, projectChecksumsFileRelativePath)

				// Attribution files can have a binary name prefix so we use a common prefix regular expression
				// and glob them to cover all possibilities.
				projectAttributionFileGlob, err := filepath.Glob(filepath.Join(projectRootFilepath, constants.AttributionsFilePattern))
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

		// Add all the updated files to the index.
		err = git.Add(worktree, updatedFiles)
		if err != nil {
			return fmt.Errorf("adding updated files to index: %v", err)
		}

		// Create a new commit including the updated files, with an appropriate commit message.
		err = git.Commit(worktree, fmt.Sprintf("Bump %s to latest release", upgradeOptions.ProjectName))
		if err != nil {
			return fmt.Errorf("committing updated project version files for [%s] project: %v", upgradeOptions.ProjectName, err)
		}

		if upgradeOptions.DryRun {
			logger.Info(fmt.Sprintf("Completed dry run of upgrade for project %s", upgradeOptions.ProjectName))
			return nil
		}

		// Push the changes to the target branch in the head repository.
		err = git.Push(repo, headRepoOwner, headBranchName, githubToken)
		if err != nil {
			return fmt.Errorf("pushing updated project version files for [%s] project: %v", upgradeOptions.ProjectName, err)
		}

		// Create a pull request from the bramch in the head repository to the target branch in the aws/eks-anywhere-build-tooling repository.
		logger.Info("Creating pull request with updated files")
		err = github.CreatePullRequest(client, projectOrg, projectRepo, baseRepoOwner, baseBranchName, headRepoOwner, headBranchName, currentRevision, latestRevision, projectHasPatches)
		if err != nil {
			return fmt.Errorf("creating pull request to %s repository: %v", constants.BuildToolingRepoName, err)
		}
	} else if latestRevision == currentRevision {
		logger.Info("Project is at the latest available version.", "Current version", currentRevision, "Latest version", latestRevision)
	}

	return nil
}

// updateProjectVersionFile updates the version information stored in a specific file.
func updateProjectVersionFile(buildToolingRepoPath, filename, projectName, value string) (string, error) {
	fileRelativepath := filepath.Join("projects", projectName, filename)
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
	err = yaml.Unmarshal(contents, &projectsList)
	if err != nil {
		return types.ProjectsList{}, types.Repo{}, fmt.Errorf("unmarshaling upstream projects tracker file to YAML: %v", err)
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
func updateUpstreamProjectsTrackerFile(projectsList *types.ProjectsList, targetRepo types.Repo, buildToolingRepoPath, upstreamProjectsTrackerFilePath, latestRevision, latestGoVersion string) error {
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

	// Create a new YAML encoder with an appropriate indentation value and encode the project list into a byte buufer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(&projectsList)

	err = os.WriteFile(upstreamProjectsTrackerFilePath, b.Bytes(), 0o644)
	if err != nil {
		return fmt.Errorf("writing upstream projects tracker file: %v", err)
	}

	return nil
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

// updateChecksumsAttributionFiles runs a script to update the version in the README file corresponding
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

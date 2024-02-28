package upgrade

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	gogithub "github.com/google/go-github/v53/github"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/git"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/github"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/command"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/file"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/slices"
)

// Run contains the business logic to execute the `upgrade` subcommand.
func Run(upgradeOptions *types.UpgradeOptions) error {
	projectName := upgradeOptions.ProjectName

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
	if slices.Contains(skippedProjects, projectName) {
		logger.Info("Project is in SKIPPED_PROJECTS list. Skipping upgrade")
		return nil
	}

	// Clone the eks-anywhere-build-tooling repository.
	buildToolingRepoPath := filepath.Join(cwd, constants.BuildToolingRepoName)
	repo, headCommit, err := git.CloneRepo(fmt.Sprintf(constants.BuildToolingRepoURL, baseRepoOwner), buildToolingRepoPath, headRepoOwner)
	if err != nil {
		return fmt.Errorf("cloning build-tooling repo: %v", err)
	}

	// Get the worktree corresponding to the cloned repository.
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("getting repo's current worktree: %v", err)
	}

	// Validate if the project name provided exists in the repository.
	projectPath := filepath.Join("projects", projectName)
	projectRootFilepath := filepath.Join(buildToolingRepoPath, projectPath)
	if _, err := os.Stat(projectRootFilepath); os.IsNotExist(err) {
		return fmt.Errorf("invalid project name %s", projectName)
	}

	// Check if project to be upgraded has patches
	projectHasPatches := false
	if _, err := os.Stat(filepath.Join(projectRootFilepath, constants.PatchesDirectory)); err == nil {
		projectHasPatches = true
	}

	// Get org and repository name from project name.
	projectOrg := strings.Split(projectName, "/")[0]
	projectRepo := strings.Split(projectName, "/")[1]

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

	currentVersion := targetRepo.Versions[0]
	// Validate whether the project builds off a commit hash instead of a tag.
	if currentVersion.Tag == "" {
		return fmt.Errorf("projects tracked with commit hashes not supported at this time")
	}
	currentRevision := currentVersion.Tag

	// Get latest revision for the project from GitHub.
	latestRevision, needsUpgrade, err := github.GetLatestRevision(client, projectOrg, projectRepo, currentRevision)
	if err != nil {
		return fmt.Errorf("getting latest revision from GitHub: %v", err)
	}

	// Upgrade project if latest commit was made after current commit and the semver of the latest revision is
	// greater than the semver of the current version.
	if needsUpgrade || slices.Contains(constants.ProjectsWithUnconventionalUpgradeFlows, projectName) {
		var updatedFiles []string
		headBranchName := fmt.Sprintf("update-%s-%s", projectOrg, projectRepo)
		baseBranchName := constants.MainBranchName
		commitMessage := fmt.Sprintf("Bump %s to latest release", projectName)

		// Checkout a new branch to keep track of version upgrade chaneges.
		err = git.Checkout(worktree, headBranchName)
		if err != nil {
			return fmt.Errorf("checking out worktree at branch %s: %v", headBranchName, err)
		}

		// Reset current worktree to get a clean index.
		err = git.ResetToMain(worktree, headCommit)
		if err != nil {
			return fmt.Errorf("resetting new branch to [origin/main] HEAD: %v", err)
		}

		if needsUpgrade {
			logger.Info("Project is out of date.", "Current version", currentRevision, "Latest version", latestRevision)

			// Reload upstream projects tracker file to get its original value instead of
			// the updated one from another project's previous upgrade
			projectsList, targetRepo, err := loadUpstreamProjectsTrackerFile(upstreamProjectsTrackerFilePath, projectOrg, projectRepo)
			if err != nil {
				return fmt.Errorf("reloading upstream projects tracker file: %v", err)
			}
			targetRepo.Versions[0].Tag = latestRevision

			// Update the Git tag file corresponding to the project
			logger.Info("Updating Git tag file corresponding to the project")
			projectGitTagRelativePath, err := updateProjectVersionFile(buildToolingRepoPath, constants.GitTagFile, projectName, latestRevision)
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
					projectGoVersionRelativePath, err := updateProjectVersionFile(buildToolingRepoPath, constants.GoVersionFile, projectName, latestGoVersion)
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
		}

		if projectName == "kubernetes-sigs/image-builder" {
			updatedBRFiles, err := updateBottlerocketVersionFiles(client, projectRootFilepath, projectPath)
			if err != nil {
				return fmt.Errorf("updating Bottlerocket version and metadata files: %v", err)
			}
			if len(updatedBRFiles) > 0 {
				updatedFiles = append(updatedFiles, updatedBRFiles...)
				if len(updatedFiles) == len(updatedBRFiles) {
					headBranchName = "update-bottlerocket-releases"
					commitMessage = "Bump Bottlerocket versions to latest release"
				} else {
					headBranchName = fmt.Sprintf("update-%s-%s-and-bottlerocket", projectOrg, projectRepo)
					commitMessage = fmt.Sprintf("Bump %s and Bottlerocket versions to latest release", projectName)
				}

				err = git.Checkout(worktree, headBranchName)
				if err != nil {
					return fmt.Errorf("checking out worktree at branch %s: %v", headBranchName, err)
				}
			}
		}

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

		if upgradeOptions.DryRun {
			logger.Info(fmt.Sprintf("Completed dry run of upgrade for project %s", projectName))
			return nil
		}

		// Push the changes to the target branch in the head repository.
		err = git.Push(repo, headRepoOwner, headBranchName, githubToken)
		if err != nil {
			return fmt.Errorf("pushing updated project version files for [%s] project: %v", projectName, err)
		}

		// Create a pull request from the bramch in the head repository to the target branch in the aws/eks-anywhere-build-tooling repository.
		logger.Info("Creating pull request with updated files")
		err = github.CreatePullRequest(client, projectOrg, projectRepo, commitMessage, baseRepoOwner, baseBranchName, headRepoOwner, headBranchName, currentRevision, latestRevision, projectHasPatches)
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
		return types.ProjectsList{}, types.Repo{}, fmt.Errorf("unmarshaling upstream projects tracker file: %v", err)
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

func updateBottlerocketVersionFiles(client *gogithub.Client, projectRootFilepath, projectPath string) ([]string, error) {
	updatedBRFiles := []string{}
	var bottlerocketReleaseMap map[string]interface{}
	bottlerocketReleasesFilePath := filepath.Join(projectRootFilepath, constants.BottlerocketReleasesFile)
	bottlerocketReleasesRelativeFilePath := filepath.Join(projectPath, constants.BottlerocketReleasesFile)
	bottlerocketReleasesFileContents, err := os.ReadFile(bottlerocketReleasesFilePath)
	if err != nil {
		return nil, fmt.Errorf("reading Bottlerocket releases file: %v", err)
	}

	err = yaml.Unmarshal(bottlerocketReleasesFileContents, &bottlerocketReleaseMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling Bottlerocket releases file: %v", err)
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

	latestBottlerocketVersion, needsUpgrade, err := github.GetLatestRevision(client, "bottlerocket-os", "bottlerocket", currentBottlerocketVersion)
	if err != nil {
		return nil, fmt.Errorf("getting latest Bottlerocket version from GitHub: %v", err)
	}

	if needsUpgrade {
		logger.Info("Bottlerocket version is out of date.", "Current version", currentBottlerocketVersion, "Latest version", latestBottlerocketVersion)

		err = updateBottlerocketReleasesFile(bottlerocketReleaseMap, bottlerocketReleasesFilePath, latestBottlerocketVersion)
		if err != nil {
			return nil, fmt.Errorf("updating Bottlerocket releases file: %v", err)
		}
		updatedBRFiles = append(updatedBRFiles, bottlerocketReleasesRelativeFilePath)

		updatedHostContainerFiles, err := updateBottlerocketHostContainerMetadata(client, projectRootFilepath, projectPath, latestBottlerocketVersion)
		if err != nil {
			return nil, fmt.Errorf("updating Bottlerocket host containers metadata files: %v", err)
		}
		updatedBRFiles = append(updatedBRFiles, updatedHostContainerFiles...)
	}

	return updatedBRFiles, nil
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
		return fmt.Errorf("marshaling Bottlerocket releases file: %v", err)
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
		imageTarget = fmt.Sprintf("bottlerocket-%s-k8s-%s-x86_64-%s.img.lz4", variant, kubeVersion, bottlerocketVersion)
	case "ova":
		variant = "vmware"
		imageTarget = fmt.Sprintf("bottlerocket-%s-k8s-%s-x86_64-%s.ova", variant, kubeVersion, bottlerocketVersion)
	case "raw":
		variant = "metal"
		imageTarget = fmt.Sprintf("bottlerocket-%s-k8s-%s-x86_64-%s.img.lz4", variant, kubeVersion, bottlerocketVersion)
	}

	timestampURL := fmt.Sprintf("https://updates.bottlerocket.aws/2020-07-07/%s-k8s-%s/x86_64/timestamp.json", variant, kubeVersion)
	timestampManifest, err := file.ReadURL(timestampURL)
	if err != nil {
		return false, fmt.Errorf("reading Bottlerocket timestamp URL: %v", err)
	}

	var timestampData interface{}
	err = json.Unmarshal(timestampManifest, &timestampData)
	if err != nil {
		return false, fmt.Errorf("unmarshaling Bottlerocket timestamp manifest: %v", err)
	}

	version := timestampData.(map[string]interface{})["signed"].(map[string]interface{})["version"].(float64)
	versionString := fmt.Sprintf("%.0f", version)

	targetsURL := fmt.Sprintf("https://updates.bottlerocket.aws/2020-07-07/%s-k8s-%s/x86_64/%s.targets.json", variant, kubeVersion, versionString)
	targetsManifest, err := file.ReadURL(targetsURL)
	if err != nil {
		return false, fmt.Errorf("reading Bottlerocket targets URL: %v", err)
	}

	var targetsData interface{}
	err = json.Unmarshal(targetsManifest, &targetsData)
	if err != nil {
		return false, fmt.Errorf("unmarshaling Bottlerocket targets manifest: %v", err)
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
	hostContainersTOMLFile := "sources/models/shared-defaults/public-host-containers.toml"
	hostContainersTOMLContents, _, _, err := client.Repositories.GetContents(context.Background(), "bottlerocket-os", "bottlerocket", hostContainersTOMLFile, &gogithub.RepositoryContentGetOptions{Ref: latestBottlerocketVersion})
	if err != nil {
		return nil, fmt.Errorf("getting contents of file [%s]: %v", hostContainersTOMLFile, err)
	}
	hostContainersTOMLContentsDecoded, err := base64.StdEncoding.DecodeString(*hostContainersTOMLContents.Content)
	if err != nil {
		return nil, fmt.Errorf("decoding contents of file [%s]: %v", hostContainersTOMLFile, err)
	}

	var hostContainersTOMLMap interface{}
	err = toml.Unmarshal(hostContainersTOMLContentsDecoded, &hostContainersTOMLMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling Bottlerocket host containers file: %v", err)
	}

	for _, container := range constants.BottlerocketHostContainers {
		var hostContainerImageMetadata types.ImageMetadata
		hostContainerMetadataFilePath := filepath.Join(projectRootFilepath, fmt.Sprintf("BOTTLEROCKET_%s_CONTAINER_METADATA", strings.ToUpper(container)))
		hostContainerMetadataRelativeFilePath := filepath.Join(projectPath, fmt.Sprintf("BOTTLEROCKET_%s_CONTAINER_METADATA", strings.ToUpper(container)))
		hostContainerMetadataFileContents, err := os.ReadFile(hostContainerMetadataFilePath)
		if err != nil {
			return nil, fmt.Errorf("reading Bottlerocket %s container metadata file: %v", container, err)
		}
		err = yaml.Unmarshal(hostContainerMetadataFileContents, &hostContainerImageMetadata)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling Bottlerocket %s container metadata file: %v", container, err)
		}

		hostContainerSourceImage := hostContainersTOMLMap.(map[string]interface{})["settings"].(map[string]interface{})["host-containers"].(map[string]interface{})[container].(map[string]interface{})["source"].(string)
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
				return nil, fmt.Errorf("marshaling updated Bottlerocket %s container file: %v", container, err)
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

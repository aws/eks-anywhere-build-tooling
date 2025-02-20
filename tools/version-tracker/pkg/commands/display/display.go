package display

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	gogithub "github.com/google/go-github/v53/github"
	"github.com/rodaine/table"
	"gopkg.in/yaml.v3"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/ecrpublic"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/git"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/github"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
)

// Run contains the business logic to execute the `display` subcommand.
func Run(displayOptions *types.DisplayOptions) error {
	// Check if branch name environment variable has been set.
	branchName, ok := os.LookupEnv(constants.BranchNameEnvVar)
	if !ok {
		branchName = constants.MainBranchName
	}

	// Check if GitHub token environment variable has been set.
	githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}
	client := gogithub.NewTokenClient(context.Background(), githubToken)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("retrieving current working directory: %v", err)
	}

	// Get base repository owner environment variable if set.
	baseRepoOwner := os.Getenv(constants.BaseRepoOwnerEnvvar)
	if baseRepoOwner == "" {
		baseRepoOwner = constants.AWSOrgName
	}

	// Clone the eks-anywhere-build-tooling repository.
	buildToolingRepoPath := filepath.Join(cwd, constants.BuildToolingRepoName)
	repo, headCommit, err := git.CloneRepo(fmt.Sprintf(constants.BuildToolingRepoURL, baseRepoOwner), buildToolingRepoPath, "", branchName)
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

	if displayOptions.ProjectName != "" {
		// Validate if the project name provided exists in the repository.
		if _, err := os.Stat(filepath.Join(buildToolingRepoPath, "projects", displayOptions.ProjectName)); os.IsNotExist(err) {
			return fmt.Errorf("invalid project name %s", displayOptions.ProjectName)
		}
	}

	// Load upstream projects tracker file.
	upstreamProjectsTrackerFilePath := filepath.Join(buildToolingRepoPath, constants.UpstreamProjectsTrackerFile)
	contents, err := os.ReadFile(upstreamProjectsTrackerFilePath)
	if err != nil {
		return fmt.Errorf("reading upstream projects tracker file: %v", err)
	}

	// Unmarshal upstream projects tracker file
	var projectsList types.ProjectsList
	err = yaml.Unmarshal(contents, &projectsList)
	if err != nil {
		return fmt.Errorf("unmarshalling upstream projects tracker file: %v", err)
	}

	var projectVersionInfoList []types.ProjectVersionInfo
	for _, project := range projectsList.Projects {
		org := project.Org
		for _, repo := range project.Repos {
			var currentVersionList, latestVersionList, upToDateList []string
			repoName := repo.Name
			fullRepoName := fmt.Sprintf("%s/%s", org, repoName)
			if displayOptions.ProjectName != "" && displayOptions.ProjectName != fullRepoName {
				continue
			}

			releaseBranched := false
			var currentVersion types.Version
			if len(repo.Versions) > 1 {
				releaseBranched = true
			}

			supportedReleaseBranches, err := getSupportedReleaseBranches(buildToolingRepoPath)
			if err != nil {
				return fmt.Errorf("getting supported EKS Distro release branches: %v", err)
			}

			for _, releaseBranch := range supportedReleaseBranches {
				err := os.Setenv(constants.ReleaseBranchEnvvar, releaseBranch)
				releaseBranchIndex := slices.Index(supportedReleaseBranches, releaseBranch)
				if !releaseBranched && releaseBranchIndex > 0 {
					break
				}
				currentVersion = repo.Versions[releaseBranchIndex]

				var isTrackedByCommitHash bool
				var currentRevision string
				if currentVersion.Tag != "" {
					currentRevision = currentVersion.Tag
				} else if currentVersion.Commit != "" {
					currentRevision = currentVersion.Commit
					isTrackedByCommitHash = true
				}
				currentVersionList = append(currentVersionList, currentRevision)

				var latestRevision string
				if fullRepoName == "cilium/cilium" || fullRepoName == "envoyproxy/envoy" {
					latestRevision, _, err = ecrpublic.GetLatestRevision(constants.ECRImageRepositories[fullRepoName], currentRevision, branchName)
					if err != nil {
						return fmt.Errorf("getting latest revision from ECR Public: %v", err)
					}
				} else {
					// Get latest revision for the project from GitHub.
					latestRevision, _, err = github.GetLatestRevision(client, org, repoName, currentRevision, branchName, isTrackedByCommitHash, releaseBranched)
					if err != nil {
						return fmt.Errorf("getting latest revision from GitHub for project %s: %v", fullRepoName, err)
					}
				}

				latestVersionList = append(latestVersionList, latestRevision)
				upToDateList = append(upToDateList, fmt.Sprintf("%t", currentRevision == latestRevision))
			}

			var paddedOrgName, paddedRepoName string
			if len(currentVersionList) > 2 {
				padding := strings.Repeat("\n", len(currentVersionList)/2)
				smallerPadding := strings.Repeat("\n", len(currentVersionList)/2-1)
				if len(currentVersionList)%2 == 0 {
					paddedOrgName = fmt.Sprintf("%s%s%s", padding, org, smallerPadding)
					paddedRepoName = fmt.Sprintf("%s%s%s", padding, repoName, smallerPadding)
				} else {
					paddedOrgName = fmt.Sprintf("%s%s%s", padding, org, padding)
					paddedRepoName = fmt.Sprintf("%s%s%s", padding, repoName, padding)
				}
			} else {
				paddedOrgName = org
				paddedRepoName = repoName
			}

			projectVersionInfoList = append(projectVersionInfoList, types.ProjectVersionInfo{Org: paddedOrgName, Repo: paddedRepoName, CurrentVersion: strings.Join(currentVersionList, "\n"), LatestVersion: strings.Join(latestVersionList, "\n"), UpToDate: strings.Join(upToDateList, "\n")})
		}
	}

	// Create a new table with the required column names in uppercase.
	tbl := table.New("Organization", "Repository", "Current Version", "Latest Version", "Up-To-Date").WithHeaderFormatter(func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	})

	// Add rows to the table for each project in the list.
	for _, versionInfo := range projectVersionInfoList {
		tbl.AddRow(versionInfo.Org, versionInfo.Repo, versionInfo.CurrentVersion, versionInfo.LatestVersion, versionInfo.UpToDate)
	}

	// Print the table contents to standard output.
	tbl.Print()

	return nil
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

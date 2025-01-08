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
	projectName := displayOptions.ProjectName

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
			repoName := repo.Name
			fullRepoName := fmt.Sprintf("%s/%s", org, repoName)
			if displayOptions.ProjectName != "" && displayOptions.ProjectName != fullRepoName {
				continue
			}
			var releaseBranched bool
			var currentVersion types.Version
			if len(repo.Versions) > 1 {
				releaseBranched = true
			}
			if releaseBranched {
				supportedReleaseBranches, err := getSupportedReleaseBranches(buildToolingRepoPath)
				if err != nil {
					return fmt.Errorf("getting supported EKS Distro release branches: %v", err)
				}
				releaseBranch := os.Getenv(constants.ReleaseBranchEnvvar)
				releaseBranchIndex := slices.Index(supportedReleaseBranches, releaseBranch)
				currentVersion = repo.Versions[releaseBranchIndex]
			} else {
				currentVersion = repo.Versions[0]
			}

			var currentRevision string
			var isTrackedByCommitHash bool
			if currentVersion.Tag != "" {
				currentRevision = currentVersion.Tag
			} else if currentVersion.Commit != "" {
				currentRevision = currentVersion.Commit
				isTrackedByCommitHash = true
			}

			var latestRevision string
			if projectName == "cilium/cilium" {
				latestRevision, _, err = ecrpublic.GetLatestRevision(constants.CiliumImageRepository, currentRevision, branchName)
				if err != nil {
					return fmt.Errorf("getting latest revision from ECR Public: %v", err)
				}
			} else {
				// Get latest revision for the project from GitHub.
				latestRevision, _, err = github.GetLatestRevision(client, org, repoName, currentRevision, branchName, isTrackedByCommitHash, releaseBranched)
				if err != nil {
					return fmt.Errorf("getting latest revision from GitHub: %v", err)
				}
			}

			// Check if we should print only the latest version of the project.
			if displayOptions.PrintLatestVersion {
				fmt.Println(latestRevision)
				return nil
			} else {
				projectVersionInfoList = append(projectVersionInfoList, types.ProjectVersionInfo{Org: org, Repo: repoName, CurrentVersion: currentRevision, LatestVersion: latestRevision})
			}
		}
	}

	// Create a new table with the required column names in uppercase.
	tbl := table.New("Organization", "Repository", "Current Version", "Latest Version").WithHeaderFormatter(func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	})

	// Add rows to the table for each project in the list.
	for _, versionInfo := range projectVersionInfoList {
		tbl.AddRow(versionInfo.Org, versionInfo.Repo, versionInfo.CurrentVersion, versionInfo.LatestVersion)
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

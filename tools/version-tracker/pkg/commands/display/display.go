package display

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gogithub "github.com/google/go-github/v53/github"
	"github.com/rodaine/table"
	"gopkg.in/yaml.v3"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/git"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/github"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
)

// Run contains the business logic to execute the `display` subcommand.
func Run(displayOptions *types.DisplayOptions) error {
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

	// Clone the eks-anywhere-build-tooling repository.
	buildToolingRepoPath := filepath.Join(cwd, constants.BuildToolingRepoName)
	_, _, err = git.CloneRepo(constants.BuildToolingRepoURL, buildToolingRepoPath, "")
	if err != nil {
		return fmt.Errorf("cloning build-tooling repo: %v", err)
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
		return fmt.Errorf("unmarshaling upstream projects tracker file to YAML: %v", err)
	}

	var projectVersionInfoList []types.ProjectVersionInfo
	for _, project := range projectsList.Projects {
		org := project.Org
		for _, repo := range project.Repos {
			repoName := repo.Name
			currentVersion := repo.Versions[len(repo.Versions)-1]
			var currentRevision string
			if currentVersion.Tag != "" {
				currentRevision = currentVersion.Tag
			} else if currentVersion.Commit != "" {
				currentRevision = currentVersion.Commit
			}
			fullRepoName := fmt.Sprintf("%s/%s", org, repoName)
			if displayOptions.ProjectName != "" && displayOptions.ProjectName != fullRepoName {
				continue
			}

			// Get latest revision for the project from GitHub.
			latestRevision, _, err := github.GetLatestRevision(client, org, repoName, currentRevision)
			if err != nil {
				return fmt.Errorf("getting latest revision from GitHub: %v", err)
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

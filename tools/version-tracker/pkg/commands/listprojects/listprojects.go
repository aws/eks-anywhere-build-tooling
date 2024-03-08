package listprojects

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rodaine/table"
	"gopkg.in/yaml.v3"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/git"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
)

// Run contains the business logic to execute the `list-projects“ subcommand.
func Run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("retrieving current working directory: %v", err)
	}

	// Get base repository owner environment variable if set.
	baseRepoOwner := os.Getenv(constants.BaseRepoOwnerEnvvar)
	if baseRepoOwner == "" {
		baseRepoOwner = constants.DefaultBaseRepoOwner
	}

	// Clone the eks-anywhere-build-tooling repository.
	buildToolingRepoPath := filepath.Join(cwd, constants.BuildToolingRepoName)
	_, _, err = git.CloneRepo(fmt.Sprintf(constants.BuildToolingRepoURL, baseRepoOwner), buildToolingRepoPath, "")
	if err != nil {
		return fmt.Errorf("cloning build-tooling repo: %v", err)
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
	var maxOrgNameLength, maxRepoNameLength int
	for _, project := range projectsList.Projects {
		org := project.Org
		// Keep track of the longest org name to align the organization column in the table.
		if len(org) > maxOrgNameLength {
			maxOrgNameLength = len(org)
		}
		var repoList []string
		for _, repo := range project.Repos {
			// Keep track of the longest repository name to align the repository column in the table.
			if len(repo.Name) > maxRepoNameLength {
				maxRepoNameLength = len(repo.Name)
			}
			repoList = append(repoList, repo.Name)
		}
		// Apply newline padding to vertically align the organization name in the table.
		if len(repoList) > 2 {
			padding := strings.Repeat("\n", len(repoList)/2)
			org = fmt.Sprintf("%s%s%s", padding, org, padding)
		}
		projectVersionInfoList = append(projectVersionInfoList, types.ProjectVersionInfo{Org: org, Repo: strings.Join(repoList, "\n")})
	}

	// Create a new table with the required column names in uppercase.
	tbl := table.New("Organization", "Repository").WithHeaderFormatter(func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	})

	// Add rows to the table for each project in the list, grouped by owner or organization.
	tbl.AddRow(strings.Repeat("-", maxOrgNameLength), strings.Repeat("-", maxRepoNameLength))
	for _, versionInfo := range projectVersionInfoList {
		tbl.AddRow(versionInfo.Org, versionInfo.Repo)
		tbl.AddRow(strings.Repeat("-", maxOrgNameLength), strings.Repeat("-", maxRepoNameLength))
	}

	// Print the table contents to standard output.
	tbl.Print()

	return nil
}

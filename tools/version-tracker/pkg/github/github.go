package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-github/v53/github"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/file"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/tar"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/version"
)

// GetReleasesForRepo retrieves the list of releases for the given GitHub repository.
func GetReleasesForRepo(client *github.Client, org, repo string) ([]*github.RepositoryRelease, error) {
	logger.V(6).Info(fmt.Sprintf("Getting releases for [%s/%s] repository\n", org, repo))
	var allReleases []*github.RepositoryRelease
	listReleasesOptions := &github.ListOptions{
		PerPage: constants.GithubPerPage,
	}

	for {
		releases, resp, err := client.Repositories.ListReleases(context.Background(), org, repo, listReleasesOptions)
		if err != nil {
			return nil, fmt.Errorf("calling ListReleases API for [%s/%s] repository: %v", org, repo, err)
		}
		for _, release := range releases {
			if !*release.Prerelease {
				allReleases = append(allReleases, release)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		listReleasesOptions.Page = resp.NextPage
	}

	return allReleases, nil
}

// GetTagsForRepo retrieves the list of tags for the given GitHub repository.
func GetTagsForRepo(client *github.Client, org, repo string) ([]*github.RepositoryTag, error) {
	logger.V(6).Info(fmt.Sprintf("Getting tags for [%s/%s] repository\n", org, repo))
	var allTags []*github.RepositoryTag
	listTagOptions := &github.ListOptions{
		PerPage: constants.GithubPerPage,
	}

	for {
		tags, resp, err := client.Repositories.ListTags(context.Background(), org, repo, listTagOptions)
		if err != nil {
			return nil, fmt.Errorf("calling ListTags API for [%s/%s] repository: %v", org, repo, err)
		}
		allTags = append(allTags, tags...)

		if resp.NextPage == 0 {
			break
		}
		listTagOptions.Page = resp.NextPage
	}

	return allTags, nil
}

// GetCommitsForRepo retrieves the list of commits for the given GitHub repository.
func GetCommitsForRepo(client *github.Client, org, repo string) ([]*github.RepositoryCommit, error) {
	logger.V(6).Info(fmt.Sprintf("Getting commits for [%s/%s] repository\n", org, repo))
	var allCommits []*github.RepositoryCommit
	listCommitOptions := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: constants.GithubPerPage,
		},
	}

	for {
		commits, resp, err := client.Repositories.ListCommits(context.Background(), org, repo, listCommitOptions)
		if err != nil {
			return nil, fmt.Errorf("calling ListCommits for [%s/%s] repository: %v", org, repo, err)
		}
		allCommits = append(allCommits, commits...)

		if resp.NextPage == 0 {
			break
		}
		listCommitOptions.ListOptions.Page = resp.NextPage
	}

	return allCommits, nil
}

// GetCommitDateEpoch gets the Unix epoch time equivalent of a given Github commit's date.
func GetCommitDateEpoch(client *github.Client, org, repo, commitSHA string) (int64, error) {
	logger.V(6).Info(fmt.Sprintf("Getting date for commit %s in [%s/%s] repository\n", commitSHA, org, repo))

	commit, _, err := client.Repositories.GetCommit(context.Background(), org, repo, commitSHA, nil)
	if err != nil {
		return 0, fmt.Errorf("getting date for commit %s in [%s/%s] repository: %v", commit, org, repo, err)
	}

	return (*commit.Commit.Author.Date).Unix(), nil
}

// GetLatestRevision returns the latest revision (GitHub release or tag) for a given GitHub repository.
func GetLatestRevision(client *github.Client, org, repo string) (string, string, []*github.RepositoryTag, error) {
	logger.V(6).Info(fmt.Sprintf("Getting latest revision for [%s/%s] repository\n", org, repo))
	var latestRevision string
	allReleases, err := GetReleasesForRepo(client, org, repo)
	if err != nil {
		return "", "", nil, fmt.Errorf("getting all releases for [%s/%s] repository: %v", org, repo, err)
	}

	allTags, err := GetTagsForRepo(client, org, repo)
	if err != nil {
		return "", "", nil, fmt.Errorf("getting all tags for [%s/%s] repository: %v", org, repo, err)
	}
	if len(allReleases) > 0 {
		latestRevision = *allReleases[0].TagName
	} else {
		if len(allTags) > 0 {
			latestRevision = *allTags[0].Name
		} else {
			allCommits, err := GetCommitsForRepo(client, org, repo)
			if err != nil {
				return "", "", nil, fmt.Errorf("getting all commits for [%s/%s] repository: %v", org, repo, err)
			}
			latestRevision = *allCommits[0].SHA
		}
	}

	var latestRevisionCommit string
	commitHashMatch, err := regexp.MatchString("^[0-9a-f]{40}$", latestRevision)
	if err != nil {
		return "", "", nil, fmt.Errorf("checking if latest revision is a commit hash: %v", err)
	}
	if commitHashMatch {
		latestRevisionCommit = latestRevision
	} else {
		latestRevisionCommit = GetCommitForTag(allTags, latestRevision)
		if latestRevisionCommit == "" {
			return "", "", nil, fmt.Errorf("empty commit hash for latest revision: %s", latestRevision)
		}
	}

	return latestRevision, latestRevisionCommit, allTags, nil
}

// GetCommitForTag returns the commit hash corresponding to the given tag.
func GetCommitForTag(allTags []*github.RepositoryTag, searchTag string) string {
	for _, tag := range allTags {
		if searchTag == *tag.Name {
			return *tag.Commit.SHA
		}
	}
	return ""
}

// GetGoVersionForLatestRevision gets the Go version used to build the latest revision of the project.
func GetGoVersionForLatestRevision(client *github.Client, org, repo, latestRevision string) (string, error) {
	logger.V(6).Info(fmt.Sprintf("Getting Go version corresponding to latest revision %s for [%s/%s] repository\n", latestRevision, org, repo))
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("retrieving current working directory: %v", err)
	}

	var goVersion string
	projectFullName := fmt.Sprintf("%s/%s", org, repo)
	if _, ok := constants.ProjectReleaseAssets[projectFullName]; ok {
		release, _, err := client.Repositories.GetReleaseByTag(context.Background(), org, repo, latestRevision)
		if err != nil {
			return "", fmt.Errorf("calling GetReleaseByTag for tag %s in [%s/%s] repository: %v", latestRevision, org, repo, err)
		}
		var tarballName, tarballUrl string
		projectReleaseAsset := constants.ProjectReleaseAssets[projectFullName]
		searchAssetName := projectReleaseAsset.AssetName
		assetVersionReplacement := latestRevision
		if constants.ProjectReleaseAssets[projectFullName].TrimLeadingVersionPrefix {
			assetVersionReplacement = latestRevision[1:]
		}
		if strings.Count(searchAssetName, "%s") > 0 {
			searchAssetName = fmt.Sprintf(searchAssetName, assetVersionReplacement)
		}
		if projectReleaseAsset.OverrideAssetURL != "" {
			tarballName = searchAssetName
			tarballUrl = projectReleaseAsset.OverrideAssetURL
			if strings.Count(tarballUrl, "%s") > 0 {
				tarballUrl = fmt.Sprintf(tarballUrl, assetVersionReplacement)
			}
		} else {
			for _, asset := range release.Assets {
				if *asset.Name == searchAssetName {
					tarballName = *asset.Name
					tarballUrl = *asset.BrowserDownloadURL
					break
				}
			}
		}

		tarballDownloadPath := filepath.Join(cwd, "github-release-downloads")
		err = os.MkdirAll(tarballDownloadPath, 0o755)
		if err != nil {
			return "", fmt.Errorf("failed to create GitHub release downloads folder: %v", err)
		}
		tarballFilePath := filepath.Join(tarballDownloadPath, tarballName)

		err = file.Download(tarballUrl, tarballFilePath)
		if err != nil {
			return "", fmt.Errorf("downloading release tarball from URL [%s]: %v", tarballUrl, err)
		}

		if projectReleaseAsset.Extract {
			tarballFile, err := os.Open(tarballFilePath)
			if err != nil {
				return "", fmt.Errorf("opening tarball filepath: %v", err)
			}

			err = tar.ExtractTarGz(tarballDownloadPath, tarballFile)
			if err != nil {
				return "", fmt.Errorf("extracting tarball file: %v", err)
			}
		}

		binaryName := projectReleaseAsset.BinaryName
		if strings.Count(binaryName, "%s") > 0 {
			binaryName = fmt.Sprintf(binaryName, assetVersionReplacement)
		}
		binaryFilePath := filepath.Join(tarballDownloadPath, binaryName)
		goVersion, err = version.GetGoVersion(binaryFilePath)
		if err != nil {
			return "", fmt.Errorf("getting Go version embedded in binary [%s]: %v", binaryFilePath, err)
		}

		err = os.RemoveAll(tarballDownloadPath)
		if err != nil {
			return "", fmt.Errorf("removing tarball download path: %v", err)
		}
	} else if _, ok := constants.ProjectGoVersionSourceOfTruth[projectFullName]; ok {
		projectGoVersionSourceOfTruthFile := constants.ProjectGoVersionSourceOfTruth[projectFullName].SourceOfTruthFile
		workflowContents, _, _, err := client.Repositories.GetContents(context.Background(), org, repo, projectGoVersionSourceOfTruthFile, &github.RepositoryContentGetOptions{Ref: latestRevision})
		if err != nil {
			return "", fmt.Errorf("getting contents of file [%s]: %v", projectGoVersionSourceOfTruthFile, err)
		}
		workflowContentsDecoded, err := base64.StdEncoding.DecodeString(*workflowContents.Content)
		if err != nil {
			return "", fmt.Errorf("decoding contents of file [%s]: %v", projectGoVersionSourceOfTruthFile, err)
		}
		pattern := regexp.MustCompile(constants.ProjectGoVersionSourceOfTruth[projectFullName].GoVersionSearchString)
		matches := pattern.FindStringSubmatch(string(workflowContentsDecoded))

		goVersion = matches[1]
	}

	return goVersion, nil
}

// CreatePullRequest creates a pull request from the head branch to the base branch on the base repository.
func CreatePullRequest(client *github.Client, org, repo, baseRepoOwner, baseBranch, headRepoOwner, headBranch, currentRevision, latestRevision string, projectHasPatches bool) error {
	logger.V(6).Info(fmt.Sprintf("Creating pull request with updated versions for [%s/%s] repository\n", org, repo))

	pullRequests, _, err := client.PullRequests.List(context.Background(), baseRepoOwner, constants.BuildToolingRepoName, &github.PullRequestListOptions{
		Head: fmt.Sprintf("%s:%s", headRepoOwner, headBranch),
	})
	if err != nil {
		return fmt.Errorf("listing pull requests with %s:%s as head branch: %v", headRepoOwner, headBranch, err)
	}
	if len(pullRequests) > 0 {
		logger.Info(fmt.Sprintf("A pull request already exists for %s:%s\n", headRepoOwner, headBranch), "Pull request", *pullRequests[0].HTMLURL)
		return nil
	}

	newPR := &github.NewPullRequest{
		Title:               github.String(fmt.Sprintf("Bump %s/%s to latest release", org, repo)),
		Head:                github.String(fmt.Sprintf("%s:%s", headRepoOwner, headBranch)),
		Base:                github.String(baseBranch),
		Body:                github.String(fmt.Sprintf(constants.PullRequestBody, org, repo, currentRevision, latestRevision)),
		MaintainerCanModify: github.Bool(true),
	}

	pullRequest, _, err := client.PullRequests.Create(context.Background(), baseRepoOwner, constants.BuildToolingRepoName, newPR)
	if err != nil {
		return fmt.Errorf("creating pull request with updated versions from %s to %s: %v", headBranch, baseBranch, err)
	}

	if projectHasPatches {
		newComment := &github.IssueComment{
			Body: github.String(constants.PatchesCommentBody),
		}

		_, _, err = client.Issues.CreateComment(context.Background(), baseRepoOwner, constants.BuildToolingRepoName, *pullRequest.Number, newComment)
		if err != nil {
			return fmt.Errorf("commenting patch warning on pull request [%s]: %v", *pullRequest.HTMLURL, err)
		}
	}

	return nil
}

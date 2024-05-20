package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/aws/eks-anywhere/pkg/semver"
	"github.com/google/go-github/v62/github"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/file"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/tar"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/version"
)

// getTagsForRepo retrieves the list of tags for the given GitHub repository.
func getTagsForRepo(client *github.Client, org, repo string) ([]*github.RepositoryTag, error) {
	logger.V(6).Info(fmt.Sprintf("Getting tags for [%s/%s] repository", org, repo))
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

// getCommitsForRepo retrieves the list of commits for the given GitHub repository.
func getCommitsForRepo(client *github.Client, org, repo string) ([]*github.RepositoryCommit, error) {
	logger.V(6).Info(fmt.Sprintf("Getting commits for [%s/%s] repository", org, repo))
	var allCommits []*github.RepositoryCommit
	listCommitOptions := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: constants.GithubPerPage,
		},
	}

	for {
		commits, resp, err := client.Repositories.ListCommits(context.Background(), org, repo, listCommitOptions)
		if err != nil {
			return nil, fmt.Errorf("calling ListCommits API for [%s/%s] repository: %v", org, repo, err)
		}
		allCommits = append(allCommits, commits...)

		if resp.NextPage == 0 {
			break
		}
		listCommitOptions.ListOptions.Page = resp.NextPage
	}

	return allCommits, nil
}

// getCommitDateEpoch gets the Unix epoch time equivalent of a given Github commit's date.
func getCommitDateEpoch(client *github.Client, org, repo, commitSHA string) (int64, error) {
	logger.V(6).Info(fmt.Sprintf("Getting date for commit %s in [%s/%s] repository", commitSHA, org, repo))

	commit, _, err := client.Repositories.GetCommit(context.Background(), org, repo, commitSHA, nil)
	if err != nil {
		return 0, fmt.Errorf("getting date for commit %s in [%s/%s] repository: %v", commit, org, repo, err)
	}

	return (*commit.Commit.Author.Date).Unix(), nil
}

func GetFileContents(client *github.Client, org, repo, filePath, ref string) ([]byte, error) {
	contents, _, _, err := client.Repositories.GetContents(context.Background(), org, repo, filePath, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return nil, fmt.Errorf("getting contents of file [%s]: %v", filePath, err)
	}
	contentsDecoded, err := base64.StdEncoding.DecodeString(*contents.Content)
	if err != nil {
		return nil, fmt.Errorf("decoding contents of file [%s]: %v", filePath, err)
	}

	return contentsDecoded, nil
}

// GetLatestRevision returns the latest revision (GitHub release or tag) for a given GitHub repository.
func GetLatestRevision(client *github.Client, org, repo, currentRevision string, isTrackedUsingCommitHash, releaseBranched bool) (string, bool, error) {
	logger.V(6).Info(fmt.Sprintf("Getting latest revision for [%s/%s] repository", org, repo))
	var currentRevisionCommit, latestRevision string
	needsUpgrade := false

	// Get all GitHub tags for this project.
	allTags, err := getTagsForRepo(client, org, repo)
	if err != nil {
		return "", false, fmt.Errorf("getting all tags for [%s/%s] repository: %v", org, repo, err)
	}

	// Get commit hash corresponding to current revision tag.
	if isTrackedUsingCommitHash {
		currentRevisionCommit = currentRevision
	} else {
		currentRevisionCommit = getCommitForTag(allTags, currentRevision)
	}

	// Get Unix timestamp for current revision's commit.
	currentRevisionCommitEpoch, err := getCommitDateEpoch(client, org, repo, currentRevisionCommit)
	if err != nil {
		return "", false, fmt.Errorf("getting epoch time corresponding to current revision commit: %v", err)
	}

	// If the project is tracked using a commit hash, upgrade to the latest commit.
	if isTrackedUsingCommitHash {
		// If the project does not have Github tags, pick the latest commit.
		allCommits, err := getCommitsForRepo(client, org, repo)
		if err != nil {
			return "", false, fmt.Errorf("getting all commits for [%s/%s] repository: %v", org, repo, err)
		}
		latestRevision = *allCommits[0].SHA
		needsUpgrade = true
	} else {
		semverRegex := regexp.MustCompile(constants.SemverRegex)
		currentRevisionForSemver := semverRegex.FindString(currentRevision)

		// Get SemVer construct corresponding to the current revision tag.
		currentRevisionSemver, err := semver.New(currentRevisionForSemver)
		if err != nil {
			return "", false, fmt.Errorf("getting semver for current version: %v", err)
		}

		for _, tag := range allTags {
			tagName := *tag.Name
			if org == "kubernetes" && repo == "autoscaler" {
				if !strings.HasPrefix(tagName, "cluster-autoscaler-") {
					continue
				}
			}
			tagNameForSemver := semverRegex.FindString(tagName)
			if tagNameForSemver == "" {
				continue
			}

			if releaseBranched {
				releaseBranch := os.Getenv(constants.ReleaseBranchEnvvar)
				releaseNumber := strings.Split(releaseBranch, "-")[1]
				tagRegex := regexp.MustCompile(fmt.Sprintf(`^v?1.%s.\d$`, releaseNumber))
				if !tagRegex.MatchString(tagNameForSemver) {
					continue
				}
			}

			revisionSemver, err := semver.New(tagNameForSemver)
			if err != nil {
				return "", false, fmt.Errorf("getting semver for the version under consideration: %v", err)
			}
			if !slices.Contains(constants.ProjectsSupportingPrereleaseTags, fmt.Sprintf("%s/%s", org, repo)) && revisionSemver.Prerelease != "" {
				continue
			}

			latestRevision = tagName

			// Determine if upgrade is required based on current and latest revisions
			upgradeRequired, shouldBreak, err := isUpgradeRequired(client, org, repo, latestRevision, currentRevisionCommitEpoch, currentRevisionSemver, allTags)
			if err != nil {
				return "", false, fmt.Errorf("determining if upgrade is required for project: %v", err)
			}
			if shouldBreak {
				needsUpgrade = upgradeRequired
				break
			}
		}
	}

	return latestRevision, needsUpgrade, nil
}

// isUpgradeRequired determines if the project requires an upgrade by comparing the current revision to the latest revision.
func isUpgradeRequired(client *github.Client, org, repo, latestRevision string, currentRevisionCommitEpoch int64, currentRevisionSemver *semver.Version, allTags []*github.RepositoryTag) (bool, bool, error) {
	needsUpgrade := false
	shouldBreak := false

	// Get commit hash corresponding to latest revision tag.
	latestRevisionCommit := getCommitForTag(allTags, latestRevision)
	if latestRevisionCommit == "" {
		return false, false, fmt.Errorf("empty commit hash for latest revision: %s", latestRevision)
	}

	// Get Unix timestamp for latest revision's commit.
	latestRevisionCommitEpoch, err := getCommitDateEpoch(client, org, repo, latestRevisionCommit)
	if err != nil {
		return false, false, fmt.Errorf("getting epoch time corresponding to latest revision commit: %v", err)
	}

	semverRegex := regexp.MustCompile(constants.SemverRegex)
	latestRevisionForSemver := semverRegex.FindString(latestRevision)

	// Get SemVer construct corresponding to the latest revision tag.
	latestRevisionSemver, err := semver.New(latestRevisionForSemver)
	if err != nil {
		return false, false, fmt.Errorf("getting semver for latest version: %v", err)
	}

	// If the latest revision comes after the current revision both chronologically and semantically, then declare that
	// an upgrade is required
	if latestRevisionSemver.GreaterThan(currentRevisionSemver) || latestRevisionCommitEpoch > currentRevisionCommitEpoch {
		needsUpgrade = true
		shouldBreak = true
	} else if latestRevisionSemver.Equal(currentRevisionSemver) {
		needsUpgrade = false
		shouldBreak = true
	}

	return needsUpgrade, shouldBreak, nil
}

// getCommitForTag returns the commit hash corresponding to the given tag.
func getCommitForTag(allTags []*github.RepositoryTag, searchTag string) string {
	for _, tag := range allTags {
		if searchTag == *tag.Name {
			return *tag.Commit.SHA
		}
	}
	return ""
}

// GetGoVersionForLatestRevision gets the Go version used to build the latest revision of the project.
func GetGoVersionForLatestRevision(client *github.Client, org, repo, latestRevision string) (string, error) {
	logger.V(6).Info(fmt.Sprintf("Getting Go version corresponding to latest revision %s for [%s/%s] repository", latestRevision, org, repo))
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("retrieving current working directory: %v", err)
	}

	var goVersion string
	projectFullName := fmt.Sprintf("%s/%s", org, repo)
	if _, ok := constants.ProjectReleaseAssets[projectFullName]; ok {
		release, _, err := client.Repositories.GetReleaseByTag(context.Background(), org, repo, latestRevision)
		if err != nil {
			return "", fmt.Errorf("calling GetReleaseByTag API for tag %s in [%s/%s] repository: %v", latestRevision, org, repo, err)
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

		binaryName := projectReleaseAsset.BinaryName
		if strings.Count(binaryName, "%s") > 0 {
			binaryName = fmt.Sprintf(binaryName, assetVersionReplacement)
		}
		if projectReleaseAsset.Extract {
			tarballFile, err := os.Open(tarballFilePath)
			if err != nil {
				return "", fmt.Errorf("opening tarball filepath: %v", err)
			}

			err = tar.ExtractFileFromTarball(tarballDownloadPath, tarballFile, binaryName)
			if err != nil {
				return "", fmt.Errorf("extracting tarball file: %v", err)
			}
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
		workflowContents, err := GetFileContents(client, org, repo, projectGoVersionSourceOfTruthFile, latestRevision)
		if err != nil {
			return "", fmt.Errorf("getting contents of file [%s]: %v", projectGoVersionSourceOfTruthFile, err)
		}

		pattern := regexp.MustCompile(constants.ProjectGoVersionSourceOfTruth[projectFullName].GoVersionSearchString)
		matches := pattern.FindStringSubmatch(string(workflowContents))

		goVersion = matches[1]
	}

	return goVersion, nil
}

// CreatePullRequest creates a pull request from the head branch to the base branch on the base repository.
func CreatePullRequest(client *github.Client, org, repo, title, body, baseRepoOwner, baseBranch, headRepoOwner, headBranch, currentRevision, latestRevision string, addPatchWarningComment bool, patchesWarningComment string) error {
	var pullRequest *github.PullRequest
	var patchWarningCommentExists bool

	// Check if there is already a pull request from the head branch to the base branch.
	pullRequests, _, err := client.PullRequests.List(context.Background(), baseRepoOwner, constants.BuildToolingRepoName, &github.PullRequestListOptions{
		Base: baseBranch,
		Head: fmt.Sprintf("%s:%s", headRepoOwner, headBranch),
	})
	if err != nil {
		return fmt.Errorf("listing pull requests from %s:%s -> %s:%s: %v", headRepoOwner, headBranch, baseRepoOwner, baseBranch, err)
	}

	if len(pullRequests) > 0 {
		pullRequest = pullRequests[0]
		logger.Info(fmt.Sprintf("A pull request already exists for %s:%s", headRepoOwner, headBranch), "Pull request", *pullRequest.HTMLURL)

		pullRequest.Body = github.String(body)
		pullRequest, _, err = client.PullRequests.Edit(context.Background(), baseRepoOwner, constants.BuildToolingRepoName, *pullRequest.Number, pullRequest)
		if err != nil {
			return fmt.Errorf("editing existing pull request [%s]: %v", *pullRequest.HTMLURL, err)
		}

		// If patches to the project failed to apply, check if the PR already has a comment warning about
		// the incomplete PR and patches needing to be regenerated.
		if addPatchWarningComment {
			pullRequestComments, _, err := client.Issues.ListComments(context.Background(), baseRepoOwner, constants.BuildToolingRepoName, *pullRequest.Number, nil)
			if err != nil {
				return fmt.Errorf("listing comments on pull request [%s]: %v", *pullRequest.HTMLURL, err)
			}

			for _, comment := range pullRequestComments {
				if comment.Body == github.String(patchesWarningComment) {
					patchWarningCommentExists = true
				}
			}
		}
	} else {
		logger.V(6).Info(fmt.Sprintf("Creating pull request with updated versions for [%s/%s] repository", org, repo))

		newPR := &github.NewPullRequest{
			Title:               github.String(title),
			Head:                github.String(fmt.Sprintf("%s:%s", headRepoOwner, headBranch)),
			Base:                github.String(baseBranch),
			Body:                github.String(body),
			MaintainerCanModify: github.Bool(true),
		}
		pullRequest, _, err = client.PullRequests.Create(context.Background(), baseRepoOwner, constants.BuildToolingRepoName, newPR)
		if err != nil {
			return fmt.Errorf("creating pull request with updated versions from %s to %s: %v", headBranch, baseBranch, err)
		}

		logger.Info(fmt.Sprintf("Created pull request: %s", *pullRequest.HTMLURL))
	}

	// If patches failed to apply and no patch warning comment exists (always the case for a new PR), then add a comment with the
	// warning.
	if addPatchWarningComment && !patchWarningCommentExists {
		patchWarningComment := &github.IssueComment{
			Body: github.String(patchesWarningComment),
		}

		_, _, err = client.Issues.CreateComment(context.Background(), baseRepoOwner, constants.BuildToolingRepoName, *pullRequest.Number, patchWarningComment)
		if err != nil {
			return fmt.Errorf("commenting failed patch apply warning on pull request [%s]: %v", *pullRequest.HTMLURL, err)
		}
	}

	return nil
}

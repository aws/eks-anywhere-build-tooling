package git

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// CloneRepo clones the remote repository to a destination folder and creates a Git remote.
func CloneRepo(cloneURL, destination, headRepoOwner string) (*git.Repository, string, error) {
	logger.V(6).Info(fmt.Sprintf("Cloning repository [%s] to %s directory\n", cloneURL, destination))
	progress := io.Discard
	if logger.Verbosity >= 6 {
		progress = os.Stdout
	}
	repo, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL:      cloneURL,
		Progress: progress,
	})
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			logger.V(6).Info(fmt.Sprintf("Repo already exists at %s\n", destination))
			repo, err = git.PlainOpen(destination)
		} else {
			return nil, "", fmt.Errorf("cloning repo %s to %s directory: %v", cloneURL, destination, err)
		}
	}

	repoHeadCommit, err := repo.ResolveRevision(plumbing.Revision(constants.BaseRepoHeadRevision))
	if err != nil {
		return nil, "", fmt.Errorf("resolving revision [%s] to commit hash: %v", constants.BaseRepoHeadRevision, err)
	}
	repoHeadCommitHash := strings.Split(repoHeadCommit.String(), " ")[0]

	if headRepoOwner != "" {
		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: headRepoOwner,
			URLs: []string{fmt.Sprintf("https://github.com/%s/%s.git", headRepoOwner, constants.BuildToolingRepoName)},
		})
		if err != nil {
			if err == git.ErrRemoteExists {
				logger.V(6).Info(fmt.Sprintf("Remote %s already exists", headRepoOwner))
			} else {
				return nil, "", fmt.Errorf("creating remote %s: %v", headRepoOwner, err)
			}
		}
	}

	return repo, repoHeadCommitHash, nil
}

// ResetToMain hard-resets the current working tree to point to the HEAD commit of the base repository.
func ResetToMain(worktree *git.Worktree, baseRepoHeadCommit string) error {
	err := worktree.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(baseRepoHeadCommit),
		Mode:   git.HardReset,
	})
	if err != nil {
		return fmt.Errorf("resetting to origin HEAD commit %s: %v", baseRepoHeadCommit, err)
	}

	return nil
}

// Checkout checks out the working tree at the given branch, creating a new branch if necessary.
func Checkout(worktree *git.Worktree, branch string) error {
	logger.V(6).Info(fmt.Sprintf("Checking out branch [%s] in local worktree\n", branch))

	err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
		Create: true,
	})
	if err != nil {
		return fmt.Errorf("checking out branch [%s]: %v", branch, err)
	}

	return nil
}

// Add adds the given paths to the Git index.
func Add(worktree *git.Worktree, paths []string) error {
	logger.V(6).Info("Adding updated files to index")
	for _, path := range paths {
		_, err := worktree.Add(path)
		if err != nil {
			return fmt.Errorf("adding file [%s] to the index: %v", path, err)
		}
	}
	return nil
}

// Commit creates a new commit with the given commit message.
func Commit(worktree *git.Worktree, commitMessage string) error {
	logger.V(6).Info("Committing file(s) in the index")
	var commitAuthorName, commitAuthorEmail string
	commitAuthorName, ok := os.LookupEnv(constants.CommitAuthorNameEnvvar)
	if !ok {
		commitAuthorName = constants.DefaultCommitAuthorName
	}

	commitAuthorEmail, ok = os.LookupEnv(constants.CommitAuthorEmailEnvvar)
	if !ok {
		commitAuthorEmail = constants.DefaultCommitAuthorEmail
	}
	_, err := worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  commitAuthorName,
			Email: commitAuthorEmail,
		},
	})
	if err != nil {
		return fmt.Errorf("committing file(s) in the index: %v", err)
	}
	return nil
}

// Push pushes changes to the given remote branch on GitHub.
func Push(repo *git.Repository, headRepoOwner, branch, githubToken string) error {
	logger.V(6).Info(fmt.Sprintf("Pushing changes to remote [%s]", headRepoOwner))
	progress := io.Discard
	if logger.Verbosity >= 6 {
		progress = os.Stdout
	}
	err := repo.Push(&git.PushOptions{
		RemoteName: headRepoOwner,
		RefSpecs:   []config.RefSpec{config.RefSpec("+refs/*:refs/*")},
		Auth: &http.BasicAuth{
			Username: headRepoOwner,
			Password: githubToken,
		},
		Progress: progress,
		Force:    true,
	})
	if err != nil {
		return fmt.Errorf("pushing changes to remote %s: %v", headRepoOwner, err)
	}
	return nil
}

package fixpatches

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// GitCommand executes a git command in the specified directory.
// Returns combined stdout/stderr output and any error.
func GitCommand(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("git %v failed: %v\nOutput: %s", args, err, string(output))
	}
	return string(output), nil
}

// GitCommandWithC executes a git command using -C flag (changes directory before running).
// This is useful when you need to run git commands on a specific repo.
func GitCommandWithC(repoPath string, args ...string) (string, error) {
	fullArgs := append([]string{"-C", repoPath}, args...)
	cmd := exec.Command("git", fullArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("git -C %s %v failed: %v\nOutput: %s", repoPath, args, err, string(output))
	}
	return string(output), nil
}

// GitAdd stages files for commit.
func GitAdd(dir string, files ...string) error {
	_, err := GitCommand(dir, append([]string{"add"}, files...)...)
	return err
}

// GitCommit creates a commit with the given message.
// Returns nil if there's nothing to commit (not an error).
func GitCommit(dir, message string) error {
	output, err := GitCommand(dir, "commit", "-m", message)
	if err != nil {
		// Check if there's nothing to commit (not an error)
		if strings.Contains(output, "nothing to commit") {
			logger.Info("No changes to commit")
			return nil
		}
		return err
	}
	return nil
}

// GitConfig sets a git configuration value.
func GitConfig(dir, key, value string) error {
	_, err := GitCommand(dir, "config", key, value)
	return err
}

// GitApply applies a patch file.
func GitApply(dir, patchFile string, reject, whitespace bool) (string, error) {
	args := []string{"apply"}
	if reject {
		args = append(args, "--reject")
	}
	if whitespace {
		args = append(args, "--whitespace=fix")
	}
	args = append(args, patchFile)
	return GitCommand(dir, args...)
}

// GitApplyWithC applies a patch file using -C flag.
func GitApplyWithC(repoPath, patchFile string, reject, whitespace bool) (string, error) {
	args := []string{"apply"}
	if reject {
		args = append(args, "--reject")
	}
	if whitespace {
		args = append(args, "--whitespace=fix")
	}
	args = append(args, patchFile)
	return GitCommandWithC(repoPath, args...)
}

// GitStatus returns the status of the repository.
func GitStatus(dir string, porcelain bool, files ...string) (string, error) {
	args := []string{"status"}
	if porcelain {
		args = append(args, "--porcelain")
	}
	args = append(args, files...)
	return GitCommand(dir, args...)
}

// GitFormatPatch generates patch files from commits.
func GitFormatPatch(dir, outputDir string, commitCount int) (string, error) {
	args := []string{"format-patch", fmt.Sprintf("-%d", commitCount), "HEAD"}
	if outputDir != "" {
		args = append(args, "-o", outputDir)
	}
	return GitCommand(dir, args...)
}

// ConfigureGitUser sets the git user name and email for a repository.
// Uses the standard patch application user from constants.
func ConfigureGitUser(dir string) error {
	if err := GitConfig(dir, "user.email", constants.PatchApplyGitUserEmail); err != nil {
		return fmt.Errorf("configuring git user.email: %v", err)
	}
	if err := GitConfig(dir, "user.name", constants.PatchApplyGitUserName); err != nil {
		return fmt.Errorf("configuring git user.name: %v", err)
	}
	return nil
}

// ResetToCleanStateWithC resets the repository to a clean state using -C flag.
func ResetToCleanStateWithC(repoPath string) error {
	logger.Info("Resetting repository to clean state", "repo", repoPath)

	if _, err := GitCommandWithC(repoPath, "reset", "--hard", "HEAD"); err != nil {
		return fmt.Errorf("git reset failed: %v", err)
	}

	if _, err := GitCommandWithC(repoPath, "clean", "-fd"); err != nil {
		return fmt.Errorf("git clean failed: %v", err)
	}

	logger.Info("Repository reset to clean state")
	return nil
}

// GetRepoPath returns the path to the cloned repository within the project directory.
// The repository is typically cloned as: projectPath/<repoName>
// For example: projects/aquasecurity/trivy/trivy
func GetRepoPath(projectPath string) (string, error) {
	repoName := filepath.Base(projectPath)
	repoPath := filepath.Join(projectPath, repoName)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("repository not found at %s", repoPath)
	}

	return repoPath, nil
}

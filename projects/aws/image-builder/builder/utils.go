package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func cloneRepo(cloneUrl, destination string) error {
	log.Print("Cloning eks-anywhere-build-tooling...")
	cloneRepoCommandSequence := fmt.Sprintf("git clone %s %s", cloneUrl, destination)
	cmd := exec.Command("bash", "-c", cloneRepoCommandSequence)
	out, err := execCommandWithStreamOutput(cmd)
	fmt.Println(out)
	return err
}

func execCommandWithStreamOutput(cmd *exec.Cmd) (string, error) {
	log.Printf("Executing command: %v\n", cmd)
	commandOutput, err := cmd.CombinedOutput()
	commandOutputStr := strings.TrimSpace(string(commandOutput))

	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}
	return commandOutputStr, nil
}

func executeMakeBuildCommand(buildCommand, releaseChannel string) error {
	cmd := exec.Command("bash", "-c", buildCommand)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("RELEASE_BRANCH=%s", releaseChannel))
	out, err := execCommandWithStreamOutput(cmd)
	fmt.Println(out)
	return err
}

func cleanup(buildToolingDir string) {
	log.Print("Cleaning up cache build files")
	err := os.RemoveAll(buildToolingDir)
	if err != nil {
		log.Fatalf("Error cleaning up build tooling dir: %v", err)
	}
}

func GetSupportedReleaseBranches() []string {
	buildToolingPath, err := getRepoRoot()
	if err != nil {
		log.Fatalf(err.Error())
	}

	supportedBranchesFile := filepath.Join(buildToolingPath, "release/SUPPORTED_RELEASE_BRANCHES")
	supportedBranchesFileData, err := os.ReadFile(supportedBranchesFile)
	supportReleaseBranches := strings.Split(string(supportedBranchesFileData), "\n")

	return supportReleaseBranches
}

func getBuildToolingPath(cwd string) string {
	buildToolingRepoPath := filepath.Join(cwd, "eks-anywhere-build-tooling")
	if codebuild == "true" {
		buildToolingRepoPath = os.Getenv("CODEBUILD_SRC_DIR")
	}
	return buildToolingRepoPath
}

func getRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	return execCommandWithStreamOutput(cmd)
}

func SliceContains(s []string, str string) bool {
	for _, elem := range s {
		if elem == str {
			return true
		}
	}
	return false
}

func getCwd() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error retrieving current working directory: %v", err)
	}
	return cwd, nil
}
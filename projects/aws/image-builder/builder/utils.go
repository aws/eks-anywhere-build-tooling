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
	cloneRepoCommandSequence := fmt.Sprintf("git clone %s %s && cd %s && git checkout cloudstack-image", cloneUrl, destination, destination)
	cmd := exec.Command("bash", "-c", cloneRepoCommandSequence)
	return execCommandWithStreamOutput(cmd)
}

func execCommandWithStreamOutput(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Executing command: %v\n", cmd)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	return nil
}

func executeMakeBuildCommand(buildCommand string, envVars ...string) error {
	cmd := exec.Command("bash", "-c", buildCommand)
	cmd.Env = os.Environ()
	for _, envVar := range envVars {
		cmd.Env = append(cmd.Env, envVar)
	}
	return execCommandWithStreamOutput(cmd)
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
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error retrieving current working directory: %v", err)
	}
	buildToolingPath := getBuildToolingPath(cwd)
	cmd := exec.Command("git", "-C", buildToolingPath, "rev-parse", "--show-toplevel")
	commandOut, err := execCommand(cmd)
	if err != nil {
		return "", err
	}
	return commandOut, nil
}

func SliceContains(s []string, str string) bool {
	for _, elem := range s {
		if elem == str {
			return true
		}
	}
	return false
}

func execCommand(cmd *exec.Cmd) (string, error) {
	log.Printf("Executing command: %v\n", cmd)
	commandOutput, err := cmd.CombinedOutput()
	commandOutputStr := strings.TrimSpace(string(commandOutput))

	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}
	return commandOutputStr, nil
}

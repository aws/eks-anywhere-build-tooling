package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func cloneRepo(cloneUrl, destination string) error {
	log.Print("Cloning eks-anywhere-build-tooling...")
	cloneRepoCommandSequence := fmt.Sprintf("git clone %s %s", cloneUrl, destination)
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

func executeMakeBuildCommand(buildCommand, releaseChannel string) error {
	cmd := exec.Command("bash", "-c", buildCommand)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("RELEASE_BRANCH=%s", releaseChannel))
	err := execCommandWithStreamOutput(cmd)
	return err
}

func cleanup(buildToolingDir string) {
	log.Print("Cleaning up cache build files")
	err := os.RemoveAll(buildToolingDir)
	if err != nil {
		log.Fatalf("Error cleaning up build tooling dir: %v", err)
	}
}

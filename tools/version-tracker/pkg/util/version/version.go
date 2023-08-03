package version

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/command"
)

// GetGoVersion gets the Go version corresponding to the given Go binary
// using the `go version -m` command.
func GetGoVersion(goBinaryLocation string) (string, error) {
	goVersionCmd := exec.Command("go", "version", "-m", goBinaryLocation)
	commandOutput, err := command.ExecCommand(goVersionCmd)
	if err != nil {
		return "", fmt.Errorf("running Go version command: %v", err)
	}
	pattern := regexp.MustCompile(fmt.Sprintf(`^%s: go(.*)\.\d+\n`, goBinaryLocation))
	matches := pattern.FindStringSubmatch(commandOutput)

	return matches[1], nil
}

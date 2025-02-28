package version

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

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
	// The first line could be a warning, so no ^ in the sprintf below
	pattern := regexp.MustCompile(fmt.Sprintf(`%s: go(.*)\.\d+`, goBinaryLocation))
	matches := pattern.FindStringSubmatch(commandOutput)

	return matches[1], nil
}

func EnsurePatchVersion(version string) string {
	hasLeadingV := (version[0] == 'v')
	// Remove 'v' prefix if present
	if hasLeadingV {
		version = strings.TrimPrefix(version, "v")
	}

	// Regular expression to match version components
	re := regexp.MustCompile(`^(\d+\.\d+)(?:\.(\d+))?(.*)$`)
	matches := re.FindStringSubmatch(version)

	if len(matches) < 2 {
		// If the version string doesn't match the expected format, return it unchanged
		return version
	}

	majorMinor := matches[1]
	patch := matches[2]
	metadata := matches[3]

	if patch == "" {
		// If patch version is missing, add ".0"
		version = fmt.Sprintf("%s.0%s", majorMinor, metadata)
	}

	if hasLeadingV {
		version = fmt.Sprintf("v%s", version)
	}

	// If patch version is present, return the original version
	return version
}

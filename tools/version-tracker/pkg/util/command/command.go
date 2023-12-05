package command

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// ExecCommand executes the given command, writing to standard output.
func ExecCommand(cmd *exec.Cmd) (string, error) {
	logger.V(6).Info(fmt.Sprintf("Executing command: %s", cmd.String()))
	commandOutput, err := cmd.CombinedOutput()
	commandOutputStr := strings.TrimSpace(string(commandOutput))
	logger.V(6).Info(commandOutputStr)
	if err != nil {
		return commandOutputStr, fmt.Errorf("executing command %s: %v", cmd.String(), err)
	}
	return commandOutputStr, nil
}

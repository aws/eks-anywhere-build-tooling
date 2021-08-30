package etcdadm

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

func runPhases(command string, phases, flags []string, args ...string) error {
	for _, phase := range phases {
		fmt.Printf("Running etcdadm %s %s phase\n", command, phase)
		cmd := exec.Command(etcdadmBinary, buildPhaseCmd(command, phase, args...).addFlags(flags...)...)
		out, err := cmd.CombinedOutput()
		fmt.Printf("Phase command output:\n--------\n%s\n--------\n", string(out))
		if err != nil {
			return errors.Wrapf(err, "error running etcdadm phase '%s %s', out:\n %s", command, phase, string(out))
		}
	}

	return nil
}

type etcdadmCommand []string

func buildPhaseCmd(command, phase string, args ...string) etcdadmCommand {
	cmd := make(etcdadmCommand, 0, len(args)+3)
	cmd = append(cmd, command, "phase", phase)
	cmd = append(cmd, args...)
	return cmd
}

func (e etcdadmCommand) addFlags(flags ...string) etcdadmCommand {
	e = append(e, flags...)
	return e
}

package utils

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

func CreateSymLink(from, to string) error {
	// Force symlink creation with cmd. sdk calls fail if the symlink dir exists
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ln -sfn %s %s", from, to))
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error creating symlink: %v", cmd)
	}

	return nil
}

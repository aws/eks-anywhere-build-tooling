package system

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/files"
	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
)

const (
	bootstrapContainerPath = "/.bottlerocket/bootstrap-containers/current"
)

var rebootedPath = filepath.Join(bootstrapContainerPath, "rebooted")

func rebootInstanceIfNeeded() error {
	// Skip reboot when instance is already rebooted once
	if files.PathExists(rebootedPath) {
		return nil
	}

	// Write rebooted file to make sure reboot only happens once
	if _, err := os.Create(rebootedPath); err != nil {
		return errors.Wrap(err, "Error writing rebooted file")
	}

	if err := reboot(); err != nil {
		return err
	}

	return nil
}

func reboot() error {
	cmd := exec.Command(utils.ApiclientBinary, "reboot")
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "Error running command: %v, output: %s\n", cmd, string(out))
	}
	return nil
}

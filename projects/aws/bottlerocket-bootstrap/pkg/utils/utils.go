package utils

import (
	"os/exec"

	"github.com/pkg/errors"
)

const (
	ApiclientBinary        = "apiclient"
	BootstrapContainerName = "kubeadm-bootstrap"
)

func DisableBootstrapContainer() error {
	cmd := exec.Command(ApiclientBinary, "set", "host-containers."+BootstrapContainerName+".enabled=false")
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Error disabling bootstrap container")
	}
	return nil
}

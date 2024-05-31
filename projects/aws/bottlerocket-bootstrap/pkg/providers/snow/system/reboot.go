package system

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/files"
)

const (
	bootstrapContainerPath = "/.bottlerocket/bootstrap-containers/current"
)

var rebootedPath = filepath.Join(bootstrapContainerPath, "rebooted")

func (s *Snow) rebootInstanceIfNeeded() error {
	// Skip reboot when instance is already rebooted once
	if files.PathExists(rebootedPath) {
		return nil
	}

	// Write rebooted file to make sure reboot only happens once
	if _, err := os.Create(rebootedPath); err != nil {
		return errors.Wrap(err, "Error writing rebooted file")
	}

	if err := s.api.Reboot(); err != nil {
		return err
	}

	return nil
}

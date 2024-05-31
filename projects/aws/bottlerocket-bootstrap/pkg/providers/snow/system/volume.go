package system

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/files"
)

var device = filepath.Join(rootfs, "/dev/vda")

func (s *Snow) mountContainerdVolume() error {
	if err := s.createPartition(); err != nil {
		return errors.Wrap(err, "error creating partition")
	}

	if err := s.mountContainerd(); err != nil {
		return errors.Wrap(err, "error mounting containerd directory")
	}

	return nil
}

func (s *Snow) createPartition() error {
	partitionCreatedPath := filepath.Join(bootstrapContainerPath, "partition.created")
	if files.PathExists(partitionCreatedPath) {
		return nil
	}

	if err := s.fs.Partition(device); err != nil {
		return err
	}

	if _, err := os.Create(partitionCreatedPath); err != nil {
		return errors.Wrapf(err, "error writing partition created file")
	}

	return nil
}

func (s *Snow) mountContainerd() error {
	containerdDir := filepath.Join(rootfs, "/var/lib/containerd")
	if err := os.MkdirAll(containerdDir, 0o640); err != nil {
		return errors.Wrap(err, "error creating directory")
	}

	return s.fs.MountVolume(device, containerdDir)
}

package system

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/files"
	"github.com/pkg/errors"
)

var device = filepath.Join(rootfs, "/dev/vda")

func mountContainerdVolume() error {
	if err := createPartition(); err != nil {
		return errors.Wrap(err, "error creating partition")
	}

	if err := mountContainerd(); err != nil {
		return errors.Wrap(err, "error mounting containerd directory")
	}

	return nil
}

func createPartition() error {
	partitionCreatedPath := filepath.Join(bootstrapContainerPath, "partition.created")
	if files.PathExists(partitionCreatedPath) {
		return nil
	}

	cmd := exec.Command("mkfs.ext4", device)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error running command: %v", cmd)
	}

	if _, err := os.Create(partitionCreatedPath); err != nil {
		return errors.Wrapf(err, "error writing partition created file")
	}

	return nil
}

func mountContainerd() error {
	containerdDir := filepath.Join(rootfs, "/var/lib/containerd")
	if err := os.MkdirAll(containerdDir, 0o640); err != nil {
		return errors.Wrap(err, "error creating directory")
	}

	cmd := exec.Command("mount", device, containerdDir)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error running command: %v", cmd)
	}
	return nil
}

package upgrade

import (
	"context"
	"fmt"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// ContainerdUpgrade upgrades containerd on the node
// As part of the upgrade:
//  1. copy over the new containerd components
//  2. reload the daemon and restart containerd with the new version
func (u *InPlaceUpgrader) ContainerdUpgrade(ctx context.Context) error {
	cmpDir, err := u.upgradeComponentsBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade components binary directory: %v", err)
	}

	containerdVersionCmd := []string{"containerd", "--version"}
	out, err := u.ExecCommand(ctx, containerdVersionCmd[0], containerdVersionCmd[1:]...)
	if err != nil {
		return execError(containerdVersionCmd, string(out))
	}

	cpCmd := []string{"cp", "-rf", fmt.Sprintf("%s/containerd/.", cmpDir), "/"}
	out, err = u.ExecCommand(ctx, cpCmd[0], cpCmd[1:]...)
	if err != nil {
		return execError(cpCmd, string(out))
	}

	version, err := u.ExecCommand(ctx, containerdVersionCmd[0], containerdVersionCmd[1:]...)
	if err != nil {
		return execError(containerdVersionCmd, string(version))
	}

	daemonReloadCmd := []string{"systemctl", "daemon-reload"}
	out, err = u.ExecCommand(ctx, daemonReloadCmd[0], daemonReloadCmd[1:]...)
	if err != nil {
		return execError(daemonReloadCmd, string(out))
	}

	containerdRestartCmd := []string{"systemctl", "restart", "containerd"}
	out, err = u.ExecCommand(ctx, containerdRestartCmd[0], containerdRestartCmd[1:]...)
	if err != nil {
		return execError(containerdRestartCmd, string(out))
	}

	logger.Info("containerd version on the node", "version", string(version))
	logger.Info("containerd upgrade successful!")
	return nil
}

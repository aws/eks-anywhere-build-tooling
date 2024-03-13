package upgrade

import (
	"context"
	"fmt"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

func (u *Upgrader) LogStatusAndCleanup(ctx context.Context) error {
	if err := u.logStatus(ctx); err != nil {
		return fmt.Errorf("logging current status of the node: %v", err)
	}

	upgCmpDir, err := u.upgradeComponentsDir()
	if err != nil {
		return fmt.Errorf("getting upgrade componenets kubernetes binary directory: %v", err)
	}

	out, err := u.ExecCommand(ctx, "rm", "-rf", upgCmpDir)
	if err != nil {
		return fmt.Errorf("executing command 'rm -rf %s': %v", upgCmpDir, string(out))
	}

	logger.Info("Cleaning up leftover upgrade components", "in-place components directory", upgCmpDir)
	return nil
}

func (u *Upgrader) logStatus(ctx context.Context) error {
	out, err := u.ExecCommand(ctx, "systemctl", "status", "containerd")
	if err != nil {
		return fmt.Errorf("executing command 'systemctl status containerd': %s", string(out))
	}
	logger.Info("Containerd status on the Node", "status", string(out))

	out, err = u.ExecCommand(ctx, "systemctl", "status", "kubelet")
	if err != nil {
		return fmt.Errorf("executing command 'systemctl status kubelet': %s", string(out))
	}
	logger.Info("Kubelet status on the Node", "status", string(out))

	out, err = u.ExecCommand(ctx, "kubeadm", "version")
	if err != nil {
		return fmt.Errorf("executing command 'kubeadm --version': %s", string(out))
	}
	logger.Info("Kubeadm Version on the Node", "Version", string(out))

	return nil
}

package upgrade

import (
	"context"
	"fmt"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

func (u *Upgrader) ContainerdUpgrade(ctx context.Context) error {
	cmpDir, err := u.upgradeComponentsBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade componenets binary directory: %v", err)
	}

	out, err := u.ExecCommand(ctx, "containerd", "--version")
	if err != nil {
		return fmt.Errorf("executing command 'containerd --version': %s", string(out))
	}

	containerdDir := fmt.Sprintf("%s/containerd/.", cmpDir)
	out, err = u.ExecCommand(ctx, "cp", "-rf", containerdDir, "/")
	if err != nil {
		return fmt.Errorf("executing command 'cp -rf %s /': %s", containerdDir, string(out))
	}

	version, err := u.ExecCommand(ctx, "containerd", "--version")
	if err != nil {
		return fmt.Errorf("executing command 'containerd --version': %s", string(out))
	}

	out, err = u.ExecCommand(ctx, "systemctl", "daemon-reload")
	if err != nil {
		return fmt.Errorf("executing command 'systemctl daemon-reload': %s", string(out))
	}

	out, err = u.ExecCommand(ctx, "systemctl", "restart", "containerd")
	if err != nil {
		return fmt.Errorf("executing command 'systemctl restart containerd': %s", string(out))
	}

	logger.Info("Containerd Version on the Node", "Version", string(version))
	logger.Info("Containerd upgrade successful!")
	return nil
}

package upgrade

import (
	"context"
	"fmt"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

func (u *Upgrader) CniPluginsUpgrade(ctx context.Context) error {
	cmpDir, err := u.upgradeComponentsBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade componenets binary directory: %v", err)
	}

	out, err := u.ExecCommand(ctx, "/opt/cni/bin/loopback", "--version")
	if err != nil {
		return fmt.Errorf("executing command '/opt/cni/bin/loopback --version': %s", string(out))
	}

	out, err = u.ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/cni-plugins/.", cmpDir), "/")
	if err != nil {
		return fmt.Errorf("executing command 'cp -rf %s /': %s", fmt.Sprintf("%s/cni-plugins/*", cmpDir), string(out))
	}

	out, err = u.ExecCommand(ctx, "/opt/cni/bin/loopback", "--version")
	if err != nil {
		return fmt.Errorf("executing command '/opt/cni/bin/loopback --version': %s", string(out))
	}

	logger.Info("Cni-Plugins Version on the Node", "Version", string(out))
	logger.Info("Cni-Plugins upgrade succesful!")
	return nil
}

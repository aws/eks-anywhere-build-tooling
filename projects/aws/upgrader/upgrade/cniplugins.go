package upgrade

import (
	"context"
	"fmt"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// CNIPluginsUpgrade upgrades cni plugins on the node
func (u *InPlaceUpgrader) CNIPluginsUpgrade(ctx context.Context) error {
	cmpDir, err := u.upgradeComponentsBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade components binary directory: %v", err)
	}

	cniVersionCmd := []string{"/opt/cni/bin/loopback", "--version"}
	out, err := u.ExecCommand(ctx, cniVersionCmd[0], cniVersionCmd[1:]...)
	if err != nil {
		return execError(cniVersionCmd, string(out))
	}

	cpCmd := []string{"cp", "-rf", fmt.Sprintf("%s/cni-plugins/.", cmpDir), "/"}
	out, err = u.ExecCommand(ctx, cpCmd[0], cpCmd[1:]...)
	if err != nil {
		return execError(cpCmd, string(out))
	}

	out, err = u.ExecCommand(ctx, cniVersionCmd[0], cniVersionCmd[1:]...)
	if err != nil {
		return execError(cniVersionCmd, string(out))
	}

	logger.Info("cni plugins version on the node", "version", string(out))
	logger.Info("cni plugins upgrade successful!")
	return nil
}

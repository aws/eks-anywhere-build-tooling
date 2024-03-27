package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/upgrader/upgrade"
	"github.com/spf13/cobra"
)

var upgradeContainerdCmd = &cobra.Command{
	Use:   "containerd",
	Short: "Upgrade containerd",
	Long:  "Use InPlace Upgrader upgrade containerd to upgrade containerd on the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := upgradeContainerd(cmd.Context())
		if err != nil {
			log.Fatalf("upgrade containerd failed: %v", err)
		}
		return nil
	},
}

func init() {
	upgradeCmd.AddCommand(upgradeContainerdCmd)
}

func upgradeContainerd(ctx context.Context) error {
	upg := upgrade.NewUpgrader()
	if err := upg.ContainerdUpgrade(ctx); err != nil {
		return fmt.Errorf("upgrading containerd on node: %v", err)
	}

	return nil
}

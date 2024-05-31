package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/upgrader/upgrade"
)

var upgradeCNIPluginsCmd = &cobra.Command{
	Use:   "cni-plugins",
	Short: "Upgrade cni-plugins",
	Long:  "Use upgrade cni-plugins command to upgrade cni plugins on the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := upgradeCNIPlugins(cmd.Context())
		if err != nil {
			log.Fatalf("upgrade cni-plugins failed: %v", err)
		}
		return nil
	},
}

func init() {
	upgradeCmd.AddCommand(upgradeCNIPluginsCmd)
}

func upgradeCNIPlugins(ctx context.Context) error {
	upg := upgrade.NewInPlaceUpgrader()
	if err := upg.CNIPluginsUpgrade(ctx); err != nil {
		return fmt.Errorf("upgrading cni-plugins on the node: %v", err)
	}
	return nil
}

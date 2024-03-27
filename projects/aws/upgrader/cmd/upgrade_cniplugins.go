package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/upgrader/upgrade"
	"github.com/spf13/cobra"
)

var upgradeCniPluginsCmd = &cobra.Command{
	Use:   "cni-plugins",
	Short: "Upgrade cni-plugins",
	Long:  "Use InPlace Upgrader upgrade cni-plugins to upgrade cni-plugins on the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := upgradeCniPlugins(cmd.Context())
		if err != nil {
			log.Fatalf("upgrade cni-plugins failed: %v", err)
		}
		return nil
	},
}

func init() {
	upgradeCmd.AddCommand(upgradeCniPluginsCmd)
}

func upgradeCniPlugins(ctx context.Context) error {
	upg := upgrade.NewUpgrader()
	if err := upg.CniPluginsUpgrade(ctx); err != nil {
		return fmt.Errorf("upgrading Cni-Plugins on the node: %v", err)
	}
	return nil
}

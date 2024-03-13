package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/upgrader/upgrade"
	"github.com/spf13/cobra"
)

var upgradeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Upgrade status",
	Long:  "Use upgrade status command to get status of upgrade components on the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := upgradeStatus(cmd.Context())
		if err != nil {
			log.Fatalf("upgrade status failed: %v", err)
		}
		return nil
	},
}

func init() {
	upgradeCmd.AddCommand(upgradeStatusCmd)
}

func upgradeStatus(ctx context.Context) error {
	upg := upgrade.NewInPlaceUpgrader()
	if err := upg.LogStatusAndCleanup(ctx); err != nil {
		return fmt.Errorf("fetching upgrade status on the node: %v", err)
	}

	return nil
}

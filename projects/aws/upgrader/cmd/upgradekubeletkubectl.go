package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/upgrader/upgrade"
)

var upgradeKubeletKubectlCmd = &cobra.Command{
	Use:   "kubelet-kubectl",
	Short: "Upgrade kubelet-kubectl",
	Long:  "Use upgrade kubelet-kubectl command to upgrade kubelet and kubectl on the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := upgradeKubeletAndKubectl(cmd.Context())
		if err != nil {
			log.Fatalf("upgrade kubelet-kubectl failed: %v", err)
		}
		return nil
	},
}

func init() {
	upgradeCmd.AddCommand(upgradeKubeletKubectlCmd)
}

func upgradeKubeletAndKubectl(ctx context.Context) error {
	upg := upgrade.NewInPlaceUpgrader()
	if err := upg.KubeletKubectlUpgrade(ctx); err != nil {
		return fmt.Errorf("upgrading kubelet and kubectl on node: %v", err)
	}

	return nil
}

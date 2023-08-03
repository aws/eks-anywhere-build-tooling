package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/commands/upgrade"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
)

var upgradeOptions = &types.UpgradeOptions{}

// upgradeCmd is the command used to upgrade versions for a particular project.
var upgradeCmd = &cobra.Command{
	Use:   "upgrade --project <project name>",
	Short: "Upgrade the version for a single project",
	Long:  "Use this command to upgrade the Git tag and related versions for a particular project in the EKS-A build-tooling repository",
	Run: func(cmd *cobra.Command, args []string) {
		err := upgrade.Run(upgradeOptions)
		if err != nil {
			log.Fatalf("Error upgrading project version: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().StringVar(&upgradeOptions.ProjectName, "project", "", "Specify the project name to upgrade versions for")
	upgradeCmd.Flags().BoolVar(&upgradeOptions.DryRun, "dry-run", false, "Upgrade the project locally but do not push changes and create PR")
	if err := upgradeCmd.MarkFlagRequired("project"); err != nil {
		log.Fatalf("Error marking flag %q as required: %v", "project", err)
	}
}

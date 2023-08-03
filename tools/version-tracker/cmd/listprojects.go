package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/commands/listprojects"
)

// listProjectsCmd is the command used to list projects built from the EKS-A build-tooling repository.
var listProjectsCmd = &cobra.Command{
	Use:   "list-projects",
	Short: "List the projects maintained in the eks-anywhere-build-tooling repository",
	Long:  "Use this command to list the upstream projects that are built from the eks-anywhere-build-tooling repository",
	Run: func(cmd *cobra.Command, args []string) {
		err := listprojects.Run()
		if err != nil {
			log.Fatalf("Error listing upstream projects: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listProjectsCmd)
}

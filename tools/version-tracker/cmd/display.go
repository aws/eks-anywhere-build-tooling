package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/commands/display"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
)

var displayOptions = &types.DisplayOptions{}

// displayCmd is the command used to display version information for projects.
var displayCmd = &cobra.Command{
	Use:   "display --project <project name>",
	Short: "Display the version information for one or all projects",
	Long:  "Use this command to display the version information for a particular project or for all projects",
	Run: func(cmd *cobra.Command, args []string) {
		err := display.Run(displayOptions)
		if err != nil {
			log.Fatalf("Error displaying version information: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)
	displayCmd.Flags().StringVar(&displayOptions.ProjectName, "project", "", "Specify the project name to track versions for")
	displayCmd.Flags().BoolVar(&displayOptions.PrintLatestVersion, "print-latest-version", false, "Flag to print only the latest version of the project")
}

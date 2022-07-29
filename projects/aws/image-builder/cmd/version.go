package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the image-builder cli version",
	Long:  "This command prints the version of image-builder cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		return printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion() error {
	fmt.Println(version)
	return nil
}

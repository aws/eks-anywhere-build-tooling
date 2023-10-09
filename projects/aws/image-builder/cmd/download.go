package cmd

import (
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads artifacts for airgapped builds",
	Long:  "Downloads manifests and other artifacts used to build EKS-A Node images in airgapped environment",
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

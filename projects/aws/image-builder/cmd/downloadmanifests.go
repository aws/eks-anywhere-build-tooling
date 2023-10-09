package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var downloadManifestsCmd = &cobra.Command{
	Use:   "manifests",
	Short: "Downloads manifests for airgapped builds",
	Long:  "Downloads EKS-D and EKS-A manifests used to build EKS-A Node images in airgapped environment",
	Run: func(cmd *cobra.Command, args []string) {
		err := bo.DownloadManifests()
		if err != nil {
			log.Fatalf("Error downloading manifests: %v", err)
		}
	},
}

func init() {
	downloadCmd.AddCommand(downloadManifestsCmd)
	downloadManifestsCmd.Flags().BoolVar(&bo.Force, "force", false, "Force flag to clean up leftover files from previous execution")
}

package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var downloadArtifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "Downloads artifacts for airgapped builds",
	Long:  "Downloads EKS-D and EKS-A artifacts used to build EKS-A Node images in airgapped environment",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate Inputs
		if bo.EKSAReleaseVersion == "" {
			log.Fatalf("Please set eksa-release flag. EKS-A release versions can be identified from https://anywhere.eks.amazonaws.com/docs/whatsnew/changelog/")
		}
		if bo.ReleaseChannel == "" {
			log.Fatalf("Please set EKS-D release-channel using --release-channel. Can be 1-24, 1-25, 1-26, 1-27 or 1-28")
		}
		err := bo.DownloadArtifacts()
		if err != nil {
			log.Fatalf("Error downloading manifests: %v", err)
		}
	},
}

func init() {
	downloadCmd.AddCommand(downloadArtifactsCmd)
	downloadArtifactsCmd.Flags().BoolVar(&bo.Force, "force", false, "Force flag to clean up leftover files from previous execution")
	downloadArtifactsCmd.Flags().StringVar(&bo.ReleaseChannel, "release-channel", "1-27", "EKS-D Release channel for node image. Can be 1-24, 1-25, 1-26, 1-27 or 1-28")
	downloadArtifactsCmd.Flags().StringVar(&bo.EKSAReleaseVersion, "eksa-release", "", "The EKS-A CLI version to build images for")
}

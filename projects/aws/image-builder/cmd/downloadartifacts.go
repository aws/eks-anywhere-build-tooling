package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var downloadArtifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "Downloads artifacts for airgapped builds",
	Long:  "Downloads EKS-D and EKS-A artifacts used to build EKS-A Node images in airgapped environment",
	Run: func(cmd *cobra.Command, args []string) {
		err := bo.DownloadArtifacts()
		if err != nil {
			log.Fatalf("Error downloading manifests: %v", err)
		}
	},
}

func init() {
	downloadCmd.AddCommand(downloadArtifactsCmd)
	downloadArtifactsCmd.Flags().BoolVar(&bo.Force, "force", false, "Force flag to clean up leftover files from previous execution")
	downloadArtifactsCmd.Flags().StringVar(&bo.ReleaseChannel, "release-channel", "1-31", "EKS-D Release channel for node image. Can be 1-28, 1-29, 1-30, 1-31, 1-32, 1-33 or 1-34")
	downloadArtifactsCmd.Flags().StringVar(&bo.EKSAReleaseVersion, "eksa-release", "", "The EKS-A CLI version to build images for")
	if err := downloadArtifactsCmd.MarkFlagRequired("eksa-release"); err != nil {
		log.Fatalf("Please set eksa-release flag. EKS-A release versions can be identified from https://anywhere.eks.amazonaws.com/docs/whatsnew/changelog/")
	}
	if err := downloadArtifactsCmd.MarkFlagRequired("release-channel"); err != nil {
		log.Fatalf("Please set EKS-D release-channel using --release-channel. Can be 1-28, 1-29, 1-30, 1-31, 1-32, 1-33 or 1-34")
	}
}

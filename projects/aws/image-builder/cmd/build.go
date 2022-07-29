package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/image-builder/builder"
)

const (
	artifactsBucket = "s3://projectbuildpipeline-857-pipelineoutputartifactsb-10ajmk30khe3f"
)

var bo = &builder.BuildOptions{}

var buildCmd = &cobra.Command{
	Use:   "build --os <image os> --hypervisor <target hypervisor>",
	Short: "Build EKS Anywhere Node Image",
	Long:  "This command is used to build EKS Anywhere node images",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating builder config")
		bo.ValidateInputs()
		bo.BuildImage()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&bo.Os, "os", "", "Operating system to use for EKS-A node image")
	buildCmd.Flags().StringVar(&bo.Hypervisor, "hypervisor", "", "Target hypervisor EKS-A node image")
	buildCmd.Flags().StringVar(&bo.VsphereConfig, "vsphere-config", "", "Path to vSphere Config file")
	buildCmd.Flags().StringVar(&bo.ReleaseChannel, "release-channel", "1-21", "EKS-D Release channel for node image. Can be 1-20, 1-21 or 1-22")
	buildCmd.Flags().StringVar(&bo.ArtifactsBucket, "artifacts-bucket", artifactsBucket, "S3 bucket storing all EKS-D and EKS-A Artifacts")
	buildCmd.Flags().BoolVar(&bo.Force, "force", false, "Force flag to clean up leftover files from previous execution")
	if err := buildCmd.MarkFlagRequired("os"); err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
	if err := buildCmd.MarkFlagRequired("hypervisor"); err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
	if err := buildCmd.MarkFlagRequired("release-channel"); err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
	// TODO validate vsphere config with hypervisor input
}

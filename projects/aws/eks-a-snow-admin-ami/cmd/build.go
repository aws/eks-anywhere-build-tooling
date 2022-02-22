package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/snow"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Snow EKS-A Admin AMI",
	SilenceUsage: true,
	RunE:         build,
}

var input = &snow.AdminAMIInput{}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVar(&input.EKSAVersion, "eksa-version", "", "EKS-A version to embed in the snow Admin AMI")
	buildCmd.Flags().StringVar(&input.EKSAReleaseURL, "eksa-release-manifest-url", "", "EKS-A release manifest URL to pull EKS-A version from")
}

var buildVersion = "v0.0"

func build(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	return snow.BuildAdminAMI(ctx, input)
}

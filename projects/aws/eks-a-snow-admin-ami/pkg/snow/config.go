package snow

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/components"
	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/pipelines"
	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/recipes"
)

func defaultSnowAdminAMIPipeline() *pipelines.Pipeline {
	return snowAdminAMIPipelineForEKSA(nil)
}

func snowAdminAMIPipelineForEKSA(input *AdminAMIInput) *pipelines.Pipeline {
	return &pipelines.Pipeline{
		Name:        "Snow EKS-A Admin AMI pipeline",
		Description: "Pipeline to build AMI to run EKS-A on Snow devices with the CAPAS provider",
		Recipe:      snowAdminAMIRecipe(input),
	}
}

func snowAdminAMIRecipe(input *AdminAMIInput) *recipes.Recipe {
	version := "0.0.0"
	description := "Base recipe for Snow EKS-A Admin AMI"
	if input != nil {
		version = strings.TrimPrefix(input.EKSAVersion, "v")
		description = fmt.Sprintf("Recipe for Snow EKS-A [%s] Admin AMI", input.EKSAVersion)
	}

	return &recipes.Recipe{
		Name:        "Snow EKS-A Admin AMI recipe",
		Description: description,
		ParentImage: "arn:aws:imagebuilder:{region}:aws:image/ubuntu-server-20-lts-x86/x.x.x",
		Version:     version,
		Components: []recipes.ComponentConfiguration{
			{
				Component: &components.External{
					Name:    "update-linux",
					Account: "aws",
				},
			},
			{
				Component: &components.External{
					Name:    "docker-ce-ubuntu",
					Account: "aws",
				},
			},
			{
				Component: &components.Component{
					Name:                "EKS Anywhere",
					YamlFilePath:        "components/install_eksa.yaml",
					SupportedOsVersions: []string{"Ubuntu 20"},
					Description:         "Install EKS Anywhere in Ubuntu",
					Platform:            "Linux",
				},
				Parameters: eksaComponentParemeters(input),
			},
			{
				Component: &components.Component{
					Name:                "Download EKS Anywhere artifacts",
					YamlFilePath:        "components/download_eksa_artifacts.yaml",
					SupportedOsVersions: []string{"Ubuntu 20"},
					Description:         "Download EKS Anywhere artifacts for disconnected use",
					Platform:            "Linux",
				},
			},
		},
	}
}

func eksaComponentParemeters(input *AdminAMIInput) []imagebuilder.ComponentParameter {
	if input == nil {
		return nil
	}

	var eksaComponentParemeters []imagebuilder.ComponentParameter

	if input.EKSAReleaseURL != "" {
		eksaComponentParemeters = append(eksaComponentParemeters, imagebuilder.ComponentParameter{
			Name:  aws.String("EksAnywhereReleaseUrl"),
			Value: aws.StringSlice([]string{input.EKSAReleaseURL}),
		})
	}

	if input.EKSAVersion != "" {
		eksaComponentParemeters = append(eksaComponentParemeters, imagebuilder.ComponentParameter{
			Name:  aws.String("EksAnywhereVersion"),
			Value: aws.StringSlice([]string{input.EKSAVersion}),
		})
	}

	return eksaComponentParemeters
}

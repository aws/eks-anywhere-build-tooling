package pipelines

import (
	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/arn"
	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/recipes"
)

type Pipeline struct {
	Name        string
	Description string
	Recipe      *recipes.Recipe
}

func (p *Pipeline) ARN(account, region string) string {
	return arn.ForImageBuilderObject(account, region, "image-pipeline", p.Name)
}

func (p *Pipeline) ValidNameForARN() string {
	return arn.NameForARN(p.Name)
}

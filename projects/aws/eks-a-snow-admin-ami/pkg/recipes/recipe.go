package recipes

import (
	"context"

	"github.com/aws/aws-sdk-go/service/imagebuilder"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/arn"
	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

type Recipe struct {
	Name        string
	Description string
	ParentImage string
	Version     string
	Components  []ComponentConfiguration
}

type Component interface {
	Deploy(ctx context.Context, session *session.Session) error
	LastVersionARN(account, region string) string
}

type ComponentConfiguration struct {
	Component
	Parameters []imagebuilder.ComponentParameter
}

func (r *Recipe) ARN(account, region string) string {
	return arn.ForVersionImageBuilderObject(account, region, "image-recipe", r.Name, r.Version)
}

func (c *ComponentConfiguration) parameterToAPI() []*imagebuilder.ComponentParameter {
	if len(c.Parameters) == 0 {
		return nil
	}

	parameters := make([]*imagebuilder.ComponentParameter, 0, len(c.Parameters))
	for i := range c.Parameters {
		parameters = append(parameters, &c.Parameters[i])
	}

	return parameters
}

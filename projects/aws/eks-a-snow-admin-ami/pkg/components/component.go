package components

import (
	"os"

	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/arn"
)

type Component struct {
	Name                string
	Description         string
	SupportedOsVersions []string
	Platform            string
	YamlFilePath        string
}

func (c *Component) Read() ([]byte, error) {
	componentYamlDefinition, err := os.ReadFile(c.YamlFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read component %s at %s", c.Name, c.YamlFilePath)
	}

	return componentYamlDefinition, nil
}

func (c *Component) LastVersionARN(account, region string) string {
	return arn.ForLastVersionImageBuilderObject(account, region, "component", c.Name)
}

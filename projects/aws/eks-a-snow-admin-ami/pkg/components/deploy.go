package components

import (
	"context"
	"log"

	"github.com/Masterminds/semver/v3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

func (c *Component) Deploy(ctx context.Context, session *session.Session) error {
	log.Printf("Deploying component [%s]\n", c.Name)
	builder := imagebuilder.New(session)

	componentARN := c.LastVersionARN(session.Account, session.Region())
	log.Printf("Checking in component [%s] exists\n", componentARN)
	componentResponse, err := builder.GetComponentWithContext(ctx, &imagebuilder.GetComponentInput{
		ComponentBuildVersionArn: aws.String(componentARN),
	})

	if isNotFound(err) {
		log.Printf("Component [%s] does not exist yet\n", c.Name)
		return c.create(ctx, builder, "1.0.0")
	}

	if err != nil {
		return errors.Wrapf(err, "failed checking if component [%s] exists", c.Name)
	}

	currentVersion := *componentResponse.Component.Version

	log.Printf("Component [%s] exists on latest version %s\n", c.Name, currentVersion)
	componentData, err := c.Read()
	if err != nil {
		return err
	}

	if string(componentData) == *componentResponse.Component.Data {
		log.Printf("Component [%s] hasn't changed, skipping deployment", c.Name)
		return nil
	}

	log.Printf("Component [%s] has changed\n", c.Name)

	nextPatchVersion, err := bumpPatchVersion(currentVersion)
	if err != nil {
		return errors.Wrapf(err, "found component with invalid version %s", currentVersion)
	}

	return c.create(ctx, builder, nextPatchVersion)
}

func (c *Component) create(ctx context.Context, builder *imagebuilder.Imagebuilder, version string) error {
	log.Printf("Creating component [%s] with version %s\n", c.Name, version)
	componentData, err := c.Read()
	if err != nil {
		return err
	}

	_, err = builder.CreateComponentWithContext(ctx, &imagebuilder.CreateComponentInput{
		Name:                aws.String(c.Name),
		Data:                aws.String(string(componentData)),
		SemanticVersion:     aws.String(version),
		Platform:            aws.String(c.Platform),
		SupportedOsVersions: aws.StringSlice(c.SupportedOsVersions),
	})
	if err != nil {
		return errors.Wrapf(err, "failed creating component %s version %s", c.Name, version)
	}

	return nil
}

func isNotFound(err error) bool {
	e := &imagebuilder.ResourceNotFoundException{}
	return errors.As(err, &e)
}

func bumpPatchVersion(version string) (string, error) {
	parsed, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	newVersion := parsed.IncPatch()
	return newVersion.String(), nil
}

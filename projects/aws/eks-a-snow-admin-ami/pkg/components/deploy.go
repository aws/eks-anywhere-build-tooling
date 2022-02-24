package components

import (
	"context"
	"log"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

const (
	maxDeplotRetries = 5
	backoffTime      = 10 * time.Second
)

var componentAlreadyExistsErr = errors.New("component already exists")

func (c *Component) Deploy(ctx context.Context, session *session.Session) error {
	var err error
	for retries := 1; retries <= maxDeplotRetries; retries++ {
		if err = c.deploy(ctx, session); errors.Is(err, componentAlreadyExistsErr) {
			log.Printf("Failed to deploy component [%s] because it already exists. This might me transient if the list API hasn't reflected the latest component version yet.", c.Name)
			log.Printf("Total retries %d. Sleeping for %s", retries, backoffTime)
			time.Sleep(backoffTime)
			continue
		}

		break
	}

	return err
}

func (c *Component) deploy(ctx context.Context, session *session.Session) error {
	log.Printf("Deploying component [%s]\n", c.Name)
	builder := imagebuilder.New(session)

	log.Printf("Checking if component [%s] exists, this might take a while\n", c.Name)

	componentLatestVersion, err := c.getLatestVersionFromList(ctx, session)
	if err != nil {
		return errors.Wrapf(err, "failed checking if component [%s] exists", c.Name)
	}

	if componentLatestVersion == nil {
		log.Printf("Component [%s] does not exist yet\n", c.Name)
		return c.create(ctx, builder, "1.0.0")
	}

	currentVersion := *componentLatestVersion.Version

	log.Printf("Component [%s] exists on latest version %s\n", c.Name, currentVersion)
	componentData, err := c.Read()
	if err != nil {
		return err
	}

	if string(componentData) == *componentLatestVersion.Data {
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
	if isAlreadyExist(err) {
		return componentAlreadyExistsErr
	}

	if err != nil {
		return errors.Wrapf(err, "failed creating component %s version %s", c.Name, version)
	}

	return nil
}

func (c *Component) getLatestVersionWithXXXARN(ctx context.Context, session *session.Session) (*imagebuilder.Component, error) {
	builder := imagebuilder.New(session)
	componentARN := c.LastVersionARN(session.Account, session.Region())
	componentResponse, err := builder.GetComponentWithContext(ctx, &imagebuilder.GetComponentInput{
		ComponentBuildVersionArn: aws.String(componentARN),
	})

	if isNotFound(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed checking if component [%s] exists", c.Name)
	}
	return componentResponse.Component, nil
}

// This method is slower than getLatestVersionWithXXXARN but its APIs calls refelect new components faster
func (c *Component) getLatestVersionFromList(ctx context.Context, session *session.Session) (*imagebuilder.Component, error) {
	componentVersions, err := c.getAllComponentVersions(ctx, session)
	if err != nil {
		return nil, err
	}

	var latestSemver *semver.Version
	var latestComponent *imagebuilder.ComponentVersion

	for i := range componentVersions {
		comp := componentVersions[i]
		version, err := semver.NewVersion(*comp.Version)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid version %s for component [%s]", *comp.Version, *comp.Arn)
		}

		if latestSemver == nil || latestSemver.LessThan(version) {
			latestSemver = version
			latestComponent = &comp
		}
	}

	if latestComponent == nil {
		return nil, nil
	}

	builder := imagebuilder.New(session)
	componentResponse, err := builder.GetComponentWithContext(ctx, &imagebuilder.GetComponentInput{
		ComponentBuildVersionArn: latestComponent.Arn,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting latest version for component [%s]", c.Name)
	}

	return componentResponse.Component, nil
}

func (c *Component) getAllComponentVersions(ctx context.Context, session *session.Session) ([]imagebuilder.ComponentVersion, error) {
	builder := imagebuilder.New(session)
	var components []imagebuilder.ComponentVersion

	hasNextPage := true
	var nextToken *string
	for hasNextPage {
		componentsList, err := builder.ListComponentsWithContext(ctx, &imagebuilder.ListComponentsInput{
			Filters: []*imagebuilder.Filter{
				{
					Name:   aws.String("name"),
					Values: aws.StringSlice([]string{c.Name}),
				},
			},
			NextToken: nextToken,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed listing all versions for component [%s]", c.Name)
		}

		for _, component := range componentsList.ComponentVersionList {
			components = append(components, *component)
		}

		nextToken = componentsList.NextToken
		hasNextPage = componentsList.NextToken != nil
	}

	return components, nil
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

func isAlreadyExist(err error) bool {
	e := &imagebuilder.ResourceAlreadyExistsException{}
	return errors.As(err, &e)
}

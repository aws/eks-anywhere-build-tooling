package recipes

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

func (r *Recipe) Create(ctx context.Context, session *session.Session) (arn string, err error) {
	componentsConfigurations := make([]*imagebuilder.ComponentConfiguration, 0, len(r.Components))
	for _, c := range r.Components {
		componentsConfigurations = append(componentsConfigurations, &imagebuilder.ComponentConfiguration{
			ComponentArn: aws.String(c.LastVersionARN(session.Account, session.Region())),
			Parameters:   c.parameterToAPI(),
		})
	}

	builder := imagebuilder.New(session)

	log.Printf("Creating recipe [%s] version [%s]\n", r.Name, r.Version)
	_, err = builder.CreateImageRecipeWithContext(ctx, &imagebuilder.CreateImageRecipeInput{
		Name:            aws.String(r.Name),
		Description:     aws.String(r.Description),
		Components:      componentsConfigurations,
		ParentImage:     aws.String(session.ARNForRegion(r.ParentImage)),
		SemanticVersion: aws.String(r.Version),
	})
	if err != nil && !isAlreadyExist(err) {
		return "", errors.Wrapf(err, "failed creating recipe %s", r.Name)
	}

	return r.ARN(session.Account, session.Region()), nil
}

func isAlreadyExist(err error) bool {
	e := &imagebuilder.ResourceAlreadyExistsException{}
	return errors.As(err, &e)
}

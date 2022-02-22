package pipelines

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

func (p *Pipeline) Deploy(ctx context.Context, session *session.Session) error {
	log.Printf("Starting pipeline [%s] deployment\n", p.Name)
	infraConfigARN, err := p.setupInfraConfig(ctx, session)
	if err != nil {
		return err
	}

	recipeArn, err := p.Recipe.Deploy(ctx, session)
	if err != nil {
		return err
	}

	builder := imagebuilder.New(session)

	log.Printf("Creating pipeline [%s]\n", p.Name)
	_, err = builder.CreateImagePipelineWithContext(ctx, &imagebuilder.CreateImagePipelineInput{
		Name:                           aws.String(p.Name),
		Description:                    aws.String(p.Description),
		ImageRecipeArn:                 aws.String(recipeArn),
		InfrastructureConfigurationArn: aws.String(infraConfigARN),
	})
	if err != nil && !isAlreadyExist(err) {
		return errors.Wrapf(err, "failed creating pipeline %s", p.Name)
	}

	return nil
}

func isAlreadyExist(err error) bool {
	e := &imagebuilder.ResourceAlreadyExistsException{}
	return errors.As(err, &e)
}

package pipelines

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

func (p *Pipeline) UpdateRecipe(ctx context.Context, session *session.Session) error {
	log.Printf("Starting pipeline [%s] update\n", p.Name)
	recipeArn, err := p.Recipe.Create(ctx, session)
	if err != nil {
		return err
	}

	builder := imagebuilder.New(session)

	pipeline, err := builder.GetImagePipeline(&imagebuilder.GetImagePipelineInput{
		ImagePipelineArn: aws.String(p.ARN(session.Account, session.Region())),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to get image pipeline [%s]", p.Name)
	}

	log.Printf("Updating pipeline [%s] with recipe [%s]\n", p.Name, recipeArn)
	_, err = builder.UpdateImagePipelineWithContext(ctx, &imagebuilder.UpdateImagePipelineInput{
		ImagePipelineArn:               pipeline.ImagePipeline.Arn,
		InfrastructureConfigurationArn: pipeline.ImagePipeline.InfrastructureConfigurationArn,
		ImageRecipeArn:                 &recipeArn,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to get update pipeline [%s]", p.Name)
	}

	return nil
}

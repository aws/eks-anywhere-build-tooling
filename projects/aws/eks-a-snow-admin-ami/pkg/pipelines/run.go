package pipelines

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

func (p *Pipeline) Run(ctx context.Context, session *session.Session) error {
	builder := imagebuilder.New(session)

	pipelineARN := p.ARN(session.Account, session.Region())

	buildResponse, err := builder.StartImagePipelineExecutionWithContext(ctx, &imagebuilder.StartImagePipelineExecutionInput{
		ImagePipelineArn: &pipelineARN,
	},
	)
	if err != nil {
		return errors.Wrapf(err, "failed to start image pipeline [%s]", p.Name)
	}

	log.Printf("Build Pipeline for image %s started", *buildResponse.ImageBuildVersionArn)

	var image *imagebuilder.GetImageOutput
	for {
		image, err = builder.GetImageWithContext(ctx, &imagebuilder.GetImageInput{
			ImageBuildVersionArn: buildResponse.ImageBuildVersionArn,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to get image pipeline [%s]", p.Name)
		}

		if isFinalStatus(*image.Image.State.Status) {
			break
		}
		log.Printf("Pipeline [%s] current status: %s\n", p.Name, *image.Image.State.Status)
		time.Sleep(30 * time.Second)
	}

	if !isStatusSuccess(*image.Image.State.Status) {
		return errors.Errorf("AMI pipeline build failed with status %s", *image.Image.State.Status)
	}

	log.Printf("Build Pipeline for image [%s] finshed successfully", p.Name)

	return nil
}

var finalStatuses = map[string]struct{}{
	imagebuilder.ImageStatusAvailable:  {},
	imagebuilder.ImageStatusCancelled:  {},
	imagebuilder.ImageStatusFailed:     {},
	imagebuilder.ImageStatusDeprecated: {},
	imagebuilder.ImageStatusDeleted:    {},
}

func isFinalStatus(status string) bool {
	_, ok := finalStatuses[status]
	return ok
}

func isStatusSuccess(status string) bool {
	return status == imagebuilder.ImageStatusAvailable
}

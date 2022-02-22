package snow

import (
	"context"
	"log"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

type AdminAMIInput struct {
	EKSAVersion    string
	EKSAReleaseURL string
}

func BuildAdminAMI(ctx context.Context, input *AdminAMIInput) error {
	log.Printf("Building AMI for EKSA %s from manifest [%s]\n", input.EKSAVersion, input.EKSAReleaseURL)
	pipeline := snowAdminAMIPipelineForEKSA(input)

	session, err := session.New(ctx)
	if err != nil {
		return err
	}

	if err := pipeline.UpdateRecipe(ctx, session); err != nil {
		return err
	}

	return pipeline.Run(ctx, session)
}

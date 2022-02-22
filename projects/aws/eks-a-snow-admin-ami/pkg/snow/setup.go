package snow

import (
	"context"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

func SetupAdminAMIPipeline(ctx context.Context) error {
	pipeline := defaultSnowAdminAMIPipeline()

	session, err := session.New(ctx)
	if err != nil {
		return err
	}

	return pipeline.Deploy(ctx, session)
}

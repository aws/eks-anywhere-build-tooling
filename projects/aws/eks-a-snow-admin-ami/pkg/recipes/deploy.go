package recipes

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

func (r *Recipe) Deploy(ctx context.Context, session *session.Session) (arn string, err error) {
	log.Printf("Starting recipe [%s] deployment\n", r.Name)
	for _, c := range r.Components {
		if err = c.Deploy(ctx, session); err != nil {
			return "", errors.Wrapf(err, "failed creating deploying component for recipe [%s]", r.Name)
		}
	}

	return r.Create(ctx, session)
}

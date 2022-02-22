package components

import (
	"context"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/arn"
	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

type External struct {
	Name    string
	Account string
}

func (e *External) Deploy(ctx context.Context, session *session.Session) error {
	return nil
}

func (e *External) LastVersionARN(account, region string) string {
	return arn.ForLastVersionImageBuilderObject(e.Account, region,  "component", e.Name)
}

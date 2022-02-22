package session

import (
	"context"

	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/arn"
)

type Session struct {
	*awsSession.Session
	Account string
}

func (s *Session) Region() string {
	return *s.Config.Region
}

func (s *Session) ARNForRegion(baseARN string) string {
	return arn.ForRegion(baseARN, s.Region())
}

func New(ctx context.Context) (*Session, error) {
	awsSession, err := awsSession.NewSession()
	if err != nil {
		return nil, err
	}
	stsService := sts.New(awsSession)
	identityResponse, err := stsService.GetCallerIdentityWithContext(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return &Session{
		Session: awsSession,
		Account: *identityResponse.Account,
	}, nil
}

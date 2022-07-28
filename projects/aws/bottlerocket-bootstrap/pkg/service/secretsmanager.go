package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"
)

type SecretsManagerService interface {
	GetSecretValue(ctx context.Context, secretName string) (*secretsmanager.GetSecretValueOutput, error)
	DeleteSecret(ctx context.Context, secretName string) (*secretsmanager.DeleteSecretOutput, error)
}

type SecretsManagerImpl struct {
	BaseClient *secretsmanager.Client
}

func NewSecretsManagerService() (SecretsManagerService, error) {
	clientImpl := new(SecretsManagerImpl)
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	imdsClient := imds.NewFromConfig(cfg)
	getRegionOutput, err := imdsClient.GetRegion(context.TODO(), &imds.GetRegionInput{})
	if err != nil {
		return nil, errors.Wrap(err, "Unable to retrieve the region from the EC2 instance")
	}
	cfg.Region = getRegionOutput.Region
	clientImpl.BaseClient = secretsmanager.NewFromConfig(cfg)
	return clientImpl, nil
}

func (client SecretsManagerImpl) GetSecretValue(ctx context.Context, secretName string) (*secretsmanager.GetSecretValueOutput, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}
	return client.BaseClient.GetSecretValue(ctx, input)
}

func (client SecretsManagerImpl) DeleteSecret(ctx context.Context, secretName string) (*secretsmanager.DeleteSecretOutput, error) {
	deleteSecretInput := &secretsmanager.DeleteSecretInput{
		SecretId: aws.String(secretName),
	}
	return client.BaseClient.DeleteSecret(ctx, deleteSecretInput)
}

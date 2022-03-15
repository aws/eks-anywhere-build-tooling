package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerService interface {
	GetSecretValue(ctx context.Context, secretName string) (*secretsmanager.GetSecretValueOutput, error)
	DeleteSecret(ctx context.Context, secretName string) (*secretsmanager.DeleteSecretOutput, error)
}

type SecretsManagerImpl struct {
	BaseClient *secretsmanager.Client
}

func NewSecretsManagerService() SecretsManagerService {
	clientImpl := new(SecretsManagerImpl)
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	cfg.Region = "us-west-2"
	clientImpl.BaseClient = secretsmanager.NewFromConfig(cfg)
	return clientImpl
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

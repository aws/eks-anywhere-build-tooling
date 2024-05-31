package utils

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang/mock/gomock"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/service"
)

// Normal UserData
const UserDataString = `
## template: jinja
#cloud-config

write_files:
-   path: /var/lib/kubeadm/pki/ca.crt
    owner: root:root
    permissions: '0640'
    content: |
        -----BEGIN CERTIFICATE-----
        ...
        -----END CERTIFICATE-----
runcmd: "ControlPlaneInit"
`

// References userdata stored in AWS SecretsManager
const AWSSecrentsManagerDataString = `
user_data_type: "AWSSecretsManager"

secrets_manager_data:
    prefix: some-prefix
    chunks: 1
`

func TestNormalUserData(t *testing.T) {
	processedUserData, err := processUserData([]byte(UserDataString))
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	if processedUserData.RunCmd != "ControlPlaneInit" {
		t.Errorf("Unexpected RunCmd: Expected: %s, Actual: %s", "ControlPlaneInit", processedUserData.RunCmd)
	}
}

func TestWithAWSSecretsManagerUserData(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	mockSecretsManagerService := service.NewMockSecretsManagerService(ctrl)
	base64UserData := base64.StdEncoding.EncodeToString([]byte(UserDataString))
	compressedUserData, _ := GzipBytes([]byte(base64UserData))
	getSecretValueOutput := secretsmanager.GetSecretValueOutput{}
	getSecretValueOutput.SecretBinary = compressedUserData

	mockSecretsManagerService.EXPECT().GetSecretValue(gomock.Any(), "some-prefix-0").Return(&getSecretValueOutput, nil)
	mockSecretsManagerService.EXPECT().DeleteSecret(gomock.Any(), "some-prefix-0").Return(&secretsmanager.DeleteSecretOutput{}, nil)

	processedUserData, err := processAWSSecretsManagerUserData([]byte(AWSSecrentsManagerDataString), mockSecretsManagerService)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	if processedUserData.RunCmd != "ControlPlaneInit" {
		t.Errorf("Unexpected RunCmd: Expected: %s, Actual: %s", "ControlPlaneInit", processedUserData.RunCmd)
	}
}

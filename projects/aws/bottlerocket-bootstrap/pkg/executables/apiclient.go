package executables

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/utils"
)

type APIClient struct {
	Executable
}

type APISetting struct {
	Kernel     *Kernel     `json:"kernel,omitempty"`
	Kubernetes *Kubernetes `json:"kubernetes,omitempty"`
}

type Kernel struct {
	Sysctl map[string]string `json:"sysctl,omitempty"`
}

type Kubernetes struct {
	AllowedUnsafeSysctls []string `json:"allowed-unsafe-sysctls,omitempty"`
}

func NewAPIClient() *APIClient {
	return &APIClient{
		Executable: NewExecutable(utils.ApiclientBinary),
	}
}

func (a *APIClient) SetKubernetesCloudProvider(cloudProvider string) error {
	_, err := a.Execute("set", fmt.Sprintf("kubernetes.cloud-provider=%q", cloudProvider))
	return err
}

func (a *APIClient) SetKubernetesNodeIP(ip string) error {
	_, err := a.Execute("set", "kubernetes.node-ip="+ip)
	return err
}

func (a *APIClient) SetKubernetesProviderID(id string) error {
	_, err := a.Execute("set", "kubernetes.provider-id="+id)
	return err
}

func (a *APIClient) Set(setting *APISetting) error {
	config, err := json.Marshal(setting)
	if err != nil {
		return errors.Wrap(err, "error json marshalling api setting")
	}
	_, err = a.Execute("set", "--json", string(config))
	return err
}

func (a *APIClient) Reboot() error {
	_, err := a.Execute("reboot")
	return err
}

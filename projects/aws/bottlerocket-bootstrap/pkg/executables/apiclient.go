package executables

import (
	"fmt"
	"strings"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
)

type APIClient struct {
	Executable
}

func NewAPIClient() *APIClient {
	return &APIClient{
		Executable: NewExecutable(utils.ApiclientBinary),
	}
}

func (a *APIClient) SetKubernetesCloudProvider(cloudProvider string) error {
	_, err := a.Execute("set", fmt.Sprintf("kubernetes.cloud-providers=%q", cloudProvider))
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

func (a *APIClient) SetKubernetesAllowUnsafeSysctls(sysctls []string) error {
	q := make([]string, len(sysctls))
	for idx, ctl := range sysctls {
		q[idx] = fmt.Sprintf("%q", ctl)
	}
	_, err := a.Execute("set", fmt.Sprintf("allowed-unsafe-sysctls=[%s]", strings.Join(q[:], ",")))
	return err
}

func (a *APIClient) SetKernelRmemMax(size string) error {
	_, err := a.Execute("set", "kernel.sysctl.net.core.rmem_max="+size)
	return err
}

func (a *APIClient) SetKernelWmemMax(size string) error {
	_, err := a.Execute("set", "kernel.sysctl.net.core.wmem_max="+size)
	return err
}

func (a *APIClient) SetKernelCorePattern(pattern string) error {
	_, err := a.Execute("set", fmt.Sprintf("kernel.core_pattern=%q", pattern))
	return err
}

func (a *APIClient) Reboot() error {
	_, err := a.Execute("reboot")
	return err
}

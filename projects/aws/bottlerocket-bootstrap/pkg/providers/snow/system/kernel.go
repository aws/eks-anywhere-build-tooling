package system

import (
	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/executables"
)

const (
	// Temporary hardcoded kernel config values for beta.
	// In the future the kernel config will be exposed for eks-a users to customize.
	maxSocketBufferSize = "8388608"
	kernelCorePattern   = "/var/corefile/core.%e.%p.%h.%t"
)

func (s *Snow) configureKernelSettings() error {
	kernelSetting := &executables.APISetting{
		Kubernetes: &executables.Kubernetes{
			AllowedUnsafeSysctls: []string{
				"net.ipv4.tcp_mtu_probing",
			},
		},
		Kernel: &executables.Kernel{
			Sysctl: map[string]string{
				"net.core.rmem_max":   maxSocketBufferSize,
				"net.core.wmem_max":   maxSocketBufferSize,
				"kernel.core_pattern": kernelCorePattern,
			},
		},
	}

	return s.api.Set(kernelSetting)
}

package system

import "github.com/pkg/errors"

const (
	// Temporary hardcoded kernel config values for beta.
	// In the future the kernel config will be exposed for eks-a users to customize.
	maxSocketBufferSize = "8388608"
	kernelCorePattern   = "/var/corefile/core.%e.%p.%h.%t"
)

func (s *Snow) configureKernelSettings() error {
	if err := s.api.SetKernelRmemMax(maxSocketBufferSize); err != nil {
		return errors.Wrap(err, "Error setting kernel maximum socket receive buffer size")
	}

	if err := s.api.SetKernelWmemMax(maxSocketBufferSize); err != nil {
		return errors.Wrap(err, "Error setting kernel maximum socket send buffer size")
	}

	if err := s.api.SetKernelCorePattern(kernelCorePattern); err != nil {
		return errors.Wrap(err, "Error setting kernel core pattern")
	}

	return nil
}

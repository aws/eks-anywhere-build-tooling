//go:generate ../../../../hack/tools/bin/mockgen -destination ./mocks/apiclient_mock.go -package mocks . APIClient
//go:generate ../../../../hack/tools/bin/mockgen -destination ./mocks/filesystem_mock.go -package mocks . FileSystem

package system

import (
	"fmt"
	"os"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/executables"
)

type APIClient interface {
	SetKubernetesCloudProvider(cloudProvider string) error
	SetKubernetesNodeIP(ip string) error
	SetKubernetesProviderID(id string) error
	Set(setting *executables.APISetting) error
	Reboot() error
}

type FileSystem interface {
	MountVolume(device, dir string) error
	Partition(device string) error
}

type Snow struct {
	api APIClient
	fs  FileSystem
}

func NewSnow() *Snow {
	return &Snow{
		api: executables.NewAPIClient(),
		fs:  executables.NewFileSystem(),
	}
}

func (s *Snow) Bootstrap() {
	if err := s.configureKubernetesSettings(); err != nil {
		fmt.Printf("Error configuring Kubernetes settings: %v\n", err)
		os.Exit(1)
	}

	if err := s.configureKernelSettings(); err != nil {
		fmt.Printf("Error configuring kernel settings: %v\n", err)
		os.Exit(1)
	}

	if err := configureDNI(); err != nil {
		fmt.Printf("Error configuring snow DNI: %v\n", err)
		os.Exit(1)
	}

	if err := s.mountContainerdVolume(); err != nil {
		fmt.Printf("Error configuring container volume: %v\n", err)
		os.Exit(1)
	}

	if err := s.rebootInstanceIfNeeded(); err != nil {
		fmt.Printf("Error rebooting instance: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("SUCCESS: Snow bootstrap tasks finished.")
	os.Exit(0)
}

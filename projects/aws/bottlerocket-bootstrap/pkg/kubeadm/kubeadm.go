package kubeadm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
	"github.com/pkg/errors"
)

const (
	kubeadmFile                      = "/tmp/kubeadm.yaml"
	kubectl                          = "/opt/bin/kubectl"
	kubeadmBinary                    = "/opt/bin/kubeadm"
	kubeconfigPath                   = "/etc/kubernetes/admin.conf"
	kubeadmJoinFile                  = "/tmp/kubeadm-join-config.yaml"
	kubeletConfigFile                = "/var/lib/kubelet/config.yaml"
	staticPodManifestsPath           = "/etc/kubernetes/manifests"
	bottlerocketRootFSKubeadmPKIPath = "/.bottlerocket/rootfs/var/lib/kubeadm/pki"
)

type kubeadm struct {
	userData *utils.UserData
}

func New(userData *utils.UserData) *kubeadm {
	return &kubeadm{userData: userData}
}

func (k *kubeadm) InitializeDirectories() error {
	fmt.Println("Initializing directories")
	err := os.MkdirAll(bottlerocketRootFSKubeadmPKIPath, 0o640)
	if err != nil {
		return errors.Wrap(err, "error creating directory")
	}
	// Force symlink creation with cmd. sdk calls fail if the symlink dir exists
	cmd := exec.Command("bash", "-c", "ln -sfn /.bottlerocket/rootfs/var/lib/kubeadm /var/lib")
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error running command: %v", cmd)
	}
	cmd = exec.Command("bash", "-c", "ln -sfn /.bottlerocket/rootfs/var/lib/kubeadm /etc/kubernetes")
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error running command: %v", cmd)
	}
	return nil
}

func (k *kubeadm) RunCmd() error {
	// Take different directions based on runCmd on node's user data
	switch k.userData.RunCmd {
	case "ControlPlaneInit":
		fmt.Println("Running controlplane init bootstrap sequence")
		if err := controlPlaneInit(); err != nil {
			return errors.Wrapf(err, "error initing controlplane")
		}
	case "ControlPlaneJoin":
		fmt.Println("Running controlplane join sequence")
		if err := controlPlaneJoin(); err != nil {
			return errors.Wrapf(err, "error joining controlplane")
		}
	case "WorkerJoin":
		fmt.Println("Running worker join sequence")
		if err := workerJoin(); err != nil {
			return errors.Wrapf(err, "error joining as worker")
		}
	}

	return nil
}

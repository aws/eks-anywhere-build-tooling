package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/files"
	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
	"github.com/pkg/errors"
)

func configureKubernetesSettings() error {
	if err := setCloudProvider(); err != nil {
		return errors.Wrap(err, "Error setting K8s cloud provider")
	}

	if err := setProviderID(); err != nil {
		return errors.Wrap(err, "Error setting K8s provider id")
	}

	if err := setNodeIP(); err != nil {
		return errors.Wrap(err, "Error setting K8s node ip")
	}

	return nil
}

func setCloudProvider() error {
	cmd := exec.Command(utils.ApiclientBinary, "set", `kubernetes.cloud-provider=""`)
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "Error running command: %v, output: %s\n", cmd, string(out))
	}
	return nil
}

func setNodeIP() error {
	nodeIPConfiguredPath := filepath.Join(bootstrapContainerPath, "nodeip.configured")

	// only set K8s node ip once after reboot
	if !files.PathExists(rebootedPath) || files.PathExists(nodeIPConfiguredPath) {
		return nil
	}

	ip, err := currentIP()
	if err != nil {
		return errors.Wrap(err, "Error getting current ip address")
	}

	cmd := exec.Command(utils.ApiclientBinary, "set", "kubernetes.node-ip="+ip)
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "Error running command: %v, output: %s\n", cmd, string(out))
	}

	// Write node ip configured file to make sure the node ip is only configured once
	if _, err := os.Create(nodeIPConfiguredPath); err != nil {
		return errors.Wrapf(err, "Error writing %s file", nodeIPConfiguredPath)
	}

	return nil
}

func setProviderID() error {
	id, err := instanceID()
	if err != nil {
		errors.Wrap(err, "error getting instance id")
	}

	cmd := exec.Command(utils.ApiclientBinary, "set", "kubernetes.provider-id=aws-snow:////"+id)
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "Error running command: %v, output: %s\n", cmd, string(out))
	}

	return nil
}

func instanceID() (string, error) {
	url := fmt.Sprintf("http://%s/latest/meta-data/instance-id", metadataServiceIP)
	body, err := httpGet(url)
	if err != nil {
		return "", errors.Wrap(err, "error requesting instance id through http")
	}

	return string(body), nil
}

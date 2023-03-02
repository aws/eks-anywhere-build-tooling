package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/files"
	"github.com/pkg/errors"
)

func (s *Snow) configureKubernetesSettings() error {
	if err := s.setCloudProvider(); err != nil {
		return errors.Wrap(err, "Error setting K8s cloud provider")
	}

	if err := s.setProviderID(); err != nil {
		return errors.Wrap(err, "Error setting K8s provider id")
	}

	if err := s.setNodeIP(); err != nil {
		return errors.Wrap(err, "Error setting K8s node ip")
	}

	return nil
}

func (s *Snow) setCloudProvider() error {
	return s.api.SetKubernetesCloudProvider("")
}

func (s *Snow) setNodeIP() error {
	nodeIPConfiguredPath := filepath.Join(bootstrapContainerPath, "nodeip.configured")

	// only set K8s node ip once after reboot
	if !files.PathExists(rebootedPath) || files.PathExists(nodeIPConfiguredPath) {
		return nil
	}

	ip, err := currentIP()
	if err != nil {
		return errors.Wrap(err, "Error getting current ip address")
	}

	if err := s.api.SetKubernetesNodeIP(ip); err != nil {
		return err
	}

	// Write node ip configured file to make sure the node ip is only configured once
	if _, err := os.Create(nodeIPConfiguredPath); err != nil {
		return errors.Wrapf(err, "Error writing %s file", nodeIPConfiguredPath)
	}

	return nil
}

func (s *Snow) setProviderID() error {
	id, err := instanceID()
	if err != nil {
		errors.Wrap(err, "error getting instance id")
	}
	ip, err := deviceIP()
	if err != nil {
		errors.Wrap(err, "error getting device ip")
	}
	return s.api.SetKubernetesProviderID(fmt.Sprintf("aws-snow:///%s/%s", ip, id))
}

func instanceID() (string, error) {
	url := fmt.Sprintf("http://%s/latest/meta-data/instance-id", metadataServiceIP)
	body, err := httpGet(url)
	if err != nil {
		return "", errors.Wrap(err, "error requesting instance id through http")
	}

	return string(body), nil
}

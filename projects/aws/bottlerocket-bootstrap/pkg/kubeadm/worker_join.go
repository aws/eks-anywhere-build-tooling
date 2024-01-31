package kubeadm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
	"github.com/pkg/errors"
)

func workerJoin() error {
	err := setHostName(kubeadmJoinFile)
	if err != nil {
		return errors.Wrap(err, "Error replacing hostname on kubeadm")
	}

	cmd := exec.Command(kubeadmBinary, utils.LogVerbosity, "join", "phase", "kubelet-start", "--config", kubeadmJoinFile)
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}

	err = killCmdAfterJoinFilesGeneration(cmd)
	if err != nil {
		return errors.Wrap(err, "Error waiting for worker join files")
	}

	apiServer, token, err := getBootstrapFromJoinConfig(kubeadmJoinFile)
	if err != nil {
		return errors.Wrap(err, "Error getting api server and token from kubeadm join config")
	}
	b64CA, err := getEncodedCA()
	if err != nil {
		return errors.Wrap(err, "Error reading the ca data")
	}
	dns, err := getDNSFromJoinConfig(kubeletConfigFile)
	if err != nil {
		return errors.Wrap(err, "Error getting dns from kubelet config")
	}
	fmt.Println(dns)

	cmd = exec.Command(utils.ApiclientBinary, "set",
		"kubernetes.api-server="+apiServer,
		"kubernetes.cluster-certificate="+b64CA,
		"kubernetes.cluster-dns-ip="+string(dns),
		"kubernetes.bootstrap-token="+token,
		"kubernetes.authentication-mode=tls",
		"kubernetes.standalone-mode=false")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Error running command: %v, output: %s\n", cmd, string(out))
	}
	fmt.Printf("Ran apiclient set call: %v\n", string(out))

	// wait for kubelet to come up before killing the bootstrap container
	err = waitForActiveKubelet()
	if err != nil {
		fmt.Printf("Error waiting for kubelet: %v", err)
		return err
	}
	return nil
}

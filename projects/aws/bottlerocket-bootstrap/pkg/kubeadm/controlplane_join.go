package kubeadm

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
	"github.com/pkg/errors"
)

const kubeconfigPath = "/etc/kubernetes/admin.conf"

func controlPlaneJoin() error {
	err := setHostName(kubeadmJoinFile)
	if err != nil {
		return errors.Wrap(err, "Error replacing hostname on kubeadm")
	}

	cmd := exec.Command(kubeadmBinary, utils.LogVerbosity, "join", "phase", "control-plane-prepare", "all", "--config", kubeadmJoinFile)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}

	cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "join", "phase", "kubelet-start", "--config", kubeadmJoinFile)
	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}

	err = killCmdAfterJoinFilesGeneration(cmd)
	if err != nil {
		return errors.Wrap(err, "Error waiting for worker join files")
	}

	dns, err := getDNSFromJoinConfig("/var/lib/kubelet/config.yaml")
	if err != nil {
		return errors.Wrap(err, "Error getting api server")
	}

	apiServer, token, err := getBootstrapFromJoinConfig(kubeadmJoinFile)
	if err != nil {
		return errors.Wrap(err, "Error getting api server")
	}

	// Get CA
	b64CA, err := getEncodedCA()
	if err != nil {
		return errors.Wrap(err, "Error reading the ca data")
	}

	cmd = exec.Command(utils.ApiclientBinary, "set",
		"kubernetes.api-server="+apiServer,
		"kubernetes.cluster-certificate="+b64CA,
		"kubernetes.cluster-dns-ip="+dns,
		"kubernetes.bootstrap-token="+token,
		"kubernetes.authentication-mode=tls",
		"kubernetes.standalone-mode=false")
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}

	err = waitForActiveKubelet()
	if err != nil {
		return errors.Wrap(err, "Error waiting for kubelet to come up")
	}

	if isEtcdExternal, err := isClusterWithExternalEtcd(kubeconfigPath); err != nil {
		return err
	} else if !isEtcdExternal {
		if err := joinLocalEtcd(); err != nil {
			return err
		}
	}

	// Migrate all static pods from this host-container to the bottlerocket host using the apiclient
	// now that etcd manifest is also created
	podDefinitions, err := utils.EnableStaticPods("/etc/kubernetes/manifests")
	if err != nil {
		return errors.Wrap(err, "Error enabling static pods")
	}

	// Now that etcd is up and running, check for other pod liveness
	err = utils.WaitForPods(podDefinitions)
	if err != nil {
		return errors.Wrapf(err, "Error waiting for static pods to be up")
	}

	// Get port number from apiServer host string
	port, err := getLocalApiBindPortFromJoinConfig(kubeadmJoinFile)
	if err != nil {
		// if we hit an error when extracting the port number, we always fallback to 6443
		fmt.Printf("unable to get local apiserver port, falling back to 6443. caused by: %s", err.Error())
		port = 6443
	}

	// Wait for Kubernetes API server to come up.
	err = utils.WaitFor200(fmt.Sprintf("https://localhost:%d/healthz", port), 30*time.Second)
	if err != nil {
		return err
	}

	err = utils.WaitFor200(string(apiServer)+"/healthz", 30*time.Second)
	if err != nil {
		return err
	}

	// finish kubeadm
	cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "join", "--skip-phases", "preflight,control-plane-prepare,kubelet-start,control-plane-join/etcd",
		"--config", kubeadmJoinFile)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}
	return nil
}

func joinLocalEtcd() error {
	// Get kubeadm to write out the manifest for etcd.
	// It will wait for etcd to start, which won't succeed because we need to set the static-pods in the BR api.
	cmd := exec.Command(kubeadmBinary, utils.LogVerbosity, "join", "phase", "control-plane-join", "etcd", "--config", kubeadmJoinFile)
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}
	etcdCheckFiles := []string{"/etc/kubernetes/manifests/etcd.yaml"}
	if err := utils.KillCmdAfterFilesGeneration(cmd, etcdCheckFiles); err != nil {
		return errors.Wrap(err, "Error waiting for etcd manifest files")
	}

	return nil
}

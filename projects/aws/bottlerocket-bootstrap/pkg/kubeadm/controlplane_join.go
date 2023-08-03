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

	// start optional EBS initialization
	ebsInitControl := startEbsInit()

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

	args := []string{
		"set",
		"kubernetes.api-server=" + apiServer,
		"kubernetes.cluster-certificate=" + b64CA,
		"kubernetes.cluster-dns-ip=" + dns,
		"kubernetes.bootstrap-token=" + token,
		"kubernetes.authentication-mode=tls",
		"kubernetes.standalone-mode=false",
	}

	kubeletTlsConfig := readKubeletTlsConfig(&RealFileReader{})
	if kubeletTlsConfig != nil {
		args = append(args, "settings.kubernetes.server-certificate="+kubeletTlsConfig.KubeletServingCert)
		args = append(args, "settings.kubernetes.server-key="+kubeletTlsConfig.KubeletServingPrivateKey)
	}

	cmd = exec.Command(utils.ApiclientBinary, args...)
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

	// Wait for Kubernetes API server to come up.
	localApiServerReadinessEndpoint, err := getLocalApiServerReadinessEndpoint()
	if err != nil {
		fmt.Printf("unable to get local apiserver readiness endpoint, falling back to localhost:6443. caused by: %s", err.Error())
		localApiServerReadinessEndpoint = "https://localhost:6443/healthz"
	}

	err = utils.WaitFor200(localApiServerReadinessEndpoint, 5*time.Minute)
	if err != nil {
		return err
	}

	err = utils.WaitFor200(string(apiServer)+"/healthz", 5*time.Minute)
	if err != nil {
		return err
	}

	// finish kubeadm
	cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "join", "--skip-phases", "preflight,control-plane-prepare,kubelet-start,control-plane-join/etcd",
		"--config", kubeadmJoinFile)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}

	if ebsInitControl != nil {
		checkEbsInit(ebsInitControl)
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

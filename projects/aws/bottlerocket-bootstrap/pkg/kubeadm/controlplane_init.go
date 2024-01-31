package kubeadm

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	versionutil "k8s.io/apimachinery/pkg/util/version"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
	"github.com/pkg/errors"
)

func controlPlaneInit() error {
	err := setHostName(kubeadmFile)
	if err != nil {
		return errors.Wrap(err, "Error replacing hostname on kubeadm")
	}

	// start optional EBS initialization
	ebsInitControl := startEbsInit()

	cmd := exec.Command(kubeadmBinary, "version", "-o=short")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Error running command: %v, Output: %s\n", cmd, string(out))
	}
	kubeadmVersion, err := versionutil.ParseSemantic(strings.TrimSuffix(string(out), "\n"))
	if err != nil {
		return errors.Wrapf(err, "%s is not a valid kubeadm version", string(out))
	}

	k8s129Compare, err := kubeadmVersion.Compare("1.29.0")
	if err != nil {
		return errors.Wrap(err, "Error comparing kubeadm version with v1.29.0")
	}

	if k8s129Compare != -1 {
		err = patchKubeVipManifest()
		if err != nil {
			return errors.Wrapf(err, "Error patching kube-vip manifest")
		}
	}

	// Generate keys and write all the manifests
	fmt.Println("Running kubeadm init commands")
	cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "init", "phase", "certs", "all", "--config", kubeadmFile)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Error running command, out: %s", string(out))
	}
	fmt.Printf("Running command: %v\n, output: %s\n", cmd, string(out))
	cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "init", "phase", "kubeconfig", "all", "--config", kubeadmFile)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}
	fmt.Printf("Running command: %v\n, output: %s\n", cmd, string(out))
	cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "init", "phase", "control-plane", "all", "--config", kubeadmFile)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}
	fmt.Printf("Running command: %v\n, output: %s\n", cmd, string(out))
	cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "init", "phase", "etcd", "local", "--config", kubeadmFile)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}
	fmt.Printf("Running command: %v\n, output: %s\n", cmd, string(out))

	// Migrate all static pods from host-container to bottlerocket host using apiclient
	podDefinitions, err := utils.EnableStaticPods(staticPodManifestsPath)
	if err != nil {
		return errors.Wrap(err, "Error enabling static pods")
	}

	// Wait for all static pods liveness probe to be up
	err = utils.WaitForPods(podDefinitions)
	if err != nil {
		return errors.Wrapf(err, "Error waiting for static pods to be up")
	}

	// Get server from admin.conf
	apiServer, err := utils.GetApiServerFromKubeConfig(kubeconfigPath)
	if err != nil {
		return errors.Wrap(err, "Error getting api server")
	}
	fmt.Printf("APIServer is %s\n", apiServer)

	// Get CA
	b64CA, err := getEncodedCA()
	if err != nil {
		return errors.Wrap(err, "Error reading the ca data")
	}

	localApiServerReadinessEndpoint, err := getLocalApiServerReadinessEndpoint()
	if err != nil {
		fmt.Printf("unable to get local apiserver readiness endpoint, falling back to localhost:6443. caused by: %s", err.Error())
		localApiServerReadinessEndpoint = "https://localhost:6443/healthz"
	}

	// Wait for Kubernetes API server to come up.
	err = utils.WaitFor200(localApiServerReadinessEndpoint, 30*time.Second)
	if err != nil {
		return err
	}

	// If the api advertise url is different than localhost, like when using kube-vip, make
	// sure it is accessible
	err = utils.WaitFor200(string(apiServer)+"/healthz", 30*time.Second)
	if err != nil {
		return err
	}

	// Set up the roles so our kubelet can bootstrap.
	cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "init", "phase", "bootstrap-token", "--config", kubeadmFile)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}

	token, err := getBootstrapToken()
	if err != nil {
		return errors.Wrap(err, "Error getting token")
	}
	fmt.Printf("Bootstrap token is: %s\n", token)
	// token string already has escaped quotes
	args := []string{
		"set",
		"kubernetes.api-server=" + apiServer,
		"kubernetes.cluster-certificate=" + b64CA,
		"kubernetes.bootstrap-token=" + string(token),
		"kubernetes.authentication-mode=tls",
		"kubernetes.standalone-mode=false",
	}

	kubeletTlsConfig := readKubeletTlsConfig(&RealFileReader{})
	if kubeletTlsConfig != nil {
		args = append(args, "settings.kubernetes.server-certificate="+kubeletTlsConfig.KubeletServingCert)
		args = append(args, "settings.kubernetes.server-key="+kubeletTlsConfig.KubeletServingPrivateKey)
	}

	cmd = exec.Command(utils.ApiclientBinary, args...)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return errors.Wrapf(err, "Error running command: %v", cmd)
	}

	err = waitForActiveKubelet()
	if err != nil {
		return errors.Wrap(err, "Error waiting for kubelet to come up")
	}

	// we compare the kubeadm version to v1.26 because a new phase "show-join-command" was introduced in that version.
	// this comparison can be removed when we deprecate v1.25.
	k8s126Compare, err := kubeadmVersion.Compare("1.26.0")
	if err != nil {
		return errors.Wrap(err, "Error comparing kubeadm version with v1.26.0")
	}

	// finish kubeadm
	if k8s126Compare == -1 {
		cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "init", "--skip-phases", "preflight,kubelet-start,certs,kubeconfig,bootstrap-token,control-plane,etcd",
			"--config", kubeadmFile)
	} else {
		cmd = exec.Command(kubeadmBinary, utils.LogVerbosity, "init", "--skip-phases", "preflight,kubelet-start,certs,kubeconfig,bootstrap-token,control-plane,etcd,show-join-command",
			"--config", kubeadmFile)
	}
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Error running command: %v, Output: %s\n", cmd, string(out))
	}

	// Now that Core DNS is installed, find the cluster DNS IP.
	dns, err := getDNS(kubeconfigPath)
	if err != nil {
		return errors.Wrap(err, "Error getting dns ip")
	}

	// set dns
	cmd = exec.Command(utils.ApiclientBinary, "set", "kubernetes.cluster-dns-ip="+string(dns))
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Error running command: %v, output: %s\n", cmd, string(out))
	}

	if ebsInitControl != nil {
		checkEbsInit(ebsInitControl)
	}
	return nil
}

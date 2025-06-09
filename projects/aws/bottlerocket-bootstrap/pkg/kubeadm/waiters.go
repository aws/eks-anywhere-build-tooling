package kubeadm

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/utils"
)

func waitForKubelet() error {
	err := utils.WaitForSystemdService(utils.KubeletService, 90*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// waitForActiveKubelet checks for the kubelet status and then sleeps for 20 seconds and
// checks again. We wait for 20 seconds to check the status just to make sure the kubelet
// hasn't crashed for some reason. There can be other reasons kubelet isn't working as expected
// but has also not crashed.
func waitForActiveKubelet() error {
	err := waitForKubelet()
	if err != nil {
		return errors.Wrap(err, "Error waiting for kubelet to come up")
	}

	// Wait for 20 seconds and check for kubelet status again
	// to check if the service has not crashed
	time.Sleep(20 * time.Second)
	err = waitForKubelet()
	if err != nil {
		return errors.Wrap(err, "Error waiting for kubelet to come up")
	}
	return nil
}

func killCmdAfterJoinFilesGeneration(cmd *exec.Cmd) error {
	checkFiles := []string{
		"/var/lib/kubelet/config.yaml",
		"/var/lib/kubelet/kubeadm-flags.env",
		"/var/lib/kubeadm/pki/ca.crt",
	}
	return utils.WaitForManifestAndOptionallyKillCmd(cmd, checkFiles, true)
}

// checkEbsInit checks the execution of ebs-init goroutine,
// the goroutine will be stopped if timeout is reached
func checkEbsInit(ctrl *EbsInitControl) {
	select {
	case <-ctrl.Timeout:
		ctrl.Cancel()
		fmt.Printf("Killing ebs-init, timeout reached \n")
		return
	case <-ctrl.OkChan:
		fmt.Printf("Finished ebs-init \n")
		return
	}
}

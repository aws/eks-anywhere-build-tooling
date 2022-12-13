package kubeadm

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
	"github.com/pkg/errors"
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
	return utils.KillCmdAfterFilesGeneration(cmd, checkFiles)
}

// checkEbsInit checks the execution of ebs-init script,
// script is killed if timeout is reached
func checkEbsInit(ctrl *EbsInitControl) {
	okChan := make(chan bool)

	// wait for script execution in another goroutine
	go func(okChan chan bool, ctrl *EbsInitControl) {

		ctrl.Command.Wait()
		// send OK when script exits
		okChan <- true
		return

	}(okChan, ctrl)

	select {
	case <-ctrl.Timeout:
		ctrl.Command.Process.Kill()
		fmt.Printf("Killing ebs-init, timeout reached \n")
		return
	case <-okChan:
		fmt.Printf("Finished ebs-init \n")
		return
	}
}

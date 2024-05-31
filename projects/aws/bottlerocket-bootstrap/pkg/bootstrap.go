package pkg

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/etcdadm"
	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/kubeadm"
	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/utils"
)

const (
	marker = "/.bottlerocket/host-containers/" + utils.BootstrapContainerName + "/.ran"
)

func acquireLock() {
	if _, err := os.Stat(marker); err == nil {
		// Lock cannot be acquired, another instance of bootstrap is running
		fmt.Println("Cannot acquire lock, another instance of bootstrap is already running")
		err = utils.DisableBootstrapContainer()
		if err != nil {
			fmt.Println("Failed to run command, set bootstrapContainer to false")
		}
		time.Sleep(100000)
	} else {
		// Create file to indicate lock acquisition
		os.Create(marker)
		fmt.Println("Acquired lock for bootstrap")
	}
}

func waitForPreReqs() error {
	fmt.Println("Waiting for bottlerocket boot to complete...")
	err := utils.WaitForSystemdService(utils.MultiUserTarget, 5*time.Minute)
	if err != nil {
		return errors.Wrapf(err, "Error waiting for multi-user.target\n")
	}

	err = utils.WaitForSystemdService(utils.KubeletService, 30*time.Second)
	if err != nil {
		return errors.Wrapf(err, "Error waiting for kubelet")
	}

	fmt.Println("Bottlerocket bootstrap pre-requisites complete")
	return nil
}

func Bootstrap() {
	err := waitForPreReqs()
	if err != nil {
		fmt.Printf("Error waiting for bottlerocket pre-reqs: %v", err)
		os.Exit(1)
	}

	fmt.Println("Initiating bottlerocket bootstrap")
	acquireLock()
	userData, err := utils.ResolveHostContainerUserData()
	if err != nil {
		fmt.Printf("Error parsing user-data: %v\n", err)
		os.Exit(1)
	}

	b := buildBootstrapper(userData)

	if err = b.InitializeDirectories(); err != nil {
		fmt.Printf("Error initializing directories: %v\n", err)
		os.Exit(1)
	}

	if err = utils.WriteUserDataFiles(userData); err != nil {
		fmt.Printf("Error writing files from user-data: %v\n", err)
		os.Exit(1)
	}

	if err = b.RunCmd(); err != nil {
		fmt.Printf("Error running bootstrapper cmd: %v\n", err)
		os.Exit(1)
	}

	// Bootstrapping done
	err = utils.DisableBootstrapContainer()
	if err != nil {
		fmt.Printf("Error disabling bootstrap container: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Bottlerocket bootstrap was successful. Disabled bootstrap container")
	os.Exit(0)
}

type bootstrapper interface {
	InitializeDirectories() error
	RunCmd() error
}

func buildBootstrapper(userData *utils.UserData) bootstrapper {
	if strings.HasPrefix(userData.RunCmd, etcdadm.CmdPrefix) {
		fmt.Println("Using etcdadm support by CAPI")
		return etcdadm.New(userData)
	} else {
		fmt.Println("Using kubeadm support by CAPI")
		return kubeadm.New(userData)
	}
}

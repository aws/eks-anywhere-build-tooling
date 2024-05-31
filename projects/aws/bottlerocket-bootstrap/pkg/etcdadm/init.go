package etcdadm

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/utils"
)

var initPreKubeletPhases = []string{
	"install",
	"certificates",
	"snapshot",
	"configure",
	"start",
}

var initPostKubeletPhases = []string{"health"}

type initCommand struct {
	repository   string
	version      string
	cipherSuites string
}

func (i *initCommand) run() error {
	flags := buildFlags(i.repository, i.version, i.cipherSuites)
	fmt.Println("Running etcdadm init phases")
	if err := runPhases("init", initPreKubeletPhases, flags); err != nil {
		return err
	}

	fmt.Println("Starting etcd static pods")
	podDefinitions, err := utils.EnableStaticPods(podSpecDir)
	if err != nil {
		return errors.Wrap(err, "error enabling etcd static pods")
	}

	fmt.Println("Waiting for etcd static pods")
	err = utils.WaitForPods(podDefinitions)
	if err != nil {
		return errors.Wrapf(err, "error waiting for etcd static pods to be up")
	}

	if err := runPhases("init", initPostKubeletPhases, flags); err != nil {
		return err
	}

	return nil
}

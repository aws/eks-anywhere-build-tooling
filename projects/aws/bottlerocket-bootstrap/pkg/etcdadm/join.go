package etcdadm

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/utils"
)

var joinPreKubeletPhases = []string{
	"stop",
	"certificates",
	"membership",
	"install",
	"configure",
	"start",
}

var joinPostKubeletPhases = []string{"health"}

type joinCommand struct {
	repository   string
	version      string
	endpoint     string
	cipherSuites string
}

func (j *joinCommand) run() error {
	flags := buildFlags(j.repository, j.version, j.cipherSuites)
	fmt.Println("Running etcdadm join phases")
	if err := runPhases("join", joinPreKubeletPhases, flags, j.endpoint); err != nil {
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

	if err := runPhases("join", joinPostKubeletPhases, flags, j.endpoint); err != nil {
		return err
	}

	return nil
}

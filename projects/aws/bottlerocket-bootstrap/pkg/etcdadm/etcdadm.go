package etcdadm

import (
	"os"
	"strings"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
	"github.com/pkg/errors"
)

const (
	CmdPrefix         = "Etcdadm"
	manifestFileName  = "etcd.manifest"
	rootfsEtcdBaseDir = "/.bottlerocket/rootfs/var/lib/etcd"
	etcdBaseDir       = "/var/lib/etcd"
	certFolder        = "/pki"
	dataFolder        = "/data"
	certDir           = etcdBaseDir + certFolder
	dataDir           = etcdBaseDir + dataFolder
	rootfsCertDir     = rootfsEtcdBaseDir + certFolder
	rootfsDataDir     = rootfsEtcdBaseDir + dataFolder
	podSpecDir        = "./manifests"

	initCmd = "EtcdadmInit"
	joinCmd = "EtcdadmJoin"

	etcdadmBinary = "/opt/bin/etcdadm"
)

var dirs = []string{podSpecDir, rootfsDataDir, rootfsCertDir}

type etcdadm struct {
	userData *utils.UserData
}

func New(userData *utils.UserData) *etcdadm {
	return &etcdadm{userData: userData}
}

func (e *etcdadm) InitializeDirectories() error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0640); err != nil {
			return errors.Wrapf(err, "error creating etcdadm directory [%s]", dir)
		}
	}

	if err := utils.CreateSymLink(rootfsEtcdBaseDir, etcdBaseDir); err != nil {
		return errors.Wrap(err, "failed init symlinks for etcdadm")
	}

	return nil
}

func (e *etcdadm) RunCmd() error {
	cmd, err := parseCmd(e.userData.RunCmd)
	if err != nil {
		return err
	}

	if err = cmd.run(); err != nil {
		return err
	}

	return nil
}

type command interface {
	run() error
}

func parseCmd(bootstrapCmd string) (command, error) {
	words := strings.Fields(bootstrapCmd)
	if len(words) == 0 {
		return nil, errors.Errorf("invalid bootstrap etcdadm command [%s]", bootstrapCmd)
	}

	cmd := words[0]

	switch cmd {
	case initCmd:
		if len(words) != 4 {
			return nil, errors.Errorf("invalid bootstrap etcdadm init command [%s]", bootstrapCmd)
		}

		return &initCommand{repository: words[1], version: words[2], cipherSuites: words[3]}, nil
	case joinCmd:
		if len(words) != 5 {
			return nil, errors.Errorf("invalid bootstrap etcdadm join command [%s]", bootstrapCmd)
		}

		return &joinCommand{repository: words[1], version: words[2], cipherSuites: words[3], endpoint: words[4]}, nil
	default:
		return nil, errors.Errorf("invalid etcadm bootstrap command %s", bootstrapCmd)
	}
}

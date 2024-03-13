package upgrade

import (
	"fmt"
	"path/filepath"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

type Upgrader struct {
	SysCalls
	// optional fields
	kubernetesVersion string
	etcdVersion       string
}

type Option func(u Upgrader) Upgrader

func WithKubernetesVersion(version string) Option {
	return func(u Upgrader) Upgrader {
		u.kubernetesVersion = version
		return u
	}
}

func WithEtcdVersion(version string) Option {
	return func(u Upgrader) Upgrader {
		u.etcdVersion = version
		return u
	}
}

func (u *Upgrader) CurrDir() (string, error) {
	ex, err := u.Executable()
	if err != nil {
		return "", fmt.Errorf("unable to get current directory")
	}

	return filepath.Dir(ex), nil
}

func (u *Upgrader) upgradeComponentsDir() (string, error) {
	scriptDir, err := u.CurrDir()
	if err != nil {
		return "", err
	}

	return filepath.Dir(scriptDir), nil
}

func (u *Upgrader) upgradeComponentsBinDir() (string, error) {
	upgCmpDir, err := u.upgradeComponentsDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/binaries", upgCmpDir), nil
}

func (u *Upgrader) upgradeComponentsKubernetesBinDir() (string, error) {
	upgCompBinDir, err := u.upgradeComponentsBinDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/kubernetes/usr/bin", upgCompBinDir), nil
}

func NewUpgrader(options ...Option) Upgrader {
	upg := Upgrader{}
	upg.SysCalls = NewSysCalls()
	for _, opt := range options {
		upg = opt(upg)
	}

	return upg
}

func (u *Upgrader) BackUpAndReplace(oldFile, backUpFolder, newFile string) error {
	fileName := filepath.Base(oldFile)
	backedUpFile := filepath.Join(backUpFolder, fmt.Sprintf("%s.bk", fileName))
	if _, err := u.Stat(backedUpFile); err == nil {
		logger.Info("BackUp File already found, skipping backup")
		return nil
	}

	if err := u.copy(oldFile, backedUpFile); err != nil {
		return copyError(oldFile, backedUpFile, err)
	}

	if err := u.copy(newFile, oldFile); err != nil {
		return copyError(newFile, oldFile, err)
	}
	logger.Info("BackUp Success", "File", oldFile, "BackedUpFile", backedUpFile)

	return nil
}

func (u *Upgrader) copy(src, dst string) error {
	data, err := u.ReadFile(src)
	if err != nil {
		return err
	}
	if err := u.WriteFile(dst, data, 0o640); err != nil {
		return err
	}

	return nil
}

func copyError(src, dst string, err error) error {
	return fmt.Errorf("copying file from %s to %s: %v", src, dst, err)
}

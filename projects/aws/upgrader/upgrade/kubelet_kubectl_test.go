package upgrade_test

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	kubeletConf    = "/etc/sysconfig/kubelet"
	kubectlBinPath = "/usr/bin/kubectl"
	kubeletBinPath = "/usr/bin/kubelet"
	fileMode416    = fs.FileMode(0o640)
	fileMode493    = fs.FileMode(0o755)
)

func Test125KubeletKubectlUpgradeBackupExist(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version", "-oshort").Return([]byte("v1.25.x"), nil)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "restart", "kubelet").Return(nil, nil).Times(1)
	err := tt.u.KubeletKubectlUpgrade(ctx)
	tt.Expect(err).To(BeNil())
}

func TestNot125KubeletKubectlUpgradeBackupExist(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version", "-oshort").Return([]byte("v1.26.x"), nil)
	tt.s.EXPECT().Stat(kubeletConf).Return(nil, os.ErrNotExist)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "restart", "kubelet").Return(nil, nil).Times(1)
	err := tt.u.KubeletKubectlUpgrade(ctx)
	tt.Expect(err).To(BeNil())
}

func Test125KubeletKubectlUpgradeBackupDoesNotExist(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	kubectlBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")
	kubeletBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")

	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(kubectlBkFile).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubectlBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubectl")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeletBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeletBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubelet")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeletBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version", "-oshort").Return([]byte("v1.25.x"), nil)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "restart", "kubelet").Return(nil, nil).Times(1)
	err := tt.u.KubeletKubectlUpgrade(ctx)
	tt.Expect(err).To(BeNil())
}

func TestNot125KubeletKubectlUpgradeBackupDoesNotExist(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	kubectlBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")
	kubeletBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")

	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubectlBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubectl")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeletBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeletBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubelet")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeletBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version", "-oshort").Return([]byte("v1.26.x"), nil)
	tt.s.EXPECT().Stat(kubeletConf).Return(nil, os.ErrNotExist)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "restart", "kubelet").Return(nil, nil).Times(1)
	tt.u.KubeletKubectlUpgrade(ctx)
}

func TestNot125KubeletKubectlUpgradeKubeletConfExist(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	kubectlBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")
	kubeletBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")
	oldKubeletConf := []byte("kubelet conf --feature-gates=KubeletCredentialProviders=true")
	newKubeletConf := []byte("kubelet conf")

	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubectlBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubectl")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeletBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeletBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubelet")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeletBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version", "-oshort").Return([]byte("v1.26.x"), nil)
	tt.s.EXPECT().Stat(kubeletConf).Return(nil, nil)
	tt.s.EXPECT().ReadFile(kubeletConf).Return(oldKubeletConf, nil)
	tt.s.EXPECT().MkdirAll(fmt.Sprintf("%s/%s", upgCompBinDir, "extra_args"), fileMode493).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeletConf).Return(oldKubeletConf, nil)
	tt.s.EXPECT().WriteFile(fmt.Sprintf("%s/%s/%s", upgCompBinDir, "extra_args", "kubelet.bk"), oldKubeletConf, fileMode416)
	tt.s.EXPECT().WriteFile(kubeletConf, newKubeletConf, fileMode416)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "restart", "kubelet").Return(nil, nil).Times(1)
	tt.u.KubeletKubectlUpgrade(ctx)
}

func TestKubeletKubectlUpgradeKubectlBackupError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubectlBinPath).Return(nil, errors.New("")).Times(1)
	err := tt.u.KubeletKubectlUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeletKubectlUpgradeStopKubeletCmdError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	kubectlBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")

	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubectlBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubectl")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeletKubectlUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeletKubectlUpgradeKubeletBackupError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	kubectlBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")

	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubectlBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubectl")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeletBinPath).Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeletKubectlUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeletKubectlUpgradeKubeadmVersionError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	kubectlBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")

	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubectlBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubectl")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version", "-oshort").Return([]byte("v1.26.x"), errors.New(""))

	err := tt.u.KubeletKubectlUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeletKubectlUpgradeDaemonReloadError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	kubectlBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")

	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubectlBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubectl")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version", "-oshort").Return([]byte("v1.25.x"), nil)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, errors.New("")).Times(1)
	err := tt.u.KubeletKubectlUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeletKubectlUpgradeKubeletRestartError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	kubectlBkFile := fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")

	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubectlBinPath).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBkFile, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(fmt.Sprintf("%s/%s", upgCompBinDir, "kubectl")).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubectlBinPath, []byte{}, fileMode416).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version", "-oshort").Return([]byte("v1.25.x"), nil)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, errors.New("")).Times(1)
	err := tt.u.KubeletKubectlUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

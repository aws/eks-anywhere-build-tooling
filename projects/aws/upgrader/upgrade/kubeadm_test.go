package upgrade_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	upgrade "github.com/aws/eks-anywhere-build-tooling/upgrader/upgrade"
)

const (
	kubeadmClusterConfigFileName = "kubeadm-cluster-config.yaml"
	kubeletConfigFileName        = "kubelet-config.yaml"
	staticKubeVipPath            = "/etc/kubernetes/manifests/kube-vip.yaml"
	kubeAdmBinDir                = "/usr/bin/kubeadm"
	clusterConfig                = `
	apiServer:
		extraArgs:
			tls-cipher-suites: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	apiVersion: kubeadm.k8s.io/v1beta3
	certificatesDir: /etc/kubernetes/pki
	clusterName: dummy-tst
	controlPlaneEndpoint: 195.17.1.74:6443
	dns:
		imageRepository: public.ecr.aws/eks-distro/coredns
		imageTag: v1.9.3-eks-1-25-34
	etcd:
		local:
		dataDir: /var/lib/etcd
		imageRepository: public.ecr.aws/eks-distro/etcd-io
		imageTag: v3.5.10-eks-1-28-11
	imageRepository: public.ecr.aws/eks-distro/kubernetes
	kind: ClusterConfiguration
	kubernetesVersion: v1.29.1-eks-1-29-6
	`
	updatedClusterConfig = `
	apiServer:
		extraArgs:
			tls-cipher-suites: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	apiVersion: kubeadm.k8s.io/v1beta3
	certificatesDir: /etc/kubernetes/pki
	clusterName: dummy-tst
	controlPlaneEndpoint: 195.17.1.74:6443
	dns:
		imageRepository: public.ecr.aws/eks-distro/coredns
		imageTag: v1.9.3-eks-1-25-34
	etcd:
		local:
		dataDir: /var/lib/etcd
		imageRepository: public.ecr.aws/eks-distro/etcd-io
		imageTag: v3.5.10-eks-1-29-6
	imageRepository: public.ecr.aws/eks-distro/kubernetes
	kind: ClusterConfiguration
	kubernetesVersion: v1.29.1-eks-1-29-6
	`
	kubeletConfig = `
	apiVersion: kubelet.config.k8s.io/v1beta1
	cgroupDriver: systemd
	clusterDNS:
	- 10.96.0.10
	clusterDomain: cluster.local
	containerRuntimeEndpoint: ""
	kind: KubeletConfiguration
	logging:
		flushFrequency: 0
	resolvConf: /run/systemd/resolve/resolv.conf
	rotateCertificates: true
	staticPodPath: /etc/kubernetes/manifests
	`
	upgCompBinDir = "/foo/binaries/kubernetes/usr/bin"
	kubeVipBackup = "/foo/kube-vip.backup.yaml"
	newKubeVip    = "/foo/kube-vip.yaml"
	coreDNSBackup = "/foo/binaries/kubernetes/usr/bin/coredns.yaml"
)

func TestKubeAdmFirstCPBackupExist(t *testing.T) {
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte(updatedClusterConfig), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath).Return(nil, nil).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).To(BeNil())
}

func TestKubeAdmFirstCPBackupDoesNotExist(t *testing.T) {
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)

	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte(updatedClusterConfig), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath).Return(nil, nil).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).To(BeNil())
}

func TestKubeAdmFirstCPBackupError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPGetKCCError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)

	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return([]byte{}, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPKCCBackupError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)

	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPKCCBackupReadError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)

	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPEtcdUpdateError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPAppendKubeletReadError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(nil, errors.New(""))

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPAppendKubeletGetCmdError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPAppendKubeletWriteError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPGetCoreDNSCMError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPCoreDNSConfWriteError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPCoreDNSCMDeleteError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPKubeAdmVersionCmdError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPKubeAdmUpgPlanError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPKubeAdmUpgApplyError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPKubeVipReadError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPKubeVipBackupError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPNewKubeVipReadError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPNewKubeVipWriteError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmFirstCPRestoreCoreDNSError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New("")).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return([]byte("coredns-conf"), nil).Times(1)
	tt.s.EXPECT().WriteFile(coreDNSBackup, []byte("coredns-conf"), fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath).Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInFirstCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmRestCPBackupExists(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath).Return(nil, nil).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).To(BeNil())
}

func TestKubeAdmRestCPBackupNotExists(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)

	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath).Return(nil, nil).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).To(BeNil())
}

func TestKubeAdmRestCPBackupError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmRestCPBackupCoreDNSError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmRestCPKubeadmVersionCmdError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmRestCPKubeadmUpgradeCmdError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmRestCPKubeVipReadError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmRestCPKubeVipBackupError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmRestCPNewKubeVipReadError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmRestCPNewKubeVipWriteError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmRestCPRestoreCoreDNSError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(kubeAdmBackup).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath).Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInRestCP(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmWorkerBackUpExist(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node").Return(nil, nil).Times(1)

	err := tt.u.KubeAdmInWorker(ctx)
	tt.Expect(err).To(BeNil())
}

func TestKubeAdmWorkerBackUpNotExist(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node").Return(nil, nil).Times(1)

	err := tt.u.KubeAdmInWorker(ctx)
	tt.Expect(err).To(BeNil())
}

func TestKubeAdmWorkerBackUpError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(errors.New("")).Times(1)

	err := tt.u.KubeAdmInWorker(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmWorkerKubeadmVersionError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInWorker(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestKubeAdmWorkerKubeadmUpgradeError(t *testing.T) {
	kubeAdmBackup := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm.bk")
	newKubeAdm := fmt.Sprintf("%s/%s", upgCompBinDir, "kubeadm")
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, errors.New(""))
	tt.s.EXPECT().ReadFile(kubeAdmBinDir).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBackup, []byte{}, fileMode640).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdm).Return([]byte{}, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmBinDir, []byte{}, fileMode640).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node").Return(nil, errors.New("")).Times(1)

	err := tt.u.KubeAdmInWorker(ctx)
	tt.Expect(err).ToNot(BeNil())
}

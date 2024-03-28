package upgrade_test

import (
	"context"
	"fmt"
	"io/fs"
	"testing"

	upgrade "github.com/aws/eks-anywhere-build-tooling/projects/aws/upgrader/upgrade"
)

const (
	kubeadmClusterConfigFileName = "kubeadm-cluster-config.yaml"
	kubeletConfigFileName        = "kubelet-config.yaml"
	staticKubeVipPath            = "/etc/kubernetes/manifests/kube-vip.yaml"
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

func TestKubeAdmFistCP(t *testing.T) {
	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", upgCompBinDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", upgCompBinDir)
	clusterConfigBytes := []byte(clusterConfig)
	appendedKubeletConfigBytes := []byte(fmt.Sprintf("%s%s%s", clusterConfig, "---\n", kubeletConfig))
	ctx := context.TODO()
	tt := newUpgraderTest(t, upgrade.WithKubernetesVersion("v1.29.1-eks-1-29-6"), upgrade.WithEtcdVersion("v3.5.10-eks-1-29-6"))
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeAdmConfigBackUp, clusterConfigBytes, fs.FileMode(0o640)).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(kubeAdmConfigBackUp).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, []byte(updatedClusterConfig), fs.FileMode(0o640)).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeAdmConfig).Return(clusterConfigBytes, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath).Return([]byte(kubeletConfig), nil).Times(1)
	tt.s.EXPECT().WriteFile(newKubeAdmConfig, appendedKubeletConfigBytes, fs.FileMode(0o640)).Return(nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig).Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fs.FileMode(0o640)).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fs.FileMode(0o640)).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath).Return(nil, nil).Times(1)
	tt.u.KubeAdmInFirstCP(ctx)
}

func TestKubeAdmRestCP(t *testing.T) {
	ctx := context.TODO()
	tt := newUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration").Return(nil, nil).Times(1)
	tt.s.EXPECT().ReadFile(staticKubeVipPath).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(kubeVipBackup, nil, fs.FileMode(0o640)).Return(nil).Times(1)
	tt.s.EXPECT().ReadFile(newKubeVip).Return(nil, nil).Times(1)
	tt.s.EXPECT().WriteFile(staticKubeVipPath, nil, fs.FileMode(0o640)).Return(nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath).Return(nil, nil).Times(1)
	tt.u.KubeAdmInRestCP(ctx)
}

func TestKubeAdmWorker(t *testing.T) {
	ctx := context.TODO()
	tt := newUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).AnyTimes()
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubeadm")).Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "upgrade", "node").Return(nil, nil).Times(1)
	tt.u.KubeAdmInWorker(ctx)
}

package kubeadm

import (
	"testing"
)

const localEtcdClusterConf = `'apiServer:
  certSANs:
  - localhost
  - 127.0.0.1
  extraArgs:
    authorization-mode: Node,RBAC
    runtime-config: ""
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta2
certificatesDir: /etc/kubernetes/pki
clusterName: eksa-test-eks-a-cluster
controlPlaneEndpoint: eksa-test-eks-a-cluster-control-plane:6443
controllerManager:
  extraArgs:
    enable-hostpath-provisioner: "true"
dns:
  imageRepository: public.ecr.aws/eks-distro/coredns
  imageTag: v1.8.3-eks-1-20-5
  type: CoreDNS
etcd:
  local:
    dataDir: /var/lib/etcd
    imageRepository: public.ecr.aws/eks-distro/etcd-io
    imageTag: v3.4.15-eks-1-20-5
featureGates:
  IPv6DualStack: true
imageRepository: public.ecr.aws/eks-distro/kubernetes
kind: ClusterConfiguration
kubernetesVersion: v1.20.7-eks-1-20-5
networking:
  dnsDomain: cluster.local
  podSubnet: 10.244.0.0/16
  serviceSubnet: 10.96.0.0/16
scheduler: {}'`

const externalEtcdClusterConf = `'apiServer:
  extraArgs:
    cloud-provider: external
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta2
certificatesDir: /var/lib/kubeadm/pki
clusterName: eksa-test
controlPlaneEndpoint: 198.18.210.144:6443
controllerManager:
  extraArgs:
    cloud-provider: external
  extraVolumes:
  - hostPath: /var/lib/kubeadm/controller-manager.conf
    mountPath: /etc/kubernetes/controller-manager.conf
    name: kubeconfig
    pathType: File
    readOnly: true
dns:
  imageRepository: public.ecr.aws/eks-distro/coredns
  imageTag: v1.8.3-eks-1-20-5
  type: CoreDNS
etcd:
  external:
    caFile: /var/lib/kubeadm/pki/etcd/ca.crt
    certFile: /var/lib/kubeadm/pki/server-etcd-client.crt
    endpoints:
    - https://198.18.138.154:2379
    - https://198.18.138.155:2379
    - https://198.18.69.78:2379
    keyFile: /var/lib/kubeadm/pki/apiserver-etcd-client.key
imageRepository: public.ecr.aws/eks-distro/kubernetes
kind: ClusterConfiguration
kubernetesVersion: v1.20.7-eks-1-20-5
networking:
  dnsDomain: cluster.local
  podSubnet: 192.168.0.0/16
  serviceSubnet: 10.96.0.0/12
scheduler:
  extraVolumes:
  - hostPath: /var/lib/kubeadm/scheduler.conf
    mountPath: /etc/kubernetes/scheduler.conf
    name: kubeconfig
    pathType: File
    readOnly: true'`

func TestIsExternalEtcd(t *testing.T) {
	tests := []struct {
		testName             string
		clusterConfiguration []byte
		wantExternal         bool
	}{
		{
			testName:             "local",
			clusterConfiguration: []byte(localEtcdClusterConf),
			wantExternal:         false,
		},
		{
			testName:             "external",
			clusterConfiguration: []byte(externalEtcdClusterConf),
			wantExternal:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gotExternal, err := isExternalEtcd(tt.clusterConfiguration)
			if err != nil {
				t.Fatalf("isExternalEtcd() -> err = %v, want err = nil", err)
			}

			if gotExternal != tt.wantExternal {
				t.Fatalf("isExternalEtcd() -> gotExternal = %t, wantExternal = %t", gotExternal, tt.wantExternal)
			}
		})
	}
}

func TestGetLocalApiBindPortFromInitConfigYaml(t *testing.T) {

	const yamlWithoutInitConfigurationKind = `
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: '10.0.20.132'
  bindPort: 443
`
	t.Run("yamlWithMultipleDoc", func(t *testing.T) {
		const yamlWithMultipleDoc = `
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
apiServer:
  certSANs:
  - 127.0.0.1
  - '78cb78ba469f0d86d443ead9b58ab978.b5005t.rv3.jnpan.people.aws.dev'
apiVersion: kubeadm.k8s.io/v1beta2
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: '10.0.20.132'
  bindPort: 443
`
		port, _ := getLocalApiBindPortFromInitConfigYaml(yamlWithMultipleDoc)
		if port != 443 {
			t.Errorf("Returned unexpected port, expected: %d, actual: %d", 443, port)
		}
	})

	t.Run("yamlWithSingleDoc", func(t *testing.T) {
		const yamlWithSingleDoc = `
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: '10.0.20.132'
  bindPort: 443
`
		port, _ := getLocalApiBindPortFromInitConfigYaml(yamlWithSingleDoc)
		if port != 443 {
			t.Errorf("Returned unexpected port, expected: %d, actual: %d", 443, port)
		}
	})

	t.Run("yamlWithoutInitConfigurationKind", func(t *testing.T) {
		const yamlWithoutInitConfigurationKind = `
apiVersion: kubeadm.k8s.io/v1beta2
kind: OtherConfiguration
localAPIEndpoint:
  advertiseAddress: '10.0.20.132'
  bindPort: 443
`
		_, err := getLocalApiBindPortFromInitConfigYaml(yamlWithoutInitConfigurationKind)
		if err == nil {
			t.Errorf("Should return error for invalid yaml")
		}
	})

	t.Run("emptyYaml", func(t *testing.T) {
		_, err := getLocalApiBindPortFromInitConfigYaml("")
		if err == nil {
			t.Errorf("Should return error for empty yaml")
		}
	})

	t.Run("invalidYaml", func(t *testing.T) {
		_, err := getLocalApiBindPortFromInitConfigYaml("SomeRandomText")
		if err == nil {
			t.Errorf("Should return error for invalid yaml")
		}
	})
}

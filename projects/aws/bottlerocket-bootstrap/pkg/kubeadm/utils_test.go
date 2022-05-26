package kubeadm

import "testing"

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

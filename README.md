## Amazon EKS Anywhere Build Tooling Repository

The EKS Anywhere Build Tooling repository contains the code to build artifacts corresponding to the various upstream dependency projects of [Amazon EKS Anywhere](https://github.com/aws/eks-anywhere). The build artifacts include container images, binary archives, and OVA image archives that will be consumed by the EKS Anywhere CLI during the cluster creation/deletion/upgrade workflow.

## Base Image Tracker

This table tracks the base images used to build the container images for the upstream dependencies of EKS Anywhere.

<details>
<summary>Click to view/hide table</summary>


| Dockerfile | Image Repo | Base image |
| --- | --- | --- |
| [EKS-A tools](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/aws/eks-anywhere-build-tooling/docker/linux/Dockerfile) | [EKS-A tools image](https://gallery.ecr.aws/eks-anywhere/cli-tools) | [EKS Distro Minimal Base Docker Client Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-docker-client) |
| [Bottlerocket bootstrap](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/aws/bottlerocket-bootstrap/docker/linux/Dockerfile) | [Bottlerocket bootstrap image](https://gallery.ecr.aws/eks-anywhere/bottlerocket-bootstrap) | [EKS Distro Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base) |
| [EKS Anywhere cluster controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/aws/eks-anywhere/docker/linux/eks-anywhere-cluster-controller/Dockerfile) | [EKS Anywhere cluster controller image](https://gallery.ecr.aws/eks-anywhere/cluster-controller) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Kube RBAC Proxy](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/brancz/kube-rbac-proxy/docker/linux/Dockerfile) | [Kube RBAC Proxy image](https://gallery.ecr.aws/eks-anywhere/brancz/kube-rbac-proxy) | [EKS Distro Minimal Base Nonroot Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nonroot) |
| [Helm controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/fluxcd/helm-controller/docker/linux/Dockerfile) | [Helm controller image](https://gallery.ecr.aws/eks-anywhere/fluxcd/helm-controller) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Kustomize controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/fluxcd/kustomize-controller/docker/linux/Dockerfile) | [Kustomize controller image](https://gallery.ecr.aws/eks-anywhere/fluxcd/kustomize-controller) | [EKS Distro Minimal Base Git Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-git) |
| [Notification controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/fluxcd/notification-controller/docker/linux/Dockerfile) | [Notification controller image](https://gallery.ecr.aws/eks-anywhere/fluxcd/notification-controller) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Source controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/fluxcd/source-controller/docker/linux/Dockerfile) | [Source controller image](https://gallery.ecr.aws/eks-anywhere/fluxcd/source-controller) | [EKS Distro Minimal Base Git Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-git) |
| [Certmanager Acmesolver](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/cert-manager/cert-manager/docker/linux/cert-manager-acmesolver/Dockerfile) | [Certmanager Acmesolver image](https://gallery.ecr.aws/eks-anywhere/cert-manager/cert-manager-acmesolver) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Certmanager CA injector](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/cert-manager/cert-manager/docker/linux/cert-manager-cainjector/Dockerfile) | [Certmanager CA Injector image](https://gallery.ecr.aws/eks-anywhere/cert-manager/cert-manager-cainjector) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Certmanager Controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/cert-manager/cert-manager/docker/linux/cert-manager-controller/Dockerfile) | [Certmanager Controller image](https://gallery.ecr.aws/eks-anywhere/cert-manager/cert-manager-controller) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Certmanager Webhook](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/cert-manager/cert-manager/docker/linux/cert-manager-webhook/Dockerfile) | [Certmanager Webhook image](https://gallery.ecr.aws/eks-anywhere/cert-manager/cert-manager-webhook) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [vSphere Cloud Provider](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/kubernetes/cloud-provider-vsphere/docker/linux/Dockerfile) | [vSphere Cloud Provider image](https://gallery.ecr.aws/eks-anywhere/kubernetes/cloud-provider-vsphere/cpi/manager) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Cluster API controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/kubernetes-sigs/cluster-api/docker/linux/cluster-api-controller/Dockerfile) | [Cluster API controller image](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/cluster-api-controller) | [EKS Distro Minimal Base Nonroot Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nonroot) |
| [Kubeadm bootstrap controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/kubernetes-sigs/cluster-api/docker/linux/kubeadm-bootstrap-controller/Dockerfile) | [Kubeadm bootstrap controller image](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller) | [EKS Distro Minimal Base Nonroot Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nonroot) |
| [Kubeadm controlplane controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/kubernetes-sigs/cluster-api/docker/linux/kubeadm-control-plane-controller/Dockerfile) | [Kubeadm controlplane controller image](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/kubeadm-control-plane-controller) | [EKS Distro Minimal Base Nonroot Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nonroot) |
| [Cluster API Docker controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/kubernetes-sigs/cluster-api/docker/linux/cluster-api-docker-controller/Dockerfile) | [Cluster API Docker controller image](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/cluster-api-docker-controller) | [EKS Distro Minimal Base Docker Client Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-docker-client) |
| [Cluster API vSphere controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/kubernetes-sigs/cluster-api-provider-vsphere/docker/linux/cluster-api-vsphere-controller/Dockerfile) | [Cluster API vSphere controller image](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api-provider-vsphere/release/manager) | [EKS Distro Minimal Base Nonroot Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nonroot) |
| [Kind node](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/kubernetes-sigs/kind/images/node/Dockerfile.squash) | [Kind node image](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/kind/node) | [EKS Distro Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base) |
| [Kindnetd](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/kubernetes-sigs/kind/images/kindnetd/Dockerfile) | [Kindnetd image](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/kind/kindnetd) | [EKS Distro Minimal Base Iptables Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-iptables) |
| [Etcdadm bootstrap provider](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/aws/etcdadm-bootstrap-provider/docker/linux/Dockerfile) | [Etcdadm bootstrap provider image](https://gallery.ecr.aws/eks-anywhere/aws/etcdadm-bootstrap-provider) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Etcdadm controller](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/aws/etcdadm-controller/docker/linux/Dockerfile) | [Etcdadm controller image](https://gallery.ecr.aws/eks-anywhere/aws/etcdadm-controller) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Kube VIP](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/plunder-app/kube-vip/docker/linux/Dockerfile) | [Kube VIP image](https://gallery.ecr.aws/eks-anywhere/plunder-app/kube-vip) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |
| [Local path provisioner](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/rancher/local-path-provisioner/docker/linux/Dockerfile) | [Local path provisioner image](https://gallery.ecr.aws/eks-anywhere/rancher/local-path-provisioner) | [EKS Distro Minimal Base Image](https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base) |

</details>

## Contribution

We appreciate your interest in contributing. Please refer to the [Amazon EKS Anywhere Contribution Guide](https://github.com/aws/eks-anywhere/blob/main/CONTRIBUTING.md) before submitting any issues or pull requests.

- [Building locally](./docs/development/building-locally.md)
- Dealing with [attribution](./docs/development/attribution-files.md) files
- [Cherry picking](./docs/development/cherry-picks.md) to release branches

## Security

If you discover a potential security issue in this project, or think you may
have discovered a security issue, we ask that you notify AWS Security via our
[vulnerability reporting
page](http://aws.amazon.com/security/vulnerability-reporting/). Please do
**not** create a public GitHub issue.

## License

This project is licensed under the Apache-2.0 License.

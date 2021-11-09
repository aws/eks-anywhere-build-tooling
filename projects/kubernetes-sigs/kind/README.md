## **Kind**
![Version](https://img.shields.io/badge/version-v0.11.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiVkgvQm93WHUvUWJ1U2ZhSG9JTUJNMFdjdGtwSkIyRCt1azM0THYxcWYweC8rM2lHRmNYMXI0QkVPUm4yZ0JZZ1c4RzdMeTJ3dGtpREdYeFpvTEhtc2FnPSIsIml2UGFyYW1ldGVyU3BlYyI6Im9GV2EzRGZQNVZ5c25kTmoiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Kind](https://github.com/kubernetes-sigs/kind) is a tool for running local Kubernetes clusters using Docker container "nodes". kind bootstraps each "node" with `kubeadm`. kind consists of:
* Go packages implementing Kubernetes cluster creation, image build, etc.
* A command-line interface (`kind`) built on these packages.
* Docker image(s) written to run systemd, Kubernetes, etc.

In EKS-A CLI, the kind node image is used for integrating Cluster API with the Docker Provider infrastructure. This image is designed to run nested containers, systemd and Kubernetes components. The Kubernetes controlplane images such as `kube-api-server`, `kube-controller-manager`, etc. are obtained from the latest [EKS Distro](https://github.com/aws/eks-distro) releases on ECR Public Gallery and packaged into the kind node image as `.tar` archives. These archives are used to create the Kubernetes control plane.

The `eks-a` CLI uses the kind binary to create the initial booststrap cluster which is integrated with the Infrastructure provider to become the management cluster. The bootstrap/management cluster is required to create and delete the workload clusters.

The Kind HA Proxy image is used internally by kind to implement kubeadm's "HA" mode, specifically to load balance the API server.
This image is built from the upstream haproxy image with addition of a minimal config, which allows it to listen on the intended port and hot-reload it at runtime with the actual desired config.

Kind also implements its own standard CNI/cluster networking configuration in the form of `kindnetd`, a simple networking daemon which performs the following:
* IP masquerade of traffic leaving the nodes that is headed out of the cluster
* Ensuring netlink routes to pod CIDRs via the host node IP for each
* Ensuring a simple CNI config based on the standard ptp / host-local plugins and the node's pod CIDR

You can find the latest versions of these images on ECR Public Gallery.

[Node](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/kind/node) | [HA Proxy](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/kind/haproxy) | [Kindnetd](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/kind/kindnetd)

## Building
This project depends on other artifacts from this repo.  To build image locally, `ARTIFACTS_BUCKET` must be supplied. For ex
the following is the presubmit bucket:

`ARTIFACTS_BUCKET=s3://projectbuildpipeline-857-pipelineoutputartifactsb-10ajmk30khe3f make build`

To avoid pushing intermediate images to a remote repo, a local registry is required
to build images.  Refer to [building locally](../../../docs/development/building-locally.md) for more instructions.

`local-path-provisioner` is required to exist in the local registry to build images

`cd projects/rancher/local-path/provisioner && IMAGE_REPO=localhost:5000 make images`

To build all images for all supported EKS-D versions and amd64 + arm64, run

`ARTIFACTS_BUCKET=<> IMAGE_REPO=localhost:5000 make release`

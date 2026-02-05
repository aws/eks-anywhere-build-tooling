## **Kind**
![Version](https://img.shields.io/badge/version-v0.31.0-blue)
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

`cd projects/rancher/local-path-provisioner && IMAGE_REPO=localhost:5000 make images`

To build all images for all supported EKS-D versions and amd64 + arm64, run:

`ARTIFACTS_BUCKET=<> IMAGE_REPO=localhost:5000 make images`

For a specific `RELEASE_BRANCH`:

`ARTIFACTS_BUCKET=<> IMAGE_REPO=localhost:5000 RELEASE_BRANCH=1-X make images`


### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/kubernetes-sigs/kind) and decide on new version. 
The maintainers are pretty good about calling breaking changes and other upgrade gotchas between release.  Please
review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
and/or automatically update between eks-anywhere version reach out to @jaxesn.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [0.17.0 compared to 0.18.0](https://github.com/kubernetes-sigs/kind/compare/v0.17.0...v0.18.0). Check the `kind` target for
any build flag changes, tag changes, dependencies, etc in the `Makefile` in the root of the repo.  Pay close attention to
`images/base/Dockerfile` for changes when updating the patch.  Update constants in [node-image-build-args.sh](./build/node-image-build-args.sh#L52).
If new yum packages are added to the base image, update the [minimal-base-kind](https://github.com/aws/eks-distro-build-tooling/blob/main/eks-distro-base/Dockerfile.minimal-base-kind)
image to include it (this is not a blocker for updating). Review changes to [buildcontext.go](https://github.com/kubernetes-sigs/kind/blob/main/pkg/build/nodeimage/buildcontext.go)
closely to ensure there are no changes neccessary in our build scripts.
1. Verify the golang version has not changed. The version specified in `.go-version` should be the source of truth.
1. Update checksums and attribution using `make attribution checksums`.
1. Validate images build locally (will take a while) using the steps above.
1. Run `make create-kind-cluster-amd64 RELEASE_BRANCH=1-X` to ensure cluster creation works with the new image.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.

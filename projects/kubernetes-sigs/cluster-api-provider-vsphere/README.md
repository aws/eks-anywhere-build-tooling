## **Cluster API Provider for vSphere**
![Version](https://img.shields.io/badge/version-v1.3.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiYm85WnJ4aDc2ZXhhVUxOWHJuUFJwN3FlQmE2L1Q4b2ZzNG91OVpjNVNGM1ZvbVBEUUM2bkdER3N5eVNrWTBKS2VSSW9Oa051aFVWS1dzVVlTOHBBZ0NRPSIsIml2UGFyYW1ldGVyU3BlYyI6IlEwOWNtd0llNXdjUGRvQWkiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [Cluster API Provider for vSphere (CAPV)](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere) is a a concrete implementation of Cluster API for vSphere, which paves the way for true vSphere hybrid deployments of Kubernetes. CAPV is designed to allow customers to use their existing vSphere infrastructure, including vCenter credentials, VMs, templates, etc. for bootstrapping and creating workload clusters.

Some of the features of Cluster API Provider vSphere include:
* Native Kubernetes manifests and API
* Manages the bootstrapping of VMs on cluster.
* Choice of Linux distribution between Ubuntu 18.04 and CentOS 7 using VM Templates based on OVA images
* Deploys Kubernetes control planes into provided clusters on vSphere.

The Cluster API Provider vSphere controller image is used in the Provider confgiration to bootstrap the vSphere Infrastructure Provider in the EKS-A CLI.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api-provider-vsphere/release/manager).

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere) and decide on the new version.
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @jaxesn, @vignesh-goutham or @g-gaston.
1. Follow these steps for changes to the patches/ folder:
   1. Fork and clone CAPV repo, and checkout the desired tag. For instance if in step 1 we decided to upgrade to v1.3.1 CAPV version, do `git checkout v1.3.1`
   on your fork.
   1. Review the patches under patches/ folder in this repo. Apply the required patches to your fork. Remove any patches that are either
   merged upstream or no longer needed. Please reach out to @jaxesn, @vignesh-goutham or @g-gaston for any questions regarding which patches to keep.
   1. Run `git format-patch <commit>`, where `<commit>` is the last upstream commit on that tag. Move the generated patches under the patches/ folder in this repo.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes.
   ex: [1.1.1 compared to 1.3.1](https://github.com/kubernetes-sigs/provider-vsphere/compare/v1.1.1...v1.3.1). Check if the [manifests](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/v1.3.1/Makefile#L341)
   target has changed in the Makefile, and make the required changes in create_manifests.sh
1. Check the go.mod file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in
   Makefile to match the version upstream.
1. Update checksums and attribution using `make update-attribution-checksums-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

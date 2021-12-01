## **Cluster API**
![Version](https://img.shields.io/badge/version-v1.0.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQVZ3TDBZZVVXZUZiVmtqLzVoOVcrV2FaMmxRRzJXRmJCRlZtQkNodXdWZ0FrNm0zQ3l5UzNqTkdsQXgwdzc0bTBZc1RIcjBhMUVFbEhIK3d2VDVPek1rPSIsIml2UGFyYW1ldGVyU3BlYyI6IkVuOGJxNXBPZEtDek81Q3giLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Cluster API](https://github.com/kubernetes-sigs/cluster-api) is a Kubernetes sub-project focused on providing declarative APIs and tooling to simplify provisioning, upgrading, and operating multiple Kubernetes clusters. It uses Kubernetes-style APIs and patterns to automate cluster lifecycle management for platform operators. The supporting infrastructure, like virtual machines, networks, load balancers, and VPCs, as well as the Kubernetes cluster configuration are all defined in the same way that application developers operate deploying and managing their workloads. This enables consistent and repeatable cluster deployments across a wide variety of infrastructure environments. Cluster API can be extended to support any infrastructure provider (AWS, Azure, vSphere, etc.) or bootstrap provider (kubeadm is default) as required by the customer.

The `eks-a` CLI uses Cluster API to generate configurations for various infrasurcture providers, and uses it to create and manage multiple workload clusters.

You can find the latest versions of these images on ECR Public Gallery.

[Cluster API Controller](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/cluster-api-controller) | 
[Cluster API Kubeadm Bootstrap Controller](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller) | 
[Cluster API Kubeadm Controlplane Controller](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/kubeadm-control-plane-controller)

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/kubernetes-sigs/cluster-api) and decide on new version.
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @jaxesn, @vignesh-goutham, @g-gaston or @mrajashree.
1. Follow these steps for changes to the patches/ folder:
   1. Checkout the desired tag on our [CAPI fork](https://github.com/mrajashree/cluster-api) and create a new branch.
   1. Review the patches under patches/ folder in this repo. Apply the required patches to the new branch created in the above step. Remove any patches that are either
   merged upstream or no longer needed. Please reach out to @jaxesn, @vignesh-goutham, @g-gaston or @mrajashree if there are any questions regarding keeping/removing patches.
   1. Run `git format-patch <commit>`, where `<commit>` is the last upstream commit on that tag. Move the generated patches under the patches/ folder in this repo.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes.
   ex: [0.3.19 compared to 1.0.1](https://github.com/kubernetes-sigs/cluster-api/compare/v0.3.19...v1.0.1). Check the targets in the Makefile
   for any changes.
   1. For instance, [the make targets for creating manifests in v1.0.1 use a different path for manager_image_patch.yaml as compared to v0.3.19](https://github.com/kubernetes-sigs/cluster-api/commit/280db9a796d5e1c2b3b75aa3036fcfe44f669909#diff-76ed074a9305c04054cdebb9e9aad2d818052b07091de1f20cad0bbac34ffb52L368-L375).
   Based on this we updated the manifest generation targets in create_manifests.sh to use the updated path same as upstream.
1. Check the go.mod file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in
   Makefile to match the version upstream.
1. Update checksums and attribution using `make update-attribution-checksums-docker PROJECT=kubernetes-sigs/cluster-api` from the root of the repo.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.
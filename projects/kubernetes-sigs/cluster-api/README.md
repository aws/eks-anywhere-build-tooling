## **Cluster API**
![Version](https://img.shields.io/badge/version-v1.10.2-blue)
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
   and/or automatically update between eks-anywhere version reach out to @jaxesn, @vignesh-goutham or @g-gaston.
1. Follow these steps for changes to the patches/ folder:
   1. Checkout the desired tag on our [CAPI fork](https://github.com/abhay-krishna/cluster-api) and create a new branch.
   1. Review the patches under patches/ folder in this repo. Apply the required patches to the new branch created in the above step.
      1. Run `git am *.patch`
      1. For patches that need some manual changes, you will see a similar error: `Patch failed at *`
      1. For that patch, run `git apply --reject --whitespace=fix *.patch`. This will apply hunks of the patch that do apply correctly, leaving
      the failing parts in a new file ending in `.rej`. This file shows what changes weren't applied and you need to manually apply.
      1. Once the changes are done, delete the `.rej` file and run `git add .` and `git am --continue`
   1. Remove any patches that are either merged upstream or no longer needed. Please reach out to @jaxesn, @vignesh-goutham or @g-gaston if there are any questions regarding keeping/removing patches.
   1. Run `git format-patch <commit>`, where `<commit>` is the last upstream commit on that tag. Move the generated patches from under the CAPI fork to the projects/kubernetes-sigs/cluster-api/patches/ folder in this repo.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes.
   ex: [1.1.3 compared to 1.2.0](https://github.com/kubernetes-sigs/cluster-api/compare/v1.1.3...v1.2.0). Check the targets in the Makefile
   for any changes.
   1. For instance, [the Makefile for the CAPD provider was removed in v1.2.0 and targets moved to the main Makefile](https://github.com/kubernetes-sigs/cluster-api/commit/88dc60e28be303d6c371a49d463f700076469c52).
   Based on this we updated the manifest generation targets in create_manifests.sh to use the targets in the main Makefile.
1. Check the golang version by checking [this file](https://github.com/kubernetes-sigs/cluster-api/blob/main/.github/workflows/release.yml#L26) Update the field `GOLANG_VERSION` in
   Makefile to match the version upstream.
1. Check default CAPI [cert-manager version]((https://github.com/kubernetes-sigs/cluster-api/blob/main/cmd/clusterctl/client/config/cert_manager_client.go#L32)) for the CAPI tag, if it has changed, then update cert-manager.
1. Update checksums and attribution using `make update-attribution-checksums-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

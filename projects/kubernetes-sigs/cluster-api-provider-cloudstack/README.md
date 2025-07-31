## **Cluster API Provider for CloudStack**
![Version](https://img.shields.io/badge/version-v0.6.1-blue)
[![Go Report Card](https://goreportcard.com/badge/kubernetes-sigs/cluster-api-provider-cloudstack)](https://goreportcard.com/report/kubernetes-sigs/cluster-api-provider-cloudstack)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiS0M4VGRyK0xWM2ZZY0pRbVMvY0pHRWlVSEJ3M1I4SXNRaVNxSnB5blVYTHpHSkNFWlpXcWhHSmdlSkhCVnVwSXJyVm16NFlSUzVSRC9vN2g2bmY5NjVnPSIsIml2UGFyYW1ldGVyU3BlYyI6ImQ4ZldMWnMweEIyTmxrTk8iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [Cluster API Provider for CloudStack (CAPC)](https://github.com/kubernetes-sigs/cluster-api-provider-cloudstack) is a concrete implementation of Cluster API for CloudStack, which paves the way for true CloudStack hybrid deployments of Kubernetes. CAPC is designed to allow customers to use their existing CloudStack infrastructure, including CloudStack credentials, VMs, templates, etc. for bootstrapping and creating workload clusters.

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/kubernetes-sigs/cluster-api-provider-cloudstack) and decide on the new version.
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @jaxesn, @vignesh-goutham, or @g-gaston
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags. **Make sure the value in GIT_TAG follows [Semantic Versioning](http://semver.org/) (e.g. v0.4.5-rc4) in order for clusterctl to work properly**
1. Compare the old tag to the new, looking specifically for Makefile changes. Check if the [manifests](https://github.com/kubernetes-sigs/cluster-api-provider-cloudstack/blob/v0.3.0/Makefile#L51)
   target has changed in the Makefile, and make the required changes in create_manifests.sh
1. Check the go.mod file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in
   Makefile to match the version upstream.
1. Update checksums using `make checksum-files-project-kubernetes-sigs_cluster-api-provider-cloudstack` from the root directory of this repo.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

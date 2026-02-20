## **harbor**
![Version](https://img.shields.io/badge/version-v2.14.2-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiU2FkKytjT1M0SXpTa3lmL3BNSFhRbWpyNVBLdVRBOHdqajI0MnB2ZnFSR2k4aVNDQ2hyS1NDTU0wdnNWT2xORVR3aWhsY29ETjBVcVB1ay9GNWpQUmlRPSIsIml2UGFyYW1ldGVyU3BlYyI6IkNJOW1HQmkzUVBzY1pVajgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [harbor project](https://github.com/goharbor/harbor) is an open source trusted cloud native registry project that stores, signs, and scans content. Harbor extends the open source Docker Distribution by adding the functionalities usually required by users such as security, identity and management. Having a registry closer to the build and run environment can improve the image transfer efficiency. Harbor supports replication of images between registries, and also offers advanced security features such as user management, access control and activity auditing.

In EKS-A, harbor offers local cloud native registry service for Kubernetes clusters on vSphere infrastructure.

You can find the latest version of its images [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/harbor/).

### Updating

1. Review releases and changelogs in upstream [code repo](https://github.com/goharbor/harbor) and [chart repo](https://github.com/goharbor/harbor-helm), and decide on new version.
1. Review the patches under `patches/` folder and remove any that are either merged upstream or no longer needed.
1. Current patch information:
* `patches`:
    1. 0001 patch includes changes from `make/photon/common/install_cert.sh`
    1. 0002 patch includes changes from `make gen_apis`
    1. 0003 patch to update tag for tencentcloud-skd-go, tag mentioned in go.mod file doesn't exist.
* `helm/patches`:
    1. 0001 patch includes changes for adding digest, namespace and imagepullsecret support


1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags, and update the [`HELM_GIT_TAG`](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/goharbor/harbor/Makefile#L57) in Makefile accordingly so the versions of referenced images in the chart match what is in `GIT_TAG`.
1. Compare the old tag to the new, looking specifically for Makefile changes. Check the `build` target for any build flag changes, tag changes, dependencies, etc. Check that the manifest target has not changed, this is called from our Makefile.
1. Check the `go.mod` file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in Makefile to match the version upstream.
1. Update checksums and attribution using make `run-attribution-checksums-in-docker`.
1. Update the version at the top of this `README`.
1. Run `make generate` to update the `UPSTREAM_PROJECTS.yaml` file.

To make changes to the patches folder, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#generate-patch-files)


To test the upgrade, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#Testing).
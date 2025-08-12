## **Metal LB**
![Version](https://img.shields.io/badge/version-v0.15.2-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQSt5WjFpTGtiSGsxTFdFLzAxakxMbU1wZUE3LzNVR0NNMWlBYjNZeDVKeFl6YWxUZ2srNmJ4YW9ST2RxOHBTOStVMnVub1FYUW1LSWF5M3RsUGx5KzhNPSIsIml2UGFyYW1ldGVyU3BlYyI6IlEzdHh1SkJJMHV5WlZXbWUiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

Official website: https://metallb.universe.tf/
Upstream repository: https://github.com/metallb/metallb

MetalLB is a virtual-ip provider for kubernetes services of type load balancer. It supports ARP, BGP and more protocols when built with [FRR](https://frrouting.org/) support.

[Upstream Configuration examples](https://metallb.universe.tf/configuration/)

### Updating

1. Review [releases notes](https://metallb.universe.tf/release-notes/)
    * Any changes to the upstream configuration needs a thorough review + testing
    * Deprecation or removal of any protocol must be considered breaking 
1. Update the `GIT_TAG` and `GOLANG_VERSION` files to have the new desired version based on the upstream release tags.
1. Review the patches under `patches/` folder and remove any that are either merged upstream or no longer needed.
1. Current patch information:
    * `helm/patches`:
        1. 0001 patch add support for packages in chart.
1. Run `make build` or `make release` to build package, if `apply patch` step fails during build follow the  steps below to update the patch and rerun build/release again.
1. Run `make generate` from the root of the repo to update the `UPSTREAM_PROJECTS.yaml` file.
1. Update the version at the top of this `README`.
1. Verify no changes have been made to the dockerfiles [speaker](https://github.com/metallb/metallb/blob/main/speaker/Dockerfile)
   [controller](https://github.com/metallb/metallb/blob/main/controller/Dockerfile) 
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.

To make changes to the patches folder, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#generate-patch-files)


To test the upgrade, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#Testing).


#### Make target changes
1. Run `make add-generated-help-block` from the project root to update available make targets.

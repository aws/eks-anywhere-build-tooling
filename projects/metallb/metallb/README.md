## **Metal LB**
![Version](https://img.shields.io/badge/version-v0.14.5-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQSt5WjFpTGtiSGsxTFdFLzAxakxMbU1wZUE3LzNVR0NNMWlBYjNZeDVKeFl6YWxUZ2srNmJ4YW9ST2RxOHBTOStVMnVub1FYUW1LSWF5M3RsUGx5KzhNPSIsIml2UGFyYW1ldGVyU3BlYyI6IlEzdHh1SkJJMHV5WlZXbWUiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

Official website: https://metallb.universe.tf/
Upstream repository: https://github.com/metallb/metallb

MetalLB is a virtual-ip provider for kubernetes services of type load balancer. It supports ARP, BGP and more protocols when built with [FRR](https://frrouting.org/) support.

[Upstream Configuration examples](https://metallb.universe.tf/configuration/)

### Updating

1. Review [releases notes](https://metallb.universe.tf/release-notes/)
    * Any changes to the upstream configuration needs a thorough review + testing
    * Deprecation or removal of any protocol must be considered breaking 
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Verify the golang version has not changed. 
1. Verify no changes have been made to the dockerfiles [speaker](https://github.com/metallb/metallb/blob/main/speaker/Dockerfile)
   [controller](https://github.com/metallb/metallb/blob/main/controller/Dockerfile) 
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.

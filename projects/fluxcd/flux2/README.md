## **Flux**
![Version](https://img.shields.io/badge/version-v2.7.5-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiYzRDM0E2d3BGeHZNenB4aVdRY0RqMkhoMUZBdjVHdjZsTSsrVEdhVEw1Sy9DREIwRUlwSEx4MFpoUVBiK2grUnhyT2JodmNVWUVaemFGR2JTOWhkWC9VPSIsIml2UGFyYW1ldGVyU3BlYyI6Im1VckJkV25QbHdyc0hRbmgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Flux](https://github.com/fluxcd/flux2) is a tool for keeping Kubernetes clusters in sync with sources of configuration (like Git repositories), and automating updates to configuration when new code is deployed.

Flux v2 is constructed with the `GitOps` Toolkit, a set of composable APIs and specialized tools for building Continuous Delivery on top of Kubernetes. In version 2, Flux supports multi-tenancy and support for syncing an arbitrary number of Git repositories, among other long-requested features. It is built from the ground up to use Kubernetes' API extension system, and to integrate with Prometheus and other core components of the Kubernetes ecosystem.

EKS-A customers can store the configurations for all their clusters under version control, and leave the heavy-lifting of syncing and reconciliation to Flux.


### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/fluxcd/flux2) and decide on new version. 
Flux maintainers are pretty good about calling breaking changes and other upgrade gotchas between release.  Please
review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
and/or automatically update between eks-anywhere version reach out to @jiayiwang7 or @danbudris
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [v2.1.2 compared to v2.7.5](https://github.com/fluxcd/flux2/compare/v2.1.2...v2.2.1). Check the `build` target for
any build flag changes, tag changes, dependencies, etc. Check that the manifest target has not changed, this is called
from our Makefile.
1. Verify the golang version has not changed. The version specified in `go.mod` seems to be kept up to date.  There is also
a github release [action](https://github.com/fluxcd/flux2/blob/main/.github/workflows/release.yaml#L18) where the golang version
is defined.
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
1. Updating flux2 usually comes with updates to the source/helm/notification/kustomize-controller, make sure and update as well.

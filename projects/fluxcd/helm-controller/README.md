## **Helm Controller**
![Version](https://img.shields.io/badge/version-v1.5.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiS045T05yUXhCRzNPeXZwczkwcjgrbm8wOWJmSXZ6dll3eHBlVTV3bERUSlhadlRyOGE1Q1AzeWpEQTlvN2RISG9MNnMrMGRmOG1FZ2N2d0Nxc0l0b2UwPSIsIml2UGFyYW1ldGVyU3BlYyI6IlpJMTJ1cUxhdzc4bWlqNFUiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [helm-controller](https://github.com/fluxcd/helm-controller) is a Kubernetes operator that allows users to declaratively manage Helm chart releases. The desired state of a Helm release is described through a Kubernetes Custom Resource named HelmRelease. Based on the creation, mutation or removal of a HelmRelease resource in the cluster, Helm actions are performed by the operator.

Some of the features of the Helm controller are:

* Watches for `HelmRelease` objects and generates `HelmChart` objects
* Supports `HelmChart` artifacts produced from `HelmRepository`, `GitRepository` and `Bucket` sources
* Fetches artifacts produced by source-controller from `HelmChart` objects
* Watches `HelmChart` objects for revision changes (including semver ranges for charts from `HelmRepository` sources)
* Performs automated Helm actions, including Helm tests, rollbacks and uninstalls
* Offers extensive configuration options for automated remediation (rollback, uninstall, retry) on failed Helm install, upgrade or test actions
* Runs Helm install/upgrade in a specific order, taking into account the depends-on relationship defined in a set of `HelmRelease` objects
* Reports Helm release statuses
* Built-in Kustomize compatible Helm post renderer, providing support for strategic merge, JSON 6902 and images patches

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/fluxcd/helm-controller).

### Updating

1. Review releases and [changelogs](https://github.com/fluxcd/helm-controller/blob/main/CHANGELOG.md) in upstream 
[repo](https://github.com/fluxcd/helm-controller) and decide on new version. Flux maintainers are pretty good 
about calling breaking changes and other upgrade gotchas between release. Please review carefully and if there are questions 
about changes necessary to eks-anywhere to support the new version and/or automatically update between 
eks-anywhere version reach out to @jiayiwang7 or @danbudris
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [v0.33.0compared to v0.34.1](https://github.com/fluxcd/helm-controller/compare/v0.33.0...v0.34.1). Check the `manager` target for
any build flag changes, tag changes, dependencies, etc.
1. Verify the golang version has not changed. The version specified in `go.mod` seems to be kept up to date.  There is also
a [dockerfile](https://github.com/fluxcd/helm-controller/blob/main/Dockerfile#L6) they use for building which has it defined.
1. Verify no changes have been made to the [dockerfile](https://github.com/fluxcd/helm-controller/blob/main/Dockerfile) looking specifically for
added runtime deps.
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
1. When upgrading helm-controller to a new version, make sure to upgrade the fluxcd/flux2 project to a release that supports this version of helm-controller.

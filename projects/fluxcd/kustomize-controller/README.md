## **Kustomize Controller**
![Version](https://img.shields.io/badge/version-v0.26.3-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoibldOWFUyd2ZXOXR1WkNhSVZDZkprbEowWi9nNEZrN2RMcCtRK3EvQW9qbWUzQjcxVEZvTEZ6VUw3M004WHNKQ0M1MGJ4SlU0RUJvVE1YQ0hFT0hzZ21nPSIsIml2UGFyYW1ldGVyU3BlYyI6Ing4cTAwdG9pc1I0Qk81MlQiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [kustomize-controller](https://github.com/fluxcd/kustomize-controller) is a Kubernetes operator, specialized in running continuous delivery pipelines for infrastructure and workloads defined with Kubernetes manifests and assembled with Kustomize.

Some of the features of the Kustomize controller are:

* Watches for `Kustomization` objects
* Fetches artifacts produced by source-controller from `Source` objects 
* Watches `Source` objects for revision changes 
* Generates the `kustomization.yaml` file if needed
* Generates Kubernetes manifests with kustomize build
* Gecrypts Kubernetes secrets with Mozilla SOPS
* Validates the build output with client-side or APIServer dry-run
* Applies the generated manifests on the cluster
* Prunes the Kubernetes objects removed from source
* Checks the health of the deployed workloads
* Runs `Kustomizations` in a specific order, taking into account the depends-on relationship 
* Notifies whenever a `Kustomization` status changes

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/fluxcd/kustomize-controller).

### Updating

1. Review releases and [changelogs](https://github.com/fluxcd/kustomize-controller/blob/main/CHANGELOG.md) in upstream 
[repo](https://github.com/fluxcd/kustomize-controller) and decide on new version. Flux maintainers are pretty good 
about calling breaking changes and other upgrade gotchas between release. Please review carefully and if there are questions 
about changes necessary to eks-anywhere to support the new version and/or automatically update between 
eks-anywhere version reach out to @jiayiwang7 or @danbudris
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [0.19.0 compared to 0.24.1](https://github.com/fluxcd/kustomize-controller/compare/v0.19.0...v0.26.3). Check the `manager` target for
any build flag changes, tag changes, dependencies, etc.
1. Verify the golang version has not changed. The version specified in `go.mod` seems to be kept up to date.  There is also
a [dockerfile](https://github.com/fluxcd/kustomize-controller/blob/main/Dockerfile#L5) they use for building which has it defined.
1. Verify no changes have been made to the [dockerfile](https://github.com/fluxcd/kustomize-controller/blob/main/Dockerfile) looking specifically for
added runtime deps.
1. Update checksums and attribution using `make update-attribution-checksums-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

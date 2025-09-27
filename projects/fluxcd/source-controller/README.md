## **Source Controller**
![Version](https://img.shields.io/badge/version-v1.7.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiS1ZJY3BFVGg0a21PUmpDVWM2T0pnc2VxV25uYWt5aGJjQktVSURIVnBsd0VBUmljSlUxTVNyeG5pSzhFbXNaMkdiUGdBRWU5L2plMG9ldVFxcHhrYjd3PSIsIml2UGFyYW1ldGVyU3BlYyI6IjgybDlDK2ZHLzJQVmNZNFoiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [source-controller](https://github.com/fluxcd/source-controller) is a Kubernetes operator specialized in artifacts acquisition from external sources such as Git, Helm repositories and S3 buckets. The controller watches for `Source` objects in a cluster and acts on them. It was designed with the goal of offloading the sources' registration, authentication, verification and resource-fetching to a dedicated controller.

Some of the features of the Source controller are:

* Authenticates to sources (SSH, user/password, API token)
* Validates source authenticity (PGP)
* Detects source changes based on update policies (semver)
* Fetches resources on-demand and on-a-schedule
* Packages the fetched resources into a well-known format (tar.gz, yaml)
* Makes the artifacts addressable by their source identifier (SHA, version, ts)
* Makes the artifacts available in-cluster to interested third-parties
* Notifies interested third-parties of source changes and availability (status conditions, events, hooks)
* Reacts to Git push and Helm chart upload events

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/fluxcd/source-controller).

### Updating

1. Review releases and [changelogs](https://github.com/fluxcd/source-controller/blob/main/CHANGELOG.md) in upstream 
[repo](https://github.com/fluxcd/source-controller) and decide on new version. Flux maintainers are pretty good 
about calling breaking changes and other upgrade gotchas between release. Please review carefully and if there are questions 
about changes necessary to eks-anywhere to support the new version and/or automatically update between 
eks-anywhere version reach out to @jiayiwang7 or @danbudris
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [v1.2.2 compared to v1.7.0](https://github.com/fluxcd/source-controller/compare/v1.2.2...v1.2.3). Check the `build` target for
any build flag changes, tag changes, dependencies, etc.
1. Verify the golang version has not changed. The version specified in `go.mod` seems to be kept up to date.  There is also
a [dockerfile](https://github.com/fluxcd/source-controller/blob/main/Dockerfile#L2) they use for building which has it defined.
1. Verify no changes have been made to the [dockerfile](https://github.com/fluxcd/source-controller/blob/main/Dockerfile) looking specifically for
added runtime deps.
1. Since source-controller requires cgo it is built in the builder base. Update checksums and attribution using `make build` from the source-controller folder.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.
1. When upgrading source-controller to a new version, make sure to upgrade the fluxcd/flux2 project to a release that supports this version of source-controller.

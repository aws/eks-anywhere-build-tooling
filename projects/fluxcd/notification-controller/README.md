## **Notification Controller**
![Version](https://img.shields.io/badge/version-v1.6.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiVDhSNCt1djNtRytnWTgvQ01BMW13b2Y1YmZPakRrSGlRWitKZ0ZLZUdaS2xxclpLOFNidnBHNjBFWjRueHpOaGRrMzV5OUhLLzhRWHgyaC85R2tET2JZPSIsIml2UGFyYW1ldGVyU3BlYyI6IlZpNGwrazFrZndNMWE4cTciLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [notification-controller](https://github.com/fluxcd/notification-controller) is a Kubernetes operator specialized in handling inbound and outbound events. The controller exposes an HTTP endpoint for receiving events from other controllers. It can be configured with Kubernetes custom resources such as `Alert`, `Event`,`Provider` and `Receiver` to define how events are processed and where to dispatch them.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/fluxcd/notification-controller).

### Updating

1. Review releases and [changelogs](https://github.com/fluxcd/notification-controller/blob/main/CHANGELOG.md) in upstream 
[repo](https://github.com/fluxcd/notification-controller) and decide on new version. Flux maintainers are pretty good 
about calling breaking changes and other upgrade gotchas between release. Please review carefully and if there are questions 
about changes necessary to eks-anywhere to support the new version and/or automatically update between 
eks-anywhere version reach out to @jiayiwang7 or @danbudris
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [v1.2.2 compared to v1.6.0](https://github.com/fluxcd/notification-controller/compare/v1.2.2...v1.2.3). Check the `manager` target for
any build flag changes, tag changes, dependencies, etc.
1. Verify the golang version has not changed. The version specified in `go.mod` seems to be kept up to date.  There is also
a [dockerfile](https://github.com/fluxcd/notification-controller/blob/main/Dockerfile#L5) they use for building which has it defined.
1. Verify no changes have been made to the [dockerfile](https://github.com/fluxcd/notification-controller/blob/main/Dockerfile) looking specifically for
added runtime deps.
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
1. When upgrading notification-controller to a new version, make sure to upgrade the fluxcd/flux2 project to a release that supports this version of notification-controller.

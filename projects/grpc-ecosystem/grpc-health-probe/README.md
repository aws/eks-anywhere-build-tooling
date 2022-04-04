## **grpc-health-probe**
![Version](https://img.shields.io/badge/version-v0.4.7-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiYzRDM0E2d3BGeHZNenB4aVdRY0RqMkhoMUZBdjVHdjZsTSsrVEdhVEw1Sy9DREIwRUlwSEx4MFpoUVBiK2grUnhyT2JodmNVWUVaemFGR2JTOWhkWC9VPSIsIml2UGFyYW1ldGVyU3BlYyI6Im1VckJkV25QbHdyc0hRbmgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [`grpc_health_probe`](https://github.com/grpc-ecosystem/grpc-health-probe) utility allows users to query health of gRPC services that expose their status through the [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md). It is meant to be used for health checking gRPC applications in Kubernetes, using the [liveness probe mechanism](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-grpc-liveness-probe). This command-line utility makes a Remote Procedure Call to `/grpc.health.v1.Health/Check`. If it responds with a SERVING status, the grpc_health_probe will exit with success, otherwise it will exit with a non-zero exit code.


### Updating

1. This project's artifacts are consumed by the [Tinkerbell PBNJ](https://github.com/tinkerbell/pbnj) project and hence build versions should ideally be aligned with the tag PBNJ uses in its [Dockerfile](https://github.com/tinkerbell/pbnj/blob/main/Dockerfile#L15).
That being said, artifact security is of utmost priority and so, whenever possible, the decision to build off later tags in favor of including security fixes and avoiding building off unsupported Go versions should override staying in alignment with the PBNJ project.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Check the [Dockerfile](https://github.com/grpc-ecosystem/grpc-health-probe/blob/master/Dockerfile#L7) for any build flag changes, tag changes, dependencies, etc.
1. Verify the golang version has not changed. The version specified in `go.mod` seems to be kept up to date.  There is also
a github release [action](https://github.com/grpc-ecosystem/grpc-health-probe/blob/master/.github/workflows/release.yml#L15) where the golang version
is defined.
1. Update checksums and attribution using `make update-attribution-checksums-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

## **Rufio**
![Version](https://img.shields.io/badge/version-126069b950a57d571df90dfec7cd98e6d64692be-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoibmZUSnF0RVBLRGhmNENKWVdNa1kzZ2V0UlFOWWJVZmM0N3UzSm12ekZkRm5KL240YmZTWXdTL2p6NXlUdnF4SUdibzFubW41dW4wTWs1c3Y1TmdSNmw0PSIsIml2UGFyYW1ldGVyU3BlYyI6IkVHRnl2M2JGVTZZSWIyZ1UiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Rufio](https://github.com/tinkerbell/rufio) is a Kubernetes controller for managing baseboard management state and actions.

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/rufio) and decide on new commit to track.
1. Update the `GIT_TAG` file to have the new desired commit based on the upstream.
1. Verify the golang version has not changed. Currently the version mentioned in a [go.mod](https://github.com/tinkerbell/rufio/blob/main/go.mod#L3) is being used to build. If it has changed, update the version in the `Makefile`: `GOLANG_VERSION?=`.
1. Verify no changes have been made to the [dockerfile](https://github.com/tinkerbell/rufio/blob/main/Dockerfile) looking specifically for added runtime deps.
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

## **Rufio**
![Version](https://img.shields.io/badge/version-v0.3.0-blue)

[Rufio](https://github.com/tinkerbell/rufio) is a Kubernetes controller for managing baseboard management state and actions.

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/rufio) and decide on new commit to track.
1. Update the `GIT_TAG` file to have the new desired commit based on the upstream.
1. Verify the golang version has not changed. Currently the version mentioned in a [go.mod](https://github.com/tinkerbell/rufio/blob/main/go.mod#L3) is being used to build.
1. Verify no changes have been made to the [dockerfile](https://github.com/tinkerbell/rufio/blob/main/Dockerfile) looking specifically for added runtime deps.
1. Update checksums and attribution using `make update-attribution-checksums-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

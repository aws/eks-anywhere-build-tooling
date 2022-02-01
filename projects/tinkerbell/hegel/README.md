## **Hegel**
![Version](https://img.shields.io/badge/version-v0.6.0-blue)

[Hegel](https://github.com/tinkerbell/hegel) is a gRPC and HTTP metadata service for Tinkerbell. Subscribes to changes in metadata, get notified when data is added/removed, etc.

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/hegel) and decide on new release tag to track.
1. Update the `GIT_TAG` file to have the new desired tag based on upstream.
1. Verify the golang version has not changed. Currently the version mentioned in a [dockerfile](https://github.com/tinkerbell/hegel/blob/main/cmd/hegel/Dockerfile#L1) is being used to build.
1. Verify no changes have been made to the [dockerfile](https://github.com/tinkerbell/hegel/blob/main/cmd/hegel/Dockerfile) looking specifically for added runtime deps.
1. Update checksums and attribution using `make update-attribution-checksums-docker PROJECT=tinkerbell/hegel` from the root of the repo.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.
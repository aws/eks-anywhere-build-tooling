## **PBNJ**
![Version](https://img.shields.io/badge/version-6b3bb36af744c896d7364cdf57b9f7f71540b573-blue)

[PBNJ](https://github.com/tinkerbell/pbnj) is a service handles BMC interactions in the Tinkerbell stack.

It is responsible for the following operations:
* machine and BMC power on/off/reset
* setting next boot device
* user management
* setting BMC network source

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/pbnj) and decide on new commit to track. PBNJ is currently [experimental](https://github.com/packethost/standards/blob/main/experimental-statement.md) and does not have a release tag.
1. Update the `GIT_TAG` file to have the new desired commit based on the upstream.
1. Verify the golang version has not changed. Currently the version mentioned in a [dockerfile](https://github.com/tinkerbell/pbnj/blob/main/cmd/pbnj/Dockerfile#L1) is being used to build.
1. Verify no changes have been made to the [dockerfile](https://github.com/tinkerbell/pbnj/blob/main/cmd/pbnj/Dockerfile) looking specifically for added runtime deps.
1. Update checksums and attribution using `make update-attribution-checksums-docker PROJECT=tinkerbell/pbnj` from the root of the repo.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.
## **Hub**
![Version](https://img.shields.io/badge/version-6c0f0d437bde2c836d90b000312c8b25fa1b65e1-blue)
[Hub](https://github.com/tinkerbell/hub) is the repository that contains reusable Tinkerbell Actions. The different images are listed under [/actions](https://github.com/tinkerbell/hub/tree/main/actions).

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/hub) and decide on release tag to track. 
1. Update the `GIT_TAG` file to have the new desired tag based on upstream.
1. Verify the golang version has not changed. Currently the version 1.15 mentioned in the [Dockerfile](https://github.com/tinkerbell/hub/blob/main/actions/cexec/v1/Dockerfile) of each action.
1. Verify no changes have been made to the Dockerfile for each action under under [actions](https://github.com/tinkerbell/hub/blob/main/actions) looking specifically for added dependencies or build 
process changes.
1. Update checksums and attribution using `make update-attribution-checksums-docker PROJECT=tinkerbell/hub` from the root of the repo.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.

## **runc**
![Version](https://img.shields.io/badge/version-v1.3.3-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQ3dHSU45Mnd3bGhzMCtlbGliWXFNcXIxbGx0VDAxVmZqaGtSQ0hXMFN2Rm1DWkNuMG5ibi9GTVRSOFVQK0ZZZW9sUEU4MGJwTzYyVUxEU0lBUG1zVlk4PSIsIml2UGFyYW1ldGVyU3BlYyI6Im5Td1JrV0NEOEh1akJWSXQiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[runc](https://github.com/opencontainers/runc) is a CLI tool for spawning and running containers on Linux according to the OCI specification.

### Updating

1. Review releases and [changelogs](https://github.com/opencontainers/runc/releases) in upstream 
[repo](https://github.com/opencontainers/runc) and decide on new version. 
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [1.1.6 compared to 1.1.7](https://github.com/opencontainers/runc/compare/v1.1.6...v1.1.15). Check the release [Makefile](https://github.com/opencontainers/runc/blob/main/Makefile)
for any build flag changes, tag changes, dependencies, etc.  The [GO_BUILD](https://github.com/opencontainers/runc/blob/main/Makefile#L27) definition should be looked at closely.
1. Verify the golang version has not changed. The version specified in the [Dockerfile](https://github.com/opencontainers/runc/blob/main/Dockerfile#L1)
should be considered the source of truth.
1. Since runc requires cgo it is built in the builder base. Update checksums and attribution using `make build` from the runc folder.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.

## **distribution**
![Version](https://img.shields.io/badge/version-v2.8.3-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoieGduSTVGQXp1STQ1b2VjY0tiZnJOVStJa1pja2pjbDJYQTdMS2V5R0lyWFJ0R1lya1lYREhuYy9xRE5sMlc2SmZVWXlNRGRJdGhwZXl5V0cwMXB2ck5nPSIsIml2UGFyYW1ldGVyU3BlYyI6IlQwNHZleTBzMzZQMjZ1VCsiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

This [distribution project](https://github.com/distribution/distribution)'s main product is the Open Source Registry implementation for storing and distributing container images using the [OCI Distribution Specification](https://github.com/opencontainers/distribution-spec). The goal of this project is to provide a simple, secure, and scalable base for building a large scale registry solution or running a simple private registry. It is a core library for many registry operators including Docker Hub, GitHub Container Registry, GitLab Container Registry and DigitalOcean Container Registry, as well as the CNCF Harbor Project, and VMware Harbor Registry.

### Updating

1. Update distribution tag when updating harbor tag if harbor is using a newer tag. Use the same tag that harbor uses by default. For instance [harbor v2.5.0 uses distribution v2.8.0 by default](https://github.com/goharbor/harbor/blob/v2.5.0/Makefile#L124) so when updating to harbor tag v2.5.0, update distribution tag to v2.8.0 or higher if security patching requires.
1. Review releases and changelogs in upstream [repo](https://github.com/distribution/distribution) and decide on new version.
1. Review the patches under `patches/` folder and remove any that are either merged upstream or no longer needed.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. Check the `build` target for any build flag changes, tag changes, dependencies, etc. Check that the manifest target has not changed, this is called from our Makefile.
1. Check the `go.mod` file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in Makefile to match the version upstream.
1. Update checksums and attribution using make `run-attribution-checksums-in-docker`.
1. Update the version at the top of this `README`.
1. Run `make generate` to update the `UPSTREAM_PROJECTS.yaml` file.
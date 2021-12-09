## **cfssl**
![Version](https://img.shields.io/badge/version-v1.6.1-blue)

[cfssl](https://github.com/cloudflare/cfssl) is both a command line tool and an HTTP API server for signing, verifying, and bundling TLS certificates.

### Updating

1. Review the changelog for the release upstream [repo](https://github.com/cloudflare/cfssl/releases) and decide on the new release to track.
1. Update the `GIT_TAG` file to have the new desired release based on the upstream.
1. Verify the golang version has not changed. Currently the version mentioned in a [go.mod](https://github.com/cloudflare/cfssl/blob/master/go.mod) is being used to build.
1. Verify no changes have been made to the [Dockerfile](https://github.com/cloudflare/cfssl/blob/master/Dockerfile) looking specifically for added runtime/build deps.
1. Update checksums and attribution using `make update-attribution-checksums-docker PROJECT=cloudflare/cfssl` from the root of the repo.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.
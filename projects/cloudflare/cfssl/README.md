## **cfssl**
![Version](https://img.shields.io/badge/version-v1.6.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiWjdCKzFMUnRRZVRXeDhpK2xMNytRQ2ZkOUlmR1F2Y2pCSk9tMUdvblcrRVpQVkZzY28vWnlMbWZxcXh5anhCUms0Qjh4aGQ3dGxsUWZ1TS9sdFNOWTBFPSIsIml2UGFyYW1ldGVyU3BlYyI6IkFWS1Q0bndnNWxIeTh2OUgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[cfssl](https://github.com/cloudflare/cfssl) is both a command line tool and an HTTP API server for signing, verifying, and bundling TLS certificates.

### Updating

1. Review the changelog for the release upstream [repo](https://github.com/cloudflare/cfssl/releases) and decide on the new release to track.
1. Update the `GIT_TAG` file to have the new desired release based on the upstream.
1. Verify the golang version has not changed. Currently the version mentioned in a [go.mod](https://github.com/cloudflare/cfssl/blob/master/go.mod) is being used to build.
1. Verify no changes have been made to the [Dockerfile](https://github.com/cloudflare/cfssl/blob/master/Dockerfile) looking specifically for added runtime/build deps.
1. Update checksums and attribution using `make update-attribution-checksums-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
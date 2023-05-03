## **Tink**
![Version](https://img.shields.io/badge/version-v0.8.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiUmxrMmd4b2N6dk02TDRPVlVXQ1N3aEhsRzAxWFBtZ1Y1VVNXWEtVZlVNS0tkQlZ4MHFuNXJiWld0ZFMvVzVmMzZxWjhKK3FERWdQeEV6RWd6WFZBcGM0PSIsIml2UGFyYW1ldGVyU3BlYyI6ImEvZEhCemJsQXJWZXVmc2kiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Tink](https://github.com/tinkerbell/tink) consists of the tink-server, tink-controller and tink-worker. The tink-worker and tink-server communicate over gRPC, and are responsible for processing workflows. Tink-controller is Kubernetes controller that is responsible for reconciling Tinkerbell hardwares, templates and workflows. The CLI is the user-interactive piece for creating workflows and their building blocks, templates and hardware data.

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/tink) and decide on release tag to track. 
1. Update the `GIT_TAG` file to have the new desired tag based on upstream.
1. Verify the golang version has not changed. Currently the version 1.17 mentioned in the github workflows [ci.yaml](https://github.com/tinkerbell/tink/blob/main/.github/workflows/ci.yaml) is being used to build.
1. Verify no changes have been made to the Dockerfile for each image under [cmd/<image-name>](https://github.com/tinkerbell/tink/tree/main/cmd) looking specifically for added dependencies.
1. Update checksums and attribution using `make run-attribution-checksums-in-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

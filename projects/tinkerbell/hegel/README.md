## **Hegel**
![Version](https://img.shields.io/badge/version-v0.10.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiUFJ0a1NyeGo1SXlHVzFMWFp3YytQTk0zeXMrSE9oYUw2VFM2MUlpa0tkbmh5S3RGYUQwTzI5VC9KVUJ6ZUJYK3NZb05ZaU15SGVMMzFNSTdmL3lzUlBjPSIsIml2UGFyYW1ldGVyU3BlYyI6IllOR29JZFNSRUFoL2ROUkIiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Hegel](https://github.com/tinkerbell/hegel) is a gRPC and HTTP metadata service for Tinkerbell. Subscribes to changes in metadata, get notified when data is added/removed, etc.

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/hegel) and decide on new release tag to track.
1. Update the `GIT_TAG` file to have the new desired tag based on upstream.
1. Verify the golang version has not changed. Currently the version mentioned in a [dockerfile](https://github.com/tinkerbell/hegel/blob/main/cmd/hegel/Dockerfile#L1) is being used to build.
1. Verify no changes have been made to the [dockerfile](https://github.com/tinkerbell/hegel/blob/main/cmd/hegel/Dockerfile) looking specifically for added runtime deps.
1. Update checksums and attribution using `make run-attribution-checksums-in-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

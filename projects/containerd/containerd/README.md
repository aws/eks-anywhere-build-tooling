## **containerd**
![Version](https://img.shields.io/badge/version-v1.7.30-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiTWhoMS9lejNIZmxuZzB2NThxU1N5VXNoVVR3MlNWYVBqajA4M3QwN3BERHRjN3oxSGxCcmk4R3pqVVU0aVVHYVVsRnVReU5pdnRRQ1FGQ2djT0pmbjVzPSIsIml2UGFyYW1ldGVyU3BlYyI6ImpGdnQ4d05CL21Lbjdsa0oiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[containerd](https://github.com/containerd/containerd) is available as a daemon for Linux and Windows. It manages the complete container lifecycle of its host system, from image transfer and storage to container execution and supervision to low-level storage to network attachments and beyond.

### Updating

1. Review releases and [changelogs](https://github.com/containerd/containerd/releases) in upstream 
[repo](https://github.com/containerd/containerd) and decide on new version. 
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [1.6.20 compared to 1.6.21](https://github.com/containerd/containerd/compare/v1.6.20...v1.7.30). Check the release [dockerfile](https://github.com/containerd/containerd/blob/main/.github/workflows/release/Dockerfile)
and [Makefile](https://github.com/containerd/containerd/blob/main/Makefile#L99) for any build flag changes, tag changes, dependencies, etc.
1. Verify the golang version has not changed. The version specified in the release github [action](https://github.com/containerd/containerd/blob/main/.github/workflows/release.yml#L16)
should be considered the source of truth.
1. Check for runc version updates via the [runc-version](https://github.com/containerd/containerd/blob/main/script/setup/runc-version) file.
If there is a new version, update runc in this repo.
1. Since containerd requires cgo it is built in the builder base. Update checksums and attribution using `make build` from the containerd folder.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.

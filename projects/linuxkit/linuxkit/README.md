## **LinuxKit**
![Version](https://img.shields.io/badge/version-v1.6.5-blue)
![Build Status]()

[Linuxkit](https://github.com/linuxkit/linuxkit) is a toolkit for building secure, portable and lean operating systems for containers

### Updating

1. The version should be based on the version used by [hook](https://github.com/tinkerbell/hook/blob/main/build.sh#L32) upstream.
1. Update the `GIT_TAG` file to have the new desired tag or commit based on upstream.
1. Verify the golang version has not changed. Use the github release [workflow](https://github.com/linuxkit/linuxkit/blob/master/.github/workflows/release.yml#L13) to determine the correct golang version used upstream.
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

### Development

We are only building the linuxkit bin as well as the images needed by the upstream hook [template](https://github.com/tinkerbell/hook/blob/main/linuxkit-templates/hook.template.yaml).
If more images are used upstream, we will need to add them here.

Currently we are building these images using usptream dockerfiles.  The plan is to update these overtime to follow our standard process.
## **Hook**
![Version](https://img.shields.io/badge/version-v5.10.57-blue)

[hook](https://github.com/tinkerbell/hook) Hook is the Tinkerbell Installation Environment for bare-metal. It runs in-memory, installs operating system, and handles deprovisioning.

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/hook) and decide on new release tag to track.
1. Update the `GIT_TAG` file to have the new desired tag based on upstream.
1. Verify the golang version has not changed. Currently for `bootkit` and `tink-docker` the version mentioned in a [dockerfile](https://github.com/tinkerbell/hook/blob/5.10.57/tink-docker/Dockerfile#L3) of the respective projects is being used to build.
1. Verify no changes have been made to the dockerfile for each image. Looking specifically for added runtime deps.
1. `tink-docker` image has docker runtime. Hence, verify no new changes have been made with docker version updates.
1. Update checksums and attribution using `make update-attribution-checksums-docker PROJECT=tinkerbell/hook` from the root of the repo.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.

### Development
1. The project consists of 3 images. `bootkit`, `tink-docker` and `kernel`.
1. For `kernel`, the image builds off upstream. The `hook` project uses the kernel.org [linux 5.10.11 kernel](https://mirrors.edge.kernel.org/pub/linux/kernel/v5.x/linux-5.10.11.tar.xz) to build an image.
1. For building the in-memory OSIE files, `hook` uses [linuxkit](https://github.com/linuxkit/linuxkit). `Linuxkit build` expects the project images to be present in the repository represented by `IMAGE_REPO` variable.
1. To build locally, we suggest using a local registry and setting `IMAGE_REPO` variable.
## **Hook**
![Version](https://img.shields.io/badge/version-9d54933a03f2f4c06322969b06caa18702d17f66-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoicVVYYXpIMzRpazNGUTBWdnY1dittK09zNDJvRmtlUlpTZUtZRFoyMkZ0YzlZT3NBMTRSSUFacFg3ZzdVNjg3SlhOZ2dZNmExOVkwaDE5U2RNQldWSTBzPSIsIml2UGFyYW1ldGVyU3BlYyI6ImdYN1lEaGZuSVpQMjhLM2EiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Hook](https://github.com/tinkerbell/hook) is the Tinkerbell Installation Environment for bare-metal. It runs in-memory, installs operating system, and handles deprovisioning.

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/hook) and decide on new release tag or commit to track.
1. Update the `GIT_TAG` file to have the new desired tag or commit based on upstream.
1. Update the 'LINUX_KERNEL_VERSION' file to the kernel version tracked in [hook.yaml](https://github.com/tinkerbell/hook/blob/029ef8f0711579717bfd14ac5eb63cdc3e658b1d/hook.yaml#L2)
1. Verify the golang version has not changed. Currently for `hook-bootkit` and `hook-docker` the version mentioned in a [dockerfile](https://github.com/tinkerbell/hook/blob/6d43b8b331c7a389f3ffeaa388fa9aa98248d7a2/hook-docker/Dockerfile#L3) of the respective projects is being used to build.
1. Verify no changes have been made to the dockerfile for each image. Looking specifically for added runtime deps.
1. `hook-docker` image has docker runtime. Hence, verify no new changes have been made with docker version updates.
1. Update checksums and attribution using `make run-attribution-checksums-in-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.

### Development
1. The project consists of 3 images. `hook-bootkit`, `hook-docker` and `kernel`.
1. For `kernel`, the image builds off upstream. The `hook` project uses the kernel.org [linux 5.10.85 kernel](https://mirrors.edge.kernel.org/pub/linux/kernel/v5.x/linux-5.10.85.tar.xz) to build an image.
1. For building the in-memory OSIE files, `hook` uses [linuxkit](https://github.com/linuxkit/linuxkit). `Linuxkit build` expects the project images to be present in the repository represented by `IMAGE_REPO` variable.
1. To build locally, we suggest using a local registry and setting `IMAGE_REPO` variable.
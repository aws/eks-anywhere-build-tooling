## **Hook**
![Version](https://img.shields.io/badge/version-v0.10.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoicVVYYXpIMzRpazNGUTBWdnY1dittK09zNDJvRmtlUlpTZUtZRFoyMkZ0YzlZT3NBMTRSSUFacFg3ZzdVNjg3SlhOZ2dZNmExOVkwaDE5U2RNQldWSTBzPSIsIml2UGFyYW1ldGVyU3BlYyI6ImdYN1lEaGZuSVpQMjhLM2EiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Hook](https://github.com/tinkerbell/hook) is the Tinkerbell Installation Environment for bare-metal. It runs in-memory, installs operating system, and handles deprovisioning.

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/hook) and decide on new release tag or commit to track.  
2. Update the `GIT_TAG` file to have the new desired tag or commit based on upstream.  
3. Update the `LINUX_KERNEL_VERSION` file to the **latest patch release** of the **currently tracked minor version**:  
   - To find the latest patch release: visit [kernel.org](https://www.kernel.org/) and check the **Longterm** or **Stable** releases for your tracked minor (e.g., update `6.6.107` → `6.6.108`).  
   - To see which minor is currently tracked: check the [hook kernel configs](https://github.com/tinkerbell/hook/tree/main/kernel/configs). Upstream supports multiple versions, consider bumping to a new supported minor
   > **Note:** If the config file names change upstream, update `HOOK_IMAGE_FILES` in the Makefile.  

4. Verify the golang version has not changed. Currently for `hook-bootkit` and `hook-docker`, the version mentioned in their respective [Dockerfiles](https://github.com/tinkerbell/hook/blob/main/images/hook-bootkit/Dockerfile) is used to build. We only support building with one golang version per project; pick the latest from these two Dockerfiles if they do not match.  
5. Verify no changes have been made to the Dockerfile for each image, especially added runtime dependencies.  
6. `hook-docker` image has Docker runtime. Verify no new changes have been made with Docker version updates.  
7. Update checksums and attribution using `make attribution checksums`.  
8. Update the version at the top of this README.  
9. Run `make generate` to update the `UPSTREAM_PROJECTS.yaml` file.


### Development

1. The project consists of 3 images. `hook-bootkit`, `hoot-containerd`, `hook-docker`, `hook-embedded`, `hook-runc` and `kernel`.
1. For building the in-memory OS files (vmlinuz and initramfs), `hook` uses [linuxkit](https://github.com/linuxkit/linuxkit). `linuxkit build` expects the project images to be present in the repository represented by `IMAGE_REPO` variable.
1. For `kernel`, the image builds from the kernel source version defined in the `LINUX_KERNEL_VERSION` file. The upstream's [Dockerfile](https://github.com/tinkerbell/hook/blob/main/kernel/Dockerfile) is patched to use AL23 instead of alpine.
This image is used by linuxkit build and ends up on the "host" via the built vmlinuz. This should be kept to the latest patch for the choosen minor.
1. The `hook-bootkit`, `hoot-containerd` and `hook-runc` images are used during the linuxkit build process and ends up on the "host" via the built initramfs.
    - `hook-containerd` and `hook-runc` use the build of containerd and runc from eks-anywhere-build-tooling
    - `hook-bootkit` is a go bin used to start the `tink-worker` (built via [tinkerbell/tink](../tink/Makefile)) image
1. The `hook-embedded` image is built using the upstream pull-images script and used by linuxkit build to "embedded" the [action](../hub/Makefile) and [tink-worker](../tink/Makefile) images in the docker-in-docker image cache.
These embedded images allow customers to point our "latest" action images without neededing to update their cluster specs and avoids pulling these images at runtime.
1. To build locally, we suggest using a local registry, or your personal public ecr repos, and setting `IMAGE_REPO` variable.
    - To assist in creating the ecr repos, you can run `make create-ecr-repos`. If using public ecr, be sure to set `IMAGE_REPO`

#### Kernel

The kernel included in the vmlinuz file built by linuxkit is built from source, using upstream hook's [Dockerfile](https://github.com/tinkerbell/hook/blob/main/kernel/Dockerfile) and kernel [config](https://github.com/tinkerbell/hook/tree/main/kernel/configs) files.
Additional config options are applied based on EKS-A customer feedback and exists as files in the [config-patches](./config-patches/) folder.  These files are merged with upstream's config during the 
docker build process using [merge-config.sh](https://github.com/torvalds/linux/blob/master/scripts/kconfig/merge_config.sh) provided in the linux source.

To create a new config patch:
1. run `make create-new-config-patch` to launch the linux `menuconfig` process.
1. using the menu, enable the new options.
1. click `save` and save the changes to `.config`.
1. click `exit`.
1. after the menuconfig is exited, `_output/kernel-config/generic-5.10.y-x86_64-eksa` will be created and `diff` will be ran to give you the config options to set in your new `config-patches` file.

Running the built kernel image locally with qemu:
1. run `make run-kernel-in-qemu`
1. this will launch a qemu vm using the built vmlinuz/initramfs files
1. `ctrl + a + z` may work to kill the qemu terminal, if not, use `kill -9` to stop the process when you are done

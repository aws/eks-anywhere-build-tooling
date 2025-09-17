BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
KERNEL_VERSION:=$(shell cat LINUX_KERNEL_VERSION)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=hook
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base

HOOK_BOOTKIT_IMAGE_COMPONENT=tinkerbell/hook-bootkit
HOOK_DOCKER_IMAGE_COMPONENT=tinkerbell/hook-docker
HOOK_DIND_IMAGE_COMPONENT=tinkerbell/hook-dind
KERNEL_IMAGE_COMPONENT=tinkerbell/hook-kernel
HOOK_RUNC_IMAGE_COMPONENT=tinkerbell/hook-runc
HOOK_CONTAINERD_IMAGE_COMPONENT=tinkerbell/hook-containerd
HOOK_EMBEDDED_IMAGE_COMPONENT=tinkerbell/hook-embedded
HOOK_UDEV_IMAGE_COMPONENT=tinkerbell/hook-udev
HOOK_ACPID_IMAGE_COMPONENT=tinkerbell/hook-acpid
MKIMAGE_ISO_EFI_INITRD_IMAGE_COMPONENT=linuxkit/mkimage-iso-efi-initrd

IMAGE_NAMES=hook-bootkit hook-docker hook-runc hook-containerd kernel hook-dind hook-embedded hook-udev hook-acpid

# We are using this pattern of setting the dockerfile_folder instead of local to the target
# because this project uses the combine-images target which overrides the dockerfile_folder
# to the standard combine-images dockerfile
# If we setup our override in this Makefile at the target it does not allow the combine-images
# to override it
# edge cases:
# - kernel is a different folder in the hook repo (hook/kernel vs hook/images/*)
# - hook-dind uses the hook-docker image with a different target
IMAGES_USING_UPSTREAM_DOCKERFILES=kernel hook-udev hook-acpid
DOCKERFILE_FOLDER=$(strip $(if $(filter $(IMAGE_NAME),$(IMAGES_USING_UPSTREAM_DOCKERFILES)),\
	$(if $(filter kernel,$(IMAGE_NAME)),$(REPO)/kernel,$(REPO)/images/$(IMAGE_NAME)),\
	$(if $(filter hook-dind,$(IMAGE_NAME)),./docker/linux/hook-docker,./docker/linux/$(IMAGE_NAME))))

# historically we have included these two images in our bundle even tho we really do not
# need them at runtime since they are built into the final vmlinuz/initramfs files
# when in codebuild we run the combine-images build, only bother merging these three
ifneq ($(findstring combine-images,$(RELEASE_TARGETS)),)
IMAGE_NAMES=hook-bootkit hook-docker kernel
endif

BINARY_TARGET_FILES=hook-bootkit hook-docker
SOURCE_PATTERNS=./ ./
GO_MOD_PATHS=images/hook-bootkit images/hook-docker
# by default the first target is assumed to be the primary and will produce a attribution with no prefix
# for this build we want to prefix all the files
HOOK_BOOTKIT_ATTRIBUTION_OVERRIDE=hook-bootkit

CREATE_HOOK_FILES_PATTERN=$(REPO)/out/hook/vmlinuz
HOOK_PLATFORMS=$(call IF_OVERRIDE_VARIABLE,IMAGE_PLATFORMS,linux/$(BUILDER_PLATFORM_ARCH))
CREATE_HOOK_FILES=$(foreach platform,$(HOOK_PLATFORMS),$(CREATE_HOOK_FILES_PATTERN)-$(if $(findstring amd64,$(platform)),x86_64,aarch64))
HOOK_IMAGE_FILES=$(OUTPUT_DIR)/kernel/config-patches $(OUTPUT_DIR)/kernel/keys.asc $(OUTPUT_DIR)/kernel/configs/generic-$(KERNEL_MAJOR_MINOR).y-aarch64 $(OUTPUT_DIR)/kernel/configs/generic-$(KERNEL_MAJOR_MINOR).y-x86_64 $(OUTPUT_DIR)/kernel/download.sh

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_TARBALLS=false

FIX_LICENSES_HOOK_BOOTKIT_TARGET=$(REPO)/images/hook-bootkit/LICENSE
FIX_LICENSES_HOOK_DOCKER_TARGET=$(REPO)/images/hook-docker/LICENSE

KERNEL_MAJOR_V=v$(firstword $(subst ., ,$(KERNEL_VERSION))).x
KERNEL_MAJOR_MINOR=$(subst $(SPACE),.,$(wordlist 1,2,$(subst ., ,$(KERNEL_VERSION))))
DEFCONFIG=generic-$(basename $(KERNEL_VERSION)).y

KERNEL_ARCH=$(if $(findstring amd64,$(BUILDER_PLATFORM_ARCH)),x86_64,aarch64)
EKSA_DEFCONFIG=$(DEFCONFIG)-$(KERNEL_ARCH)-eksa
USER_DEFCONFIG=$(DEFCONFIG)-$(KERNEL_ARCH)
KERNEL_CONFIG_HOST_PATH=$(MAKE_ROOT)/$(OUTPUT_DIR)/kernel-config

HOOK_EMBEDDED_FOLDER=$(REPO)/images/hook-embedded
EMBEDDED_IMAGES=$(foreach platform,$(HOOK_PLATFORMS),$(OUTPUT_DIR)/hook-embedded/$(subst /,-,$(platform))/images)

LINUXKIT_VERSION=$(shell cat $(MAKE_ROOT)/../../linuxkit/linuxkit/GIT_TAG | cut -c2-)
LINUXKIT_CACHE_FILE=$(REPO)/cache/linuxkit-linux-$(BUILDER_PLATFORM_ARCH)-$(LINUXKIT_VERSION)
LINUXKIT_IMAGE_REPO?=$(IMAGE_REPO)

# we need to set IMAGE_BUILD_ARGS here even though its the same as the default. 
# it is set in Common.mk on the images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L799)
# and the combine-images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L846)
# since combine-images has images as prereq target, the ?= does not really behavior as one might expect.
# the images target being the actual action, its version of the set takes prioirty and resets to empty
# setting it explicitly to empty here takes allows the combine-images override to take proirty
IMAGE_BUILD_ARGS=

BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

PROJECT_DEPENDENCIES=eksa/containerd/containerd eksa/tinkerbell/actions eksa/tinkerbell/tink eksa/linuxkit/linuxkit

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_HOOK_BOOTKIT_TARGET) $(FIX_LICENSES_HOOK_DOCKER_TARGET)

# For local image builds we override to push to local registry.
%/images/amd64 %/images/arm64: IMAGE_OUTPUT=push=true
%/images/amd64 %/images/arm64: IMAGE_OUTPUT_TYPE=image

# hook-docker image required docker runtime.
# We are using eks-distro-minimal-base-glibc as the base and builder to install docker.
hook-docker/images/% hook-dind/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-iptables

hook-dind/images/%: IMAGE_TARGET=hook-dind

kernel-config/images/% kernel/images/%: IMAGE_BUILD_ARGS+=KERNEL_MAJOR_V KERNEL_VERSION DEFCONFIG
kernel-config/images/% kernel/images/%: IMAGE_CONTEXT_DIR=$(OUTPUT_DIR)/kernel
kernel-config/images/% kernel/images/%: BASE_IMAGE=public.ecr.aws/amazonlinux/amazonlinux:2023

$(call IMAGE_TARGETS_FOR_NAME, kernel): $(HOOK_IMAGE_FILES)
$(call IMAGE_TARGETS_FOR_NAME, kernel-config): $(HOOK_IMAGE_FILES)

kernel-config/images/%: IMAGE_TARGET=kernel-configurator
kernel-config/images/%: IMAGE_OUTPUT_TYPE=docker
kernel-config/images/%: IMAGE_IMPORT_CACHE=
kernel-config/images/%: LATEST_IMAGE=hook-kernel-builder:$(KERNEL_VERSION)

$(call IMAGE_TARGETS_FOR_NAME,hook-embedded): $(EMBEDDED_IMAGES)

# upstream hook builds these instead of using the linuxkit provided images
# building this from their dockerfile
hook-udev/images/%: IMAGE_CONTEXT_DIR=$(REPO)/images/$(IMAGE_NAME)
hook-acpid/images/%: IMAGE_CONTEXT_DIR=$(REPO)/images/$(IMAGE_NAME)

s3-artifacts: $(CREATE_HOOK_FILES)

$(OUTPUT_DIR)/kernel/config-patches:
	@mkdir -p $(OUTPUT_DIR)/kernel
	cp -rf config-patches $@

$(OUTPUT_DIR)/kernel/%:
	@mkdir -p $(@D)
	cp -f $(REPO)/kernel/$* $@

$(LINUXKIT_CACHE_FILE): $(PROJECT_DEPENDENCIES_TARGETS)
	@mkdir -p $(@D)
	cp $(BINARY_DEPS_DIR)/linux-$(BUILDER_PLATFORM_ARCH)/eksa/linuxkit/linuxkit/linuxkit $@

# We use the linuxkit built from eks-a build tooling
# This includes the linuxkit bin and the images
# The hook linuxkit-template has two sets of images it pulls in
# - from linuxkit upstream (dhcp/init/sysctl/etc)
# - from hook builds (ip/mdev)
# Before calling the hook target to build the OSIE files we replace the images it is
# pulling from upstream linuxkit with those built via eksa
# The images provided by hook are handled by the hook support for overriding those
# We use the sed approach instead of a patch to avoid having to update conflicts with each update
# After the images are replaced in the hook.template.yaml
# Using sed they are pulled back out so that we can run cache pull with retries since we tend to see
# issues from time to time pulling from ecr
# We skip ENABLE_DOCKER here because for some of the image builds in linuxkit such as the ISO
# linuxkit pulls in the base image and does a docker run.
# Since the docker daemon isn't running on the build container this would not work.
# We build on host to circumvent the above problem, this wouldn't ever occur on our CI builds.
$(CREATE_HOOK_FILES_PATTERN)-%: BUILD_ARCH=$(if $(findstring x86_64,$*),amd64,arm64)
$(CREATE_HOOK_FILES_PATTERN)-%: $(LINUXKIT_CACHE_FILE)
	@source $(BUILD_LIB)/common.sh && build::common::use_go_version $(GOLANG_VERSION); \
	if [ -f $(REPO)/out/linuxkit-hook-default-$(BUILD_ARCH)/hook-cmdline ]; then echo "remove hook/out before continuing"; exit 1; fi; \
	sed -i 's/[^ ]*linuxkit\/\([^:]*\)\(:.*$$\)/$(subst /,\/,$(LINUXKIT_IMAGE_REPO))\/linuxkit\/\1:$(LATEST_TAG)/' $(REPO)/linuxkit-templates/hook.template.yaml; \
	images=(); \
	while read -r image; do \
		if [[ "$$image" == *"sshd"* ]]; then continue; fi; \
		images+=($$image); \
	done < <(sed -n 's/.*\($(subst /,\/,$(LINUXKIT_IMAGE_REPO)).*$$\)/\1/p' $(REPO)/linuxkit-templates/hook.template.yaml); \
	for image in hook-bootkit hook-docker hook-runc hook-containerd hook-kernel hook-embedded hook-acpid hook-udev; do \
		images+=($(IMAGE_REPO)/tinkerbell/$$image:$(LATEST_TAG)$(IMAGE_TAG_SUFFIX)); \
	done; \
	for image in "$${images[@]}"; do build::common::echo_and_run retry $(LINUXKIT_CACHE_FILE) cache pull --cache $(REPO)/cache/linuxkit $$image; done; \
	cd $(REPO); \
	DEBUG=yes LINUXKIT_VERSION_DEFAULT=$(LINUXKIT_VERSION) HOOK_KERNEL_OCI_BASE=$(IMAGE_REPO)/$(KERNEL_IMAGE_COMPONENT) \
		HOOK_LK_CONTAINERS_OCI_BASE=$(IMAGE_REPO)/tinkerbell/ \
		HOOK_KERNEL_POINT_RELEASE=$(lastword $(subst ., ,$(KERNEL_VERSION))) \
		HOOK_KERNEL_OCI_VERSION=$(LATEST_TAG)$(IMAGE_TAG_SUFFIX) \
		HOOK_LK_CONTAINERS_OCI_VERSION=$(LATEST_TAG)$(IMAGE_TAG_SUFFIX) \
		EKS_A_BUILD_ISO_EFI_IMAGE=$(IMAGE_REPO)/$(MKIMAGE_ISO_EFI_INITRD_IMAGE_COMPONENT):$(LATEST_TAG)$(IMAGE_TAG_SUFFIX) \
		LINUXKIT_ISO=true \
		./build.sh build hook-default-$(BUILD_ARCH); \
	DEBUG=yes LINUXKIT_VERSION_DEFAULT=$(LINUXKIT_VERSION) HOOK_KERNEL_OCI_BASE=$(IMAGE_REPO)/$(KERNEL_IMAGE_COMPONENT) \
		HOOK_LK_CONTAINERS_OCI_BASE=$(IMAGE_REPO)/tinkerbell/ \
		HOOK_KERNEL_POINT_RELEASE=$(lastword $(subst ., ,$(KERNEL_VERSION))) \
		HOOK_KERNEL_OCI_VERSION=$(LATEST_TAG)$(IMAGE_TAG_SUFFIX) \
		HOOK_LK_CONTAINERS_OCI_VERSION=$(LATEST_TAG)$(IMAGE_TAG_SUFFIX) \
		./build.sh build hook-default-$(BUILD_ARCH); \
	cd ..; \
	mkdir -p $(OUTPUT_DIR)/hook/$(GIT_TAG); \
	cp $(REPO)/out/hook/vmlinuz-$* $(OUTPUT_DIR)/hook/$(GIT_TAG)/; \
	cp $(REPO)/out/hook/initramfs-$* $(OUTPUT_DIR)/hook/$(GIT_TAG)/; \
	cp $(REPO)/out/hook-$*-efi-initrd.iso $(OUTPUT_DIR)/hook/$(GIT_TAG)/; \
	mkdir -p $(ARTIFACTS_PATH); \
	cp -rf $(OUTPUT_DIR)/hook/* $(ARTIFACTS_PATH)

$(REPO)/%/LICENSE: | $(GO_MOD_DOWNLOAD_TARGETS)
#go-licenses requires a LICENSE file in each folder with the go.mod
	cp $(REPO)/LICENSE $@

$(OUTPUT_DIR)/hook-embedded/linux-%/images: MAKEFLAGS=--no-print-directory
# in codebuild this is set to arm64/amd64, for the images pulled during this target
# we never want to use the -<arch>
$(OUTPUT_DIR)/hook-embedded/linux-%/images: IMAGE_TAG_SUFFIX=
$(OUTPUT_DIR)/hook-embedded/linux-%/images: DIND_IMAGE=$(shell $(MAKE) var-value-LATEST_IMAGE IMAGE_NAME=hook-dind)
$(OUTPUT_DIR)/hook-embedded/linux-%/images: $(GIT_PATCH_TARGET)
	mkdir -p $(@D)
	echo "$$($(MAKE) -C $(MAKE_ROOT)/../tink var-value-LATEST_IMAGE IMAGE_NAME=tink-worker) 127.0.0.1/embedded/tink-worker true" > $@.txt
	for action in cexec kexec image2disk oci2disk writefile reboot; do \
		echo "$$($(MAKE) -C $(MAKE_ROOT)/../actions var-value-LATEST_IMAGE IMAGE_NAME=$$action) 127.0.0.1/embedded/$$action true" >> $@.txt; \
	done
	cat $@.txt
	@if [[ "$(JOB_TYPE)" == "presubmit" ]] && [[ "$(DOCKER_AVAILABLE)" == "false" ]]; then \
		mkdir -p $@; \
	else \
		source $(BUILD_LIB)/common.sh; \
		function docker() { \
			local -r subcommand="$${1:-}"; \
			if [[ "$$subcommand" == "pull" ]]; then \
				retry command docker "$$@"; \
			else \
				command docker "$$@"; \
			fi; \
		}; \
		docker pull $(DIND_IMAGE); \
		cd $(@D) && source $(MAKE_ROOT)/$(HOOK_EMBEDDED_FOLDER)/pull-images.sh $* $(DIND_IMAGE); \
		sudo chown -R 1000:1000 ./images; \
	fi

.PHONY: run-kernel-in-qemu
run-kernel-in-qemu: $(CREATE_HOOK_FILES) $(LINUXKIT_CACHE_FILE)
	mkdir -p $(REPO)/out/linuxkit-hook-default-$(BUILDER_PLATFORM_ARCH)
	cp $(REPO)/out/hook/vmlinuz-$(KERNEL_ARCH) $(REPO)/out/linuxkit-hook-default-$(BUILDER_PLATFORM_ARCH)/hook-kernel
	cp $(REPO)/out/hook/initramfs-$(KERNEL_ARCH) $(REPO)/out/linuxkit-hook-default-$(BUILDER_PLATFORM_ARCH)/hook-initrd.img
	cd $(REPO); \
	DEBUG=yes TINK_WORKER_IMAGE=127.0.0.1/embedded/tink-worker LINUXKIT_VERSION_DEFAULT=$(LINUXKIT_VERSION) \
		sudo -E ./build.sh run hook-default-$(BUILDER_PLATFORM_ARCH)

kernel-config-%: | ensure-docker
	@$(MAKE) kernel-config/images/$(BUILDER_PLATFORM_ARCH) IMAGE_NAMES=kernel-config DOCKERFILE_FOLDER=$(REPO)/kernel
	@mkdir -p $(KERNEL_CONFIG_HOST_PATH)
	docker run -it -v $(KERNEL_CONFIG_HOST_PATH):/host hook-kernel-builder:$(KERNEL_VERSION) bash -c "cp .config /host/$(EKSA_DEFCONFIG) && make $* && cp -v .config /host/$(USER_DEFCONFIG)"

menuconfig: kernel-config-menuconfig

create-new-config-patch: kernel-config-menuconfig
	@diff $(KERNEL_CONFIG_HOST_PATH)/$(EKSA_DEFCONFIG) $(KERNEL_CONFIG_HOST_PATH)/$(USER_DEFCONFIG) || true


clean-repo: clean-extra

# the extract container images in the _output directory
# will have root permissions and require sudo to delete
.PHONY: clean-extra
clean-extra:	
	sudo rm -rf $(EMBEDDED_IMAGES) $(REPO)/out


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
LINUX_KERNEL_VERSION:=$(shell cat LINUX_KERNEL_VERSION)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=hook
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

HOOK_BOOTKIT_IMAGE_COMPONENT=tinkerbell/hook-bootkit
HOOK_DOCKER_IMAGE_COMPONENT=tinkerbell/hook-docker
KERNEL_IMAGE_COMPONENT=tinkerbell/hook-kernel

IMAGE_NAMES=hook-bootkit hook-docker kernel

BINARY_TARGET_FILES=hook-bootkit hook-docker
SOURCE_PATTERNS=./ ./
GO_MOD_PATHS=hook-bootkit hook-docker
# by default the first target is assumed to be the primary and will produce a attribution with no prefix
# for this build we want to prefix all the files
HOOK_BOOTKIT_ATTRIBUTION_OVERRIDE=hook-bootkit

CREATE_HOOK_FILES=$(REPO)/out/sha-6d43b8b/rel/vmlinuz-x86_64

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_TARBALLS=false

FIX_LICENSES_HOOK_BOOTKIT_TARGET=$(REPO)/hook-bootkit/LICENSE
FIX_LICENSES_HOOK_DOCKER_TARGET=$(REPO)/hook-docker/LICENSE

# we need to set IMAGE_BUILD_ARGS here even though its the same as the default. 
# it is set in Common.mk on the images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L799)
# and the combine-images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L846)
# since combine-images has images as prereq target, the ?= does not really behavior as one might expect.
# the images target being the actual action, its version of the set takes prioirty and resets to empty
# setting it explicitly to empty here takes allows the combine-images override to take proirty
IMAGE_BUILD_ARGS=

BUILDSPECS=buildspec.yml buildspecs/combine-images.yml
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE
BUILDSPEC_1_VARS_KEYS=IMAGE_PLATFORMS
BUILDSPEC_1_VARS_VALUES=IMAGE_PLATFORMS
BUILDSPEC_1_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_HOOK_BOOTKIT_TARGET) $(FIX_LICENSES_HOOK_DOCKER_TARGET)

# For local image builds we override to push to local registry.
%/images/amd64 %/images/arm64: IMAGE_OUTPUT=push=true
%/images/amd64 %/images/arm64: IMAGE_OUTPUT_TYPE=image

# hook-docker image required docker runtime.
# We are using eks-distro-minimal-base-glibc as the base and builder to install docker.
hook-docker/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-glibc

# Currently the kernel image is being built off upstream.
# TODO: Setup a build environment to build the linux kernel.
kernel/images/%: BASE_IMAGE=quay.io/tinkerbell/hook-kernel:$(LINUX_KERNEL_VERSION)

s3-artifacts: $(CREATE_HOOK_FILES)

$(CREATE_HOOK_FILES): MAKEFLAGS=
$(CREATE_HOOK_FILES): tarballs | $$(ENABLE_DOCKER)
# Modify the linuxkit config file hook/hook.yaml to point to IMAGE_REPO.
# Upstream make target `dist` is triggered to perform linuxkit build and generate OSIE files ().
	@source $(BUILD_LIB)/common.sh && build::common::use_go_version "1.19" && make dist IMAGE_REPO=$(IMAGE_REPO) LATEST_TAG=$(LATEST_TAG) -C $(REPO); \
	mkdir -p $(OUTPUT_DIR)/hook/$(GIT_TAG); \
	cp $(REPO)/out/sha-*/rel/vmlinuz-* $(OUTPUT_DIR)/hook/$(GIT_TAG)/; \
	cp $(REPO)/out/sha-*/rel/initramfs-* $(OUTPUT_DIR)/hook/$(GIT_TAG)/; \
	mkdir -p $(ARTIFACTS_PATH); \
	cp -rf $(OUTPUT_DIR)/hook/* $(ARTIFACTS_PATH)

$(REPO)/%/LICENSE: | $(GO_MOD_DOWNLOAD_TARGETS)
#go-licenses requires a LICENSE file in each folder with the go.mod
	cp $(REPO)/LICENSE $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

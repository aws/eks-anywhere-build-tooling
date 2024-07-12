BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
KERNEL_VERSION:=$(shell cat LINUX_KERNEL_VERSION)
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

CREATE_HOOK_FILES_PATTERN=$(REPO)/out/sha-6d43b8b/rel/vmlinuz
HOOK_PLATFORMS=$(call IF_OVERRIDE_VARIABLE,IMAGE_PLATFORMS,linux/$(BUILDER_PLATFORM_ARCH))
CREATE_HOOK_FILES=$(foreach platform,$(HOOK_PLATFORMS),$(CREATE_HOOK_FILES_PATTERN)-$(if $(findstring amd64,$(platform)),x86_64,aarch64))
HOOK_IMAGE_FILES=$(OUTPUT_DIR)/kernel/config-patches $(OUTPUT_DIR)/kernel/keys.asc $(OUTPUT_DIR)/kernel/configs/config-5.10.x-aarch64 $(OUTPUT_DIR)/kernel/configs/config-5.10.x-x86_64 

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_TARBALLS=false

FIX_LICENSES_HOOK_BOOTKIT_TARGET=$(REPO)/hook-bootkit/LICENSE
FIX_LICENSES_HOOK_DOCKER_TARGET=$(REPO)/hook-docker/LICENSE

KERNEL_MAJOR_V=v$(firstword $(subst ., ,$(KERNEL_VERSION))).x
DEFCONFIG=config-$(basename $(KERNEL_VERSION)).x

KERNEL_ARCH=$(if $(findstring amd64,$(BUILDER_PLATFORM_ARCH)),x86_64,aarch64)
EKSA_DEFCONFIG=$(DEFCONFIG)-$(KERNEL_ARCH)-eksa
USER_DEFCONFIG=$(DEFCONFIG)-$(KERNEL_ARCH)
KERNEL_CONFIG_HOST_PATH=$(MAKE_ROOT)/$(OUTPUT_DIR)/kernel-config

BUILDSPECS=buildspec.yml buildspecs/combine-images.yml
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE
BUILDSPEC_1_VARS_KEYS=IMAGE_PLATFORMS
BUILDSPEC_1_VARS_VALUES=IMAGE_PLATFORMS
BUILDSPEC_1_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER
BUILDSPEC_2_DEPENDS_ON_OVERRIDE=tinkerbell_hook_linux_amd64 tinkerbell_hook_linux_arm64

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
hook-docker/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-glibc

kernel-config/images/% kernel/images/%: IMAGE_BUILD_ARGS=KERNEL_MAJOR_V KERNEL_VERSION DEFCONFIG
kernel-config/images/% kernel/images/%: IMAGE_CONTEXT_DIR=$(OUTPUT_DIR)/kernel
kernel-config/images/% kernel/images/%: BASE_IMAGE=public.ecr.aws/amazonlinux/amazonlinux:2023
kernel-config/images/% kernel/images/%: DOCKERFILE_FOLDER=$(REPO)/kernel

$(call IMAGE_TARGETS_FOR_NAME, kernel): $(HOOK_IMAGE_FILES)
$(call IMAGE_TARGETS_FOR_NAME, kernel-config): $(HOOK_IMAGE_FILES)

kernel-config/images/%: IMAGE_TARGET=kernel-configurator
kernel-config/images/%: IMAGE_OUTPUT_TYPE=docker
kernel-config/images/%: IMAGE_IMPORT_CACHE=
kernel-config/images/%: LATEST_IMAGE=hook-kernel-builder:$(KERNEL_VERSION)

s3-artifacts: $(CREATE_HOOK_FILES)

$(OUTPUT_DIR)/kernel/config-patches:
	@mkdir -p $(OUTPUT_DIR)/kernel
	cp -rf config-patches $@

$(OUTPUT_DIR)/kernel/%:
	@mkdir -p $(@D)
	cp -f $(REPO)/kernel/$(@F) $@

$(CREATE_HOOK_FILES_PATTERN)-%: MAKEFLAGS=
$(CREATE_HOOK_FILES_PATTERN)-%: tarballs | $$(ENABLE_DOCKER)
# Modify the linuxkit config file hook/hook.yaml to point to IMAGE_REPO.
# Upstream make target `dist-<arch>` is triggered to perform linuxkit build and generate OSIE files ().
	@source $(BUILD_LIB)/common.sh && build::common::use_go_version "1.19" && make dist-$* IMAGE_REPO=$(IMAGE_REPO) LATEST_TAG=$(LATEST_TAG) -C $(REPO); \
	mkdir -p $(OUTPUT_DIR)/hook/$(GIT_TAG); \
	cp $(REPO)/out/sha-*/rel/vmlinuz-$* $(OUTPUT_DIR)/hook/$(GIT_TAG)/; \
	cp $(REPO)/out/sha-*/rel/initramfs-$* $(OUTPUT_DIR)/hook/$(GIT_TAG)/; \
	mkdir -p $(ARTIFACTS_PATH); \
	cp -rf $(OUTPUT_DIR)/hook/* $(ARTIFACTS_PATH)

$(REPO)/%/LICENSE: | $(GO_MOD_DOWNLOAD_TARGETS)
#go-licenses requires a LICENSE file in each folder with the go.mod
	cp $(REPO)/LICENSE $@

kernel-config-%: | ensure-docker
	@$(MAKE) kernel-config/images/$(BUILDER_PLATFORM_ARCH) IMAGE_NAMES=kernel-config
	@mkdir -p $(KERNEL_CONFIG_HOST_PATH)
	docker run -it -v $(KERNEL_CONFIG_HOST_PATH):/host hook-kernel-builder:$(KERNEL_VERSION) bash -c "cp .config /host/$(EKSA_DEFCONFIG) && make $* && cp -v .config /host/$(USER_DEFCONFIG)"

menuconfig: kernel-config-menuconfig

create-new-config-patch: kernel-config-menuconfig
	@diff $(KERNEL_CONFIG_HOST_PATH)/$(EKSA_DEFCONFIG) $(KERNEL_CONFIG_HOST_PATH)/$(USER_DEFCONFIG) || true


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

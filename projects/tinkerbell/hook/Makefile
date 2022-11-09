BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
LINUX_KERNEL_VERSION:=$(shell cat LINUX_KERNEL_VERSION)
GOLANG_VERSION?="1.17"

REPO=hook
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

BOOTKIT_IMAGE_COMPONENT=tinkerbell/hook-bootkit
TINK_DOCKER_IMAGE_COMPONENT=tinkerbell/hook-docker
KERNEL_IMAGE_COMPONENT=tinkerbell/hook-kernel

IMAGE_NAMES=hook-bootkit hook-docker kernel

BINARY_TARGET_FILES=hook-bootkit hook-docker
SOURCE_PATTERNS=./ ./
GO_MOD_PATHS=hook-bootkit hook-docker
# by default the first target is assumed to be the primary and will produce a attribution with no prefix
# for this build we want to prefix all the files
BOOTKIT_ATTRIBUTION_OVERRIDE=hook-bootkit

CREATE_HOOK_FILES=$(REPO)/dist/vmlinuz-x86_64

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_TARBALLS=false

FIX_LICENSES_BOOTKIT_TARGET=$(REPO)/hook-bootkit/LICENSE
FIX_LICENSES_TINK_DOCKER_TARGET=$(REPO)/hook-docker/LICENSE

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_BOOTKIT_TARGET) $(FIX_LICENSES_TINK_DOCKER_TARGET)

# For local image builds we override to push to local registry.
%/images/amd64: IMAGE_OUTPUT=push=true
%/images/amd64: IMAGE_OUTPUT_TYPE=image

# hook-docker image required docker runtime.
# We are using eks-distro-minimal-base-glibc as the base and builder to install docker.
hook-docker/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-glibc

# Currently the kernel image is being built off upstream.
# TODO: Setup a build environment to build the linux kernel.
kernel/images/%: BASE_IMAGE=quay.io/tinkerbell/hook-kernel:$(LINUX_KERNEL_VERSION)

s3-artifacts: $(CREATE_HOOK_FILES)

$(CREATE_HOOK_FILES): tarballs
# Modify the linuxkit config file hook/hook.yaml to point to IMAGE_REPO.
# Upstream make target `dist-existing-images` is triggered to perform linuxkit build and generate OSIE files.
	sed -i -E -e 's,quay.io/tinkerbell/hook-(bootkit|docker|kernel).*,$(IMAGE_REPO)/tinkerbell/hook-\1:$(LATEST_TAG),g' $(REPO)/hook.yaml
	source $(BUILD_LIB)/common.sh && build::common::use_go_version "1.16" && make dist-existing-images -C $(REPO)
	mkdir -p $(OUTPUT_DIR)/hook/$(GIT_TAG)
	cp $(REPO)/dist/* $(OUTPUT_DIR)/hook/$(GIT_TAG)/
	mkdir -p $(ARTIFACTS_PATH)
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

BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION:=$(shell cat GOLANG_VERSION)
REPO=hub
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base
BINARY_PLATFORMS=linux/amd64
IMAGE_PLATFORMS=linux/amd64

CEXEC_IMAGE_COMPONENT=tinkerbell/hub/cexec
KEXEC_IMAGE_COMPONENT=tinkerbell/hub/kexec
IMAGE2DISK_IMAGE_COMPONENT=tinkerbell/hub/image2disk
OCI2DISK_IMAGE_COMPONENT=tinkerbell/hub/oci2disk
WRITEFILE_IMAGE_COMPONENT=tinkerbell/hub/writefile
REBOOT_IMAGE_COMPONENT=tinkerbell/hub/reboot

BINARY_TARGET_FILES=cexec kexec image2disk oci2disk writefile
GO_MOD_PATHS=actions/cexec/v1 actions/kexec/v1 actions/image2disk/v1 actions/oci2disk/v1 actions/writefile/v1
# by default the first target is assumed to be the primary and will produce a attribution with no prefix
# for this build we want to prefix all the files
CEXEC_ATTRIBUTION_OVERRIDE=cexec

IMAGE_NAMES=$(BINARY_TARGET_FILES) reboot

DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

IMAGE_BUILD_ARGS=ATTRIBUTION_FILE
ATTRIBUTION_FILE=$(call TO_UPPER,$(IMAGE_NAME))_ATTRIBUTION.txt

PROJECT_DEPENDENCIES=eksa/torvalds/linux

FIX_LICENSES=$(foreach path,$(GO_MOD_PATHS),$(REPO)/$(path)/LICENSE)

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES)

# reboot action image requires touch.
# We are using eks-distro-minimal-base-glibc as the base to install touch
reboot/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-glibc

# bootconfig is dependency included in the writefile image, and bootconfig
# is not a static binary. It requires libs on the glibc image
writefile/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-glibc

$(REPO)/actions/%/LICENSE: | $(GO_MOD_DOWNLOAD_TARGETS)
#go-licenses requires a LICENSE file in each folder with the go.mod
	cp $(REPO)/LICENSE $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

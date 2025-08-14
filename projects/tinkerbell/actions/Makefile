BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=actions
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base

CEXEC_IMAGE_COMPONENT=tinkerbell/actions/cexec
KEXEC_IMAGE_COMPONENT=tinkerbell/actions/kexec
IMAGE2DISK_IMAGE_COMPONENT=tinkerbell/actions/image2disk
OCI2DISK_IMAGE_COMPONENT=tinkerbell/actions/oci2disk
WRITEFILE_IMAGE_COMPONENT=tinkerbell/actions/writefile
REBOOT_IMAGE_COMPONENT=tinkerbell/actions/reboot

BINARY_TARGET_FILES=cexec kexec image2disk oci2disk writefile
SOURCE_PATTERNS=./cexec ./kexec ./image2disk ./oci2disk ./writefile
GO_MOD_PATHS=. . . . .

IMAGE_NAMES=$(BINARY_TARGET_FILES) reboot

DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

PROJECT_DEPENDENCIES=eksa/torvalds/linux

EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk


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

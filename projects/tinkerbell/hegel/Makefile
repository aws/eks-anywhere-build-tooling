BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=hegel
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=hegel
SOURCE_PATTERNS=./cmd/hegel

IMAGE_USERADD_USER_NAME=tinkerbell

include $(BASE_DIRECTORY)/Common.mk

$(call IMAGE_TARGETS_FOR_NAME, hegel): hegel-useradd/images/export


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

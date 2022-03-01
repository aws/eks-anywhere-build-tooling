BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=eks-anywhere-packages
REPO_OWNER=aws
IMAGE_COMPONENT=eks-anywhere-packages
IMAGE_NAMES=$(IMAGE_COMPONENT)

BINARY_TARGET_FILES=package-manager
BASE_IMAGE_NAME=eks-distro-minimal-base

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

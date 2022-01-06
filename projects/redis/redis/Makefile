BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=redis
REPO_OWNER=redis
IMAGE_COMPONENT=redis/redis
IMAGE_BUILD_ARGS?=IMAGE_TAG

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc
SIMPLE_CREATE_BINARIES=false
GO_MOD_DOWNLOAD_TARGETS=$(GIT_PATCH_TARGET)
HAS_GO_LICENSES=false

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

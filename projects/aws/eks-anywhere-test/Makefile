BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=eks-anywhere-test
REPO_OWNER=aws
IMAGE_COMPONENT=eks-anywhere-test
IMAGE_BUILD_ARGS?=IMAGE_TAG
HAS_HELM_CHART=true

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc
SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false
REPO_NO_CLONE=true

BUILD_TARGETS=local-images
RELEASE_TARGETS=images

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

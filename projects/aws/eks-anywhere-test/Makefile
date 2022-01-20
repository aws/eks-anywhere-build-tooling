BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=eks-anywhere-test
REPO_OWNER=aws
IMAGE_COMPONENT=eks-anywhere-test
IMAGE_BUILD_ARGS?=IMAGE_TAG

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc
SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false
REPO_NO_CLONE=true

BUILD_TARGETS=local-images
RELEASE_TARGETS=images

HAS_HELM_CHART=true
HELM_ADDITIONAL_KEY_VALUES=\
  EKS_ANYWHERE_TAG=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} eks-anywhere-test $(LATEST)) \

include $(BASE_DIRECTORY)/Common.mk

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

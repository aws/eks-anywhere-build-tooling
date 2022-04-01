BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=packages/hello-eks-anywhere
REPO_OWNER=jonahjon
IMAGE_COMPONENT=hello-eks-anywhere
IMAGE_BUILD_ARGS?=IMAGE_TAG

BASE_IMAGE_NAME?=eks-distro-minimal-base-nginx
SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false
REPO_NO_CLONE=true

BUILD_TARGETS=local-images
RELEASE_TARGETS=images

HAS_HELM_CHART=true
HELM_IMAGE_LIST="$(REPO)"

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

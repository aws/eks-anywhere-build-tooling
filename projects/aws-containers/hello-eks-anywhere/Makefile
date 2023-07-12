BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION:=$(shell cat GOLANG_VERSION)
REPO=hello-eks-anywhere
REPO_OWNER=aws-containers
IMAGE_COMPONENT=$(REPO)
IMAGE_BUILD_ARGS?=IMAGE_TAG

BASE_IMAGE_NAME?=eks-distro-minimal-base-nginx
SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false

HAS_HELM_CHART=true
HELM_DIRECTORY=charts/hello-eks-anywhere
HELM_DESTINATION_REPOSITORY=$(REPO)
HELM_IMAGE_LIST=$(REPO)

EXCLUDE_FROM_STAGING_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

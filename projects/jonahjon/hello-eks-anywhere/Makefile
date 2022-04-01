BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=packages/hello-eks-anywhere
REPO_OWNER=jonahjon
IMAGE_COMPONENT=hello-eks-anywhere
IMAGE_NAMES=$(IMAGE_COMPONENT)

BASE_IMAGE_NAME?=eks-distro-minimal-base-nginx
SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false

HAS_HELM_CHART=true
HELM_DIRECTORY=charts/hello-eks-anywhere
HELM_IMAGE_LIST=packages/hello-eks-anywhere

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=eks-anywhere-packages
REPO_OWNER=aws
IMAGE_COMPONENT=eks-anywhere-packages
IMAGE_NAMES=$(IMAGE_COMPONENT)

BINARY_TARGET_FILES=package-manager
BASE_IMAGE_NAME=eks-distro-minimal-base

HAS_HELM_CHART=true
HELM_DIRECTORY=charts/eks-anywhere-packages
HELM_IMAGE_LIST=eks-anywhere-packages brancz/kube-rbac-proxy

# Needed to be built first for final helm/bundle builds
PROJECT_DEPENDENCIES=eksa/goharbor/harbor eksa/aws-containers/hello-eks-anywhere

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

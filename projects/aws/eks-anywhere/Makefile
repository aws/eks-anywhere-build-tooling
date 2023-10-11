BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)

REPO=eks-anywhere
REPO_OWNER=aws

BASE_IMAGE_NAME?=eks-distro-base

SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false

IMAGE_COMPONENT=eks-anywhere-diagnostic-collector
IMAGE_NAMES=$(IMAGE_COMPONENT)

PROJECT_DEPENDENCIES=eksd/kubernetes/client

REPO_NO_CLONE=true

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

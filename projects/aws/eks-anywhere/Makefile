BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=eks-anywhere
REPO_OWNER=aws

BASE_IMAGE_NAME?=eks-distro-base

SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false

IMAGE_COMPONENT=eks-anywhere-diagnostic-collector
IMAGE_NAMES=$(IMAGE_COMPONENT)

BUILD_TARGETS=local-images
RELEASE_TARGETS=images

FETCH_BINARIES_TARGETS=eksd/kubernetes/client

REPO_NO_CLONE=true

include $(BASE_DIRECTORY)/Common.mk


$(call IMAGE_TARGETS_FOR_NAME, eks-anywhere-diagnostic-collector): $(call FULL_FETCH_BINARIES_TARGETS, $(FETCH_BINARIES_TARGETS))


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

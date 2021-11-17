BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=eks-anywhere
REPO_OWNER=aws

BASE_IMAGE_NAME?=eks-distro-base

SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false

IMAGE_COMPONENT=eks-anywhere-diagnostic-collector

BUILD_TARGETS=local-images
RELEASE_TARGETS=images

include $(BASE_DIRECTORY)/Common.mk


$(call IMAGE_TARGETS_FOR_NAME, eks-anywhere-diagnostic-collector): fetch-binaries

.PHONY: local-images
local-images: eks-anywhere-diagnostic-collector/images/amd64
	
.PHONY: images
images: eks-anywhere-diagnostic-collector/images/push

.PHONY: fetch-binaries
fetch-binaries: kubernetes

kubernetes:
	build/fetch_binaries.sh

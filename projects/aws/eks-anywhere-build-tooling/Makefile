BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=eks-anwyhere-build-tooling
REPO_OWNER=aws

IMAGE_COMPONENT=eks-anywhere-cli-tools

BASE_IMAGE_NAME?=eks-distro-minimal-base-docker-client

SIMPLE_CREATE_BINARIES=false
HAS_S3_ARTIFACTS=false
HAS_LICENSES=false

BUILD_TARGETS=local-images
RELEASE_TARGETS=images

FETCH_BINARIES_TARGET=$(OUTPUT_DIR)/linux-amd64/eks-a-tools

include $(BASE_DIRECTORY)/Common.mk


$(call IMAGE_TARGETS_FOR_NAME, eks-anywhere-cli-tools): fetch-binaries

.PHONY: local-images
local-images: eks-anywhere-cli-tools/images/amd64

.PHONY: images
images: eks-anywhere-cli-tools/images/push

.PHONY: fetch-binaries
fetch-binaries: $(FETCH_BINARIES_TARGET)

$(FETCH_BINARIES_TARGET):
	build/fetch_binaries.sh $(ARTIFACTS_BUCKET)

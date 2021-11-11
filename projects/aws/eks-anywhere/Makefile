BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=eks-anywhere
REPO_OWNER=aws

BASE_IMAGE_NAME?=eks-distro-base

HAS_S3_ARTIFACTS=false
SIMPLE_CREATE_BINARIES=false
IMAGE_COMPONENT=eks-anywhere-diagnostic-collector


include $(BASE_DIRECTORY)/Common.mk


# overriding targets called by build + release to be no-ops 
# these do not apply here since there is no binary/checksums/attribution/licenses
validate-checksums $(GATHER_LICENSES_TARGET) generate-attribution $(OUTPUT_DIR)/ATTRIBUTION.txt: ;

$(call IMAGE_TARGETS_FOR_NAME, eks-anywhere-diagnostic-collector): fetch-binaries

.PHONY: local-images
local-images: eks-anywhere-diagnostic-collector/images/amd64
	
.PHONY: images
images: eks-anywhere-diagnostic-collector/images/push

.PHONY: fetch-binaries
fetch-binaries: kubernetes

kubernetes:
	build/fetch_binaries.sh

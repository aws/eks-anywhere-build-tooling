BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=eks-anywhere-test
REPO_OWNER=aws
IMAGE_COMPONENT=eks-anywhere-test

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc

include $(BASE_DIRECTORY)/Common.mk

# overriding targets called by build + release to be no-ops 
# these do not apply here since there is no binary/checksums/attribution/licenses
$(BINARY_TARGET) validate-checksums $(GATHER_LICENSES_TARGET) $(ATTRIBUTION_TARGET) $(OUTPUT_DIR)/ATTRIBUTION.txt: ;

.PHONY: local-images
local-images: eks-anywhere-test/images/amd64
	
.PHONY: images
images: eks-anywhere-test/images/push

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
HELM_REPO_VERSION?=$(shell cat HELM_REPO_VERSION)

REPO=cilium
REPO_OWNER=cilium

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk

# overriding targets called by build + release to be no-ops 
# these do not apply here since there is no binary/checksums/attribution/licenses
local-images images tarballs validate-checksums $(GATHER_LICENSES_TARGET) generate-attribution $(OUTPUT_DIR)/ATTRIBUTION.txt: ;

s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests:
	build/create_manifests.sh $(GIT_TAG) $(HELM_REPO_VERSION) $(ARTIFACTS_PATH)

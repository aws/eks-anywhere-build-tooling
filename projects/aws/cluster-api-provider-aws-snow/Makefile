BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=cluster-api-provder-aws-snow
REPO_OWNER=aws

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false
IMAGE_NAMES=

include $(BASE_DIRECTORY)/Common.mk


s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests:
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)

REPO=cluster-api-provider-aws-snow
REPO_OWNER=aws

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false
IMAGE_COMPONENT=aws/cluster-api-provider-aws-snow/manager

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk


s3-artifacts: create-manifests

cluster-api-provider-aws-snow/images/%: BASE_IMAGE=382577505035.dkr.ecr.us-west-2.amazonaws.com/cluster-api-provider-aws-snow-ecr:$(GIT_TAG)

.PHONY: create-manifests
create-manifests:
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

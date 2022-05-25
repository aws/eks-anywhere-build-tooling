BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.18"

REPO=rufio
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc

BINARY_TARGET_FILES=manager
SOURCE_PATTERNS=.

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk

s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests: tarballs
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(GOLANG_VERSION)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

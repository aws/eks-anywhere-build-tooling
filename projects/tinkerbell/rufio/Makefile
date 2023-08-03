BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=rufio
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc

BINARY_TARGET_FILES=manager
SOURCE_PATTERNS=.

HAS_S3_ARTIFACTS=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests: tarballs | ensure-yq $$(ENABLE_DOCKER)
	@build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(GOLANG_VERSION)

FIX_LICENSE_TARGET=$(REPO)/vendor/github.com/bmc-toolbox/common/LICENSE

# Ensure we download the license for the bmc-toolbox/common dependency ahead of gathering licenses.
$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSE_TARGET)

$(FIX_LICENSE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
	dest=$(REPO)/vendor/github.com/bmc-toolbox/common/LICENSE; \
	mkdir -p $$(dirname $$dest); \
	wget -q --retry-connrefused -O $$dest \
		https://raw.githubusercontent.com/bmc-toolbox/common/main/LICENSE;


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

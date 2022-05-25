BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
IPMITOOL_CLONE_URL=$(call GET_CLONE_URL,ipmitool,ipmitool)
IPMITOOL_GIT_TAG:=$(shell cat IPMITOOL_GIT_TAG)
GOLANG_VERSION?="1.18"

REPO=rufio
REPO_OWNER=tinkerbell

IPMITOOL_CLONE_TARGETS=$(MAKE_ROOT)/ipmitool/eks-anywhere-checkout-$(IPMITOOL_GIT_TAG)

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc

BINARY_TARGET_FILES=manager
SOURCE_PATTERNS=.

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk

$(call IMAGE_TARGETS_FOR_NAME, rufio): $(IPMITOOL_CLONE_TARGETS)

$(IPMITOOL_CLONE_TARGETS):
	git clone $(IPMITOOL_CLONE_URL) $(@D)
	git -C $(@D) checkout $(IPMITOOL_GIT_TAG)
	touch $@

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

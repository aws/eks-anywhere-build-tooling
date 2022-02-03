BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
BMCLIB_GIT_COMMIT=$(shell cat BMCLIB_GIT_COMMIT)
IPMITOOL_CLONE_URL=$(call GET_CLONE_URL,ipmitool,ipmitool)
BMCLIB_CLONE_URL=$(call GET_CLONE_URL,bmc-toolbox,bmclib)
IPMITOOL_GIT_TAG=$(shell cat IPMITOOL_GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=pbnj
REPO_OWNER=tinkerbell

IPMITOOL_CLONE_TARGETS=$(MAKE_ROOT)/ipmitool/eks-anywhere-checkout-$(IPMITOOL_GIT_TAG)
FIX_BMCLIB_TARGETS=$(REPO)/bmclib

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc

BINARY_TARGET_FILES=pbnj
SOURCE_PATTERNS=.

IMAGE_USERADD_USER_ID=10001
IMAGE_USERADD_USER_NAME=pbnj

include $(BASE_DIRECTORY)/Common.mk

FETCH_BINARIES_TARGETS=eksa/grpc-ecosystem/grpc-health-probe

$(GO_MOD_DOWNLOAD_TARGETS): $(FIX_BMCLIB_TARGETS)

$(call IMAGE_TARGETS_FOR_NAME, pbnj): $(IPMITOOL_CLONE_TARGETS) $(call FULL_FETCH_BINARIES_TARGETS, $(FETCH_BINARIES_TARGETS)) pbnj-useradd/images/export

$(IPMITOOL_CLONE_TARGETS):
	git clone $(IPMITOOL_CLONE_URL) $(@D)
	git -C $(@D) checkout $(IPMITOOL_GIT_TAG)
	touch $@

# PBNJ had a dependency on an older commit of bmclib module, that
# had an LGPL dependency (gebn/bmc). This was removed in 
# https://github.com/bmc-toolbox/bmclib/pull/257. However bmclib
# also bumped the versions of some logging libraries like logr,
# zap, logusr so directly using the PR's commit caused issues
# with other go mods in PBNJ. So cloning bmclib and patching it
# with that PR's changes
$(FIX_BMCLIB_TARGETS): | $(GIT_PATCH_TARGETS)
	cd $(REPO)
	git clone $(BMCLIB_CLONE_URL) $@
	git -C $@ checkout $(BMCLIB_GIT_COMMIT)
	git -C $@ config user.email prow@amazonaws.com
	git -C $@ config user.name "Prow Bot"
	git -C $@ am --committer-date-is-author-date $(MAKE_ROOT)/patches/bmclib/*
	touch $@/eks-anywhere-patched


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

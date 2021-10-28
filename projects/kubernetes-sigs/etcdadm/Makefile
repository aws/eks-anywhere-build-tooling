BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=etcdadm
REPO_OWNER=kubernetes-sigs

LICENSE_PACKAGE_FILTER?=.

GIT_PATCH_TARGET=$(REPO)/eks-anywhere-patched
BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/etcdadm

include $(BASE_DIRECTORY)/Common.mk


build: s3-artifacts
release: upload-artifacts

# Patch source before building
$(BINARY_TARGET): | $(GIT_PATCH_TARGET)

$(GIT_PATCH_TARGET): $(GIT_CHECKOUT_TARGET)
	git -C $(REPO) apply --verbose $(MAKE_ROOT)/patches/*
	@touch $@ 

.PHONY: s3-artifacts
s3-artifacts: tarballs
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)

.PHONY: local-images
local-images: #noop

.PHONY: images
images: #noop

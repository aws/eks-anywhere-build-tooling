BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=flux2
REPO_OWNER=fluxcd

LICENSE_PACKAGE_FILTER?=./cmd/flux

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/flux
GIT_PATCH_TARGET=$(REPO)/eks-anywhere-patched
FLUX_MANIFESTS_TARGET=$(REPO)/cmd/flux/manifests

BINARY_PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

include $(BASE_DIRECTORY)/Common.mk


build: s3-artifacts
release: upload-artifacts

# Patch source before building
$(BINARY_TARGET): $(FLUX_MANIFESTS_TARGET) | $(GIT_PATCH_TARGET)

$(GIT_PATCH_TARGET): $(GIT_CHECKOUT_TARGET)
	git -C $(REPO) apply --verbose $(MAKE_ROOT)/patches/*
	@touch $@ 

$(FLUX_MANIFESTS_TARGET): export PATH:=$(MAKE_ROOT)/$(OUTPUT_DIR):$(PATH)
$(FLUX_MANIFESTS_TARGET): $(KUSTOMIZE_TARGET)
	make -C $(REPO) cmd/flux/manifests

.PHONY: s3-artifacts
s3-artifacts: tarballs
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)

.PHONY: local-images
local-images: #noop

.PHONY: images
images: #noop

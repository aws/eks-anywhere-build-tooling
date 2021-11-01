BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=flux2
REPO_OWNER=fluxcd

LICENSE_PACKAGE_FILTER?=./cmd/flux

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/flux
FLUX_MANIFESTS_TARGET=$(REPO)/cmd/flux/manifests

BINARY_PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


# Patch source before building
$(BINARY_TARGET): $(FLUX_MANIFESTS_TARGET)

$(FLUX_MANIFESTS_TARGET): export PATH:=$(MAKE_ROOT)/$(OUTPUT_DIR):$(PATH)
$(FLUX_MANIFESTS_TARGET): $(KUSTOMIZE_TARGET)
	make -C $(REPO) cmd/flux/manifests

.PHONY: local-images
local-images: #noop

.PHONY: images
images: #noop

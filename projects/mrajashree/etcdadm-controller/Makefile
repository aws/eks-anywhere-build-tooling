BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=etcdadm-controller
REPO_OWNER=mrajashree

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/manager

BASE_IMAGE_NAME?=eks-distro-minimal-base

include $(BASE_DIRECTORY)/Common.mk


build: s3-artifacts
release: upload-artifacts

.PHONY: local-images
local-images: etcdadm-controller/images/amd64

.PHONY: images
images: etcdadm-controller/images/push

.PHONY: s3-artifacts
s3-artifacts: tarballs create-manifests
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)

.PHONY: create-manifests
create-manifests: export PATH:=$(MAKE_ROOT)/$(OUTPUT_DIR):$(PATH)
create-manifests: $(KUSTOMIZE_TARGET)
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

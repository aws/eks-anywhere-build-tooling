BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=etcdadm-controller
REPO_OWNER=mrajashree

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/manager

BASE_IMAGE_NAME?=eks-distro-minimal-base

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


s3-artifacts: create-manifests

.PHONY: local-images
local-images: etcdadm-controller/images/amd64

.PHONY: images
images: etcdadm-controller/images/push

.PHONY: create-manifests
create-manifests: export PATH:=$(MAKE_ROOT)/$(OUTPUT_DIR):$(PATH)
create-manifests: $(KUSTOMIZE_TARGET) tarballs
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

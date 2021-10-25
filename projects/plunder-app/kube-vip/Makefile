BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=kube-vip
REPO_OWNER=plunder-app

BASE_IMAGE_NAME?=eks-distro-minimal-base

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/kube-vip

include $(BASE_DIRECTORY)/Common.mk


build: s3-artifacts
release: upload-artifacts

.PHONY: local-images
local-images: kube-vip/images/amd64

.PHONY: images
images: kube-vip/images/push

.PHONY: s3-artifacts
s3-artifacts: tarballs
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)

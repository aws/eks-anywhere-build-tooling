BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=kube-rbac-proxy
REPO_OWNER=brancz

BASE_IMAGE_NAME?=eks-distro-minimal-base-nonroot

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/kube-rbac-proxy

KUBE_RBAC_PROXY_IMAGE=$(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

include $(BASE_DIRECTORY)/Common.mk


build: s3-artifacts
release: upload-artifacts

.PHONY: local-images
local-images: kube-rbac-proxy/images/amd64

.PHONY: images
images: kube-rbac-proxy/images/push

.PHONY: s3-artifacts
s3-artifacts: tarballs
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)



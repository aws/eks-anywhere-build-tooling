BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=cluster-api-provider-vsphere
REPO_OWNER=kubernetes-sigs

BASE_IMAGE_NAME?=eks-distro-minimal-base-nonroot

BINARY_TARGET_FILES=manager
EXTRA_GO_LDFLAGS=$(shell $(BUILD_LIB)/version.sh $(REPO) sigs.k8s.io/cluster-api-provider-vsphere/pkg/version)

IMAGE_COMPONENT=kubernetes-sigs/cluster-api-provider-vsphere/release/manager

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests: tarballs
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

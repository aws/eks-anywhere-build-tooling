BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=cluster-api-provider-tinkerbell
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base-nonroot

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET_FILES=manager
SOURCE_PATTERNS=.
EXTRA_GO_LDFLAGS=$(shell $(BUILD_LIB)/version.sh $(REPO) sigs.k8s.io/cluster-api-provider-packet/version)

IMAGE_COMPONENT=tinkerbell/cluster-api-provider-tinkerbell/release/manager

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


s3-artifacts: create-manifests

.PHONY: local-images
local-images: cluster-api-tinkerbell-controller/images/amd64

.PHONY: images
images: cluster-api-tinkerbell-controller/images/push

.PHONY: create-manifests
create-manifests: tarballs
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

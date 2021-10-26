BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.14"

REPO=cri-tools
REPO_OWNER=kubernetes-sigs

LICENSE_PACKAGE_FILTER?=./cmd/critest ./cmd/crictl

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/crictl

include $(BASE_DIRECTORY)/Common.mk


build: s3-artifacts
release: upload-artifacts

.PHONY: s3-artifacts
s3-artifacts: tarballs
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)

.PHONY: local-images
local-images: #noop

.PHONY: images
images: #noop

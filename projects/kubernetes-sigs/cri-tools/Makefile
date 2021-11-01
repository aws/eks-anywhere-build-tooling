BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.14"

REPO=cri-tools
REPO_OWNER=kubernetes-sigs

LICENSE_PACKAGE_FILTER?=./cmd/critest ./cmd/crictl

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/crictl

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


.PHONY: local-images
local-images: #noop

.PHONY: images
images: #noop

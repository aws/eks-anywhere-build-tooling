BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=govmomi
REPO_OWNER=vmware

LICENSE_PACKAGE_FILTER?=./govc

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/govc

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


.PHONY: local-images
local-images: #noop

.PHONY: images
images: #noop

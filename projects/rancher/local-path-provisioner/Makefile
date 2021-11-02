BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=local-path-provisioner
REPO_OWNER=rancher

BASE_IMAGE_NAME?=eks-distro-minimal-base

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET_FILES=local-path-provisioner

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


.PHONY: local-images
local-images: local-path-provisioner/images/amd64

.PHONY: images
images: local-path-provisioner/images/push

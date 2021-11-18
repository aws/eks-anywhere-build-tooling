BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=kube-vip
REPO_OWNER=plunder-app

BASE_IMAGE_NAME?=eks-distro-minimal-base

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET_FILES=kube-vip

include $(BASE_DIRECTORY)/Common.mk


.PHONY: local-images
local-images: kube-vip/images/amd64

.PHONY: images
images: kube-vip/images/push

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=local-path-provisioner
REPO_OWNER=rancher

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=local-path-provisioner

include $(BASE_DIRECTORY)/Common.mk

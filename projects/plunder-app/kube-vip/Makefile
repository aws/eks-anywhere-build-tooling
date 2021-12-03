BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=kube-vip
REPO_OWNER=plunder-app

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=kube-vip

include $(BASE_DIRECTORY)/Common.mk

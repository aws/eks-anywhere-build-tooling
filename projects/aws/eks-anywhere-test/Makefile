BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=eks-anywhere-test
REPO_OWNER=aws
IMAGE_COMPONENT=eks-anywhere-test
IMAGE_BUILD_ARGS?=IMAGE_TAG
HAS_HELM_CHART=true

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc
SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false

BUILD_TARGETS=local-images
RELEASE_TARGETS=images

include $(BASE_DIRECTORY)/Common.mk

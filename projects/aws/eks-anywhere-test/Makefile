BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=eks-anywhere-test
REPO_OWNER=aws
IMAGE_COMPONENT=eks-anywhere-test
IMAGE_DESCRIPTION="EKS Anywhere Add-ons Test Application"
IMAGE_PLATFORMS?=linux/amd64
IMAGE_BUILD_ARGS?=IMAGE_TAG

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc
HAS_S3_ARTIFACTS=false
SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false

BUILD_TARGETS=local-images
RELEASE_TARGETS=images

include $(BASE_DIRECTORY)/Common.mk

.PHONY: local-images
local-images: eks-anywhere-test/images/amd64 helm/build
	
.PHONY: images
images: eks-anywhere-test/images/push helm/push

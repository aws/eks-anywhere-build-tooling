BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=non-existent

REPO=ubuntu
REPO_OWNER=canonical

RTOS_BUCKET_NAME?=
RTOS_IMAGE_DATE=$(shell cat RTOS_IMAGE_DATE)
RELEASE_BRANCH=1-29

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false
IMAGE_NAMES=

REPO_NO_CLONE=true
GIT_CHECKOUT_TARGET=RTOS_IMAGE_DATE
FAKE_GIT_REPO_TARGET=.git

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

BUILD_TARGETS=upload-artifacts
RELEASE_TARGETS=upload-artifacts

include $(BASE_DIRECTORY)/Common.mk

$(FAKE_GIT_REPO_TARGET):
	@git init
	@git remote add origin https://github.com/aws/eks-anywhere-build-tooling.git 

s3-artifacts: download-rtos-image

.PHONY: download-rtos-image
download-rtos-image:
	build/download_rtos_image.sh $(RTOS_BUCKET_NAME) $(RTOS_IMAGE_DATE) $(ARTIFACTS_PATH) $(RELEASE_BRANCH)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

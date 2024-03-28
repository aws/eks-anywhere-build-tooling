BASE_DIRECTORY:=$(abspath ../../../)

REPO=ubuntu
REPO_OWNER=canonical
RTOS_SOURCE_URL=

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false
IMAGE_NAMES=
EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

BUILD_TARGETS=upload-artifacts
RELEASE_TARGETS=upload-artifacts

include $(BASE_DIRECTORY)/Common.mk

s3-artifacts: download-rtos-image

.PHONY: download-rtos-image
download-rtos-image:
	build/download_rtos_image.sh $(RTOS_SOURCE_URL) $(ARTIFACTS_PATH)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

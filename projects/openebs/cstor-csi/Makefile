BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=cstor-csi
REPO_OWNER=openebs

BINARY_TARGET_FILES=cstor-csi
SOURCE_PATTERNS=./cmd

CSTOR_CSI_IMAGE_COMPONENT=cstor-csi/cstor-csi

IMAGE_NAMES=cstor-csi
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################
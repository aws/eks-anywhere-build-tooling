BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=chartmuseum
REPO_OWNER=helm

BASE_IMAGE_NAME?=eks-distro-minimal-base
SOURCE_PATTERNS=./cmd/chartmuseum

BINARY_TARGET_FILES=chartm

HAS_S3_ARTIFACTS=true
IMAGE_NAMES=

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

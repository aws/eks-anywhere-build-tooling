BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=grpc-health-probe
REPO_OWNER=grpc-ecosystem

BINARY_TARGET_FILES=grpc_health_probe
SOURCE_PATTERNS=.

EXTRA_GOBUILD_FLAGS=-tags netgo

HAS_S3_ARTIFACTS=true

IMAGE_NAMES=

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

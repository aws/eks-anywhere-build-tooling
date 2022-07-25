BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=cstor-operators
REPO_OWNER=openebs

BINARY_TARGET_FILES=pool-manager cvc-operator volume-manager
SOURCE_PATTERNS=./cmd/pool-manager ./cmd/cvc-operator ./cmd/volume-manager

CVC_OPERATOR_IMAGE_COMPONENT=cstor-operators/cvc-operator
POOL_MANAGER_IMAGE_COMPONENT=cstor-operators/pool-manager
VOLUME_MANAGER_IMAGE_COMPONENT=cstor-operators/volume-manager

IMAGE_NAMES=pool-manager cvc-operator volume-manager
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################
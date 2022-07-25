BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=maya
REPO_OWNER=openebs

BINARY_TARGET_FILES=admission-server cstor-pool-mgmt cstor-volume-mgmt maya-apiserver
SOURCE_PATTERNS=./cmd/admission-server ./cmd/cstor-pool-mgmt ./cmd/cstor-volume-mgmt ./cmd/maya-apiserver

ADMISSION_SERVER_IMAGE_COMPONENT=maya/admission-server
CSTOR_POOL_MGMT_IMAGE_COMPONENT=maya/cstor-pool-mgmt
CSTOR_VOLUME_MGMT_IMAGE_COMPONENT=maya/cstor-volume-mgmt
MAYA_APISERVER_IMAGE_COMPONENT=maya/maya-apiserver

IMAGE_NAMES=admission-server cstor-pool-mgmt cstor-volume-mgmt maya-apiserver
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################
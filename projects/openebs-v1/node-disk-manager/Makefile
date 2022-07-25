BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=node-disk-manager
REPO_OWNER=openebs

# BINARY_TARGET_FILES=manager ndm-exporter
# SOURCE_PATTERNS=./cmd/manager ./cmd/ndm-exporter

# SIMPLE_CREATE_BINARIES=false

MANAGER_IMAGE_COMPONENT=node-disk-manager/manager
NDM_EXPORTER_IMAGE_COMPONENT=node-disk-manager/ndm-exporter

IMAGE_NAMES=manager ndm-exporter
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################
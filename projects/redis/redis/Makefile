BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=redis
REPO_OWNER=redis

SIMPLE_CREATE_BINARIES=false
GO_MOD_DOWNLOAD_TARGETS=
HAS_LICENSES=false
BINARY_TARGET_FILES=redis-server

BASE_IMAGE_NAME?=eks-distro-minimal-base

include $(BASE_DIRECTORY)/Common.mk

$(OUTPUT_BIN_DIR)/linux-amd64/redis-server:
	$(MAKE_ROOT)/build/create_binaries.sh $(MAKE_ROOT) $(REPO)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################


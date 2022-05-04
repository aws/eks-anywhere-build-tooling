BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)

REPO=linux
REPO_OWNER=torvalds

SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false
TAR_FILE_PREFIX=bootconfig
HAS_S3_ARTIFACTS=true

IMAGE_NAMES=

include $(BASE_DIRECTORY)/Common.mk

$(BINARY_TARGETS): create-binaries

.PHONY: create-binaries
create-binaries: $(GIT_CHECKOUT_TARGET)
	build/create_binaries.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH)

.PHONY: tarballs
tarballs: create-binaries
	build/create_tarballs.sh $(REPO) $(REPO_OWNER) $(GIT_TAG) $(TAR_FILE_PREFIX) $(OUTPUT_DIR) $(OUTPUT_DIR)/bin/linux $(ARTIFACTS_PATH) $(GIT_HASH)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=boots
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=boots
SOURCE_PATTERNS=./cmd/boots

GITREV?=$(shell git -C $(REPO) rev-parse --short HEAD)

EXTRA_GO_LDFLAGS=-X main.GitRev=${GITREV}

IPXE_DEPENDENCY_TARGET=$(REPO)/tftp/ipxe/undionly.kpxe

include $(BASE_DIRECTORY)/Common.mk

$(OUTPUT_BIN_DIR)/linux-amd64/boots $(OUTPUT_BIN_DIR)/linux-arm64/boots $(GATHER_LICENSES_TARGETS): $(IPXE_DEPENDENCY_TARGET)

$(IPXE_DEPENDENCY_TARGET):
# Boots has a dependency on 4 ipxe binaries that get embedded in the code.
# Currently we have these binaries under /ipxe folder.
# This is temporary, till we have a good way to build these binaries.
	mkdir $(REPO)/tftp/ipxe
	cp -a ./ipxe/. $(REPO)/tftp/ipxe/


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

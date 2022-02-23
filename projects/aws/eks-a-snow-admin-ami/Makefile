BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GOLANG_VERSION?=1.16
GIT_TAG=0.0.0-dev # This value does not matter but needs to be set

REPO=eks-a-snow-admin-ami
REPO_OWNER=aws

# This project is a bit of odd ball when it comes to using the standard makefile targets
# - there is no clone of source since the source for the bootstrap exist in this repo already
# - REPO target is overridden to create empty folder
# - the binary is built for the current os/arch so it can be run immediately 

BINARY_PLATFORMS=$(BUILDER_PLATFORM)
BINARY_TARGET_FILES=eks-a-snow-admin-ami
REPO_SUBPATH=..

BUILD_TARGETS=binaries
RELEASE_TARGETS=binaries

# fake git checkout
REPO_NO_CLONE=true
GIT_CHECKOUT_TARGET=main.go
HAS_LICENSES=false
IMAGE_NAMES=

EKSA_VERSION?=v0.7.0
EKSA_MANIFEST?=https://anywhere-assets.eks.amazonaws.com/releases/eks-a/manifest.yaml

include $(BASE_DIRECTORY)/Common.mk

.PHONY: setup
setup: binaries
	$(OUTPUT_BIN_DIR)/$(subst /,-,$(BINARY_PLATFORMS))/eks-a-snow-admin-ami setup

.PHONY: build-ami
build-ami: binaries
	$(OUTPUT_BIN_DIR)/$(subst /,-,$(BINARY_PLATFORMS))/eks-a-snow-admin-ami build --eksa-version $(EKSA_VERSION) --eksa-release-manifest-url $(EKSA_MANIFEST)

$(REPO):
	@mkdir $@

.PHONY: clean-extra
clean-extra:
	@rm -rf vendor

clean: clean-extra


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

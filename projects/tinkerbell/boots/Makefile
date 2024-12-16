BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=boots
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=smee
SOURCE_PATTERNS=./cmd/smee

GITREV?=$(shell git -C $(REPO) rev-parse --short HEAD)

EXTRA_GO_LDFLAGS=-X main.GitRev=${GITREV}

EXCLUDE_FROM_UPGRADE_BUILDSPEC=true
SKIP_CHECKSUM_VALIDATION=true

PROJECT_DEPENDENCIES=eksa/tinkerbell/ipxedust

OVERRIDE_VENDOR_BINARIES_TARGET=$(REPO)/vendor/github.com/tinkerbell/ipxedust/binary/eks-a-overridden

include $(BASE_DIRECTORY)/Common.mk


# smee embeds the ipxe binarys into the final go bin
# overriding them with the ones built via eks-a build tooling
$(OVERRIDE_VENDOR_BINARIES_TARGET): $(PROJECT_DEPENDENCIES_TARGETS) | $(GO_MOD_DOWNLOAD_TARGETS)
	for bin in "ipxe-efi.img" "ipxe.efi" "ipxe.iso" "undionly.kpxe"; do cp $(OUTPUT_DIR)/dependencies/linux-amd64/eksa/tinkerbell/ipxedust/$$bin $(@D); done
	for bin in "snp.efi"; do cp $(OUTPUT_DIR)/dependencies/linux-arm64/eksa/tinkerbell/ipxedust/$$bin $(@D); done
	touch $(@)

$(BINARY_TARGETS): $(OVERRIDE_VENDOR_BINARIES_TARGET)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

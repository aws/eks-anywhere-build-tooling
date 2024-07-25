BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)

REPO=ipxedust
REPO_OWNER=tinkerbell

SIMPLE_CREATE_BINARIES=false # not a go project, can't use this script "simple_create_binaries.sh"


# this is a native build which we can not cross compile
# ensure we are building to a directory based on the current
# host platform to avoid ever creating a different arch'd
# binary in the wrong folder
BINARY_PLATFORMS?=linux/$(BUILDER_PLATFORM_ARCH)
BINARY_TARGETS_AMD64=$(OUTPUT_BIN_DIR)/linux-amd64/ipxe.efi $(OUTPUT_BIN_DIR)/linux-amd64/undionly.kpxe
BINARY_TARGETS_ARM64=$(OUTPUT_BIN_DIR)/linux-arm64/snp.efi

# we are not building these because we do not need them
# but the smee import is going to reguire them so we 
# create empty files
FAKE_BINARY_TARGETS=$(OUTPUT_BIN_DIR)/linux-amd64/ipxe.iso $(OUTPUT_BIN_DIR)/linux-amd64/ipxe-efi.img
BINARY_TARGETS=$(if $(findstring amd64,$(BINARY_PLATFORMS)),$(BINARY_TARGETS_AMD64),) $(if $(findstring arm64,$(BINARY_PLATFORMS)),$(BINARY_TARGETS_ARM64),) $(FAKE_BINARY_TARGETS)

HAS_LICENSES=false #not go 

HAS_S3_ARTIFACTS=true
SKIP_CHECKSUM_VALIDATION=true

IMAGE_NAMES=

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

BUILDSPEC_VARS_KEYS=BINARY_PLATFORMS
BUILDSPEC_VARS_VALUES=BINARY_PLATFORMS
BUILDSPEC_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true


include $(BASE_DIRECTORY)/Common.mk

tarballs: gather-non-golang-licenses $(OUTPUT_DIR)/ATTRIBUTION.txt

$(BINARY_TARGETS_AMD64): PLATFORM=linux/amd64
$(BINARY_TARGETS_ARM64): PLATFORM=linux/arm64

$(FAKE_BINARY_TARGETS):
	touch $@

# since we are yum installing, we need to run the docker container as root
$(OUTPUT_BIN_DIR)/%: MAKEFLAGS=
$(OUTPUT_BIN_DIR)/%: $(GIT_PATCH_TARGET) | $$(call ENABLE_DOCKER_PLATFORM,$$(PLATFORM))
	sudo yum install -y genisoimage perl syslinux xz-devel; \
	mkdir -p $(@D); \
	$(MAKE) -C $(REPO) binary/$(@F); \
	cp $(REPO)/binary/$(@F) $@


.PHONY: gather-non-golang-licenses
gather-non-golang-licenses: $(GIT_PATCH_TARGET)
	source $(BUILD_LIB)/common.sh && build::non-golang::copy_licenses $(REPO)/upstream-$(shell cat $(REPO)/binary/script/ipxe.commit) $(OUTPUT_DIR)/LICENSES/github.com/ipxe/ipxe


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

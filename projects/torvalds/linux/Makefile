BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)

REPO=linux
REPO_OWNER=torvalds
REPO_SPARSE_CHECKOUT=tools lib/bootconfig.c include/linux/bootconfig.h

SIMPLE_CREATE_BINARIES=false

# this is a native build which we can not cross compile
# ensure we are building to a directory based on the current
# host platform to avoid ever creating a different arch'd
# binary in the wrong folder
BINARY_PLATFORMS?=linux/$(BUILDER_PLATFORM_ARCH)
BINARY_TARGETS=$(foreach platform,$(BINARY_PLATFORMS),$(OUTPUT_BIN_DIR)/$(subst /,-,$(platform))/tools/bootconfig)

HAS_LICENSES=false

TAR_FILE_PREFIX=bootconfig
HAS_S3_ARTIFACTS=true
SKIP_CHECKSUM_VALIDATION=true

IMAGE_NAMES=

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true

BUILDSPEC_VARS_KEYS=BINARY_PLATFORMS
BUILDSPEC_VARS_VALUES=BINARY_PLATFORMS
BUILDSPEC_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true


include $(BASE_DIRECTORY)/Common.mk


tarballs: gather-non-golang-licenses $(OUTPUT_DIR)/ATTRIBUTION.txt

$(OUTPUT_BIN_DIR)/%: MAKEFLAGS=
$(OUTPUT_BIN_DIR)/%: PLATFORM=$(subst -,/,$(*D:/tools=))
$(OUTPUT_BIN_DIR)/%: $(GIT_CHECKOUT_TARGET) | $$(call ENABLE_DOCKER_PLATFORM,$$(PLATFORM))
	@mkdir -p $(MAKE_ROOT)/$(@D); \
	$(MAKE) -C $(REPO)/tools/bootconfig OUTPUT=$(MAKE_ROOT)/$(@D)/

.PHONY: gather-non-golang-licenses
gather-non-golang-licenses: $(GIT_CHECKOUT_TARGET)
	source $(BUILD_LIB)/common.sh && build::non-golang::copy_licenses $(REPO) $(OUTPUT_DIR)/LICENSES/github.com/$(REPO_OWNER)/$(REPO)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

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

HAS_LICENSES=false

HAS_S3_ARTIFACTS=true
SKIP_CHECKSUM_VALIDATION=true

IMAGE_NAMES=

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

BUILDSPECS=buildspec.yml projects/$(COMPONENT)/buildspecs/combine-artifacts.yml
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE
BUILDSPEC_1_VARS_KEYS=BINARY_PLATFORMS
BUILDSPEC_1_VARS_VALUES=BINARY_PLATFORMS
BUILDSPEC_1_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER
BUILDSPEC_2_DEPENDS_ON_OVERRIDE=tinkerbell_ipxedust_linux_amd64 tinkerbell_ipxedust_linux_arm64


# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
# this is overridden to false when running the combine-artifacts target
UPLOAD_DO_NOT_DELETE=true

MODULE_SRC_TARBALL_TARGET=$(ARTIFACTS_PATH)/ipxedust-src.tar.gz


include $(BASE_DIRECTORY)/Common.mk

tarballs: gather-non-golang-licenses $(OUTPUT_DIR)/ATTRIBUTION.txt

$(BINARY_TARGETS_AMD64): PLATFORM=linux/amd64
$(BINARY_TARGETS_ARM64): PLATFORM=linux/arm64

$(FAKE_BINARY_TARGETS):
	mkdir -p $(@D)
	touch $@

# since we are yum installing, we need to run the docker container as root
$(OUTPUT_BIN_DIR)/%: MAKEFLAGS=
$(OUTPUT_BIN_DIR)/%: $(GIT_PATCH_TARGET) | $$(call ENABLE_DOCKER_PLATFORM,$$(PLATFORM))
	sudo yum install -y genisoimage perl syslinux xz-devel; \
	mkdir -p $(@D); \
	$(MAKE) -C $(REPO) binary/$(@F); \
	cp $(REPO)/binary/$(@F) $@


# download the binary tars, checkout and patch repo, copy eksa built ipxe on top of upstream vendor, tar rep as src
.PHONY: combine-artifacts
combine-artifacts: $(GIT_PATCH_TARGET) $(foreach arch,amd64 arm64,$(BINARY_DEPS_DIR)/linux-$(arch)/eksa/tinkerbell/ipxedust) overwrite-ipxe $(MODULE_SRC_TARBALL_TARGET)
	echo $(GIT_HASH) > $(ARTIFACTS_PATH)/githash
	$(MAKE) s3-artifacts SIMPLE_CREATE_TARBALLS=false EXPECTED_FILES_PATH=expected_artifacts_src

.PHONY: upload-src-artifacts
upload-src-artifacts:
	$(MAKE) combine-artifacts SIMPLE_CREATE_TARBALLS=false EXPECTED_FILES_PATH=expected_artifacts_src

.PHONY: overwrite-ipxe
overwrite-ipxe:
	for bin in "ipxe-efi.img" "ipxe.efi" "ipxe.iso" "undionly.kpxe"; do cp $(OUTPUT_DIR)/dependencies/linux-amd64/eksa/tinkerbell/ipxedust/$$bin $(REPO)/binary/; done
	for bin in "snp.efi"; do cp $(OUTPUT_DIR)/dependencies/linux-arm64/eksa/tinkerbell/ipxedust/$$bin $(REPO)/binary/; done
	cp -rf $(OUTPUT_DIR)/dependencies/linux-amd64/eksa/tinkerbell/ipxedust/ATTRIBUTION.txt $(REPO)/binary/IPXE_ATTRIBUTION.txt

$(MODULE_SRC_TARBALL_TARGET):
	mkdir -p $(@D)
	git -C $(REPO) archive --format=tar.gz --output=$(@) HEAD

.PHONY: gather-non-golang-licenses
gather-non-golang-licenses: $(GIT_PATCH_TARGET)
	$(MAKE) -C $(REPO) extract-ipxe
	mkdir -p $(OUTPUT_DIR)/LICENSES/github.com/ipxe/ipxe
	cp $(REPO)/upstream-$(shell cat $(REPO)/binary/script/ipxe.commit)/{COPYING,COPYING.GPLv2} $(OUTPUT_DIR)/LICENSES/github.com/ipxe/ipxe


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

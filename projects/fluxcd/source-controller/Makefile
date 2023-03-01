BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=source-controller
REPO_OWNER=fluxcd

BINARY_TARGET_FILES=source-controller

FIX_LICENSES_XEIPUUV_TARGET=$(REPO)/vendor/github.com/xeipuuv/gojsonpointer/LICENSE.txt
FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/source-controller/api/LICENSE
LIBGIT2_LICENSES_TARGET=$(OUTPUT_DIR)/LICENSES/github.com/libgit2/libgit2/COPYING
BASE_IMAGE_NAME?=eks-distro-minimal-base-git

CGO_DEPS_TARGET=$(CGO_SOURCE)/linux-%/lib64/libgit2.so
CGO_CREATE_BINARIES=true
EXTRA_GOBUILD_FLAGS=-tags netgo,osusergo,static_build

# These artifacts are only used to coordinate between code build batch builds
# they will not exist in the final eks-a bundle
HAS_S3_ARTIFACTS?=true

BUILDSPECS=projects/$(COMPONENT)/buildspecs/binaries.yml projects/$(COMPONENT)/buildspecs/images.yml
BUILDSPEC_1_VARS_KEYS=BINARY_PLATFORMS
BUILDSPEC_1_VARS_VALUES=BINARY_PLATFORMS
BUILDSPEC_1_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER
BUILDSPEC_2_DEPENDS_ON_OVERRIDE=fluxcd_source_controller_linux_amd64 fluxcd_source_controller_linux_arm64

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true

# the checksum periodic job in prow has space issues due to using buildctl to build the arm binary
# skipping for now in this build until we can improve the job
BINARY_PLATFORMS?=linux/amd64 linux/arm64
BINARY_PLATFORMS:=$(if $(findstring periodic,$(JOB_TYPE)),linux/amd64,$(BINARY_PLATFORMS))

include $(BASE_DIRECTORY)/Common.mk


# source-controller depends on building with cgo to include git2go
$(OUTPUT_BIN_DIR)/linux-amd64/source-controller: | $(CGO_SOURCE)/linux-amd64/usr/lib64/libgit2.so
$(OUTPUT_BIN_DIR)/linux-arm64/source-controller: | $(CGO_SOURCE)/linux-arm64/usr/lib64/libgit2.so

$(CGO_SOURCE)/linux-%/usr/lib64/libgit2.so:
	$(MAKE) binary-builder-deps/cgo/$* IMAGE_TARGET=deps IMAGE_OUTPUT=dest=$(CGO_SOURCE)/linux-$* DOCKERFILE_FOLDER=./docker/build

$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_XEIPUUV_TARGET) $(FIX_LICENSES_API_LICENSE_TARGET) $(LIBGIT2_LICENSES_TARGET)

$(FIX_LICENSES_XEIPUUV_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The xeipuuv dependency github repos all have licenses however they all do not have go.mod files
# checked in to the repo. Hence we need to manually download licenses from Github for each of them 
# and place them in the respective folders under vendor directory so that they is available for 
# go-licenses to pick up	
	for package in gojsonpointer gojsonreference gojsonschema ; do \
		dest=$(REPO)/vendor/github.com/xeipuuv/$$package/LICENSE.txt; \
		mkdir -p $$(dirname $$dest); \
		wget --retry-connrefused -q https://raw.githubusercontent.com/xeipuuv/$$package/master/LICENSE-APACHE-2.0.txt -O \
			$$dest; \
	done;

$(FIX_LICENSES_API_LICENSE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# Internal go.mod under /api directory
	@mkdir -p $(@D)
	cp $(REPO)/LICENSE $@

$(LIBGIT2_LICENSES_TARGET):
	$(BASE_DIRECTORY)/build/lib/gather_non_golang_licenses.sh libgit2/libgit2 v1.3.1 $(MAKE_ROOT)/$(OUTPUT_DIR)

# in codebuild batch builds we seperate binary builds from image builds to work around
# issue building arm on an amd instance
# use the standard handle-dependencies process and extract to $(OUTPUT_BIN_DIR) and validate checksums match what is checked in
.PHONY: download-binaries-from-bucket
download-binaries-from-bucket: fake-binary-prereqs $(foreach arch,amd64 arm64,$(BINARY_DEPS_DIR)/linux-$(arch)/eksa/fluxcd/source-controller) move-downloaded-artifacts validate-checksums

# reset the binary timestamps to make sure make doesnt think they are older than the pre-req targets, like go-mod-download
.PHONY: move-downloaded-artifacts
move-downloaded-artifacts:
	for arch in amd64 arm64; do \
		mkdir -p $(OUTPUT_BIN_DIR)/linux-$$arch; \
	 	mv $(BINARY_DEPS_DIR)/linux-$$arch/eksa/fluxcd/source-controller/source-controller $(OUTPUT_BIN_DIR)/linux-$$arch; \
		touch $(OUTPUT_BIN_DIR)/linux-$$arch/source-controller; \
	done
	rm -rf $(BINARY_DEPS_DIR)

# When downloading the binaries vs building the standard make targets (release/images/etc) will still try to rebuild since
# the binary targets have prereqs on the libgit2 cgo deps target fake these during this flow
.PHONY: fake-binary-prereqs
fake-binary-prereqs: $(GO_MOD_DOWNLOAD_TARGETS)
	mkdir -p $(CGO_SOURCE)/linux-amd64/usr/lib64 $(CGO_SOURCE)/linux-arm64/usr/lib64
	touch $(CGO_SOURCE)/linux-amd64/usr/lib64/libgit2.so $(CGO_SOURCE)/linux-arm64/usr/lib64/libgit2.so
	

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

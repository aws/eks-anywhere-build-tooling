BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
LIBYANG_GIT_TAG=$(shell cat LIBYANG_GIT_TAG)
TINI_GIT_TAG=$(shell cat TINI_GIT_TAG)

REPO=frr
REPO_OWNER=frrouting

# All sources are fetched as checksum-verified archives rather than Git clones.
REPO_NO_CLONE=true
SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false
LICENSES_TARGETS_FOR_PREREQ=gather-non-golang-licenses $(OUTPUT_DIR)/ATTRIBUTION.txt validate-source-checksums

IMAGE_NAMES=frr
# Keep these global so combine-images can override them for its manifest stage.
DOCKERFILE_FOLDER=./docker/linux
IMAGE_BUILD_ARGS=
BASE_IMAGE_NAME?=eks-distro-minimal-base-python
BASE_IMAGE_OS_VERSION=al2023
COMPILER_IMAGE_VERSION=3.11

# Compile on an instance of each architecture, then combine the image manifests.
BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE
IMAGE_PLATFORMS?=linux/$(BUILDER_PLATFORM_ARCH)
NATIVE_BUILD?=true
EXCLUDE_FROM_STAGING_BUILDSPEC=true
EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true
SKIP_ON_RELEASE_BRANCH=true

SOURCE_VERSION_FILES=GIT_TAG LIBYANG_GIT_TAG TINI_GIT_TAG
SOURCES_TARGET=$(OUTPUT_DIR)/src/.sources-ready

include $(BASE_DIRECTORY)/Common.mk

combine-images: IMAGE_PLATFORMS=linux/amd64,linux/arm64
combine-images: NATIVE_BUILD=false

# This is a checked-in static notice, not a Go-generated attribution target.
ATTRIBUTION.txt: ;
$(OUTPUT_DIR)/ATTRIBUTION.txt: ATTRIBUTION.txt
	@mkdir -p $(@D)
	cp $< $@

$(SOURCES_TARGET): $(SOURCE_VERSION_FILES) CHECKSUMS build/fetch_sources.sh
	./build/fetch_sources.sh "$(GIT_TAG)" "$(LIBYANG_GIT_TAG)" "$(TINI_GIT_TAG)" "$(OUTPUT_DIR)"
	@touch $@

.PHONY: fetch-sources validate-source-checksums gather-non-golang-licenses validate-native-platform
fetch-sources: $(SOURCES_TARGET) # Download and verify the three pinned source archives

checkout-repo: fetch-sources
validate-checksums: validate-source-checksums
validate-source-checksums: $(SOURCES_TARGET) # Verify the pinned source archive checksums
	sha256sum -c CHECKSUMS

# Common.mk's Go attribution generator does not apply to these C projects.
# Preserve FRR's complete license directory as well as the root notices.
gather-non-golang-licenses: $(SOURCES_TARGET) # Collect FRR, libyang and tini license notices
	source $(BUILD_LIB)/common.sh && \
		build::non-golang::copy_licenses $(OUTPUT_DIR)/src/frr $(OUTPUT_DIR)/LICENSES/github.com/frrouting/frr && \
		build::non-golang::copy_licenses $(OUTPUT_DIR)/src/libyang $(OUTPUT_DIR)/LICENSES/github.com/CESNET/libyang && \
		build::non-golang::copy_licenses $(OUTPUT_DIR)/src/tini $(OUTPUT_DIR)/LICENSES/github.com/krallin/tini
	cp -R $(OUTPUT_DIR)/src/frr/doc/licenses $(OUTPUT_DIR)/LICENSES/github.com/frrouting/frr/

gather-licenses: gather-non-golang-licenses
attribution: gather-non-golang-licenses $(OUTPUT_DIR)/ATTRIBUTION.txt

frr/images/amd64: IMAGE_PLATFORMS=linux/amd64
frr/images/arm64: IMAGE_PLATFORMS=linux/arm64
frr/images/amd64 frr/images/arm64 frr/images/push: validate-native-platform
validate-native-platform:
	@if [ "$(NATIVE_BUILD)" = "true" ] && [ "$(IMAGE_PLATFORMS)" != "linux/$(BUILDER_PLATFORM_ARCH)" ]; then \
		echo "FRR must be compiled on a native $(IMAGE_PLATFORMS) builder; this host is linux/$(BUILDER_PLATFORM_ARCH)."; \
		exit 1; \
	fi


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

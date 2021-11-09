BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=source-controller
REPO_OWNER=fluxcd

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET_FILES=source-controller
SOURCE_PATTERNS=.

FIX_LICENSES_XEIPUUV_TARGET=$(REPO)/vendor/github.com/xeipuuv/gojsonpointer/LICENSE.txt
FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/source-controller/api/LICENSE
LIBGIT2_LICENSES_TARGET=$(OUTPUT_DIR)/LICENSES/github.com/libgit2/libgit2/COPYING
BASE_IMAGE_NAME?=eks-distro-minimal-base-git
BUILDER_IMAGE?=$(BASE_IMAGE_REPO)/builder-base:latest

CGO_DEPS_TARGET=$(CGO_SOURCE)/linux-%/lib64/libgit2.so
CGO_CREATE_BINARIES=true

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


# source-controller depends on building with cgo to include git2go
$(OUTPUT_BIN_DIR)/linux-amd64/source-controller: | $(CGO_SOURCE)/linux-amd64/usr/lib64/libgit2.so
$(OUTPUT_BIN_DIR)/linux-arm64/source-controller: | $(CGO_SOURCE)/linux-arm64/usr/lib64/libgit2.so

$(CGO_SOURCE)/linux-%/usr/lib64/libgit2.so:
	$(MAKE) binary-builder-deps/cgo/$* IMAGE_TARGET=deps IMAGE_OUTPUT=dest=$(CGO_SOURCE)/linux-$*


# $(BINARY_TARGET):
# 	for platform in $(CGO_PLATFORMS); do \
# 		$(MAKE) source-controller-deps/cgo/$${platform#linux/} IMAGE_TARGET=deps IMAGE_OUTPUT=dest=$(CGO_SOURCE); \
# 		$(MAKE) source-controller/cgo/$${platform#linux/} IMAGE_OUTPUT=dest=$(OUTPUT_DIR); \
# 	done;
# 	@cp $(OUTPUT_DIR)/attribution/ATTRIBUTION.txt ./

$(GATHER_LICENSES_TARGET): | $(FIX_LICENSES_XEIPUUV_TARGET) $(FIX_LICENSES_API_LICENSE_TARGET) $(LIBGIT2_LICENSES_TARGET)

.PHONY: local-images
local-images: source-controller/images/amd64

.PHONY: images
images: source-controller/images/push

$(FIX_LICENSES_XEIPUUV_TARGET): | $(GIT_CHECKOUT_TARGET)
# The xeipuuv dependency github repos all have licenses however they all do not have go.mod files
# checked in to the repo. Hence we need to manually download licenses from Github for each of them 
# and place them in the respective folders under vendor directory so that they is available for 
# go-licenses to pick up	
	for package in gojsonpointer gojsonreference gojsonschema ; do \
		dest=$(REPO)/vendor/github.com/xeipuuv/$$package/LICENSE.txt; \
		mkdir -p $$(dirname $$dest); \
		wget -q https://raw.githubusercontent.com/xeipuuv/$$package/master/LICENSE-APACHE-2.0.txt -O \
			$$dest; \
	done;

$(FIX_LICENSES_API_LICENSE_TARGET): | $(GIT_CHECKOUT_TARGET)
# Internal go.mod under /api directory
	@mkdir -p $(@D)
	cp $(REPO)/LICENSE $@

$(LIBGIT2_LICENSES_TARGET):
	$(BASE_DIRECTORY)/build/lib/gather_non_golang_licenses.sh libgit2/libgit2 v1.1.0 $(MAKE_ROOT)/$(OUTPUT_DIR)

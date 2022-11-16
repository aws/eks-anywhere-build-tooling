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

# This is an artificial dep so that during the staging batch build the source-controller
# falls into the second "batch" of builds, not blocking the first batch from finishing
# since it takes significantly longer
BUILDSPEC_DEPENDS_ON_OVERRIDE=eksa/kubernetes-sigs/etcdadm

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


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

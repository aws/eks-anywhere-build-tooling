BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.20"

REPO=source-controller
REPO_OWNER=fluxcd

BINARY_TARGET_FILES=source-controller

FIX_LICENSES_XEIPUUV_TARGET=$(REPO)/vendor/github.com/xeipuuv/gojsonpointer/LICENSE.txt
FIX_LICENSES_ALIBABACLOUD_TARGET=$(REPO)/vendor/github.com/alibabacloud-go/cr-20160607/LICENSE
FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/source-controller/api/LICENSE
FIX_LICENSES_SHA1CD=$(REPO)/vendor/github.com/pjbgf/sha1cd/LICENSE
BASE_IMAGE_NAME?=eks-distro-minimal-base

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_XEIPUUV_TARGET) $(FIX_LICENSES_API_LICENSE_TARGET) $(FIX_LICENSES_SHA1CD) $(FIX_LICENSES_ALIBABACLOUD_TARGET)

$(FIX_LICENSES_XEIPUUV_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The alibabacloud repos includes a LICENSE (apache2) in a later version than the currently depended on version
# Copying in manually for now until the dep is updated upstream	
	for package in gojsonpointer gojsonreference gojsonschema ; do \
		dest=$(REPO)/vendor/github.com/xeipuuv/$$package/LICENSE.txt; \
		mkdir -p $$(dirname $$dest); \
		wget --retry-connrefused -q https://raw.githubusercontent.com/xeipuuv/$$package/master/LICENSE-APACHE-2.0.txt -O \
			$$dest; \
	done;

$(FIX_LICENSES_SHA1CD): | $(GO_MOD_DOWNLOAD_TARGETS)
# pjbgf/sha1cd includes a LICENSE (apache2) in a later version than the currently depended on version
# Copying in manually for now until the dep is updated upstream
	wget -q --retry-connrefused -O $@ https://raw.githubusercontent.com/pjbgf/sha1cd/main/LICENSE

$(FIX_LICENSES_ALIBABACLOUD_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The xeipuuv dependency github repos all have licenses however they all do not have go.mod files
# checked in to the repo. Hence we need to manually download licenses from Github for each of them 
# and place them in the respective folders under vendor directory so that they is available for 
# go-licenses to pick up	
	for package in cr-20160607 tea-xml  ; do \
		dest=$(REPO)/vendor/github.com/alibabacloud-go/$$package/LICENSE.txt; \
		mkdir -p $$(dirname $$dest); \
		wget --retry-connrefused -q https://raw.githubusercontent.com/alibabacloud-go/$$package/master/LICENSE -O \
			$$dest; \
	done;


$(FIX_LICENSES_API_LICENSE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# Internal go.mod under /api directory
	@mkdir -p $(@D)
	cp $(REPO)/LICENSE $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=helm-controller
REPO_OWNER=fluxcd

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET_FILES=helm-controller

FIX_LICENSES_XEIPUUV_TARGET=$(REPO)/vendor/github.com/xeipuuv/gojsonpointer/LICENSE.txt
FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/helm-controller/api/LICENSE

BASE_IMAGE_NAME=eks-distro-minimal-base

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGET): $(FIX_LICENSES_XEIPUUV_TARGET) $(FIX_LICENSES_API_LICENSE_TARGET)

.PHONY: local-images
local-images: helm-controller/images/amd64

.PHONY: images
images: helm-controller/images/push

$(FIX_LICENSES_XEIPUUV_TARGET): | $(GIT_CHECKOUT_TARGET)
# The xeipuuv dependency github repos all have licenses however they all do not have go.mod files
# checked in to the repo. Hence we need to manually download licenses from Github for each of them 
# and place them in the respective folders under vendor directory so that they is available for 
# go-licenses to pick up	
	for package in gojsonpointer gojsonreference gojsonschema ; do \
		wget -q https://raw.githubusercontent.com/xeipuuv/$$package/master/LICENSE-APACHE-2.0.txt -O \
			$(REPO)/vendor/github.com/xeipuuv/$$package/LICENSE.txt; \
	done;

$(FIX_LICENSES_API_LICENSE_TARGET): | $(GIT_CHECKOUT_TARGET)
# Internal go.mod under /api directory
	cp $(REPO)/LICENSE $@

BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=helm-controller
REPO_OWNER=fluxcd

BINARY_TARGET_FILES=helm-controller

FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/helm-controller/api/LICENSE

BASE_IMAGE_NAME=eks-distro-minimal-base

include $(BASE_DIRECTORY)/Common.mk

$(FIX_LICENSES_API_LICENSE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# Internal go.mod under /api directory
	cp $(REPO)/LICENSE $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

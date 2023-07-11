BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.20"

REPO=notification-controller
REPO_OWNER=fluxcd

BINARY_TARGET_FILES=notification-controller

FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/notification-controller/api/LICENSE

BASE_IMAGE_NAME=eks-distro-minimal-base

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_API_LICENSE_TARGET)

$(FIX_LICENSES_API_LICENSE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# Internal go.mod under /api directory
	cp $(REPO)/LICENSE $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

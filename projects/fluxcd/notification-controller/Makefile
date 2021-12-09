BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=notification-controller
REPO_OWNER=fluxcd

BINARY_TARGET_FILES=notification-controller

FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/notification-controller/api/LICENSE
FIX_LICENSES_K0KUBUN_TARGET=$(REPO)/vendor/github.com/k0kubun/pp/LICENSE.txt

BASE_IMAGE_NAME=eks-distro-minimal-base
IMAGE_USERADD_USER_NAME=controller

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_K0KUBUN_TARGET) $(FIX_LICENSES_API_LICENSE_TARGET)

$(call IMAGE_TARGETS_FOR_NAME, notification-controller): notification-controller-useradd/images/export

$(FIX_LICENSES_API_LICENSE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# Internal go.mod under /api directory
	cp $(REPO)/LICENSE $@

$(FIX_LICENSES_K0KUBUN_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The notification-controller project consumes an older version (v2.3.0) of the k0kubun/pp module 
# which does not have a LICENSE file committed to the repo. Hence we are fetching the license file 
# from the master branch so that it is available in the vendor directory for go-licenses to pick up
	wget -q -O $@  https://raw.githubusercontent.com/k0kubun/pp/master/LICENSE.txt


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=kustomize-controller
REPO_OWNER=fluxcd

BINARY_TARGET_FILES=kustomize-controller

DEFAULT_KUBECONFIG_TARGET=$(OUTPUT_DIR)/files/kubeconfig
FIX_LICENSES_MOZILLA_TARGET=$(REPO)/vendor/go.mozilla.org/gopgagent/LICENSE.txt
FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/kustomize-controller/api/LICENSE

FETCH_BINARIES_TARGETS=eksd/kubernetes/client

BASE_IMAGE_NAME=eks-distro-minimal-base-git
IMAGE_USERADD_USER_NAME=controller

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_MOZILLA_TARGET) $(FIX_LICENSES_API_LICENSE_TARGET)

$(call IMAGE_TARGETS_FOR_NAME, kustomize-controller): $(call FULL_FETCH_BINARIES_TARGETS, $(FETCH_BINARIES_TARGETS)) $(DEFAULT_KUBECONFIG_TARGET) kustomize-controller-useradd/images/export

$(DEFAULT_KUBECONFIG_TARGET): | $(GIT_CHECKOUT_TARGET)
	@mkdir $(@D)
	@cp $(REPO)/config/kubeconfig $@ 

$(FIX_LICENSES_MOZILLA_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The mozilla-services/gopgagent does not have a license file checked into the repo, but there is a currently open PR
# https://github.com/mozilla-services/gopgagent/pull/4 which adds it. Until this is merged, we need to fetch the license
# file from the commit and not from main/master.
	wget -q -O $@ \
		https://raw.githubusercontent.com/mozilla-services/gopgagent/39936d55b621318e919509000af38573d91c42ad/LICENSE.txt

$(FIX_LICENSES_API_LICENSE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# Internal go.mod under /api directory
	cp $(REPO)/LICENSE $@

BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=kustomize-controller
REPO_OWNER=fluxcd

BINARY_TARGET_FILES=kustomize-controller

FIX_LICENSES_GETSOPS_TARGET=$(REPO)/vendor/github.com/getsops/gopgagent/LICENSE.txt
FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/kustomize-controller/api/LICENSE

BASE_IMAGE_NAME=eks-distro-minimal-base-git

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_GETSOPS_TARGET) $(FIX_LICENSES_API_LICENSE_TARGET)

$(FIX_LICENSES_GETSOPS_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# In version v1.2.0, kuztomize-controller moved to using the getsops/gopgagent module instead of mozilla-services/gopgagent.
# The new module is a fork of the old one and neither of them have a license file checked into the repo, but there is an
# open PR https://github.com/mozilla-services/gopgagent/pull/4 which adds the liense to the parent repo. Until getsops/gopgagent
# adds a license file, we need to fetch the license from this pull request.
	wget -q --retry-connrefused -O $@ \
		https://raw.githubusercontent.com/mozilla-services/gopgagent/39936d55b621318e919509000af38573d91c42ad/LICENSE.txt

$(FIX_LICENSES_API_LICENSE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# Internal go.mod under /api directory
	cp $(REPO)/LICENSE $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

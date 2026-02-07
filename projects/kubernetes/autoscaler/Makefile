BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat ./$(RELEASE_BRANCH)/GIT_TAG)
GOLANG_VERSION=$(shell cat ./$(RELEASE_BRANCH)/GOLANG_VERSION)

REPO=autoscaler
REPO_OWNER=kubernetes

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=cluster-autoscaler
SOURCE_PATTERNS=.
GO_MOD_PATHS=cluster-autoscaler

FIX_LICENSES=$(REPO)/cluster-autoscaler/LICENSE

HAS_RELEASE_BRANCHES=true
EXCLUDE_FROM_STAGING_BUILDSPEC=true
SKIP_ON_RELEASE_BRANCH=true

HAS_HELM_CHART=true

HELM_GIT_TAG=$(shell cat ./$(RELEASE_BRANCH)/HELM_GIT_TAG)
HELM_TAG=$(subst cluster-autoscaler-chart-,,$(HELM_GIT_TAG))-$(subst -,.,$(RELEASE_BRANCH))-$(GIT_HASH)
HELM_SOURCE_OWNER=kubernetes
HELM_SOURCE_REPOSITORY=autoscaler

# cluster-autoscaler has migrated its chart since v0.55.0
HELM_DIRECTORY=cluster-autoscaler/charts/cluster-autoscaler
HELM_IMAGE_LIST=kubernetes/autoscaler

HELM_CHART_NAMES=cluster-autoscaler/charts/cluster-autoscaler

BUILDSPEC_VARS_KEYS=RELEASE_BRANCH
BUILDSPEC_VARS_VALUES=SUPPORTED_K8S_VERSIONS
BUILDSPEC_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

REMOVE_CLOUD_PROVIDERS_TARGET=$(REPO)/cluster-autoscaler/cloudprovider/eks-anywhere-cloud-providers-removed-$(GIT_TAG)

include $(BASE_DIRECTORY)/Common.mk


$(GO_MOD_DOWNLOAD_TARGETS): $(REMOVE_CLOUD_PROVIDERS_TARGET)
attributions: $(REMOVE_CLOUD_PROVIDERS_TARGET)

$(REMOVE_CLOUD_PROVIDERS_TARGET): $(GIT_PATCH_TARGET)
	@echo "Remove cloud providers"
	cd ./autoscaler/cluster-autoscaler/cloudprovider && ls . | grep -v -e builder -e clusterapi -e mocks -e test -e .go | xargs rm -rf
	touch $@


.PHONY: update-go-mod-and-create-patch
update-go-mod-and-create-patch: $(REMOVE_CLOUD_PROVIDERS_TARGET)
	@echo "Running 'go mod tidy'"
	cd $(REPO)/cluster-autoscaler && go mod tidy
	@echo "Creating patch..."
	git -C  $(REPO) add cluster-autoscaler/go.mod cluster-autoscaler/go.sum
	git -C  $(REPO) commit -m "Update go.mod Dependencies"
	git -C  $(REPO) format-patch -1
	@echo "Copying patch to release branch patches directory..."
	mv $(REPO)/*.patch $(BASE_DIRECTORY)/projects/kubernetes/autoscaler/$(RELEASE_BRANCH)/patches/


$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES)

$(FIX_LICENSES): | $(GO_MOD_DOWNLOAD_TARGETS)
#go-licenses requires a LICENSE file in each folder with the go.mod
	cp $(REPO)/LICENSE $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

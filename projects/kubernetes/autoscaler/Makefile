BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG=$(shell cat ./$(RELEASE_BRANCH)/GIT_TAG)
GOLANG_VERSION=$(shell cat ./$(RELEASE_BRANCH)/GOLANG_VERSION)

REPO=autoscaler
REPO_OWNER=kubernetes

BASE_IMAGE_NAME?=eks-distro-minimal-base
SKIPPED_K8S_VERSIONS=1-20

BINARY_TARGET_FILES=cluster-autoscaler
SOURCE_PATTERNS=.
GO_MOD_PATHS=cluster-autoscaler

FIX_LICENSES=$(REPO)/cluster-autoscaler/LICENSE
FIX_LICENSES_AZURE_TARGET=$(REPO)/cluster-autoscaler/vendor/github.com/Azure/azure-sdk-for-go/LICENSE.txt

HAS_RELEASE_BRANCHES=true
EXCLUDE_FROM_STAGING_BUILDSPEC=true

HAS_HELM_CHART=true

HELM_GIT_TAG=$(shell cat ./$(RELEASE_BRANCH)/HELM_GIT_TAG)
HELM_TAG=$(subst cluster-autoscaler-chart-,,$(HELM_GIT_TAG))-$(subst -,.,$(RELEASE_BRANCH))-$(GIT_HASH)
HELM_SOURCE_OWNER=kubernetes
HELM_SOURCE_REPOSITORY=autoscaler
HELM_DIRECTORY=charts/cluster-autoscaler
HELM_DESTINATION_REPOSITORY=cluster-autoscaler/charts/cluster-autoscaler
HELM_IMAGE_LIST=kubernetes/autoscaler

BUILDSPEC_VARS_KEYS=RELEASE_BRANCH
BUILDSPEC_VARS_VALUES=SUPPORTED_K8S_VERSIONS
BUILDSPEC_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES) $(FIX_LICENSES_AZURE_TARGET) 

$(FIX_LICENSES): | $(GO_MOD_DOWNLOAD_TARGETS)
#go-licenses requires a LICENSE file in each folder with the go.mod
	cp $(REPO)/LICENSE $@

$(FIX_LICENSES_AZURE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The azure sdk dependency github repo has a LICENSE however it does not have a go.mod
# checked in to the repo. Hence we need to manually download LICENSE from Github
# and place them in the vendor directory so that they is available for go-licenses to pick up
	wget -q --retry-connrefused https://raw.githubusercontent.com/Azure/azure-sdk-for-go/main/LICENSE.txt -O $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

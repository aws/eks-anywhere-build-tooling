BASE_DIRECTORY:=$(abspath ../../../)
EKS_ADDON_IMAGE_TAG=$(shell cat EKS_ADDON_IMAGE_TAG)
GIT_TAG=$(EKS_ADDON_IMAGE_TAG)
HELM_GIT_TAG=$(shell cat HELM_GIT_TAG)

 # Upstream images are used directly without re-building and re-tagging in build
 # tooling, so the IMAGE_TAG needs to be overwritten to match upstream ECR tags.
IMAGE_TAG=$(EKS_ADDON_IMAGE_TAG)

REPO_OWNER=kubernetes-sigs
REPO=metrics-server

SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false
IMAGE_NAMES=

BUILD_TARGETS=helm/build
RELEASE_TARGETS=helm/push
EXCLUDE_FROM_STAGING_BUILDSPEC=true
SKIP_ON_RELEASE_BRANCH=true

HAS_HELM_CHART=true

HELM_SOURCE_IMAGE_REPO=602401143452.dkr.ecr.us-west-2.amazonaws.com
HELM_SOURCE_OWNER=kubernetes-sigs
HELM_SOURCE_REPOSITORY=metrics-server
HELM_DIRECTORY=charts/metrics-server
HELM_IMAGE_LIST=eks/metrics-server
HELM_USE_UPSTREAM_IMAGE=true

HELM_CHART_NAMES=metrics-server/charts/metrics-server

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

# To avoid issues cloning/checking out the repo, git_tag is set to
# a valid branch/tag
$(GIT_CHECKOUT_TARGET): GIT_TAG=$(word 1,$(subst -, ,$(EKS_ADDON_IMAGE_TAG)))

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

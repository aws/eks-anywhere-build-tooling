BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG=$(shell source $(BUILD_LIB)/common.sh && build::eksd_releases::get_eksd_component_asset_image_tag metrics-server metrics-server-image $(RELEASE_BRANCH))
HELM_GIT_TAG:=$(shell cat $(RELEASE_BRANCH)/HELM_GIT_TAG)
 # Upstream images are used directly without re-building and re-tagging in build
 # tooling, so the IMAGE_TAG needs to be overwritten to match upstream ECR tags.
IMAGE_TAG=$(GIT_TAG)

HAS_RELEASE_BRANCHES=true
SKIPPED_K8S_VERSIONS=1-20

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

HELM_SOURCE_IMAGE_REPO=public.ecr.aws
HELM_SOURCE_OWNER=kubernetes-sigs
HELM_SOURCE_REPOSITORY=metrics-server
HELM_DIRECTORY=charts/metrics-server
HELM_DESTINATION_REPOSITORY=metrics-server/charts/metrics-server
HELM_IMAGE_LIST=eks-distro/kubernetes-sigs/metrics-server
HELM_USE_UPSTREAM_IMAGE=true

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

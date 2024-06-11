BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=charts
REPO_OWNER=tinkerbell

HAS_HELM_CHART=true
HELM_DIRECTORY=tinkerbell

BUILD_TARGETS=helm/build
RELEASE_TARGETS=helm/push

HELM_IMAGE_LIST=
HELM_CHART_NAMES=$(REPO_OWNER)/$(REPO)
HELM_CHART_FOLDER=charts/stack

EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

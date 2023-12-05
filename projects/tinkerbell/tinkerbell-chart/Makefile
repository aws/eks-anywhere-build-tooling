BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)

REPO_OWNER=tinkerbell
REPO=tinkerbell-chart

IMAGE_NAMES=
REPO_NO_CLONE=true
BUILD_TARGETS=helm/build
RELEASE_TARGETS=helm/push

HAS_HELM_CHART=true
HELM_DIRECTORY=chart
HELM_SOURCE_REPOSITORY=.
HELM_IMAGE_LIST=
HELM_CHART_NAMES=$(REPO_OWNER)/$(REPO)

EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

verify:
	scripts/verify.sh $(BASE_DIRECTORY)/projects/$(HELM_CHART_NAMES)/chart $(GIT_TAG)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

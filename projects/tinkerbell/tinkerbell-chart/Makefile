BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)

REPO_OWNER=tinkerbell
REPO=tinkerbell-chart

REPO_NO_CLONE=true
BUILD_TARGETS=helm/build
RELEASE_TARGETS=helm/push

HAS_HELM_CHART=true
HELM_DIRECTORY=chart
HELM_SOURCE_REPOSITORY=.
HELM_DESTINATION_REPOSITORY=$(REPO_OWNER)/$(REPO)

include $(BASE_DIRECTORY)/Common.mk

verify:
	scripts/verify.sh $(BASE_DIRECTORY)/projects/$(HELM_DESTINATION_REPOSITORY)/chart $(GIT_TAG)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

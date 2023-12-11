BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=metallb
REPO_OWNER=metallb

BINARY_TARGET_FILES=controller speaker
SOURCE_PATTERNS=go.universe.tf/metallb/controller go.universe.tf/metallb/speaker


CONTROLLER_IMAGE_COMPONENT=metallb/controller
SPEAKER_IMAGE_COMPONENT=metallb/speaker

IMAGE_NAMES=controller speaker
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

HAS_HELM_CHART=true
HELM_DIRECTORY=charts/metallb
HELM_IMAGE_LIST=metallb/controller metallb/speaker

HELM_CHART_NAMES=metallb/crds metallb/metallb

EXCLUDE_FROM_STAGING_BUILDSPEC=true
SKIP_ON_RELEASE_BRANCH=true

include $(BASE_DIRECTORY)/Common.mk


metallb/metallb/helm/%: PACKAGE_DEPENDENCIES=metallb-crds

metallb/crds/helm/%: HELM_DIRECTORY=charts/crds

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

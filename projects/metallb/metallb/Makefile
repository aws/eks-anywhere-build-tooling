BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

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
HELM_DESTINATION_REPOSITORY=metallb/metallb
HELM_IMAGE_LIST=metallb/controller metallb/speaker

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

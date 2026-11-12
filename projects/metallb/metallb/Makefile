BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=metallb
REPO_OWNER=metallb
BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=controller speaker
SOURCE_PATTERNS=go.universe.tf/metallb/controller go.universe.tf/metallb/speaker


CONTROLLER_IMAGE_COMPONENT=metallb/controller
SPEAKER_IMAGE_COMPONENT=metallb/speaker

IMAGE_NAMES=controller speaker
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

HAS_HELM_CHART=true
HELM_DIRECTORY=charts/metallb
FRR_K8S_GIT_TAG=$(shell cat ../frr-k8s/GIT_TAG)
FRR_GIT_TAG=$(shell cat ../../frrouting/frr/GIT_TAG)
HELM_IMAGE_LIST=metallb/controller metallb/speaker metallb/frr-k8s frrouting/frr
HELM_IMAGE_TAG_LIST=$(IMAGE_TAG) $(IMAGE_TAG) $(FRR_K8S_GIT_TAG)-$(GIT_HASH) $(FRR_GIT_TAG)-$(GIT_HASH)

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

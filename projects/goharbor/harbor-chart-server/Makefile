BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=chartmuseum
REPO_OWNER=helm
REPO_SUBPATH=./cmd/chartmuseum

BASE_IMAGE_NAME?=eks-distro-base

IMAGE_COMPONENT=harbor/harbor-chart-server
BINARY_TARGET_FILES=harbor-chart-server
EXTRA_GOBUILD_FLAGS=-mod=vendor

$(call IMAGE_TARGETS_FOR_NAME, harbor-chart-server): harbor-chart-server-useradd/images/export

harbor-chart-server-useradd/images/export: IMAGE_USERADD_USER_NAME=chart IMAGE_USERADD_USER_ID=10000

PRECOMPILE_PREPARE=chartmuseum_prepare

harbor_base_namespace=goharbor
IMAGE_BUILD_ARGS=harbor_base_namespace

COMPONENT=goharbor/harbor-chart-server
CLONE_URL=https://github.com/$(REPO_OWNER)/$(REPO).git

HARBOR_VERSION=v2.4.1

chartmuseum_prepare: 
	cd $(PROJECT_ROOT) \
	&& git clone https://github.com/goharbor/harbor.git harbor \
	&& git -C harbor checkout -f $(HARBOR_VERSION) \
	&& cd harbor && git apply $(PROJECT_ROOT)/../harbor/patches/0001-remove-base-image-check-for-cert-installation.patch && cd .. \
	&& cd $(REPO) && git apply ../harbor/make/photon/chartserver/redis.patch

include $(BASE_DIRECTORY)/Common.mk

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

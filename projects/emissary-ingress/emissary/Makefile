BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

APPMESH_VERSION:=$(shell cat APPMESH_VERSION)
IMAGE_BUILD_ARGS=APPMESH_VERSION COMPILER_IMAGE

REPO=emissary
REPO_OWNER=emissary-ingress

EXCLUDE_FROM_STAGING_BUILDSPEC=true
BINARY_TARGET_FILES=busyambassador capabilities_wrapper example-envoy-metrics-sink k8sregistryctl kat-client kat-server
SOURCE_PATTERNS=./cmd/busyambassador ./cmd/capabilities_wrapper ./cmd/example-envoy-metrics-sink ./cmd/k8sregistryctl ./cmd/kat-client ./cmd/kat-server

EMISSARY_IMAGE_COMPONENT=emissary-ingress/emissary

IMAGE_NAMES=emissary
DOCKERFILE_FOLDER=./docker/linux/emissary-image

BASE_IMAGE_NAME?=eks-distro-minimal-base-python
BASE_IMAGE_OS_VERSION=al2023

COMPILER_IMAGE_VERSION=3.9

HAS_HELM_CHART=true
HELM_DIRECTORY=charts/emissary-ingress
HELM_DESTINATION_REPOSITORY=emissary-ingress/emissary
HELM_IMAGE_LIST=emissary-ingress/emissary

AMBASSADOR_VERSION_TARGET=$(REPO)/python/ambassador.version

include $(BASE_DIRECTORY)/Common.mk

$(GO_MOD_DOWNLOAD_TARGETS): $(AMBASSADOR_VERSION_TARGET)

# do not pass our makeflags down to emissary make since there will be unncessary warnings
$(AMBASSADOR_VERSION_TARGET): MAKEFLAGS=
$(AMBASSADOR_VERSION_TARGET):
	@export VERSION="v3.6.0" && export CHART_VERSION=$(CHART_VERSION) && cd emissary && GOPATH=/var && make python/ambassador.version

# Build helm chart
.PHONY: helm/build
helm/build: $(if $(wildcard $(PROJECT_ROOT)/helm/patches),$(HELM_GIT_PATCH_TARGET),)
	# Emissary requires 2 charts. Consider this a temporary work-around.
	$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) emissary/crds charts/crds $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_require.sh $(HELM_SOURCE_IMAGE_REPO) emissary/crds  $(OUTPUT_DIR) $(IMAGE_TAG) $(HELM_TAG) $(PROJECT_ROOT) $(LATEST) $(HELM_USE_UPSTREAM_IMAGE) $(HELM_IMAGE_LIST)
	$(BUILD_LIB)/helm_replace.sh emissary/crds $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) emissary/crds
	$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) $(HELM_DESTINATION_REPOSITORY) $(HELM_DIRECTORY) $(OUTPUT_DIR)
	PACKAGE_DEPENDENCIES=emissary-crds $(BUILD_LIB)/helm_require.sh $(HELM_SOURCE_IMAGE_REPO) $(HELM_DESTINATION_REPOSITORY) $(OUTPUT_DIR) $(IMAGE_TAG) $(HELM_TAG) $(PROJECT_ROOT) $(LATEST) $(HELM_USE_UPSTREAM_IMAGE) $(HELM_IMAGE_LIST)
	$(BUILD_LIB)/helm_replace.sh $(HELM_DESTINATION_REPOSITORY) $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) $(HELM_DESTINATION_REPOSITORY)

.PHONY: helm/push
helm/push: helm/build
	$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) $(HELM_DESTINATION_REPOSITORY) $(HELM_TAG) $(GIT_TAG) $(OUTPUT_DIR) $(LATEST)
	$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) emissary-ingress/crds $(HELM_TAG) $(GIT_TAG) $(OUTPUT_DIR) $(LATEST)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=flux2
REPO_OWNER=fluxcd
IMAGE_COMPONENT=fluxcd/flux-cli
IMAGE_NAMES=flux-cli

BASE_IMAGE_NAME=eks-distro-minimal-base

BINARY_TARGET_FILES=flux
SOURCE_PATTERNS=./cmd/flux
EXTRA_GO_LDFLAGS=-X main.VERSION=$(GIT_TAG:v%=%)

FLUX_MANIFESTS_TARGET=$(REPO)/cmd/flux/manifests

BINARY_PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
FETCH_BINARIES_TARGETS=eksd/kubernetes/client

HAS_S3_ARTIFACTS=true
HAS_HELM_CHART=true
HELM_DESTINATION_REPOSITORY=fluxcd/flux2
HELM_SOURCE_OWNER=fluxcd-community
HELM_SOURCE_REPOSITORY=helm-charts
HELM_GIT_TAG?=$(shell cat HELM_GIT_TAG)
HELM_DIRECTORY=charts/flux2
HELM_ADDITIONAL_KEY_VALUES=\
  FLUX_CLI_TAG=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} fluxcd/flux-cli $(LATEST)) \
  HELM_CONTROLLER_TAG=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} fluxcd/helm-controller $(LATEST)) \
  KUSTOMIZE_CONTROLLER_TAG=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} fluxcd/kustomize-controller $(LATEST)) \
  NOTIFICATION_CONTROLLER_TAG=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} fluxcd/notification-controller $(LATEST)) \
  SOURCE_CONTROLLER_TAG=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} fluxcd/source-controller $(LATEST))

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS) $(OUTPUT_BIN_DIR)/linux-amd64/flux $(OUTPUT_BIN_DIR)/linux-arm64/flux: $(FLUX_MANIFESTS_TARGET) $(call FULL_FETCH_BINARIES_TARGETS, $(FETCH_BINARIES_TARGETS))

$(FLUX_MANIFESTS_TARGET): export PATH:=$(MAKE_ROOT)/$(OUTPUT_DIR):$(PATH)
$(FLUX_MANIFESTS_TARGET): $(KUSTOMIZE_TARGET)
	# EMBEDDED_MANIFESTS_TARGET from https://github.com/fluxcd/flux2/blob/main/Makefile#L2
	make -C $(REPO) cmd/flux/.manifests.done


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

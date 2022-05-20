BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=flux2
REPO_OWNER=fluxcd
IMAGE_COMPONENT=fluxcd/flux-cli
IMAGE_NAMES=flux-cli

BASE_IMAGE_NAME=eks-distro-minimal-base

BINARY_TARGET_FILES=flux
SOURCE_PATTERNS=./cmd/flux
EXTRA_GO_LDFLAGS=-X main.VERSION=$(GIT_TAG:v%=%)

FLUX_MANIFESTS_TARGET=$(REPO)/cmd/flux/manifests

HAS_S3_ARTIFACTS=true
# HAS_HELM_CHART=true
# HELM_DESTINATION_REPOSITORY=fluxcd/flux2
# HELM_SOURCE_OWNER=fluxcd-community
# HELM_SOURCE_REPOSITORY=helm-charts
# HELM_GIT_TAG?=$(shell cat HELM_GIT_TAG)
# HELM_DIRECTORY=charts/flux2
# HELM_IMAGE_LIST="fluxcd/flux-cli fluxcd/helm-controller fluxcd/kustomize-controller fluxcd/notification-controller fluxcd/source-controller"

# Needed to be built for the helm chart build, which happens in this project
PROJECT_DEPENDENCIES=eksa/fluxcd/source-controller eksd/kubernetes/client

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS) $(OUTPUT_BIN_DIR)/linux-amd64/flux $(OUTPUT_BIN_DIR)/linux-arm64/flux: $(FLUX_MANIFESTS_TARGET)

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

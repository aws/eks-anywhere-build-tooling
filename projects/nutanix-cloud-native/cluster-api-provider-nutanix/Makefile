BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=cluster-api-provider-nutanix
REPO_OWNER=nutanix-cloud-native

BASE_IMAGE_NAME?=eks-distro-minimal-base-nonroot

BINARY_TARGET_FILES=manager
SOURCE_PATTERNS=.
EXTRA_GO_LDFLAGS=$(shell $(BUILD_LIB)/version.sh $(REPO) sigs.k8s.io/cluster-api-provider-nutanix/pkg/version)

HAS_S3_ARTIFACTS=true

FIX_LICENSES_STORAGE_GO_CLIENT_TARGET=$(REPO)/vendor/github.com/nutanix/ntnx-api-golang-clients/storage-go-client/v4/LICENSE.txt

TOOLS_BIN_DIR := $(shell pwd)/hack/tools/bin
CONTROLLER_GEN := $(TOOLS_BIN_DIR)/controller-gen
GO_BIN_PATH=$(shell source $(BUILD_LIB)/common.sh && build::common::get_go_path $(GOLANG_VERSION))
GO=$(GO_BIN_PATH)/go

include $(BASE_DIRECTORY)/Common.mk

$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_STORAGE_GO_CLIENT_TARGET)

.PHONY: $(FIX_LICENSES_STORAGE_GO_CLIENT_TARGET)
$(FIX_LICENSES_STORAGE_GO_CLIENT_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)	
	wget -q --retry-connrefused https://raw.githubusercontent.com/nutanix/ntnx-api-golang-clients/refs/heads/main/LICENSE.md -O \
		$(REPO)/vendor/github.com/nutanix/ntnx-api-golang-clients/storage-go-client/v4/LICENSE.txt; \

s3-artifacts: create-manifests

$(TOOLS_BIN_DIR):
	@mkdir -p $(TOOLS_BIN_DIR)

$(CONTROLLER_GEN): $(TOOLS_BIN_DIR)
	GOBIN=$(TOOLS_BIN_DIR) $(GO) install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.14.0

.PHONY: create-manifests
create-manifests: export PATH:=$(TOOLS_BIN_DIR):$(MAKE_ROOT)/$(OUTPUT_DIR):$(PATH)
create-manifests: $(CONTROLLER_GEN) $(KUSTOMIZE_TARGET) tarballs | ensure-yq $$(ENABLE_DOCKER)
	@build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG) $(GOLANG_VERSION)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

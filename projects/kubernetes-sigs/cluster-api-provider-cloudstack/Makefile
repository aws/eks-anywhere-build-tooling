BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=cluster-api-provider-cloudstack
REPO_OWNER=kubernetes-sigs

BASE_IMAGE_NAME?=eks-distro-minimal-base-nonroot

BINARY_TARGET_FILES=manager
EXTRA_GO_LDFLAGS=$(shell $(BUILD_LIB)/version.sh $(REPO) kubernetes-sigs/cluster-api-provider-cloudstack/pkg/version)

IMAGE_COMPONENT=kubernetes-sigs/cluster-api-provider-cloudstack/release/manager
CAPC_MANIFESTS_TARGET=$(REPO)/pkg/mocks/mock_client.go

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk

$(REPO)/eks-anywhere-go-mod-download: $(CAPC_MANIFESTS_TARGET)

$(CAPC_MANIFESTS_TARGET): $(GIT_CHECKOUT_TARGET) | $$(ENABLE_DOCKER)
	@source $(BUILD_LIB)/common.sh && build::common::use_go_version $(GOLANG_VERSION) && $(MAKE) -C $(REPO) generate-mocks

s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests: export PATH:=$(MAKE_ROOT)/$(OUTPUT_DIR):$(PATH)
create-manifests: $(KUSTOMIZE_TARGET) tarballs | $$(call ENABLE_DOCKER)
	@build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG) $(GOLANG_VERSION)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

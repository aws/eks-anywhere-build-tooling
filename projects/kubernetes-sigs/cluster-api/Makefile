BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=cluster-api
REPO_OWNER=kubernetes-sigs

LICENSE_PACKAGE_FILTER?=./cmd/clusterctl ./bootstrap/kubeadm ./controlplane/kubeadm ./

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/clusterctl

GIT_PATCH_TARGET=$(REPO)/eks-anywhere-patched

BASE_IMAGE_NAME=eks-distro-minimal-base-nonroot

CAPI_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/cluster-api-controller
CAPI_KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller
CAPI_KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-control-plane-controller
CAPD_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/capd-manager

CAPD_GATHER_LICENSES_TARGET=$(OUTPUT_DIR)/capd/attribution/go-license.csv
CAPD_ATTRIBUTION_TARGET=CAPD_ATTRIBUTION.txt

include $(BASE_DIRECTORY)/Common.mk


build: $(CAPD_ATTRIBUTION_TARGET) s3-artifacts
release: upload-artifacts
clean: clean-binary-deps
attribution: $(CAPD_ATTRIBUTION_TARGET)

# Patch source before building
$(BINARY_TARGET): | $(GIT_PATCH_TARGET)

$(GIT_PATCH_TARGET): $(GIT_CHECKOUT_TARGET)
	git -C $(REPO) apply --verbose $(MAKE_ROOT)/patches/*
	@touch $@ 

$(CAPD_GATHER_LICENSES_TARGET): $(BINARY_TARGET)
	$(BASE_DIRECTORY)/build/lib/gather_licenses.sh $(REPO) $(MAKE_ROOT)/$(OUTPUT_DIR)/capd . test/infrastructure/docker

$(CAPD_ATTRIBUTION_TARGET): $(CAPD_GATHER_LICENSES_TARGET)
	$(BASE_DIRECTORY)/build/lib/create_attribution.sh $(MAKE_ROOT) $(GOLANG_VERSION) $(MAKE_ROOT)/$(OUTPUT_DIR)/capd CAPD_ATTRIBUTION.txt

.PHONY: local-images
local-images: capi-controller/images/amd64 capi-kubeadm-bootstrap-controller/images/amd64 capi-kubeadm-control-plane-controller/images/amd64 capd-controller/images/amd64 

.PHONY: images
images: capi-controller/images/push capi-kubeadm-bootstrap-controller/images/push capi-kubeadm-control-plane-controller/images/push capd-controller/images/push 

capd-controller/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-docker-client
capd-controller/images/%: DOCKERFILE_FOLDER=./docker/linux/cluster-api-docker-controller
$(call IMAGE_TARGETS_FOR_NAME, capd-controller): fetch-binaries $(CAPD_GATHER_LICENSES_TARGET)

capi-controller/images/%: DOCKERFILE_FOLDER=./docker/linux/cluster-api-controller

capi-kubeadm-bootstrap-controller/images/%: DOCKERFILE_FOLDER=./docker/linux/kubeadm-bootstrap-controller

capi-kubeadm-control-plane-controller/images/%: DOCKERFILE_FOLDER=./docker/linux/kubeadm-control-plane-controller

.PHONY: s3-artifacts
s3-artifacts: tarballs create-manifests
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)

.PHONY: create-manifests
create-manifests:
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

kubernetes:
	build/fetch_binaries.sh

.PHONY: fetch-binaries
fetch-binaries: kubernetes

.PHONY: clean-binary-deps
clean-binary-deps:
	@rm -rf kubernetes
	@rm -f kubernetes-client-linux-amd64.tar.gz

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=cluster-api
REPO_OWNER=kubernetes-sigs

LICENSE_PACKAGE_FILTER?=./cmd/clusterctl ./bootstrap/kubeadm ./controlplane/kubeadm ./

BINARY_TARGET_FILES=clusterctl manager kubeadm-bootstrap-manager kubeadm-control-plane-manager cluster-api-provider-docker-manager
SOURCE_PATTERNS=sigs.k8s.io/cluster-api/cmd/clusterctl sigs.k8s.io/cluster-api sigs.k8s.io/cluster-api/bootstrap/kubeadm sigs.k8s.io/cluster-api/controlplane/kubeadm sigs.k8s.io/cluster-api/test/infrastructure/docker
EXTRA_GO_LDFLAGS=$(shell $(BUILD_LIB)/version.sh $(REPO) sigs.k8s.io/cluster-api/cmd/version)

GIT_PATCH_TARGET=$(REPO)/eks-anywhere-patched

BASE_IMAGE_NAME=eks-distro-minimal-base-nonroot
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CLUSTER_API_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/cluster-api-controller
KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller
KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-control-plane-controller
CLUSTER_API_DOCKER_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/capd-manager

CAPD_GATHER_LICENSES_TARGET=$(OUTPUT_DIR)/capd/attribution/go-license.csv
CAPD_ATTRIBUTION_TARGET=CAPD_ATTRIBUTION.txt

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk

$(OUTPUT_BIN_DIR)/linux-%/cluster-api-provider-docker-manager: REPO_SUBPATH=test/infrastructure/docker
$(OUTPUT_BIN_DIR)/linux-%/cluster-api-provider-docker-manager: EXTRA_GO_LDFLAGS=
	
s3-artifacts: create-manifests

clean: clean-binary-deps
generate-attribution: $(CAPD_ATTRIBUTION_TARGET)

$(CAPD_GATHER_LICENSES_TARGET): $(BINARY_TARGETS)
	$(BASE_DIRECTORY)/build/lib/gather_licenses.sh $(REPO) $(MAKE_ROOT)/$(OUTPUT_DIR)/capd . test/infrastructure/docker

$(CAPD_ATTRIBUTION_TARGET): $(CAPD_GATHER_LICENSES_TARGET)
	$(BASE_DIRECTORY)/build/lib/create_attribution.sh $(MAKE_ROOT) $(GOLANG_VERSION) $(MAKE_ROOT)/$(OUTPUT_DIR)/capd CAPD_ATTRIBUTION.txt

.PHONY: local-images
local-images: cluster-api-controller/images/amd64 kubeadm-bootstrap-controller/images/amd64 kubeadm-control-plane-controller/images/amd64 cluster-api-docker-controller/images/amd64 

.PHONY: images
images: cluster-api-controller/images/push kubeadm-bootstrap-controller/images/push kubeadm-control-plane-controller/images/push cluster-api-docker-controller/images/push 

cluster-api-docker-controller/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-docker-client
$(call IMAGE_TARGETS_FOR_NAME, cluster-api-docker-controller): fetch-binaries $(CAPD_GATHER_LICENSES_TARGET)

.PHONY: create-manifests
create-manifests: tarballs
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

kubernetes:
	build/fetch_binaries.sh

.PHONY: fetch-binaries
fetch-binaries: kubernetes

.PHONY: clean-binary-deps
clean-binary-deps:
	@rm -rf kubernetes
	@rm -f kubernetes-client-linux-amd64.tar.gz

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=cluster-api
REPO_OWNER=kubernetes-sigs

LICENSE_PACKAGE_FILTER?=./cmd/clusterctl ./bootstrap/kubeadm ./controlplane/kubeadm ./

BINARY_TARGET_FILES=clusterctl manager kubeadm-bootstrap-manager kubeadm-control-plane-manager cluster-api-provider-docker-manager
SOURCE_PATTERNS=sigs.k8s.io/cluster-api/cmd/clusterctl sigs.k8s.io/cluster-api sigs.k8s.io/cluster-api/bootstrap/kubeadm sigs.k8s.io/cluster-api/controlplane/kubeadm sigs.k8s.io/cluster-api/test/infrastructure/docker
EXTRA_GO_LDFLAGS=$(shell $(BUILD_LIB)/version.sh $(REPO) sigs.k8s.io/cluster-api/cmd/version)

BASE_IMAGE_NAME=eks-distro-minimal-base-nonroot
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CLUSTER_API_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/cluster-api-controller
KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller
KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-control-plane-controller
CLUSTER_API_DOCKER_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/capd-manager

IMAGE_NAMES=cluster-api-controller kubeadm-bootstrap-controller kubeadm-control-plane-controller cluster-api-docker-controller

GO_MOD_DOWNLOAD_TARGETS=$(REPO)/eks-anywhere-go-mod-download $(REPO)/test/infrastructure/docker/eks-anywhere-go-mod-download
GATHER_LICENSES_TARGETS=$(OUTPUT_DIR)/attribution/go-license.csv $(OUTPUT_DIR)/capd/attribution/go-license.csv
ATTRIBUTION_TARGETS=ATTRIBUTION.txt CAPD_ATTRIBUTION.txt

FETCH_BINARIES_TARGETS=eksd/kubernetes/client

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk

$(OUTPUT_BIN_DIR)/linux-%/cluster-api-provider-docker-manager: EXTRA_GO_LDFLAGS=

$(OUTPUT_BIN_DIR)/linux-%/cluster-api-provider-docker-manager \
	$(REPO)/test/infrastructure/docker/eks-anywhere-go-mod-download \
	$(OUTPUT_DIR)/capd/attribution/go-license.csv: REPO_SUBPATH=test/infrastructure/docker

$(OUTPUT_DIR)/capd/attribution/go-license.csv CAPD_ATTRIBUTION.txt: LICENSE_PACKAGE_FILTER=.
$(OUTPUT_DIR)/capd/attribution/go-license.csv CAPD_ATTRIBUTION.txt: LICENSES_OUTPUT_DIR=$(OUTPUT_DIR)/capd


s3-artifacts: create-manifests

cluster-api-docker-controller/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-docker-client
$(call IMAGE_TARGETS_FOR_NAME, cluster-api-docker-controller): $(call FULL_FETCH_BINARIES_TARGETS, $(FETCH_BINARIES_TARGETS))

.PHONY: create-manifests
create-manifests: tarballs
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

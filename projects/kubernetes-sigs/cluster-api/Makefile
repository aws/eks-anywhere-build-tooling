BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=cluster-api
REPO_OWNER=kubernetes-sigs

LICENSE_PACKAGE_FILTER?=./cmd/clusterctl ./bootstrap/kubeadm ./controlplane/kubeadm ./

BINARY_TARGET_FILES=clusterctl manager kubeadm-bootstrap-manager kubeadm-control-plane-manager cluster-api-provider-docker-manager
SOURCE_PATTERNS=sigs.k8s.io/cluster-api/cmd/clusterctl sigs.k8s.io/cluster-api sigs.k8s.io/cluster-api/bootstrap/kubeadm sigs.k8s.io/cluster-api/controlplane/kubeadm sigs.k8s.io/cluster-api/test/infrastructure/docker
GO_MOD_PATHS=. . . . test/infrastructure/docker
EXTRA_GO_LDFLAGS=$(shell $(BUILD_LIB)/version.sh $(REPO) sigs.k8s.io/cluster-api/version)
CLUSTER_API_PROVIDER_DOCKER_MANAGER_ATTRIBUTION_OVERRIDE=capd

BASE_IMAGE_NAME=eks-distro-minimal-base-nonroot
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CLUSTER_API_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/cluster-api-controller
KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller
KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-control-plane-controller
CLUSTER_API_DOCKER_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/capd-manager

IMAGE_NAMES=cluster-api-controller kubeadm-bootstrap-controller kubeadm-control-plane-controller cluster-api-docker-controller

FETCH_BINARIES_TARGETS=eksd/kubernetes/client

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


$(OUTPUT_BIN_DIR)/linux-%/cluster-api-provider-docker-manager: EXTRA_GO_LDFLAGS=

FIX_LICENSES_GO_JSON_TARGET=$(REPO)/vendor/github.com/github.com/ajeddeloh/go-json/LICENSE.txt

s3-artifacts: create-manifests

$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_GO_JSON_TARGET)

cluster-api-docker-controller/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base
$(call IMAGE_TARGETS_FOR_NAME, cluster-api-docker-controller): $(call FULL_FETCH_BINARIES_TARGETS, $(FETCH_BINARIES_TARGETS))

.PHONY: create-manifests
create-manifests: tarballs
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

$(FIX_LICENSES_GO_JSON_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The ajeddeloh/go-json dependency github repo does not have a license, however it redirects to coreos/go-json repo
# that does have its own license. Hence we need to manually download license from coreos/go-json repo
# and place it in the ajeddeloh/go-json folder under vendor directory so that it is available for
# go-licenses to pick up
	wget -q https://raw.githubusercontent.com/coreos/go-json/main/LICENSE -O \
			$(REPO)/vendor/github.com/ajeddeloh/go-json/LICENSE.txt;


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

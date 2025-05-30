BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=cluster-api
REPO_OWNER=kubernetes-sigs

BINARY_TARGET_FILES=clusterctl manager kubeadm-bootstrap-manager kubeadm-control-plane-manager cluster-api-provider-docker-manager
SOURCE_PATTERNS=sigs.k8s.io/cluster-api/cmd/clusterctl sigs.k8s.io/cluster-api sigs.k8s.io/cluster-api/bootstrap/kubeadm sigs.k8s.io/cluster-api/controlplane/kubeadm sigs.k8s.io/cluster-api/test/infrastructure/docker
GO_MOD_PATHS=. . . . test
EXTRA_GO_LDFLAGS=$(shell $(BUILD_LIB)/version.sh $(REPO) sigs.k8s.io/cluster-api/version)
CLUSTER_API_PROVIDER_DOCKER_MANAGER_ATTRIBUTION_OVERRIDE=capd

BASE_IMAGE_NAME=eks-distro-minimal-base-nonroot
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CLUSTER_API_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/cluster-api-controller
KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller
KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/kubeadm-control-plane-controller
CLUSTER_API_DOCKER_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api/capd-manager

IMAGE_NAMES=cluster-api-controller kubeadm-bootstrap-controller kubeadm-control-plane-controller cluster-api-docker-controller

PROJECT_DEPENDENCIES=eksd/kubernetes/client

HAS_S3_ARTIFACTS=true

BUILDSPEC_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

FIX_LICENSES_GO_JSON_TARGET=$(REPO)/vendor/github.com/github.com/ajeddeloh/go-json/LICENSE.txt
FIX_LICENSES_TEST_CONTAINER_TARGET=$(REPO)/test/infrastructure/docker/LICENSE
FIX_LICENSES_TEST_DOCKER_TARGET=$(REPO)/test/infrastructure/container/LICENSE
FIX_LICENSES_TEST_KIND_TARGET=$(REPO)/test/infrastructure/kind/LICENSE
FIX_LICENSES_TEST_INMEMORY_TARGET=$(REPO)/test/infrastructure/inmemory/LICENSE

include $(BASE_DIRECTORY)/Common.mk


$(OUTPUT_BIN_DIR)/linux-%/cluster-api-provider-docker-manager: EXTRA_GO_LDFLAGS=

s3-artifacts: create-manifests

$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_GO_JSON_TARGET) $(FIX_LICENSES_TEST_CONTAINER_TARGET) $(FIX_LICENSES_TEST_DOCKER_TARGET) $(FIX_LICENSES_TEST_KIND_TARGET) $(FIX_LICENSES_TEST_INMEMORY_TARGET)

cluster-api-docker-controller/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base

.PHONY: create-manifests
create-manifests: MAKEFLAGS=
create-manifests: tarballs | ensure-yq $$(ENABLE_DOCKER)
	@build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG) $(GOLANG_VERSION)

$(FIX_LICENSES_GO_JSON_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The ajeddeloh/go-json dependency github repo does not have a license, however it redirects to coreos/go-json repo
# that does have its own license. Hence we need to manually download license from coreos/go-json repo
# and place it in the ajeddeloh/go-json folder under vendor directory so that it is available for
# go-licenses to pick up
	wget --retry-connrefused -q https://raw.githubusercontent.com/coreos/go-json/main/LICENSE -O \
			$(REPO)/vendor/github.com/ajeddeloh/go-json/LICENSE.txt;

$(FIX_LICENSES_TEST_CONTAINER_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# CAPD pulls in a dep in test/infrastructure/container which does not have a LICENSE file and go-licenses does look it up
# We use capi license instead
	cp $(REPO)/LICENSE $@

$(FIX_LICENSES_TEST_DOCKER_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# CAPD is a separate module but it doesn't have its own license, it inherits the one at the REPO top level.
	cp $(REPO)/LICENSE $@

$(FIX_LICENSES_TEST_KIND_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# CAPD pulls in a dep in test/infrstructure/kind which does not have a LICENSE file and go-licenses does look it up
# We use capi license instead
	cp $(REPO)/LICENSE $@

$(FIX_LICENSES_TEST_INMEMORY_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# CAPD pulls in a dep in test/infrstructure/inmemory which does not have a LICENSE file and go-licenses does look it up
# We use capi license instead
	cp $(REPO)/LICENSE $@	


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

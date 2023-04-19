BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.19"

REPO=kind
REPO_OWNER=kubernetes-sigs

BINARY_TARGET_FILES=kind kindnetd
SOURCE_PATTERNS=. ./cmd/kindnetd
GO_MOD_PATHS=. ./images/kindnetd

# force binariess to go to non-release branch bin folder
BINARIES_ARE_RELEASE_BRANCHED=false

# do not look for checksums in release branch folder, instead use project root
PROJECT_ROOT=$(MAKE_ROOT)

MINIMAL_BASE_IPTABLES_IMAGE_NAME?=eks-distro-minimal-base-iptables
MINIMAL_BASE_HAPROXY_IMAGE_NAME?=eks-distro-minimal-base-haproxy
MINIMAL_BASE_KIND_IMAGE_NAME?=eks-distro-minimal-base-kind
MINIMAL_BASE_KIND_IMAGE_NAME_OS_VERSION=al2023

KIND_BASE_IMAGE_COMPONENT=kubernetes-sigs/kind/base
KIND_BASE_VERSIONED_IMAGE_COMPONENT=kubernetes-sigs/kind/base
KIND_NODE_IMAGE_COMPONENT=kubernetes-sigs/kind/node
KINDNETD_IMAGE_COMPONENT=kubernetes-sigs/kind/kindnetd
HAPROXY_IMAGE_COMPONENT=kubernetes-sigs/kind/haproxy

KIND_BASE_IMAGE_BUILD_ARGS=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/kind-base-image-build-args
KIND_BASE_KUBEADM_OVERRIDE=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/images/base/files/usr/local/bin/kubeadm
KIND_NODE_IMAGE_BUILD_ARGS=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/kind-node-image-build-args
EKSD_KUBE_VERSION=

# Overriding since kind is released branch it will try to add git_tag to latest tag
LATEST_TAG=$(LATEST)

UNVERSIONED_BASE_IMAGE=$(IMAGE_REPO)/$(KIND_BASE_IMAGE_COMPONENT):$(GIT_TAG)-$(GIT_HASH)
VERSIONED_BASE_IMAGE=$(IMAGE_REPO)/$(KIND_BASE_IMAGE_COMPONENT):$(EKSD_KUBE_VERSION)-$(GIT_TAG)-$(GIT_HASH)

HAS_S3_ARTIFACTS?=true
HAS_RELEASE_BRANCHES=true

PROJECT_DEPENDENCIES=eksd/kubernetes/client eksa/kubernetes-sigs/etcdadm eksd/cni-plugins eksa/kubernetes-sigs/cri-tools eksd/etcd/etcd.tar.gz eksa/rancher/local-path-provisioner
ORGANIZE_BINARIES_AMD64_TARGET=$(BINARY_DEPS_DIR)/linux-amd64/files/rootfs/LICENSES
ORGANIZE_BINARIES_ARM64_TARGET=$(BINARY_DEPS_DIR)/linux-arm64/files/rootfs/LICENSES
ARM_ENV_CONF_TARGET=$(BINARY_DEPS_DIR)/linux-arm64/files/rootfs/etc
KIND_NODE_BUILD_AMD64_TARGET=$(BINARY_DEPS_DIR)/linux-amd64/files/io.containerd.metadata.v1.bolt/meta.db
KIND_NODE_BUILD_ARM64_TARGET=$(BINARY_DEPS_DIR)/linux-arm64/files/io.containerd.metadata.v1.bolt/meta.db

HAPROXY_CFG=$(OUTPUT_DIR)/files/haproxy/usr/local/etc/haproxy/haproxy.cfg

IMAGE_NAMES?=haproxy kindnetd kind-base 

FIX_LICENSES_KINDNETD_TARGET=$(REPO)/images/kindnetd/LICENSE

BUILDSPECS=buildspecs/binaries.yml buildspecs/images.yml
BUILDSPEC_1_PLATFORM=ARM_CONTAINER
BUILDSPEC_1_DEPENDS_ON_OVERRIDE=none
BUILDSPEC_2_PLATFORM=LINUX_CONTAINER
BUILDSPEC_2_COMPUTE_TYPE=BUILD_GENERAL1_LARGE
BUILDSPEC_2_VARS_KEYS=RELEASE_BRANCH
BUILDSPEC_2_VARS_VALUES=SUPPORTED_K8S_VERSIONS



include $(BASE_DIRECTORY)/Common.mk

# Do not need to import if running clean or any var-value-% targets
ifneq ($(and $(RELEASE_BRANCH),$(filter-out clean,$(MAKECMDGOALS_WITHOUT_VAR_VALUE))),)
-include $(KIND_BASE_IMAGE_BUILD_ARGS)
-include $(KIND_NODE_IMAGE_BUILD_ARGS)
endif

$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_KINDNETD_TARGET)

# During build the node image is never built so we can extract this file
# instead fake it to keep s3-artifacts happy
s3-artifacts: $(ARTIFACTS_PATH)/manifests/kindnetd/$(GIT_TAG)/kindnetd.yaml

images: kind-base-versioned/images/push kind-node/images/push
create-ecr-repos: kind-node/create-ecr-repo

$(OUTPUT_BIN_DIR)/%/kind: EXTRA_GO_LDFLAGS=-X=sigs.k8s.io/kind/pkg/cmd/kind/version.gitCommit=$(shell git -C $(REPO) rev-list -n 1  "${GIT_TAG}")

# haproxy
haproxy/images/%: DOCKERFILE_FOLDER=./images/haproxy
haproxy/images/%: BASE_IMAGE_NAME=$(MINIMAL_BASE_HAPROXY_IMAGE_NAME)
$(call IMAGE_TARGETS_FOR_NAME, haproxy): $(HAPROXY_CFG)

$(HAPROXY_CFG): $(GIT_PATCH_TARGET)
	@mkdir -p $(@D)
	cp $(REPO)/images/haproxy/haproxy.cfg $@

# Kindnetd
kindnetd/images/%: DOCKERFILE_FOLDER=./images/kindnetd
kindnetd/images/%: BASE_IMAGE_NAME=$(MINIMAL_BASE_IPTABLES_IMAGE_NAME)

# Kind-base
kind-base/images/% kind-base-versioned/images/%: DOCKERFILE_FOLDER=$(REPO)/images/base
kind-base/images/% kind-base-versioned/images/%: BASE_IMAGE_NAME=$(MINIMAL_BASE_KIND_IMAGE_NAME)
kind-base/images/% kind-base-versioned/images/%: BASE_IMAGE_OS_VERSION=$(MINIMAL_BASE_KIND_IMAGE_NAME_OS_VERSION)
kind-base/images/% kind-base-versioned/images/%: IMAGE_BUILD_ARGS=CNI_PLUGINS_AMD64_URL CNI_PLUGINS_ARM64_URL CNI_PLUGINS_AMD64_SHA256SUM \
	CNI_PLUGINS_ARM64_SHA256SUM CRICTL_AMD64_SHA256SUM_URL CRICTL_ARM64_SHA256SUM_URL CRICTL_AMD64_URL CRICTL_ARM64_URL

kind-base/images/%: IMAGE_TARGET=base
kind-base/images/%: IMAGE_CONTEXT_DIR=$(REPO)/images/base
$(call IMAGE_TARGETS_FOR_NAME, kind-base): $(KIND_BASE_IMAGE_BUILD_ARGS)

kind-base-versioned/images/%: IMAGE_TARGET=base-versioned
kind-base-versioned/images/%: IMAGE_TAG:=$(EKSD_KUBE_VERSION)-$(IMAGE_TAG)
kind-base-versioned/images/%: LATEST_TAG:=$(EKSD_KUBE_VERSION)-$(LATEST_TAG)
kind-base-versioned/images/%: IMAGE_CONTEXT_DIR=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/images/base
kind-base-versioned/images/%: BUILDER_IMAGE=$(UNVERSIONED_BASE_IMAGE)
$(call IMAGE_TARGETS_FOR_NAME, kind-base-versioned): $(KIND_BASE_KUBEADM_OVERRIDE)

# Kind-node
kind-node/images/push: DOCKERFILE_FOLDER=./images/node
kind-node/images/push: IMAGE_CONTEXT_DIR=$(BINARY_DEPS_DIR)
kind-node/images/push: BASE_IMAGE=$(VERSIONED_BASE_IMAGE)
# We are fully overriding the base_image above, setting the name is to ensure the builder image is set correctly
kind-node/images/push: BASE_IMAGE_NAME=$(MINIMAL_BASE_KIND_IMAGE_NAME)
kind-node/images/push: BASE_IMAGE_OS_VERSION=$(MINIMAL_BASE_KIND_IMAGE_NAME_OS_VERSION)
kind-node/images/push: IMAGE_TAG=$(NODE_IMAGE_TAG)
kind-node/images/push: LATEST_TAG=$(NODE_IMAGE_LATEST_TAG)
kind-node/images/push: IMAGE_TARGET=node
kind-node/images/push: $(KIND_NODE_BUILD_AMD64_TARGET) $(KIND_NODE_BUILD_ARM64_TARGET) $(ARM_ENV_CONF_TARGET)

$(KIND_BASE_IMAGE_BUILD_ARGS):
	@echo -e $(call TARGET_START_LOG)
	build/base-image-build-args.sh $(RELEASE_BRANCH) $(ARTIFACTS_BUCKET) $@ $(LATEST)
	@echo -e $(call TARGET_END_LOG)

$(KIND_NODE_IMAGE_BUILD_ARGS):
	@echo -e $(call TARGET_START_LOG)
	build/node-image-build-args.sh $(RELEASE_BRANCH) $(KINDNETD_IMAGE_COMPONENT) $(IMAGE_REPO) "$(ARTIFACTS_BUCKET)" $(IMAGE_TAG) $(LATEST) $@
	@echo -e $(call TARGET_END_LOG)

# Tweak the kind/base image to have a hardcode kubeadm config
# so that during the image pull phase it pulls eks-d images
# vs upstream images
# kubeadm-override and config are copied into kind/images/base/files/etc
# so they are automatically added into the image by the dockerfile
$(KIND_BASE_KUBEADM_OVERRIDE): KUBEADM_CONFIG_DIR=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/images/base/files/etc
$(KIND_BASE_KUBEADM_OVERRIDE):
	@mkdir -p $(@D)
	@mkdir -p $(KUBEADM_CONFIG_DIR)
	set -a && source $(KIND_BASE_IMAGE_BUILD_ARGS) && envsubst '$$COREDNS_VERSION:$$ETCD_VERSION:$$EKSD_KUBE_VERSION:$$EKSD_IMAGE_REPO' \
		< ./images/base/kubeadm.config.tmpl \
		> $(KUBEADM_CONFIG_DIR)/kubeadm.config
	cp ./images/base/kubeadm-override.sh $@

# kind build node-image process, also extracts kindnetd manifest to push to s3
# kind cli built for current host is needed to create the node image
$(KIND_NODE_BUILD_AMD64_TARGET): KIND_CLI=./$(OUTPUT_BIN_DIR)/$(subst /,-,$(BUILDER_PLATFORM))/kind
$(KIND_NODE_BUILD_AMD64_TARGET): $(KIND_NODE_IMAGE_BUILD_ARGS) $(ORGANIZE_BINARIES_AMD64_TARGET)
	$(MAKE) $(KIND_CLI) BINARY_PLATFORMS=$(BUILDER_PLATFORM)
	build/build-kind-node-image.sh $(RELEASE_BRANCH) $(VERSIONED_BASE_IMAGE) amd64
	@mkdir -p $(ARTIFACTS_PATH)/manifests/kindnetd/$(GIT_TAG)
	cp $(BINARY_DEPS_DIR)/linux-amd64/files/rootfs/kind/manifests/default-cni.yaml $(ARTIFACTS_PATH)/manifests/kindnetd/$(GIT_TAG)/kindnetd.yaml
	sed -i -e 's/{{ .PodSubnet }}/192.168.0.0\/16/' $(ARTIFACTS_PATH)/manifests/kindnetd/$(GIT_TAG)/kindnetd.yaml 

$(KIND_NODE_BUILD_ARM64_TARGET): $(KIND_NODE_IMAGE_BUILD_ARGS) $(ORGANIZE_BINARIES_ARM64_TARGET)
	build/build-kind-node-image.sh $(RELEASE_BRANCH) $(VERSIONED_BASE_IMAGE) arm64

$(ORGANIZE_BINARIES_AMD64_TARGET): $(HANDLE_DEPENDENCIES_TARGET)
	build/organize_binaries.sh $(RELEASE_BRANCH) amd64

$(ORGANIZE_BINARIES_ARM64_TARGET): $(HANDLE_DEPENDENCIES_TARGET)
	build/organize_binaries.sh $(RELEASE_BRANCH) arm64

$(ARM_ENV_CONF_TARGET):
	@cp -rf images/node/files/* $(BINARY_DEPS_DIR)/linux-arm64/files/rootfs/

$(ARTIFACTS_PATH)/manifests/kindnetd/$(GIT_TAG)/kindnetd.yaml:
	@mkdir -p $(@D)
	@touch $@

$(FIX_LICENSES_KINDNETD_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
#go-licenses requires a LICENSE file in each folder with the go.mod
	cp $(REPO)/LICENSE $@

.PHONY: create-kind-cluster-%
create-kind-cluster-%:
	build/create-kind-cluster.sh $(IMAGE_REPO)/$(KIND_NODE_IMAGE_COMPONENT):$(NODE_IMAGE_LATEST_TAG) $* $(RELEASE_BRANCH)

create-kind-cluster-amd64: # Create local kind cluster using built amd64 image
create-kind-cluster-arm64: # Create local kind cluster using built arm64 image


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

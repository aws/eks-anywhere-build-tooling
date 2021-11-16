BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=kind
REPO_OWNER=kubernetes-sigs

LICENSE_PACKAGE_FILTER?=./cmd/kind

BINARY_PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
BINARY_TARGETS=$(call BINARY_TARGETS_FROM_FILES_PLATFORMS,kind,$(BINARY_PLATFORMS))

BINARY_PLATFORMS_KINDNETD=linux/amd64 linux/arm64
BINARY_TARGETS+=$(call BINARY_TARGETS_FROM_FILES_PLATFORMS,kindnetd,$(BINARY_PLATFORMS_KINDNETD))

# force binarys to go to non-release branch bin folder
OUTPUT_DIR=_output

# do not look for checksums in release branch folder, instead use project root
PROJECT_ROOT=$(MAKE_ROOT)

BASE_IMAGE_NAME?=eks-distro-base
MINIMAL_IPTABLES_BASE_IMAGE_NAME?=eks-distro-minimal-base-iptables

KINDNETD_GATHER_LICENSES_TARGET=$(OUTPUT_DIR)/kindnetd/attribution/go-license.csv
KINDNETD_ATTRIBUTION_TARGET=KINDNETD_ATTRIBUTION.txt

KIND_BASE_IMAGE_COMPONENT=kindest/base
KIND_BASE_VERSIONED_IMAGE_COMPONENT=kindest/base
KIND_NODE_IMAGE_COMPONENT=kubernetes-sigs/kind/node
KINDNETD_IMAGE_COMPONENT=kubernetes-sigs/kind/kindnetd

KIND_BASE_IMAGE_BUILD_ARGS=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/kind-base-image-build-args
KIND_BASE_KUBEADM_OVERRIDE=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/images/base/files/usr/local/bin/kubeadm
KIND_NODE_IMAGE_BUILD_ARGS=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/kind-node-image-build-args

LOCAL_REGISTRY=localhost:5000
INTERMEDIATE_BASE_IMAGE=$(LOCAL_REGISTRY)/$(KIND_BASE_IMAGE_COMPONENT):$(EKSD_KUBE_VERSION)
INTERMEDIATE_NODE_IMAGE=$(LOCAL_REGISTRY)/$(KIND_NODE_IMAGE_COMPONENT):$(EKSD_KUBE_VERSION)

HAS_S3_ARTIFACTS=true
HAS_RELEASE_BRANCHES=true

# Binary targets are added manually in this makefile instead of relying on simple_create logic
SIMPLE_CREATE_BINARIES=false

ifneq ($(and $(RELEASE_BRANCH),$(filter-out clean,$(MAKECMDGOALS))),)
-include $(KIND_BASE_IMAGE_BUILD_ARGS)
-include $(KIND_NODE_IMAGE_BUILD_ARGS)
endif

include $(BASE_DIRECTORY)/Common.mk


$(OUTPUT_BIN_DIR)/%/kindnetd: REPO_SUBPATH=images/kindnetd

$(OUTPUT_BIN_DIR)/%/kind: EXTRA_GO_LDFLAGS=-X=sigs.k8s.io/kind/pkg/cmd/kind/version.GitCommit=$(shell git -C $(REPO) rev-list -n 1  "${GIT_TAG}")

# Since there are different platforms for the different binaries, calling helper from common.mk directly
# instead of relying on it automatically
$(eval $(foreach platform, $(BINARY_PLATFORMS), $(call BINARY_TARGET_BODY,$(platform),kind,.)))
$(eval $(foreach platform, $(BINARY_PLATFORMS_KINDNETD), $(call BINARY_TARGET_BODY,$(platform),kindnetd,./cmd/kindnetd)))

# Kindnetd
kindnetd/images/%: DOCKERFILE_FOLDER=./images/kindnetd
kindnetd/images/%: BASE_IMAGE_NAME=$(MINIMAL_IPTABLES_BASE_IMAGE_NAME)
$(call IMAGE_TARGETS_FOR_NAME, kindnetd): $(KINDNETD_GATHER_LICENSES_TARGET)

# Kind-base
kind-base/images/% kind-base-versioned/images/%: DOCKERFILE_FOLDER=$(REPO)/images/base
kind-base/images/% kind-base-versioned/images/%: BASE_IMAGE_NAME=eks-distro-base
kind-base/images/% kind-base-versioned/images/%: IMAGE_BUILD_ARGS=CNI_PLUGINS_AMD64_URL CNI_PLUGINS_ARM64_URL CNI_PLUGINS_AMD64_SHA256SUM \
	CNI_PLUGINS_ARM64_SHA256SUM CRICTL_AMD64_SHA256SUM CRICTL_ARM64_SHA256SUM CRICTL_AMD64_URL CRICTL_ARM64_URL

kind-base/images/%: IMAGE_TARGET=base
kind-base/images/%: IMAGE_TAG=latest
kind-base/images/%: IMAGE_CONTEXT_DIR=$(REPO)/images/base
$(call IMAGE_TARGETS_FOR_NAME, kind-base): $(KIND_BASE_IMAGE_BUILD_ARGS)

kind-base-versioned/images/%: IMAGE_TARGET=base-versioned
kind-base-versioned/images/%: IMAGE_TAG=$(EKSD_KUBE_VERSION)
kind-base-versioned/images/%: IMAGE_CONTEXT_DIR=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/images/base
kind-base-versioned/images/%: BUILDER_IMAGE=$(LOCAL_REGISTRY)/$(KIND_BASE_IMAGE_COMPONENT):latest
$(call IMAGE_TARGETS_FOR_NAME, kind-base-versioned): $(KIND_BASE_KUBEADM_OVERRIDE)

# Avoid pushing to remote repo, this image is used by the `kind build node-image`
kind-base/images/push kind-base-versioned/images/push: IMAGE_REPO=$(LOCAL_REGISTRY)

# Kind-node
kind-node/images/%: DOCKERFILE_FOLDER=./images/node
kind-node/images/%: IMAGE_CONTEXT_DIR=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/dependencies
kind-node/images/%: BUILDER_IMAGE=$(INTERMEDIATE_NODE_IMAGE)
kind-node/images/%: BASE_IMAGE=$(LOCAL_REGISTRY)/$(KIND_BASE_IMAGE_COMPONENT):latest
kind-node/images/%: IMAGE_TAG=$(EKSD_KUBE_VERSION)
kind-node/images/%: IMAGE_BUILD_ARGS=IMAGE_REPO AL2_HELPER_IMAGE DEBIAN_BASE_IMAGE_TAG LOCAL_PATH_PROVISONER_IMAGE_TAG LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE \
	PAUSE_IMAGE_TAG_OVERRIDE PAUSE_IMAGE_TAG KIND_KINDNETD_IMAGE_OVERRIDE KINDNETD_IMAGE_TAG
$(call IMAGE_TARGETS_FOR_NAME, kind-node): $(KIND_NODE_IMAGE_BUILD_ARGS)

generate-attribution: $(KINDNETD_ATTRIBUTION_TARGET)

$(KINDNETD_GATHER_LICENSES_TARGET): $(BINARY_TARGETS)
	$(BASE_DIRECTORY)/build/lib/gather_licenses.sh $(REPO) $(MAKE_ROOT)/$(OUTPUT_DIR)/kindnetd ./cmd/kindnetd images/kindnetd 

$(KINDNETD_ATTRIBUTION_TARGET): $(KINDNETD_GATHER_LICENSES_TARGET)
	$(BASE_DIRECTORY)/build/lib/create_attribution.sh $(MAKE_ROOT) $(GOLANG_VERSION) $(MAKE_ROOT)/$(OUTPUT_DIR)/kindnetd $@

.PHONY: local-images
local-images: kindnetd/images/amd64 kind-base/images/amd64

.PHONY: images
images: kindnetd/images/push kind-base/images/push kind-base-versioned/images/push build-kind-node-amd64 build-kind-node-arm64 kind-node/images/push

$(KIND_BASE_IMAGE_BUILD_ARGS):
	build/base-image-build-args.sh $(RELEASE_BRANCH) $(ARTIFACTS_BUCKET) $@

$(KIND_NODE_IMAGE_BUILD_ARGS):
	build/node-image-build-args.sh $(RELEASE_BRANCH) $(KINDNETD_IMAGE_COMPONENT) $(IMAGE_REPO) $(ARTIFACTS_BUCKET) $(LATEST_TAG) $@

# Tweak the kind/base image to have a hardcode kubeadm config
# so that during the image pull phase it pulls eks-d images
# vs upstream images
# kubeadm-override and config are copied into kind/images/base/files/etc
# so they are automatically added into the image by the dockerfile
$(KIND_BASE_KUBEADM_OVERRIDE): KUBEADM_CONFIG_DIR=$(OUTPUT_DIR)/$(RELEASE_BRANCH)/images/base/files/etc
$(KIND_BASE_KUBEADM_OVERRIDE):
	@mkdir -p $(@D)
	@mkdir -p $(KUBEADM_CONFIG_DIR)
	set -a && source $(KIND_BASE_IMAGE_BUILD_ARGS) && envsubst '$$COREDNS_VERSION:$$ETCD_VERSION:$$EKSD_KUBE_VERSION' \
		< ./images/base/kubeadm.config.tmpl \
		> $(KUBEADM_CONFIG_DIR)/kubeadm.config
	cp ./images/base/kubeadm-override.sh $@

.PHONY: build-kind-node-%
build-kind-node-%: $(KIND_NODE_IMAGE_BUILD_ARGS)
	build/build-kind-node-image.sh $(RELEASE_BRANCH) $(INTERMEDIATE_BASE_IMAGE) \
		$(INTERMEDIATE_NODE_IMAGE)-$* $(ARTIFACTS_BUCKET) $*

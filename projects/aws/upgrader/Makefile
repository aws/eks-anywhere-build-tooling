BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)

REPO=upgrader
REPO_OWNER=aws

GIT_TAG=v$(RELEASE_BRANCH)-$(shell yq e ".releases[] | select(.branch==\"${RELEASE_BRANCH}\").number" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)

HAS_RELEASE_BRANCHES=true

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true
# for the staging buildspec generation
BUILDSPEC_DEPENDS_ON_OVERRIDE=containerd_containerd_linux_amd64 containerd_containerd_linux_arm64 kubernetes_sigs_cri_tools

BINARY_TARGET_FILES=upgrader
SOURCE_PATTERNS=. 
GO_MOD_PATHS=..

GIT_CHECKOUT_TARGET=main.go
FAKE_GIT_REPO_TARGET=.git
REPO_NO_CLONE=true
HAS_LICENSES=false
SKIP_CHECKSUM_VALIDATION=true

BUILDSPECS=buildspec.yml buildspecs/combine-images.yml
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE
BUILDSPEC_1_VARS_KEYS=IMAGE_PLATFORMS
BUILDSPEC_1_VARS_VALUES=IMAGE_PLATFORMS
BUILDSPEC_1_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER
BUILDSPEC_2_DEPENDS_ON_OVERRIDE=aws_upgrader_linux_amd64 aws_upgrader_linux_arm64


BASE_IMAGE_NAME=eks-distro-minimal-base-nsenter
IMAGE_NAMES=upgrader
# its not clear why, but if we do not set this here, even though this is the default
# when we try and override it in the combine-images target the override will not take
DOCKERFILE_FOLDER=./docker/linux
IMAGE_BUILD_ARGS=

PROJECT_DEPENDENCIES=eksd/kubernetes/client eksd/kubernetes/server eksd/cni-plugins eksa/containerd/containerd eksa/kubernetes-sigs/cri-tools


include $(BASE_DIRECTORY)/Common.mk

TOOLS_BIN_DIR := $(shell pwd)/hack/tools/bin
MOCKGEN := $(TOOLS_BIN_DIR)/mockgen

build: unit-test

$(REPO):
	@mkdir $@

$(GATHER_LICENSES_TARGETS): | $(FAKE_GIT_REPO_TARGET)

$(ATTRIBUTION_TARGETS): GIT_TAG

$(FAKE_GIT_REPO_TARGET):
	@git init
	@git remote add origin https://github.com/aws/eks-anywhere-build-tooling.git 

$(TOOLS_BIN_DIR):
	@mkdir -p $(TOOLS_BIN_DIR)

$(MOCKGEN): $(TOOLS_BIN_DIR)
	GOBIN=$(TOOLS_BIN_DIR) build::common::use_go_version $(GOLANG_VERSION) install github.com/golang/mock/mockgen@v1.6.0

unit-test: | $$(ENABLE_DOCKER)
	build::common::use_go_version $(GOLANG_VERSION) test ./...

.PHONY: mocks
mocks: $(MOCKGEN)
	build::common::use_go_version $(GOLANG_VERSION) generate ./...

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

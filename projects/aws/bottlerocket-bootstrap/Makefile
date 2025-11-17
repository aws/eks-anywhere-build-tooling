# avoid finding the local fake git repo by looking one level up when looking for base_directory
BASE_DIRECTORY:=$(shell git -C .. rev-parse --show-toplevel)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=bottlerocket-bootstrap
REPO_OWNER=aws

BASE_IMAGE_NAME?=eks-distro-base

# we need to set IMAGE_BUILD_ARGS here even though its the same as the default. 
# it is set in Common.mk on the images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L799)
# and the combine-images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L846)
# since combine-images has images as prereq target, the ?= does not really behavior as one might expect.
# the images target being the actual action, its version of the set takes prioirty and resets to empty
# setting it explicitly to empty here takes allows the combine-images override to take proirty
IMAGE_BUILD_ARGS=

IMAGE_NAMES=bottlerocket-bootstrap bottlerocket-bootstrap-snow bottlerocket-bootstrap-multi-network
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

BOTTLEROCKET_BOOTSTRAP_IMAGE_COMPONENT=bottlerocket-bootstrap
BOTTLEROCKET_BOOTSTRAP_SNOW_IMAGE_COMPONENT=bottlerocket-bootstrap-snow
BOTTLEROCKET_BOOTSTRAP_MULTI_NETWORK_IMAGE_COMPONENT=bottlerocket-bootstrap-multi-network

BINARY_TARGET_FILES=bottlerocket-bootstrap bottlerocket-bootstrap-snow
SOURCE_PATTERNS=. ./pkg/providers/snow
GO_MOD_PATHS=.. ..

GIT_TAG=v$(RELEASE_BRANCH)-$(shell yq e ".releases[] | select(.branch==\"${RELEASE_BRANCH}\").number" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)

HAS_RELEASE_BRANCHES=true

BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

# Go 1.18 and newer automatically embeds vcs information if the go module being built
# belongs (upwards) to a git repo. In this case, that repo is the buildtooling repo, 
# which is problematic. Given these binaries are built from different banches,
# the local dev branhc on PRs and `main` later, the VCS info will change between
# builds, making builds quite difficult to reproduce. That would require that the presubmits 
# knew the future commit that might be created if the PR gets merged.
# Ideally this go module would live on its own repo, so all builds can be undoubtely
# linked to a particular commit of that repo and become independent from the build tooling repo.
# As an alternative, we don't include VCS info in the binary.
EXTRA_GOBUILD_FLAGS=-buildvcs=false

# This project is a bit of odd ball when it comes to using the standard makefile targets
# - the image builds are release branch aware but the binary is the same across all versions
# - there is no clone of source since the source for the bootstrap exist in this repo already
# - go-licenses supports only HTTPS URLs so a temporary git repo is created in this project root
#   to add the https origin
# - REPO target is overridden to create empty folder

# force binaries to go to non-release branch bin folder
BINARIES_ARE_RELEASE_BRANCHED=false

# do not look for checksums in release branch folder, instead use project root
PROJECT_ROOT=$(MAKE_ROOT)

# fake git checkout
REPO_NO_CLONE=true
GIT_CHECKOUT_TARGET=main.go
FAKE_GIT_REPO_TARGET=.git

PROJECT_DEPENDENCIES=eksd/kubernetes/client eksd/kubernetes/server eksa/kubernetes-sigs/etcdadm

TOOLS_BIN_DIR := $(shell pwd)/hack/tools/bin
MOCKGEN := $(TOOLS_BIN_DIR)/mockgen
GO_BIN_PATH=$(shell source $(BUILD_LIB)/common.sh && build::common::get_go_path $(GOLANG_VERSION))
GO=$(GO_BIN_PATH)/go

include $(BASE_DIRECTORY)/Common.mk


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
	GOBIN=$(TOOLS_BIN_DIR) $(GO) install github.com/golang/mock/mockgen@v1.6.0

GIT_TAG:
	echo 0.0.0 > $@

.PHONY: clean-extra
clean-extra:
	@rm -rf vendor .git GIT_TAG $(GO_MOD_DOWNLOAD_TARGETS) $(REPO)

clean: clean-extra

unit-test: | $$(ENABLE_DOCKER)
	$(GO) test ./...

.PHONY: mocks
mocks: $(MOCKGEN)
	$(GO) generate ./...


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

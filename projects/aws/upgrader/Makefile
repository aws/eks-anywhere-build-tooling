BASE_DIRECTORY:=$(shell git -C .. rev-parse --show-toplevel)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)

REPO=upgrader
REPO_OWNER=aws

GIT_TAG=v$(RELEASE_BRANCH)-$(shell yq e ".releases[] | select(.branch==\"${RELEASE_BRANCH}\").number" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)

HAS_RELEASE_BRANCHES=true

# Go 1.18 and newer automatically embeds vcs information if the go module being built
# belongs (upwards) to a git repo. In this case, that repo is the buildtooling repo, 
# which is problematic. Given these binaries are built from different banches,
# the local dev branch on PRs and `main` later, the VCS info will change between
# builds, making builds quite difficult to reproduce. That would require that the presubmits 
# knew the future commit that might be created if the PR gets merged.
# Ideally this go module would live on its own repo, so all builds can be undoubtely
# linked to a particular commit of that repo and become independent from the build tooling repo.
# As an alternative, we don't include VCS info in the binary.
EXTRA_GOBUILD_FLAGS=-buildvcs=false

BINARY_TARGET_FILES=upgrader
SOURCE_PATTERNS=. 
GO_MOD_PATHS=..

# force binaries to go to non-release branch bin folder
BINARIES_ARE_RELEASE_BRANCHED=false
GIT_CHECKOUT_TARGET=main.go
FAKE_GIT_REPO_TARGET=.git
REPO_NO_CLONE=true

BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

BASE_IMAGE_NAME=eks-distro-minimal-base-nsenter
IMAGE_NAMES=upgrader
# its not clear why, but if we do not set this here, even though this is the default
# when we try and override it in the combine-images target the override will not take
DOCKERFILE_FOLDER=./docker/linux
IMAGE_BUILD_ARGS=

PROJECT_DEPENDENCIES=eksd/kubernetes/client eksd/kubernetes/server eksd/cni-plugins eksa/containerd/containerd eksa/kubernetes-sigs/cri-tools
TOOLS_BIN_DIR=$(MAKE_ROOT)/hack/tools/bin
MOCKGEN=$(TOOLS_BIN_DIR)/mockgen
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

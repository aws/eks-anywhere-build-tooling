BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION="1.19"

REPO=containerd
REPO_OWNER=containerd

BINARY_TARGET_FILES=containerd containerd-shim containerd-shim-runc-v1 containerd-shim-runc-v2 ctr
SOURCE_PATTERNS=$(foreach target,$(BINARY_TARGET_FILES),./cmd/$(target))

GITVERSION=$(shell git -C $(REPO) describe --dirty --long --always)
COMMIT=$(shell git -C $(REPO) rev-parse HEAD)
EXTRA_GO_LDFLAGS=-X github.com/containerd/containerd/version.Version=$(GITVERSION) -X github.com/containerd/containerd/version.Revision=$(COMMIT) -X github.com/containerd/containerd/version.Package=github.com/containerd/containerd -extldflags -static
EXTRA_GOBUILD_FLAGS=-tags urfave_cli_no_docs,osusergo,netgo,static_build

CGO_CREATE_BINARIES=true

HAS_S3_ARTIFACTS=true

IMAGE_NAMES=

BUILDSPEC_VARS_KEYS=BINARY_PLATFORMS
BUILDSPEC_VARS_VALUES=BINARY_PLATFORMS
BUILDSPEC_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER


include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

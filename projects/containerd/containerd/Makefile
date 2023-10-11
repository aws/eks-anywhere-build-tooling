BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION:=$(shell cat GOLANG_VERSION)
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
BUILDSPEC_DEPENDS_ON_OVERRIDE=opencontainers_runc_linux_amd64 opencontainers_runc_linux_arm64

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true

PROJECT_DEPENDENCIES=eksa/opencontainers/runc
SIMPLE_CREATE_TARBALLS=false

include $(BASE_DIRECTORY)/Common.mk

.PHONY: tarballs
tarballs: $(LICENSES_TARGETS_FOR_PREREQ) $(HANDLE_DEPENDENCIES_TARGET)
	@echo -e $(call TARGET_START_LOG)
	build/create_tarballs.sh $(TAR_FILE_PREFIX) $(MAKE_ROOT)/$(OUTPUT_DIR) $(MAKE_ROOT)/$(OUTPUT_BIN_DIR) $(GIT_TAG) "$(BINARY_PLATFORMS)" $(ARTIFACTS_PATH) $(GIT_HASH) $(MAKE_ROOT)/$(BINARY_DEPS_DIR)
	@echo -e $(call TARGET_END_LOG)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

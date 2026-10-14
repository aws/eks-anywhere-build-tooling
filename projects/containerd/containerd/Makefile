BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat ./$(RELEASE_BRANCH)/GIT_TAG)
GOLANG_VERSION=$(shell cat ./$(RELEASE_BRANCH)/GOLANG_VERSION)
REPO=containerd
REPO_OWNER=containerd

BINARY_TARGET_FILES=$(if $(RELEASE_BRANCH),$(shell cat ./$(RELEASE_BRANCH)/BINARY_TARGET_FILES),)
SOURCE_PATTERNS=$(foreach target,$(BINARY_TARGET_FILES),./cmd/$(target))

GITVERSION=$(shell git -C $(REPO) describe --dirty --long --always)
COMMIT=$(shell git -C $(REPO) rev-parse HEAD)
CONTAINERD_VERSION_PACKAGE=github.com/containerd/containerd/version
CONTAINERD_MODULE=github.com/containerd/containerd
ifeq ($(RELEASE_BRANCH),1-37)
CONTAINERD_VERSION_PACKAGE=github.com/containerd/containerd/v2/version
CONTAINERD_MODULE=github.com/containerd/containerd/v2
endif
EXTRA_GO_LDFLAGS=-X $(CONTAINERD_VERSION_PACKAGE).Version=$(GITVERSION) -X $(CONTAINERD_VERSION_PACKAGE).Revision=$(COMMIT) -X $(CONTAINERD_VERSION_PACKAGE).Package=$(CONTAINERD_MODULE) -extldflags -static
EXTRA_GOBUILD_FLAGS=-tags urfave_cli_no_docs,osusergo,netgo,static_build

CGO_CREATE_BINARIES=true

HAS_S3_ARTIFACTS=true
HAS_RELEASE_BRANCHES=true

IMAGE_NAMES=

BUILDSPEC_VARS_KEYS=RELEASE_BRANCH BINARY_PLATFORMS
BUILDSPEC_VARS_VALUES=SUPPORTED_K8S_VERSIONS BINARY_PLATFORMS

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true

PROJECT_DEPENDENCIES=eksa/opencontainers/runc
SIMPLE_CREATE_TARBALLS=false

include $(BASE_DIRECTORY)/Common.mk

ifeq ($(RELEASE_BRANCH),1-37)
.PHONY: seed-go-license-cache
1-37/ATTRIBUTION.txt: | seed-go-license-cache
seed-go-license-cache: $(GO_MOD_DOWNLOAD_TARGETS) | $$(ENABLE_DOCKER)
	@bash $(MAKE_ROOT)/1-37/seed-go-license-cache.sh $(MAKE_ROOT) $(GOLANG_VERSION)
endif

.PHONY: tarballs
tarballs: $(LICENSES_TARGETS_FOR_PREREQ) $(HANDLE_DEPENDENCIES_TARGET) | $$(ENABLE_LOGGING)
	@build/create_tarballs.sh $(TAR_FILE_PREFIX) $(MAKE_ROOT)/$(OUTPUT_DIR) $(MAKE_ROOT)/$(OUTPUT_BIN_DIR) $(GIT_TAG) "$(BINARY_PLATFORMS)" $(ARTIFACTS_PATH) $(GIT_HASH) $(MAKE_ROOT)/$(BINARY_DEPS_DIR)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

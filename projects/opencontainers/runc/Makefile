BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=runc
REPO_OWNER=opencontainers

GITCOMMIT=$(shell git -C $(REPO) describe --dirty --long --always)
EXTRA_GO_LDFLAGS=-X main.VERSION=$(GIT_TAG:v%=%) -X main.gitCommit=$(GITCOMMIT) -extldflags -static
EXTRA_GOBUILD_FLAGS=-tags seccomp,netgo,osusergo
BINARY_TARGET_FILES=runc

CGO_CREATE_BINARIES=true

HAS_S3_ARTIFACTS=true

IMAGE_NAMES=

BUILDSPEC_VARS_KEYS=BINARY_PLATFORMS
BUILDSPEC_VARS_VALUES=BINARY_PLATFORMS
BUILDSPEC_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

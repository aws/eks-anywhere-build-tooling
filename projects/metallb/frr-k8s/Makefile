BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=frr-k8s
REPO_OWNER=metallb
BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=frr-k8s frr-metrics frr-status statuscleaner
SOURCE_PATTERNS=./cmd/frr-k8s-controller ./cmd/metrics ./cmd/status ./cmd/statuscleaner
EXTRA_GO_LDFLAGS=-X github.com/metallb/frr-k8s/internal/version.version=$(GIT_TAG) \
	-X github.com/metallb/frr-k8s/internal/version.gitCommit=$(shell git -C $(REPO) rev-list -n 1 "$(GIT_TAG)") \
	-X github.com/metallb/frr-k8s/internal/version.gitBranch=$(GIT_TAG)

FRR_K8S_IMAGE_COMPONENT=metallb/frr-k8s

# Build the RPM layer natively to avoid QEMU failures.
BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)

EXCLUDE_FROM_STAGING_BUILDSPEC=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true
SKIP_ON_RELEASE_BRANCH=true

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

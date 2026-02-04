BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=tinkerbell
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base

# The mono-repo builds two main binaries:
# - tinkerbell: main binary containing smee, tink-server, tink-controller, rufio, tootles (services we use in EKS-A)
# - tink-agent: worker agent that runs on provisioned machines (replaces tink-worker)
# NOTE: tootles is the renamed hegel metadata service - EKS-A uses it.
# NOTE: secondstar (SSH over serial) is excluded via patch as EKS-A does not use it.
BINARY_TARGET_FILES=tinkerbell tink-agent
SOURCE_PATTERNS=./cmd/tinkerbell ./cmd/agent

# Image components for the mono-repo binaries
TINKERBELL_IMAGE_COMPONENT=tinkerbell/tinkerbell/tinkerbell
TINK_AGENT_IMAGE_COMPONENT=tinkerbell/tinkerbell/tink-agent

IMAGE_NAMES=tinkerbell tink-agent

DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

# Version info - matches pattern from tink/Makefile
VERSION?=$(shell git -C $(REPO) rev-parse --short HEAD)
EXTRA_GO_LDFLAGS=-X main.version=$(VERSION)

# we need to set IMAGE_BUILD_ARGS here even though its the same as the default.
# it is set in Common.mk on the images target and the combine-images target
# since combine-images has images as prereq target, the ?= does not really behavior as one might expect.
# the images target being the actual action, its version of the set takes priority and resets to empty
# setting it explicitly to empty here allows the combine-images override to take priority
IMAGE_BUILD_ARGS=

# Build configuration
BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

HAS_S3_ARTIFACTS=true

# Since we build the arm and amd binaries on different instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
UPLOAD_DO_NOT_DELETE=true

# Helm chart configuration - chart is in the mono-repo at helm/tinkerbell
HAS_HELM_CHART=true
HELM_DIRECTORY=helm/tinkerbell
HELM_SOURCE_REPOSITORY=tinkerbell
HELM_CHART_NAMES=tinkerbell
HELM_DESTINATION_REPOSITORY=tinkerbell/tinkerbell
HELM_IMAGE_LIST=tinkerbell/tinkerbell/tinkerbell tinkerbell/tinkerbell/tink-agent

# Exclude from upgrade buildspec until migration is complete
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

# tinkerbell requires glibc base for ipmitool (bmclib BMC operations)
tinkerbell/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-glibc

s3-artifacts:


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

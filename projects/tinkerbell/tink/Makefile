BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=tink
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

PROJECT_DEPENDENCIES=eksa/isc-projects/dhcp

BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

# we need to set IMAGE_BUILD_ARGS here even though its the same as the default. 
# it is set in Common.mk on the images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L799)
# and the combine-images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L846)
# since combine-images has images as prereq target, the ?= does not really behavior as one might expect.
# the images target being the actual action, its version of the set takes prioirty and resets to empty
# setting it explicitly to empty here takes allows the combine-images override to take proirty
IMAGE_BUILD_ARGS=

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true

TINK_CONTROLLER_IMAGE_COMPONENT=tinkerbell/tink/tink-controller
TINK_SERVER_IMAGE_COMPONENT=tinkerbell/tink/tink-server
TINK_WORKER_IMAGE_COMPONENT=tinkerbell/tink/tink-worker
NGINX_IMAGE_COMPONENT=tinkerbell/tink/nginx
TINK_RELAY_IMAGE_COMPONENT=tinkerbell/tink/tink-relay
TINK_RELAY_INIT_IMAGE_COMPONENT=tinkerbell/tink/tink-relay-init

IMAGE_NAMES=tink-controller tink-server tink-worker nginx tink-relay tink-relay-init

BINARY_TARGET_FILES=tink-controller tink-server tink-worker
SOURCE_PATTERNS=./cmd/tink-controller ./cmd/tink-server ./cmd/tink-worker

VERSION?=$(shell git -C $(REPO) rev-parse --short HEAD)

# When we track a tag instead of commit, this would become $(GIT_TAG)-$(VERSION)
EXTRA_GO_LDFLAGS=-X main.version=$(VERSION)

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk

s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests: tarballs | $$(ENABLE_DOCKER)
	@build/create_manifests.sh $(REPO) $(ARTIFACTS_PATH) $(IMAGE_REPO) $(TINK_SERVER_IMAGE_COMPONENT) ${TINK_CONTROLLER_IMAGE_COMPONENT} $(IMAGE_TAG) $(GOLANG_VERSION)

nginx/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-nginx
tink-relay/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-glibc
tink-relay-init/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-nsenter


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=tink
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

TINK_CONTROLLER_IMAGE_COMPONENT=tinkerbell/tink/tink-controller
TINK_SERVER_IMAGE_COMPONENT=tinkerbell/tink/tink-server
TINK_WORKER_IMAGE_COMPONENT=tinkerbell/tink/tink-worker

IMAGE_NAMES=tink-controller tink-server tink-worker

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


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

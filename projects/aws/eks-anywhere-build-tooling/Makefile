BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib
LATEST_RELEASE_BRANCH=$(shell yq e ".latest" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)
LATEST_RELEASE=$(shell yq e ".$(LATEST_RELEASE_BRANCH).number" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)
GIT_TAG=v$(LATEST_RELEASE_BRANCH)-$(LATEST_RELEASE)

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-minimal-base-docker-client
BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_MINIMAL_BASE_DOCKER_CLIENT_TAG_FILE)
BASE_IMAGE?=$(BASE_REPO)/$(BASE_IMAGE_NAME):$(BASE_TAG)

IMAGE_REPO?=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_NAME=eks-anywhere-cli-tools
IMAGE_TAG?=$(shell git rev-parse HEAD)
IMAGE=$(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

# These will go away after the full release stuff merges
# This is to not block anyone using the images
SECONDARY_IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
SECONDARY_IMAGE=$(SECONDARY_IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
SECONDARY_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/$(IMAGE_NAME):latest

.PHONY: fetch-binaries
fetch-binaries:
	build/fetch_binaries.sh $(RELEASE_BRANCH) $(RELEASE) $(EKS_DISTRO_RELEASE_VERSION)

.PHONY: local-images
local-images: fetch-binaries
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(IMAGE),$(SECONDARY_IMAGE),$(SECONDARY_LATEST_IMAGE)\",dest=/tmp/eks-anywhere-cli-tools.tar

.PHONY: images
images: fetch-binaries
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(IMAGE),$(SECONDARY_IMAGE),$(SECONDARY_LATEST_IMAGE)\",push=true

.PHONY: build
build: local-images
	echo "Done $(COMPONENT)"

.PHONY: release
release: images
	echo "Done $(COMPONENT)"

.PHONY: all
all: release

.PHONY: clean
clean:
	rm -rf $(REPO)
	rm -rf "_output"

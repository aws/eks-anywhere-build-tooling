BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib
GIT_TAG?=$(shell cat GIT_TAG)
COMPONENT=eks-anywhere-test

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/nginx
BASE_IMAGE_NAME?=nginx
BASE_TAG?=$(shell cat $(MAKE_ROOT)/BASE_TAG_FILE)
BASE_IMAGE?=$(BASE_REPO)/$(BASE_IMAGE_NAME):$(BASE_TAG)

IMAGE_REPO?=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_NAME=$(COMPONENT)
IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
IMAGE=$(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
LATEST_IMAGE=$(IMAGE_REPO)/$(IMAGE_NAME):latest

.PHONY: local-images
local-images:
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(IMAGE),$(LATEST_IMAGE)\",dest=/tmp/$(COMPONENT).tar

.PHONY: images
images:
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(IMAGE),$(LATEST_IMAGE)\",push=true

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
	echo "Done $(COMPONENT)"

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib
GIT_TAG?=$(shell cat GIT_TAG)
COMPONENT=test-eks-anywhere

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
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
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(IMAGE),$(LATEST_IMAGE)\",dest=/tmp/$(COMPONENT).tar

.PHONY: images
images:
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
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

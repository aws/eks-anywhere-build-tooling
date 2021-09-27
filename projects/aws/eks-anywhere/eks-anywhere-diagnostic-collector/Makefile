BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib
GIT_TAG=$(shell cat GIT_TAG)

GOPROXY_DNS?=https://proxy.golang.org
export GOPROXY=$(GOPROXY_DNS)

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/eks-distro-build-tooling
COLLECTOR_BASE_IMAGE_NAME?=eks-distro-base
COLLECTOR_BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../../EKS_DISTRO_BASE_TAG_FILE)
COLLECTOR_BASE_IMAGE=$(BASE_REPO)/$(COLLECTOR_BASE_IMAGE_NAME):$(COLLECTOR_BASE_TAG)

IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
COLLECTOR_IMAGE_NAME=eks-anywhere-diagnostic-collector
COLLECTOR_IMAGE=$(IMAGE_REPO)/$(COLLECTOR_IMAGE_NAME):$(IMAGE_TAG)
COLLECTOR_LATEST_IMAGE=$(IMAGE_REPO)/$(COLLECTOR_IMAGE_NAME):latest

.PHONY: local-images
local-images:
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(COLLECTOR_BASE_IMAGE) \
		--local dockerfile=./docker/linux \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(COLLECTOR_IMAGE),$(COLLECTOR_LATEST_IMAGE)\",dest=/tmp/eksa-diagnostic-collector.tar

.PHONY: images
images:
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(COLLECTOR_BASE_IMAGE) \
		--local dockerfile=./docker/linux \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(COLLECTOR_IMAGE),$(COLLECTOR_LATEST_IMAGE)\",push=true

.PHONY: build
build: local-images
	echo "Done"

.PHONY: release
release: images
	echo "Done"

.PHONY: release-upload
release-upload: release
	echo "Done upload"

.PHONY: all
all: release

.PHONY: clean
clean:
	rm -rf "_output"

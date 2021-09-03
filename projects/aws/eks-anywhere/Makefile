BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib
ARTIFACTS_BUCKET?=my-s3-bucket
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

REPO=eks-anywhere
BINARY_NAME=eks-anywhere-cluster-controller
COMPONENT=aws/$(REPO)
CLONE_URL?=git@github.com:$(COMPONENT).git
ifdef CODEBUILD_SRC_DIR
	TAR_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	CLONE_URL=https://git-codecommit.us-west-2.amazonaws.com/v1/repos/aws.eks-anywhere
else
	TAR_PATH?="_output/tar"
endif

GOPROXY_DNS?=https://proxy.golang.org
export GOPROXY=$(GOPROXY_DNS)

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-minimal-base
BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_MINIMAL_BASE_TAG_FILE)
BASE_IMAGE?=$(BASE_REPO)/$(BASE_IMAGE_NAME):$(BASE_TAG)

IMAGE_REPO?=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_NAME=eks-anywhere-cluster-controller
IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
IMAGE=$(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

COLLECTOR_IMAGE_NAME=eks-anywhere-diagnostic-collector
COLLECTOR_IMAGE=$(IMAGE_REPO)/$(COLLECTOR_IMAGE_NAME):$(IMAGE_TAG)

# These will go away after the full release stuff merges
# This is to not block anyone using the images
SECONDARY_IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
SECONDARY_IMAGE=$(SECONDARY_IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
SECONDARY_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/$(IMAGE_NAME):latest

SECONDARY_COLLECTOR_IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
SECONDARY_COLLECTOR_IMAGE=$(SECONDARY_COLLECTOR_IMAGE_REPO)/$(COLLECTOR_IMAGE_NAME):$(IMAGE_TAG)
SECONDARY_COLLECTOR_LATEST_IMAGE=$(SECONDARY_COLLECTOR_IMAGE_REPO)/$(COLLECTOR_IMAGE_NAME):latest

.PHONY: upload-artifacts
upload-artifacts:
	build/upload_artifacts.sh $(TAR_PATH) $(ARTIFACTS_BUCKET) $(PROJECT_PATH) $(CODEBUILD_BUILD_NUMBER) $(CODEBUILD_RESOLVED_SOURCE_VERSION)

.PHONY: binaries
binaries:
	build/create_binaries.sh $(REPO) $(CLONE_URL) $(GIT_TAG) $(GOLANG_VERSION) $(BINARY_NAME) $(IMAGE_REPO) $(IMAGE_TAG)

.PHONY: generate-attribution
generate-attribution:
	build/create_attribution.sh $(GOLANG_VERSION)

.PHONY: tarballs
tarballs:  binaries
	build/create_tarballs.sh $(BINARY_NAME) $(GIT_TAG) $(TAR_PATH)

.PHONY: local-images
local-images: binaries
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/eks-anywhere-cluster-controller \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(IMAGE),$(SECONDARY_IMAGE),$(SECONDARY_LATEST_IMAGE)\",dest=/tmp/eks-anywhere-cluster-controller.tar
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/eksa-diagnostic-collector \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(IMAGE),$(SECONDARY_IMAGE),$(SECONDARY_LATEST_IMAGE)\",dest=/tmp/eksa-collector.tar

.PHONY: images
images: binaries
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/eks-anywhere-cluster-controller \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(IMAGE),$(SECONDARY_IMAGE),$(SECONDARY_LATEST_IMAGE)\",push=true
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/eksa-diagnostic-collector \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(IMAGE),$(SECONDARY_IMAGE),$(SECONDARY_LATEST_IMAGE)\",push=true

.PHONY: build
build: local-images tarballs
	echo "Done $(COMPONENT)"

.PHONY: release
release: images tarballs
	echo "Done $(COMPONENT)"

.PHONY: release-upload
release-upload: release upload-artifacts
	echo "Done upload $(COMPONENT)"

.PHONY: all
all: release

.PHONY: clean
clean:
	rm -rf $(REPO)
	rm -rf "_output"

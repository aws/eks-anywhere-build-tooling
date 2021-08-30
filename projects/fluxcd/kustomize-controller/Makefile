BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
ARTIFACT_BUCKET?=my-s3-bucket
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=kustomize-controller
COMPONENT=fluxcd/$(REPO)
ifdef CODEBUILD_SRC_DIR
	TAR_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	CLONE_URL?=https://git-codecommit.us-west-2.amazonaws.com/v1/repos/fluxcd.kustomize-controller
else
	TAR_PATH?="_output/tar"
	CLONE_URL?=https://github.com/$(COMPONENT).git
endif

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-minimal-base-git
BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_MINIMAL_BASE_GIT_TAG_FILE)
BASE_IMAGE?=$(BASE_REPO)/$(BASE_IMAGE_NAME):$(BASE_TAG)
BUILDER_IMAGE=$(BASE_REPO)/eks-distro-base:$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_BASE_TAG_FILE)

IMAGE_REPO?=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_NAME=$(COMPONENT)
IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
IMAGE=$(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

# These will go away after the full release stuff merges
# This is to not block anyone using the images
SECONDARY_IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
SECONDARY_IMAGE=$(SECONDARY_IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
SECONDARY_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/$(IMAGE_NAME):latest

.PHONY: create-binaries
create-binaries:
	build/create_binaries.sh $(REPO) $(CLONE_URL) $(GIT_TAG) $(GOLANG_VERSION)

.PHONY: fetch-binaries
fetch-binaries:
	build/fetch_binaries.sh

.PHONY: binaries
binaries: create-binaries fetch-binaries

.PHONY: generate-attribution
generate-attribution:
	build/create_attribution.sh $(GOLANG_VERSION)

.PHONY: tarballs
tarballs:  binaries
	build/create_tarballs.sh $(REPO) $(GIT_TAG) $(TAR_PATH)

.PHONY: local-images
local-images: binaries
	buildctl \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(IMAGE),$(SECONDARY_IMAGE),$(SECONDARY_LATEST_IMAGE)\",dest=/tmp/kustomize-controller.tar

.PHONY: images
images: binaries
	/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(IMAGE),$(SECONDARY_IMAGE),$(SECONDARY_LATEST_IMAGE)\",push=true

.PHONY: build
build: local-images tarballs
	echo "Done $(COMPONENT)"

.PHONY: release
release: images tarballs
	echo "Done $(COMPONENT)"

.PHONY: release-upload
release-upload: release
	echo "Done upload $(COMPONENT)"

.PHONY: all
all: release

.PHONY: clean
clean:
	rm -rf $(REPO)
	rm -rf "_output"

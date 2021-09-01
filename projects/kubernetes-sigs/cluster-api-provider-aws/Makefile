BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib
ARTIFACTS_BUCKET?=my-s3-bucket
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=cluster-api-provider-aws
COMPONENT=kubernetes-sigs/$(REPO)
ifdef CODEBUILD_SRC_DIR
	TAR_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	CLONE_URL?=https://git-codecommit.us-west-2.amazonaws.com/v1/repos/kubernetes-sigs.cluster-api-provider-aws
else
	TAR_PATH?="_output/tar"
	CLONE_URL?=https://github.com/$(COMPONENT).git
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
IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
CAPA_CONTROLLER_IMAGE=$(IMAGE_REPO)/kubernetes-sigs/cluster-api-provider-aws/cluster-api-aws-controller:$(IMAGE_TAG)
CAPA_EKS_BOOTSTRAP_CONTROLLER_IMAGE=$(IMAGE_REPO)/kubernetes-sigs/cluster-api-provider-aws/eks-bootstrap-controller:$(IMAGE_TAG)
CAPA_EKS_CONTROL_PLANE_CONTROLLER_IMAGE=$(IMAGE_REPO)/kubernetes-sigs/cluster-api-provider-aws/eks-control-plane-controller:$(IMAGE_TAG)

# These will go away after the full release stuff merges
# This is to not block anyone using the images
SECONDARY_IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
SECONDARY_CAPA_CONTROLLER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api-provider-aws/cluster-api-aws-controller:$(IMAGE_TAG)
SECONDARY_CAPA_CONTROLLER_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api-provider-aws/cluster-api-aws-controller:latest

SECONDARY_CAPA_EKS_BOOTSTRAP_CONTROLLER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api-provider-aws/eks-bootstrap-controller:$(IMAGE_TAG)
SECONDARY_CAPA_EKS_BOOTSTRAP_CONTROLLER_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api-provider-aws/eks-bootstrap-controller:latest

SECONDARY_CAPA_EKS_CONTROL_PLANE_CONTROLLER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api-provider-aws/eks-control-plane-controller:$(IMAGE_TAG)
SECONDARY_CAPA_EKS_CONTROL_PLANE_CONTROLLER_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api-provider-aws/eks-control-plane-controller:latest

.PHONY: upload-artifacts
upload-artifacts:
	build/upload_artifacts.sh $(TAR_PATH) $(ARTIFACTS_BUCKET) $(PROJECT_PATH) $(CODEBUILD_BUILD_NUMBER) $(CODEBUILD_RESOLVED_SOURCE_VERSION)

.PHONY: binaries
binaries:
	build/create_binaries.sh $(REPO) $(CLONE_URL) $(GIT_TAG) $(GOLANG_VERSION) $(IMAGE_REPO) $(IMAGE_TAG)

.PHONY: generate-attribution
generate-attribution:
	build/create_attribution.sh $(GOLANG_VERSION)

.PHONY: tarballs
tarballs:  binaries
	build/create_tarballs.sh $(REPO) $(GIT_TAG) $(TAR_PATH)

.PHONY: local-images
local-images: binaries
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/cluster-api-aws-controller \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CAPA_CONTROLLER_IMAGE),$(SECONDARY_CAPA_CONTROLLER_IMAGE),$(SECONDARY_CAPA_CONTROLLER_LATEST_IMAGE)\",dest=/tmp/capa-controller.tar
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/eks-bootstrap-controller \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CAPA_EKS_BOOTSTRAP_CONTROLLER_IMAGE),$(SECONDARY_CAPA_EKS_BOOTSTRAP_CONTROLLER_IMAGE),$(SECONDARY_CAPA_EKS_BOOTSTRAP_CONTROLLER_LATEST_IMAGE)\",dest=/tmp/capa-eks-bootstrap-controller.tar
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/eks-control-plane-controller \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CAPA_EKS_CONTROL_PLANE_CONTROLLER_IMAGE),$(SECONDARY_CAPA_EKS_CONTROL_PLANE_CONTROLLER_IMAGE),$(SECONDARY_CAPA_EKS_CONTROL_PLANE_CONTROLLER_LATEST_IMAGE)\",dest=/tmp/capa-eks-control-plane-controller.tar

.PHONY: images
images: binaries
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/cluster-api-aws-controller \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CAPA_CONTROLLER_IMAGE),$(SECONDARY_CAPA_CONTROLLER_IMAGE),$(SECONDARY_CAPA_CONTROLLER_LATEST_IMAGE)\",push=true
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/eks-bootstrap-controller \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CAPA_EKS_BOOTSTRAP_CONTROLLER_IMAGE),$(SECONDARY_CAPA_EKS_BOOTSTRAP_CONTROLLER_IMAGE),$(SECONDARY_CAPA_EKS_BOOTSTRAP_CONTROLLER_LATEST_IMAGE)\",push=true
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/eks-control-plane-controller \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CAPA_EKS_CONTROL_PLANE_CONTROLLER_IMAGE),$(SECONDARY_CAPA_EKS_CONTROL_PLANE_CONTROLLER_IMAGE),$(SECONDARY_CAPA_EKS_CONTROL_PLANE_CONTROLLER_LATEST_IMAGE)\",push=true

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

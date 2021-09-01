BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
ARTIFACT_BUCKET?=my-s3-bucket
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=cert-manager
COMPONENT=jetstack/$(REPO)
ifdef CODEBUILD_SRC_DIR
	TAR_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	CLONE_URL?=https://git-codecommit.us-west-2.amazonaws.com/v1/repos/jetstack.cert-manager
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
BUILDER_IMAGE=$(BASE_REPO)/eks-distro-base:$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_BASE_TAG_FILE)

IMAGE_REPO?=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
CERT_MANAGER_ACMESOLVER_IMAGE=$(IMAGE_REPO)/jetstack/cert-manager-acmesolver:$(IMAGE_TAG)
CERT_MANAGER_CAINJECTOR_IMAGE=$(IMAGE_REPO)/jetstack/cert-manager-cainjector:$(IMAGE_TAG)
CERT_MANAGER_CONTROLLER_IMAGE=$(IMAGE_REPO)/jetstack/cert-manager-controller:$(IMAGE_TAG)
CERT_MANAGER_WEBHOOK_IMAGE=$(IMAGE_REPO)/jetstack/cert-manager-webhook:$(IMAGE_TAG)

# These will go away after the full release stuff merges
# This is to not block anyone using the images
SECONDARY_IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
SECONDARY_CERT_MANAGER_ACMESOLVER_IMAGE=$(SECONDARY_IMAGE_REPO)/jetstack/cert-manager-acmesolver:$(IMAGE_TAG)
SECONDARY_CERT_MANAGER_CAINJECTOR_IMAGE=$(SECONDARY_IMAGE_REPO)/jetstack/cert-manager-cainjector:$(IMAGE_TAG)
SECONDARY_CERT_MANAGER_CONTROLLER_IMAGE=$(SECONDARY_IMAGE_REPO)/jetstack/cert-manager-controller:$(IMAGE_TAG)
SECONDARY_CERT_MANAGER_WEBHOOK_IMAGE=$(SECONDARY_IMAGE_REPO)/jetstack/cert-manager-webhook:$(IMAGE_TAG)
SECONDARY_LATEST_CERT_MANAGER_ACMESOLVER_IMAGE=$(SECONDARY_IMAGE_REPO)/jetstack/cert-manager-acmesolver:latest
SECONDARY_LATEST_CERT_MANAGER_CAINJECTOR_IMAGE=$(SECONDARY_IMAGE_REPO)/jetstack/cert-manager-cainjector:latest
SECONDARY_LATEST_CERT_MANAGER_CONTROLLER_IMAGE=$(SECONDARY_IMAGE_REPO)/jetstack/cert-manager-controller:latest
SECONDARY_LATEST_CERT_MANAGER_WEBHOOK_IMAGE=$(SECONDARY_IMAGE_REPO)/jetstack/cert-manager-webhook:latest

.PHONY: binaries
binaries:
	build/create_binaries.sh $(REPO) $(CLONE_URL) $(GIT_TAG) $(GOLANG_VERSION)

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
		--local dockerfile=./docker/linux/cert-manager-acmesolver \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CERT_MANAGER_ACMESOLVER_IMAGE),$(SECONDARY_CERT_MANAGER_ACMESOLVER_IMAGE),$(SECONDARY_LATEST_CERT_MANAGER_ACMESOLVER_IMAGE)\",dest=/tmp/cert-manager-acmesolver.tar
	buildctl \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--local dockerfile=./docker/linux/cert-manager-cainjector \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CERT_MANAGER_CAINJECTOR_IMAGE),$(SECONDARY_CERT_MANAGER_CAINJECTOR_IMAGE),$(SECONDARY_LATEST_CERT_MANAGER_CAINJECTOR_IMAGE)\",dest=/tmp/cert-manager-cainjector.tar
	buildctl \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--local dockerfile=./docker/linux/cert-manager-controller \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CERT_MANAGER_CONTROLLER_IMAGE),$(SECONDARY_CERT_MANAGER_CONTROLLER_IMAGE),$(SECONDARY_LATEST_CERT_MANAGER_CONTROLLER_IMAGE)\",dest=/tmp/cert-manager-controller.tar
	buildctl \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/cert-manager-webhook \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CERT_MANAGER_WEBHOOK_IMAGE),$(SECONDARY_CERT_MANAGER_WEBHOOK_IMAGE),$(SECONDARY_LATEST_CERT_MANAGER_WEBHOOK_IMAGE)\",dest=/tmp/cert-manager-webhook.tar

.PHONY: images
images: binaries
	/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--local dockerfile=./docker/linux/cert-manager-acmesolver \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CERT_MANAGER_ACMESOLVER_IMAGE),$(SECONDARY_CERT_MANAGER_ACMESOLVER_IMAGE),$(SECONDARY_LATEST_CERT_MANAGER_ACMESOLVER_IMAGE)\",push=true
	/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--local dockerfile=./docker/linux/cert-manager-cainjector \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CERT_MANAGER_CAINJECTOR_IMAGE),$(SECONDARY_CERT_MANAGER_CAINJECTOR_IMAGE),$(SECONDARY_LATEST_CERT_MANAGER_CAINJECTOR_IMAGE)\",push=true
	/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--local dockerfile=./docker/linux/cert-manager-controller \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CERT_MANAGER_CONTROLLER_IMAGE),$(SECONDARY_CERT_MANAGER_CONTROLLER_IMAGE),$(SECONDARY_LATEST_CERT_MANAGER_CONTROLLER_IMAGE)\",push=true
	/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/cert-manager-webhook \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CERT_MANAGER_WEBHOOK_IMAGE),$(SECONDARY_CERT_MANAGER_WEBHOOK_IMAGE),$(SECONDARY_LATEST_CERT_MANAGER_WEBHOOK_IMAGE)\",push=true
	
.PHONY: build
build: local-images tarballs
	echo "Done $(COMPONENT)"

.PHONY: release
release: images tarballs
	echo "Done $(COMPONENT)"

.PHONY: all
all: release

.PHONY: release-upload
release-upload: release
	echo "Done upload $(COMPONENT)"

.PHONY: clean
clean:
	rm -rf $(REPO)
	rm -rf _output
	

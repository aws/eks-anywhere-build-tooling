BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib

ARTIFACTS_BUCKET?=my-s3-bucket

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION?=us-west-2

REPO=cloud-provider-vsphere
COMPONENT=kubernetes/$(REPO)
CLONE_URL?=https://github.com/$(COMPONENT).git
ifdef CODEBUILD_SRC_DIR
	TAR_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	CLONE_URL?=https://git-codecommit.$(AWS_REGION).amazonaws.com/v1/repos/kubernetes.cloud-provider-vsphere
else
	TAR_PATH?="_output/tar"
	CLONE_URL?=https://github.com/$(COMPONENT).git
endif

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-minimal-base
BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_MINIMAL_BASE_TAG_FILE)
BASE_IMAGE?=$(BASE_REPO)/$(BASE_IMAGE_NAME):$(BASE_TAG)

IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_NAME=kubernetes/cloud-provider-vsphere/cpi/manager
IMAGE=$(IMAGE_REPO)/$(IMAGE_NAME)

SUPPORTED_K8S_VERSIONS=$(shell yq e 'keys | .[]' $(BASE_DIRECTORY)/projects/kubernetes-sigs/image-builder/BOTTLEROCKET_OVA_RELEASES)
BINARY_TARGETS=$(addprefix binaries-, $(SUPPORTED_K8S_VERSIONS))
BUILD_TARGETS=$(addprefix local-images-, $(SUPPORTED_K8S_VERSIONS)) $(addprefix tarballs-, $(SUPPORTED_K8S_VERSIONS))
RELEASE_TARGETS=$(addprefix images-, $(SUPPORTED_K8S_VERSIONS)) $(addprefix tarballs-, $(SUPPORTED_K8S_VERSIONS))
ATTRIBUTION_TARGETS=$(addprefix generate-attribution-, $(SUPPORTED_K8S_VERSIONS))
## --------------------------------------
## Help
## --------------------------------------
##@ Helpers
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-35s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


.PHONY: upload-artifacts
upload-artifacts:
	build/upload_artifacts.sh $(TAR_PATH) $(ARTIFACTS_BUCKET) $(PROJECT_PATH) $(CODEBUILD_BUILD_NUMBER) $(CODEBUILD_RESOLVED_SOURCE_VERSION)

.PHONY: binaries
binaries: $(BINARY_TARGETS)

.PHONY: binaries-%
binaries-%: GOLANG_VERSION=$(shell cat ./$*/GOLANG_VERSION)
binaries-%: GIT_TAG=$(shell cat ./$*/GIT_TAG)
binaries-%:
	build/create_binaries.sh $(REPO) $* $(CLONE_URL) $(GIT_TAG) $(GOLANG_VERSION)

.PHONY: generate-attribution-%
generate-attribution-%: GOLANG_VERSION=$(shell cat ./$*/GOLANG_VERSION)
generate-attribution-%: binaries-%
	build/create_attribution.sh $* $(GOLANG_VERSION)

.PHONY: generate-attribution
generate-attribution: $(ATTRIBUTION_TARGETS)

.PHONY: tarballs-%
tarballs-%: GIT_TAG=$(shell cat ./$*/GIT_TAG)
tarballs-%: binaries-%
	build/create_tarballs.sh $(REPO) $* $(GIT_TAG) $(TAR_PATH)

.PHONY: local-images-%
local-images-%: GIT_TAG=$(shell cat ./$*/GIT_TAG)
local-images-%: IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
local-images-%: RELEASE_LATEST=$(GIT_TAG)-latest
local-images-%: binaries-%
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:RELEASE_BRANCH=$* \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(IMAGE):$(IMAGE_TAG),$(IMAGE):$(RELEASE_LATEST)\",dest=/tmp/cloud-provider-vsphere.tar

.PHONY: images-%
images-%: GIT_TAG=$(shell cat ./$*/GIT_TAG)
images-%: IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
images-%: RELEASE_LATEST=$(GIT_TAG)-latest
images-%: binaries-%
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:RELEASE_BRANCH=$* \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(IMAGE):$(IMAGE_TAG),$(IMAGE):$(RELEASE_LATEST)\",push=true

.PHONY: build
build: $(BUILD_TARGETS)
	echo "Done $(COMPONENT)"

.PHONY: release
release: $(RELEASE_TARGETS)
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

## --------------------------------------
## Document dynamic build targets
## --------------------------------------
##@ Builds
binaries-1-18: ## Build binaries for 1-18 cloud-provider
binaries-1-19: ## Build binaries for 1-19 cloud-provider
binaries-1-20: ## Build binaries for 1-20 cloud-provider
binaries-1-21: ## Build binaries for 1-21 cloud-provider

local-images-1-18: ## Build image with 1-18 cloud-provider
local-images-1-19: ## Build image with 1-19 cloud-provider
local-images-1-20: ## Build image with 1-20 cloud-provider
local-images-1-21: ## Build image with 1-21 cloud-provider

##@ Releases
images-1-18: ## Build and push image with 1-18 cloud-provider
images-1-19: ## Build and push image with 1-19 cloud-provider
images-1-20: ## Build and push image with 1-20 cloud-provider
images-1-21: ## Build and push image with 1-21 cloud-provider

tarballs-1-18: ## Build tarballs with 1-18 cloud-provider
tarballs-1-19: ## Build tarballs with 1-19 cloud-provider
tarballs-1-20: ## Build tarballs with 1-20 cloud-provider
tarballs-1-21: ## Build tarballs with 1-21 cloud-provider
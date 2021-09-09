BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib
AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-base
BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_BASE_TAG_FILE)
BASE_IMAGE?=$(BASE_REPO)/$(BASE_IMAGE_NAME):$(BASE_TAG)

IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_NAME=bottlerocket-bootstrap
IMAGE=$(IMAGE_REPO)/$(IMAGE_NAME)

SUPPORTED_K8S_VERSIONS=$(shell yq e 'keys | .[]' $(BASE_DIRECTORY)/projects/kubernetes-sigs/image-builder/BOTTLEROCKET_OVA_RELEASES)
BUILD_TARGETS=$(addprefix local-images-, $(SUPPORTED_K8S_VERSIONS))
RELEASE_TARGETS=$(addprefix images-, $(SUPPORTED_K8S_VERSIONS))

## --------------------------------------
## Help
## --------------------------------------
##@ Helpers
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-35s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: bootstrap-binary
bootstrap-binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bottlerocket-bootstrap github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap

.PHONY: fetch-binaries-%
fetch-binaries-%:
	build/fetch_binaries.sh $*

.PHONY: local-images-%
local-images-%: GIT_TAG=v$(*)-$(shell yq e ".$(*).number" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)
local-images-%: IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
local-images-%: RELEASE_LATEST=$(GIT_TAG)-latest
local-images-%: fetch-binaries-% bootstrap-binary
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:RELEASE_BRANCH=$(*) \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(IMAGE):$(IMAGE_TAG),$(IMAGE):$(RELEASE_LATEST)\",dest=/tmp/bottlerocket-bootstrap.tar

.PHONY: images-%
images-%: GIT_TAG=v$(*)-$(shell yq e ".$(*).number" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)
images-%: IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
images-%: RELEASE_LATEST=$(GIT_TAG)-latest
images-%: fetch-binaries-% bootstrap-binary
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:RELEASE_BRANCH=$(*) \
		--local dockerfile=./docker/linux/ \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(IMAGE):$(IMAGE_TAG),$(IMAGE):$(RELEASE_LATEST)\",push=true

.PHONY: build
build: ## Build images for all Kubernetes versions to ECR
build: $(BUILD_TARGETS)
	echo "Done $(COMPONENT)"

.PHONY: release
release: ## Build and push images for all Kubernetes versions to ECR
release: $(RELEASE_TARGETS)
	echo "Done $(COMPONENT)"

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
local-images-1-18: ## Build image with Kubernetes 1-18
local-images-1-19: ## Build image with Kubernetes 1-19
local-images-1-20: ## Build image with Kubernetes 1-20

##@ Releases
images-1-18: ## Build and push image with Kubernetes 1-18
images-1-19: ## Build and push image with Kubernetes 1-19
images-1-20: ## Build and push image with Kubernetes 1-20

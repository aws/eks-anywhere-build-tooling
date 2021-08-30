REPO_ROOT="$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")/../../../" && pwd -P)"

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
ARTIFACTS_BUCKET?=my-s3-bucket
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.16"

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

REPO=kind
COMPONENT=kubernetes-sigs/$(REPO)
ifdef CODEBUILD_SRC_DIR
	TAR_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	CLONE_URL?=https://git-codecommit.us-west-2.amazonaws.com/v1/repos/kubernetes-sigs.kind
else
	TAR_PATH?="_output/tar"
	CLONE_URL?=https://github.com/$(COMPONENT).git
endif

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-base
BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_BASE_TAG_FILE)
EKS_DISTRO_BASE_IMAGE?=$(BASE_REPO)/$(BASE_IMAGE_NAME):$(BASE_TAG)
MINIMAL_IPTABLES_BASE_IMAGE_NAME?=eks-distro-minimal-base-iptables
MINIMAL_IPTABLES_BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_MINIMAL_BASE_IPTABLES_TAG_FILE)
EKS_DISTRO_MINIMAL_BASE_IPTABLES_IMAGE?=$(BASE_REPO)/$(MINIMAL_IPTABLES_BASE_IMAGE_NAME):$(MINIMAL_IPTABLES_BASE_TAG)

IMAGE_REPO?=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
KIND_BASE_IMAGE_NAME=kindest/base
KIND_NODE_IMAGE_NAME=kubernetes-sigs/kind/node
KIND_KINDNETD_IMAGE_NAME=kubernetes-sigs/kind/kindnetd

KIND_BASE_IMAGE=$(IMAGE_REPO)/$(KIND_BASE_IMAGE_NAME)
KIND_NODE_IMAGE=$(IMAGE_REPO)/${KIND_NODE_IMAGE_NAME}
KIND_KINDNETD_IMAGE=$(IMAGE_REPO)/$(KIND_KINDNETD_IMAGE_NAME):$(IMAGE_TAG)

# These will go away after the full release stuff merges
# This is to not block anyone using the images
SECONDARY_IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
SECONDARY_KIND_NODE_IMAGE=$(SECONDARY_IMAGE_REPO)/$(KIND_NODE_IMAGE_NAME)
SECONDARY_KIND_KINDNETD_IMAGE=$(SECONDARY_IMAGE_REPO)/$(KIND_KINDNETD_IMAGE_NAME):$(IMAGE_TAG)
SECONDARY_KIND_KINDNETD_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/$(KIND_KINDNETD_IMAGE_NAME):latest

PUSH=false

SUPPORTED_K8S_VERSIONS=$(shell yq e 'keys | .[]' $(REPO_ROOT)/projects/kubernetes-sigs/image-builder/BOTTLEROCKET_OVA_RELEASES)
BUILD_TARGETS=$(addprefix build-kind-node-, $(SUPPORTED_K8S_VERSIONS))

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
binaries:
	build/create_binaries.sh $(REPO) $(CLONE_URL) $(GIT_TAG) $(GOLANG_VERSION)

.PHONY: generate-attribution
generate-attribution:
	build/create_attribution.sh $(GOLANG_VERSION)

.PHONY: tarballs
tarballs:  binaries
	build/create_tarballs.sh $(REPO) $(GIT_TAG) $(TAR_PATH)

.PHONY: build-kind-node-%
build-kind-node-%:
	# base image is still tied to a specific eks version because the cni_plugins are included
	build/build-kind-base-image.sh $* $(KIND_BASE_IMAGE) $(ARTIFACTS_BUCKET) $(EKS_DISTRO_BASE_IMAGE) 
	build/build-kind-node-image.sh $* $(KIND_BASE_IMAGE) $(KIND_NODE_IMAGE) $(SECONDARY_KIND_NODE_IMAGE) $(KIND_KINDNETD_IMAGE) $(IMAGE_REPO) $(IMAGE_TAG) $(ARTIFACTS_BUCKET) $(PUSH)

.PHONY: build-kindnetd-image
build-kindnetd-image:
	build/build-kindnetd-image.sh $(EKS_DISTRO_MINIMAL_BASE_IPTABLES_IMAGE) $(KIND_KINDNETD_IMAGE) $(SECONDARY_KIND_KINDNETD_IMAGE) $(SECONDARY_KIND_KINDNETD_LATEST_IMAGE) $(PUSH)


##@ Builds

.PHONY: build
build: ## Build binaries + kindnetd + node images for every kubernetes version
build: tarballs build-kindnetd-image $(BUILD_TARGETS)
	echo "Done $(COMPONENT)"

.PHONY: release
release: ## Build binaries + kindnetd + node images for every kubernetes version and push to ecr
release: PUSH=true
release: build
	echo "Done $(COMPONENT)"

.PHONY: release-upload
release-upload: release upload-artifacts

.PHONY: clean
clean:
	rm -rf $(REPO)
	rm -rf "_output"

## --------------------------------------
## Document dynamic build targets
## --------------------------------------
build-kind-node-1-18: ## Build kind node image for kubernetes 1.18
build-kind-node-1-19: ## Build kind node image for kubernetes 1.19
build-kind-node-1-20: ## Build kind node image for kubernetes 1.20
build-kind-node-1-21: ## Build kind node image for kubernetes 1.21

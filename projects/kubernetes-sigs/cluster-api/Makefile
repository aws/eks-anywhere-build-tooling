BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
ARTIFACTS_BUCKET?=my-s3-bucket
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=cluster-api
COMPONENT=kubernetes-sigs/$(REPO)
ifdef CODEBUILD_SRC_DIR
	TAR_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	CLONE_URL?=https://git-codecommit.us-west-2.amazonaws.com/v1/repos/kubernetes-sigs.cluster-api
else
	TAR_PATH?="_output/tar"
	CLONE_URL?=https://github.com/$(COMPONENT).git
endif

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-minimal-base-nonroot
BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_MINIMAL_BASE_NONROOT_TAG_FILE)
BASE_IMAGE?=$(BASE_REPO)/$(BASE_IMAGE_NAME):$(BASE_TAG)
CAPD_BASE_IMAGE_NAME?=eks-distro-minimal-base-docker-client
CAPD_BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_MINIMAL_BASE_DOCKER_CLIENT_TAG_FILE)
CAPD_BASE_IMAGE?=$(BASE_REPO)/$(CAPD_BASE_IMAGE_NAME):$(CAPD_BASE_TAG)

IMAGE_REPO?=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
CAPI_CONTROLLER_IMAGE=$(IMAGE_REPO)/kubernetes-sigs/cluster-api/cluster-api-controller:$(IMAGE_TAG)
CAPI_KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE=$(IMAGE_REPO)/kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller:$(IMAGE_TAG)
CAPI_KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE=$(IMAGE_REPO)/kubernetes-sigs/cluster-api/kubeadm-control-plane-controller:$(IMAGE_TAG)
CAPD_CONTROLLER_IMAGE=$(IMAGE_REPO)/kubernetes-sigs/cluster-api/capd-manager:$(IMAGE_TAG)

# These will go away after the full release stuff merges
# This is to not block anyone using the images
SECONDARY_IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
SECONDARY_CAPI_CONTROLLER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api/cluster-api-controller:$(IMAGE_TAG)
SECONDARY_CAPI_CONTROLLER_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api/cluster-api-controller:latest

SECONDARY_KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller:$(IMAGE_TAG)
SECONDARY_KUBEADM_BOOTSTRAP_CONTROLLER_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller:latest

SECONDARY_KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api/kubeadm-control-plane-controller:$(IMAGE_TAG)
SECONDARY_KUBEADM_CONTROL_PLANE_CONTROLLER_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api/kubeadm-control-plane-controller:latest

SECONDARY_CAPD_CONTROLLER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api/capd-manager:$(IMAGE_TAG)
SECONDARY_CAPD_CONTROLLER_LATEST_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/cluster-api/capd-manager:latest

.PHONY: upload-artifacts
upload-artifacts:
	build/upload_artifacts.sh $(TAR_PATH) $(ARTIFACTS_BUCKET) $(PROJECT_PATH) $(CODEBUILD_BUILD_NUMBER) $(CODEBUILD_RESOLVED_SOURCE_VERSION)

.PHONY: create-binaries
create-binaries:
	build/create_binaries.sh $(REPO) $(CLONE_URL) $(GIT_TAG) $(GOLANG_VERSION) $(IMAGE_REPO) $(IMAGE_TAG)

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
		--local dockerfile=./docker/linux/cluster-api-controller \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CAPI_CONTROLLER_IMAGE),$(SECONDARY_CAPI_CONTROLLER_IMAGE),$(SECONDARY_CAPI_CONTROLLER_LATEST_IMAGE)\",dest=/tmp/capi-controller.tar
	buildctl \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/kubeadm-bootstrap-controller \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CAPI_KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE),$(SECONDARY_KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE),$(SECONDARY_KUBEADM_BOOTSTRAP_CONTROLLER_LATEST_IMAGE)\",dest=/tmp/capi-kubeadm-bootstrap-controller.tar
	buildctl \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/kubeadm-control-plane-controller \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CAPI_KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE),$(SECONDARY_KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE),$(SECONDARY_KUBEADM_CONTROL_PLANE_CONTROLLER_LATEST_IMAGE)\",dest=/tmp/capi-kubeadm-control-plane-controller.tar
	buildctl \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(CAPD_BASE_IMAGE) \
		--local dockerfile=./docker/linux/cluster-api-docker-controller \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CAPD_CONTROLLER_IMAGE),$(SECONDARY_CAPD_CONTROLLER_IMAGE),$(SECONDARY_CAPD_CONTROLLER_LATEST_IMAGE)\",dest=/tmp/capd-controller.tar

.PHONY: images
images: binaries
	/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/cluster-api-controller \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CAPI_CONTROLLER_IMAGE),$(SECONDARY_CAPI_CONTROLLER_IMAGE),$(SECONDARY_CAPI_CONTROLLER_LATEST_IMAGE)\",push=true
	/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/kubeadm-bootstrap-controller \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CAPI_KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE),$(SECONDARY_KUBEADM_BOOTSTRAP_CONTROLLER_IMAGE),$(SECONDARY_KUBEADM_BOOTSTRAP_CONTROLLER_LATEST_IMAGE)\",push=true
	/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/kubeadm-control-plane-controller \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CAPI_KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE),$(SECONDARY_KUBEADM_CONTROL_PLANE_CONTROLLER_IMAGE),$(SECONDARY_KUBEADM_CONTROL_PLANE_CONTROLLER_LATEST_IMAGE)\",push=true
	/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(CAPD_BASE_IMAGE) \
		--local dockerfile=./docker/linux/cluster-api-docker-controller \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CAPD_CONTROLLER_IMAGE),$(SECONDARY_CAPD_CONTROLLER_IMAGE),$(SECONDARY_CAPD_CONTROLLER_LATEST_IMAGE)\",push=true

.PHONY: build
build: local-images tarballs
	echo "Done $(COMPONENT)"

.PHONY: release
release: images tarballs
	echo "Done $(COMPONENT)"

.PHONY: all
all: release

.PHONY: release-upload
release-upload: release upload-artifacts
	echo "Done upload $(COMPONENT)"

.PHONY: clean
clean:
	rm -rf $(REPO)
	rm -rf _output
	rm -rf kubernetes
	rm -f kubernetes-client-linux-amd64.tar.gz

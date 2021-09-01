BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib
ARTIFACTS_BUCKET?=my-s3-bucket
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=vsphere-csi-driver
COMPONENT=kubernetes-sigs/$(REPO)
ifdef CODEBUILD_SRC_DIR
	TAR_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	CLONE_URL?=https://git-codecommit.us-west-2.amazonaws.com/v1/repos/kubernetes-sigs.vsphere-csi-driver
else
	TAR_PATH?="_output/tar"
	CLONE_URL?=https://github.com/$(COMPONENT).git
endif
GOPROXY_DNS?=direct
export GOPROXY=$(GOPROXY_DNS)

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION=us-west-2

MAKE_ROOT=$(shell cd "$(shell dirname "${BASH_SOURCE[0]}")" && pwd -P)
BASE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-minimal-base
BASE_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_MINIMAL_BASE_TAG_FILE)
BASE_IMAGE?=$(BASE_REPO)/$(BASE_IMAGE_NAME):$(BASE_TAG)
BASE_CSI_IMAGE_NAME?=eks-distro-minimal-base-csi
BASE_CSI_TAG?=$(shell cat $(MAKE_ROOT)/../../../EKS_DISTRO_MINIMAL_BASE_CSI_TAG_FILE)
BASE_CSI_IMAGE?=$(BASE_REPO)/$(BASE_CSI_IMAGE_NAME):$(BASE_CSI_TAG)

IMAGE_REPO?=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_TAG?=$(GIT_TAG)-$(shell git rev-parse HEAD)
CSI_DRIVER_IMAGE=$(IMAGE_REPO)/kubernetes-sigs/vsphere-csi-driver/csi/driver:$(IMAGE_TAG)
CSI_SYNCER_IMAGE=$(IMAGE_REPO)/kubernetes-sigs/vsphere-csi-driver/csi/syncer:$(IMAGE_TAG)

# These will go away after the full release stuff merges
# This is to not block anyone using the images
SECONDARY_IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
SECONDARY_DRIVER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/vsphere-csi-driver/csi/driver:$(IMAGE_TAG)
SECONDARY_LATEST_DRIVER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/vsphere-csi-driver/csi/driver:latest
SECONDARY_SYNCER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/vsphere-csi-driver/csi/syncer:$(IMAGE_TAG)
SECONDARY_LATEST_SYNCER_IMAGE=$(SECONDARY_IMAGE_REPO)/kubernetes-sigs/vsphere-csi-driver/csi/syncer:latest

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

.PHONY: local-images
local-images: binaries
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_CSI_IMAGE) \
		--local dockerfile=./docker/linux/csi-driver \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CSI_DRIVER_IMAGE),$(SECONDARY_DRIVER_IMAGE),$(SECONDARY_LATEST_DRIVER_IMAGE)\",dest=/tmp/csi-driver.tar
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/csi-syncer \
		--local context=. \
		--output type=oci,oci-mediatypes=true,\"name=$(CSI_SYNCER_IMAGE),$(SECONDARY_SYNCER_IMAGE),$(SECONDARY_LATEST_SYNCER_IMAGE)\",dest=/tmp/csi-syncer.tar

.PHONY: images
images: binaries
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_CSI_IMAGE) \
		--local dockerfile=./docker/linux/csi-driver \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CSI_DRIVER_IMAGE),$(SECONDARY_DRIVER_IMAGE),$(SECONDARY_LATEST_DRIVER_IMAGE)\",push=true
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=linux/amd64 \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--local dockerfile=./docker/linux/csi-syncer \
		--local context=. \
		--output type=image,oci-mediatypes=true,\"name=$(CSI_SYNCER_IMAGE),$(SECONDARY_SYNCER_IMAGE),$(SECONDARY_LATEST_SYNCER_IMAGE)\",push=true

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

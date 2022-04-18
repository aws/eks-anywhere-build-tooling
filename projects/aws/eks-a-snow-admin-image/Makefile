BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG=0.0.0-dev # This value does not matter but needs to be set

SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false
REPO_NO_CLONE=true
IMAGE_NAMES=

REPO=eks-a-snow-admin-image
REPO_OWNER=aws

BUILD_TARGETS=validate
RELEASE_TARGETS=build-raw-image

PACKER_VERSION=1.8.0

EKSA_VERSION?=v0.8.0
EKSA_RELEASE_MANIFEST_URL?=https://anywhere-assets.eks.amazonaws.com/releases/eks-a/manifest.yaml

ARTIFACTS_UPLOAD_PATH?=$(PROJECT_PATH)
EXPORT_AMI_BUCKET?=$(ARTIFACTS_BUCKET)
FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ami/snow
MANIFEST_OUTPUT?=$(FULL_OUTPUT_DIR)/manifest.json
AMI_S3_DST=$(EXPORT_AMI_BUCKET)/$(ARTIFACTS_UPLOAD_PATH)
EXPORT_AMI_DST=$(AMI_S3_DST)/$(EKSA_VERSION)/$(GIT_HASH)/
LATEST_AMI_S3_URL=$(AMI_S3_DST)/$(LATEST)/snow-admin.raw

DEPENDENCIES_FOLDER=$(MAKE_ROOT)/$(OUTPUT_DIR)/dependencies
PACKER=$(DEPENDENCIES_FOLDER)/packer

include $(BASE_DIRECTORY)/Common.mk


$(PACKER): OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
$(PACKER):
	mkdir -p $(@D)
	curl --silent -L https://releases.hashicorp.com/packer/$(PACKER_VERSION)/packer_$(PACKER_VERSION)_$(OS)_amd64.zip > packer.zip
	unzip packer.zip -d $(@D)
	rm packer.zip
	$@ init .	

$(FULL_OUTPUT_DIR):
	mkdir -p $@

.PHONY: build-raw-image
build-raw-image: build-ami export-to-raw

.PHONY: build-ami
build-ami: export PKR_VAR_eks-a-version = $(EKSA_VERSION)
build-ami: export PKR_VAR_eks-a-release-manifest-url = $(EKSA_RELEASE_MANIFEST_URL)
build-ami: export PKR_VAR_build-version = $(shell echo $(GIT_HASH) | cut -c1-7)
build-ami: export PKR_VAR_manifest-output = $(MANIFEST_OUTPUT)
build-ami: KUBECTL_URL?=$(shell ./build/get_kubectl_url.sh)
build-ami: export PKR_VAR_kubectl-url = $(KUBECTL_URL)
build-ami: $(PACKER) $(FULL_OUTPUT_DIR)
	$(PACKER) build snow-admin.pkr.hcl

.PHONY: export-to-raw
export-to-raw:
	build/export-ami-to-s3.sh $(MANIFEST_OUTPUT) raw $(EXPORT_AMI_DST) $(LATEST_AMI_S3_URL)

.PHONY: validate
validate: $(PACKER)
	$(PACKER) fmt -check .
	$(PACKER) validate .

.PHONY: format
format: $(PACKER)
	$(PACKER) fmt .


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

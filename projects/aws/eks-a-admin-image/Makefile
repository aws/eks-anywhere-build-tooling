-include .env

BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG=0.0.0-dev # This value does not matter but needs to be set

SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false
REPO_NO_CLONE=true
IMAGE_NAMES=

REPO=eks-a-admin-image
REPO_OWNER=aws

BUILD_TARGETS=validate
RELEASE_TARGETS=build-raw-image

PACKER_VERSION=1.8.0

EKSA_VERSION?=v0.0.0
EKSA_RELEASE_MANIFEST_URL?=https://dev-release-prod-pdx.s3.us-west-2.amazonaws.com/eks-a-release.yaml

FULL_OUTPUT_DIR?=$(MAKE_ROOT)/$(OUTPUT_DIR)
AMI_OUTPUT_DIR?=$(FULL_OUTPUT_DIR)/ami
OVA_OUTPUT_DIR?=$(FULL_OUTPUT_DIR)/ova
MANIFEST_FILE?=manifest.json
AMI_MANIFEST_PATH?=$(AMI_OUTPUT_DIR)/$(MANIFEST_FILE)
OVA_MANIFEST_PATH?=$(OVA_OUTPUT_DIR)/$(MANIFEST_FILE)

ARTIFACTS_UPLOAD_PATH?=$(PROJECT_PATH)
EXPORT_AMI_BUCKET?=$(ARTIFACTS_BUCKET)
AMI_S3_DST=$(EXPORT_AMI_BUCKET)/$(ARTIFACTS_UPLOAD_PATH)
EXPORT_AMI_DST=$(AMI_S3_DST)/$(EKSA_VERSION)/$(GIT_HASH)/
LATEST_AMI_S3_URL=$(AMI_S3_DST)/$(LATEST)/snow-admin.raw

EXPORT_OVA_BUCKET?=$(ARTIFACTS_BUCKET)
EKSA_ADMIN_PROJECT_PATH?=projects/aws/eks-a-admin-image
OVA_PROJECT_PATH?=$(EKSA_ADMIN_PROJECT_PATH)/ova


include $(BASE_DIRECTORY)/Common.mk


PACKER=$(BINARY_DEPS_DIR)/packer

build-%: export PKR_VAR_eks-a-version=$(EKSA_VERSION)
build-%: export PKR_VAR_eks-a-release-manifest-url=$(EKSA_RELEASE_MANIFEST_URL)
build-%: export PKR_VAR_build-version=$(shell echo $(GIT_HASH) | cut -c1-7)
build-%: KUBECTL_URL?=$(shell ./build/get_kubectl_url.sh)
build-%: export PKR_VAR_kubectl-url=$(KUBECTL_URL)


$(PACKER): OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
$(PACKER):
	mkdir -p $(@D)
	curl --silent -L https://releases.hashicorp.com/packer/$(PACKER_VERSION)/packer_$(PACKER_VERSION)_$(OS)_amd64.zip > packer.zip
	unzip packer.zip -d $(@D)
	rm packer.zip

$(FULL_OUTPUT_DIR):
	mkdir -p $@

$(AMI_OUTPUT_DIR): $(FULL_OUTPUT_DIR)
	mkdir -p $@

$(OVA_OUTPUT_DIR): $(FULL_OUTPUT_DIR)
	mkdir -p $@

.PHONY: build-raw-image
build-raw-image: build-ami export-to-raw

.PHONY: build-ami
build-ami: export PKR_VAR_manifest-output=$(AMI_MANIFEST_PATH)
build-ami: $(PACKER) init clean-ami $(AMI_OUTPUT_DIR)
	$(PACKER) build -only='*.amazon-ebs.*' .

.PHONY: export-to-raw
export-to-raw:
	build/export-ami-to-s3.sh $(AMI_OUTPUT_DIR)/$(MANIFEST_FILE) raw $(EXPORT_AMI_DST) $(LATEST_AMI_S3_URL)

.PHONY: build-ova
build-ova: export PKR_VAR_manifest-output=$(OVA_MANIFEST_PATH)
build-ova: export PKR_VAR_vsphere_password = $(VSPHERE_PASSWORD)
build-ova: export PKR_VAR_vsphere_username = $(VSPHERE_USERNAME)
build-ova: export PKR_VAR_vsphere_endpoint = $(VSPHERE_ENDPOINT)
build-ova: export PKR_VAR_vsphere_datacenter = $(VSPHERE_DATACENTER)
build-ova: export PKR_VAR_vsphere_datastore = $(VSPHERE_DATASTORE)
build-ova: export PKR_VAR_vsphere_cluster = $(VSPHERE_CLUSTER)
build-ova: export PKR_VAR_vsphere_network = $(VSPHERE_NETWORK)
build-ova: export PKR_VAR_build-password = $(BUILD_PASSWORD)
build-ova: export PKR_VAR_build-output = $(OVA_OUTPUT_DIR)
build-ova: KIND_URL=$(shell ./build/get_kind_url.sh $(ARTIFACTS_BUCKET) )
build-ova: export PKR_VAR_kind-url=$(KIND_URL)
build-ova: $(PACKER) init clean-ova $(OVA_OUTPUT_DIR)
	$(PACKER) build -only='*.vsphere-iso.*' -var-file=ova/linux/ubuntu/vsphere.ubuntu.pkrvars.hcl .

.PHONY: release-ova
release-ova: build-ova publish-ova-to-s3

.PHONY: clean-ova
clean-ova:
	rm -rf $(OVA_OUTPUT_DIR)

.PHONY: clean-ami
clean-ami:
	rm -rf $(AMI_OUTPUT_DIR)

.PHONY: clean
clean: clean-ova clean-ami

.PHONY: publish-ova-to-s3
publish-ova-to-s3:
	aws s3 cp $(OVA_OUTPUT_DIR)/*.ova $(EXPORT_OVA_BUCKET)/$(OVA_PROJECT_PATH)/

.PHONY: init
init: $(PACKER)
	$(PACKER) init .

.PHONY: validate
validate: $(PACKER) init
	$(PACKER) fmt -check .
	$(PACKER) validate .

.PHONY: format
format: $(PACKER) init
	$(PACKER) fmt .


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

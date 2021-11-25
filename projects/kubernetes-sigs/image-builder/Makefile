BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=image-builder
REPO_OWNER=kubernetes-sigs

IMAGE_BUILDER_DIR=$(REPO)/images/capi

OVA_PATH=$(IMAGE_BUILDER_DIR)/output
ifdef CODEBUILD_SRC_DIR
	OVA_PATH=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/image-builder/images/capi/output
endif

BOTTLEROCKET_DOWNLOAD_PATH?=$(OUTPUT_DIR)/bottlerocket/downloads
CARGO_HOME=$(OUTPUT_DIR)/cargo

VSPHERE_CONNECTION_DATA?={}
# Aws accounts to share built AMI with
DEV_ACCOUNTS?=

PACKER_AMI_VAR_FILES=$(MAKE_ROOT)/packer/ami/packer.json
PACKER_AMI_SHARE_FILE=$(MAKE_ROOT)/packer/ami/share-ami.json
PACKER_AMI_VAR_FILES+=$(PACKER_AMI_SHARE_FILE)

PACKER_OVA_CONF_FILE=$(MAKE_ROOT)/packer/ova/vsphere.json
PACKER_OVA_VAR_FILES=$(PACKER_OVA_CONF_FILE)
PACKER_OVA_VAR_FILES+=$(MAKE_ROOT)/packer/ova/vsphere-library.json

PACKER_TYPE_VAR_FILES?=
PACKER_VAR_FILES=$(MAKE_ROOT)/$(OUTPUT_DIR)/config/kubernetes.json $(PACKER_TYPE_VAR_FILES) $(MAKE_ROOT)/$(OUTPUT_DIR)/config/common.json $(MAKE_ROOT)/$(OUTPUT_DIR)/config/cni.json $(MAKE_ROOT)/$(OUTPUT_DIR)/config/additional_components.json

BUILD_TARGETS=build-ami-ubuntu-2004 release-ova-bottlerocket
RELEASE_TARGETS= release-ami-ubuntu-2004 release-ova-ubuntu-2004 upload-artifacts release-ova-bottlerocket

# Kubeadm hardcodes the version of the pause image it pulls, even tho containerd config
# overrides which image is actually used
# Setting up this mapping so we can build the additional image into the image
ADDITIONAL_PAUSE_1-18_FROM=1-18
ADDITIONAL_PAUSE_1-19_FROM=1-18
ADDITIONAL_PAUSE_1-20_FROM=1-19
ADDITIONAL_PAUSE_1-21_FROM=1-20

HAS_RELEASE_BRANCHES=true

include $(BASE_DIRECTORY)/Common.mk


.PHONY: setup-ami-share
setup-ami-share:
	jq -n '{"ami_users": "$(DEV_ACCOUNTS)", "snapshot_users": "$(DEV_ACCOUNTS)"}' > $(PACKER_AMI_SHARE_FILE)

.PHONY: setup-vsphere
setup-vsphere:
	echo $(VSPHERE_CONNECTION_DATA) > $(PACKER_OVA_CONF_FILE)

.PHONY: setup-bottlerocket
setup-bottlerocket:
	build/setup_bottlerocket.sh $(BOTTLEROCKET_DOWNLOAD_PATH) $(CARGO_HOME)

.PHONY: deps-ami
deps-ami: $(GIT_PATCH_TARGET)
	make -C $(IMAGE_BUILDER_DIR) deps-ami

.PHONY: deps-ova
deps-ova: $(GIT_PATCH_TARGET)
	make -C $(IMAGE_BUILDER_DIR) deps-ova

.PHONY: setup-packer-configs
setup-packer-configs:
	build/setup_packer_configs.sh $(RELEASE_BRANCH) $(ARTIFACTS_BUCKET) $(OVA_PATH) $(ADDITIONAL_PAUSE_$(RELEASE_BRANCH)_FROM)

.PHONY: build-ami-ubuntu-2004
build-ami-ubuntu-2004: PACKER_TYPE_VAR_FILES=$(PACKER_AMI_VAR_FILES)
build-ami-ubuntu-2004: setup-ami-share deps-ami setup-packer-configs
	PACKER_VAR_FILES="$(PACKER_VAR_FILES)" make -C $(IMAGE_BUILDER_DIR) validate-ami-ubuntu-2004

.PHONY: release-ami-ubuntu-2004
release-ami-ubuntu-2004: PACKER_TYPE_VAR_FILES=$(PACKER_AMI_VAR_FILES)
release-ami-ubuntu-2004-%: setup-ami-share deps-ami setup-packer-configs
	PACKER_VAR_FILES="$(PACKER_VAR_FILES)" make -C $(IMAGE_BUILDER_DIR) build-ami-ubuntu-2004

.PHONY: release-ova-ubuntu-2004
build-ova-ubuntu-2004: PACKER_TYPE_VAR_FILES=$(PACKER_OVA_VAR_FILES)
release-ova-ubuntu-2004: deps-ova setup-vsphere setup-packer-configs
	PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(OVA_PATH)/packer.log PACKER_VAR_FILES="$(PACKER_VAR_FILES)" \
		OVF_CUSTOM_PROPERTIES="$(MAKE_ROOT)/$(OUTPUT_DIR)/config/ovf_custom_properties.json" \
		make -C $(IMAGE_BUILDER_DIR) build-node-ova-vsphere-ubuntu-2004
	build/get_artifacts.sh $(IMAGE_BUILDER_DIR) ubuntu

.PHONY: release-ova-bottlerocket
release-ova-bottlerocket: setup-bottlerocket
	build/upload_bottlerocket_artifacts.sh $(RELEASE_BRANCH) bottlerocket $(BOTTLEROCKET_DOWNLOAD_PATH) $(CARGO_HOME) $(PROJECT_PATH)/$(RELEASE_BRANCH) $(LATEST_TAG)

.PHONY: validate-ubuntu-2004
validate-ubuntu-2004: check-env-validation $(GIT_PATCH_TARGET) setup-packer-configs
	sh validate/validate.sh $(MAKE_ROOT)/$(OUTPUT_DIR)/config/validate_goss_inline_vars.json

.PHONY: check-env-validation
check-env-validation:
	ifndef TEMPLATE
		$(error Environment var TEMPLATE not set. Example: TEMPLATE=<ubuntu.livecd> make ...)
	endif

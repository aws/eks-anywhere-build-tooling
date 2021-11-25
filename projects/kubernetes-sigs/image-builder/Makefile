BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=image-builder
REPO_OWNER=kubernetes-sigs

IMAGE_BUILDER_DIR=$(REPO)/images/capi

PROJECT_PATH?=projects/$(REPO_OWNER)/$(REPO)
OVA_PATH=$(IMAGE_BUILDER_DIR)/output
ifdef CODEBUILD_SRC_DIR
	OVA_PATH=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/image-builder/images/capi/output
endif

FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)
# Force non release-branched
CARGO_HOME=$(MAKE_ROOT)/_output/cargo
RUSTUP_HOME=$(MAKE_ROOT)/_output/rustup
BOTTLEROCKET_DOWNLOAD_PATH?=$(FULL_OUTPUT_DIR)/bottlerocket/downloads

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
PACKER_VAR_FILES=$(FULL_OUTPUT_DIR)/config/kubernetes.json $(PACKER_TYPE_VAR_FILES) $(FULL_OUTPUT_DIR)/config/common.json $(FULL_OUTPUT_DIR)/config/cni.json $(FULL_OUTPUT_DIR)/config/additional_components.json

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

TUFTOOL_TARGET=$(CARGO_HOME)/bin/tuftool
BOOTLEROCKET_SETUP_TARGET=$(BOTTLEROCKET_DOWNLOAD_PATH)/bottlerocket-root-json-checksum

include $(BASE_DIRECTORY)/Common.mk

export PATH:=$(CARGO_HOME)/bin:$(MAKE_ROOT)/$(IMAGE_BUILDER_DIR)/.local/bin:$(PATH)

.PHONY: setup-ami-share
setup-ami-share:
	jq -n '{"ami_users": "$(DEV_ACCOUNTS)", "snapshot_users": "$(DEV_ACCOUNTS)"}' > $(PACKER_AMI_SHARE_FILE)

.PHONY: setup-vsphere
setup-vsphere:
	echo $(VSPHERE_CONNECTION_DATA) > $(PACKER_OVA_CONF_FILE)

$(TUFTOOL_TARGET):
	# This code installs the Rust toolchain manager called rustup along
	# with other Rust binaries such as rustc, rustfmt. It also installs Cargo,
	# the Rust package manager which is then used to install Tuftool.
	@mkdir -p $(CARGO_HOME)
	@mkdir -p $(RUSTUP_HOME)
	curl https://sh.rustup.rs -sSf | CARGO_HOME=$(CARGO_HOME) RUSTUP_HOME=$(RUSTUP_HOME) sh -s -- -y
	$(CARGO_HOME)/bin/rustup default stable
	CARGO_NET_GIT_FETCH_WITH_CLI=true $(CARGO_HOME)/bin/cargo install --force --root $(CARGO_HOME) tuftool

$(BOOTLEROCKET_SETUP_TARGET): export BOTTLEROCKET_ROOT_JSON_PATH=$(BOTTLEROCKET_DOWNLOAD_PATH)/root.json
$(BOOTLEROCKET_SETUP_TARGET):
	@mkdir -p $(BOTTLEROCKET_DOWNLOAD_PATH)
	# This configuration supports local installations and checksum validations
	# of root.json file
	envsubst '$$BOTTLEROCKET_ROOT_JSON_PATH' \
		< $(MAKE_ROOT)/bottlerocket-root-json-checksum \
		> $(BOOTLEROCKET_SETUP_TARGET)
	curl https://cache.bottlerocket.aws/root.json -o $$BOTTLEROCKET_ROOT_JSON_PATH
	sha512sum -c $(BOOTLEROCKET_SETUP_TARGET)

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
		OVF_CUSTOM_PROPERTIES="$(FULL_OUTPUT_DIR)/config/ovf_custom_properties.json" \
		make -C $(IMAGE_BUILDER_DIR) build-node-ova-vsphere-ubuntu-2004
	build/get_artifacts.sh $(IMAGE_BUILDER_DIR) ubuntu

.PHONY: release-ova-bottlerocket
release-ova-bottlerocket: $(TUFTOOL_TARGET) $(BOOTLEROCKET_SETUP_TARGET)
	build/upload_bottlerocket_artifacts.sh $(RELEASE_BRANCH) bottlerocket $(BOTTLEROCKET_DOWNLOAD_PATH) $(CARGO_HOME) $(PROJECT_PATH)/$(RELEASE_BRANCH) $(LATEST_TAG)

.PHONY: validate-ubuntu-2004
validate-ubuntu-2004: check-env-validation $(GIT_PATCH_TARGET) setup-packer-configs
	sh validate/validate.sh $(FULL_OUTPUT_DIR)/config/validate_goss_inline_vars.json

.PHONY: check-env-validation
check-env-validation:
	ifndef TEMPLATE
		$(error Environment var TEMPLATE not set. Example: TEMPLATE=<ubuntu.livecd> make ...)
	endif

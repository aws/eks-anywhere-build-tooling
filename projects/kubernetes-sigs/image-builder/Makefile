BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=image-builder
REPO_OWNER=kubernetes-sigs

IMAGE_BUILDER_DIR=$(REPO)/images/capi
RAW_IMAGE_BUILD_AMI?=ami-0bd99362b6ccd24ea
RAW_IMAGE_BUILD_INSTANCE_TYPE?=t1.micro
RAW_IMAGE_BUILD_KEY_NAME?=raw-image-build

FULL_OUTPUT_DIR=

# Force non release-branched
CARGO_HOME=$(MAKE_ROOT)/_output/cargo
RUSTUP_HOME=$(MAKE_ROOT)/_output/rustup
export BOTTLEROCKET_DOWNLOAD_PATH?=$(FULL_OUTPUT_DIR)/bottlerocket/downloads

export BUILDER_ROOT=$(MAKE_ROOT)

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

# Kubeadm hardcodes the version of the pause image it pulls, even tho containerd config
# overrides which image is actually used
# Setting up this mapping so we can build the additional image into the image
# The versions are dependent on the tag set for the Kubernetes version we are consuming vs what is hardcoded in the
# kubeadm constants for that specific version:
# Example for v1.22.4 on why we rely on the pause tag from v1.21.5, which is the previous version that we support:
# Kubernetes pause tag in v1.22.4: https://github.com/kubernetes/kubernetes/blob/v1.22.4/build/pause/Makefile#L20
# Kubernetes pause tag in v1.21.5: https://github.com/kubernetes/kubernetes/blob/v1.21.5/build/pause/Makefile#L20
# Kubeadm constants: https://github.com/kubernetes/kubernetes/blob/v1.22.4/cmd/kubeadm/app/constants/constants.go#L397
ADDITIONAL_PAUSE_1-20_FROM=1-19
ADDITIONAL_PAUSE_1-21_FROM=1-20
ADDITIONAL_PAUSE_1-22_FROM=1-21

HAS_RELEASE_BRANCHES=true
HAS_LICENSES=false
SIMPLE_CREATE_TARBALLS=false

# to support a no op attribution target
TARGETS_ALLOWED_WITH_NO_RELEASE_BRANCH=build release clean help binaries checksums attribution release-ami-ubuntu-2004
TUFTOOL_TARGET=$(CARGO_HOME)/bin/tuftool
BOTTLEROCKET_SETUP_TARGET=$(BOTTLEROCKET_DOWNLOAD_PATH)/bottlerocket-root-json-checksum

FINAL_UBUNTU_OVA_PATH=$(ARTIFACTS_PATH)/ova/ubuntu.ova
FINAL_UBUNTU_RAW_IMAGE_PATH=$(ARTIFACTS_PATH)/raw/ubuntu.gz
FINAL_BOTTLEROCKET_OVA_PATH=$(ARTIFACTS_PATH)/ova/bottlerocket.ova
FAKE_UBUNTU_OVA_PATH=$(IMAGE_BUILDER_DIR)/output/fake-ubuntu.ova
FAKE_UBUNTU_RAW_PATH=$(IMAGE_BUILDER_DIR)/output/fake-ubuntu.gz

BUILD_AMI_TARGETS=build-ami-ubuntu-2004
BUILD_OVA_TARGETS=setup-packer-configs-ova release-ova-bottlerocket $(FAKE_UBUNTU_OVA_PATH) $(FINAL_UBUNTU_OVA_PATH) $(FINAL_BOTTLEROCKET_OVA_PATH) upload-artifacts-ova
BUILD_RAW_TARGETS=release-raw-ubuntu-2004 $(FAKE_UBUNTU_RAW_PATH) $(FINAL_UBUNTU_RAW_IMAGE_PATH) upload-artifacts-raw
BUILD_TARGETS=$(BUILD_RAW_TARGETS) $(BUILD_AMI_TARGETS) $(BUILD_OVA_TARGETS)
ifeq ($(IMAGE_FORMAT),ova)
	S3_TARGET_PREREQUISITES=$(FINAL_UBUNTU_OVA_PATH) $(FINAL_BOTTLEROCKET_OVA_PATH)
	RELEASE_TARGETS=release-ova-ubuntu-2004 release-ova-bottlerocket upload-artifacts-ova
else ifeq ($(IMAGE_FORMAT),raw)
	S3_TARGET_PREREQUISITES=$(FINAL_UBUNTU_RAW_IMAGE_PATH)
	RELEASE_TARGETS=release-raw-ubuntu-2004 upload-artifacts-raw
endif

include $(BASE_DIRECTORY)/Common.mk

export PATH:=$(CARGO_HOME)/bin:$(MAKE_ROOT)/$(IMAGE_BUILDER_DIR)/.local/bin:$(PATH)

# Since we do not build the ova in presubmit but want to validate upload-artifacts behavior
$(FAKE_UBUNTU_OVA_PATH):
	@mkdir -p $(@D)
	touch $@
	touch $(ARTIFACTS_PATH)/ova/packer.log

$(FAKE_UBUNTU_RAW_PATH):
	@mkdir -p $(@D)
	touch $@
	touch $(@D)/packer.log

$(FINAL_UBUNTU_OVA_PATH):
	mv $(IMAGE_BUILDER_DIR)/output/*.ova $@

$(FINAL_UBUNTU_RAW_IMAGE_PATH):
	mv $(IMAGE_BUILDER_DIR)/output/*.gz $@
	mv $(IMAGE_BUILDER_DIR)/output/packer.log $(@D)/packer.log

$(FINAL_BOTTLEROCKET_OVA_PATH): FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ova
$(FINAL_BOTTLEROCKET_OVA_PATH):
	mv $(BOTTLEROCKET_DOWNLOAD_PATH)/ova/*.ova $@

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

$(BOTTLEROCKET_SETUP_TARGET): FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ova
$(BOTTLEROCKET_SETUP_TARGET): export BOTTLEROCKET_ROOT_JSON_PATH=$(BOTTLEROCKET_DOWNLOAD_PATH)/root.json
$(BOTTLEROCKET_SETUP_TARGET):
	@mkdir -p $(BOTTLEROCKET_DOWNLOAD_PATH)
	# This configuration supports local installations and checksum validations
	# of root.json file
	envsubst '$$BOTTLEROCKET_ROOT_JSON_PATH' \
		< $(MAKE_ROOT)/bottlerocket-root-json-checksum \
		> $(BOTTLEROCKET_SETUP_TARGET)
	curl https://cache.bottlerocket.aws/root.json -o $$BOTTLEROCKET_ROOT_JSON_PATH
	sha512sum -c $(BOTTLEROCKET_SETUP_TARGET)

.PHONY: deps-%
deps-%: MAKEFLAGS=
deps-%: $(GIT_PATCH_TARGET)
	$(MAKE) -C $(IMAGE_BUILDER_DIR) deps-$*

.PHONY: setup-packer-configs-%
setup-packer-configs-%:
	build/setup_packer_configs.sh $(RELEASE_BRANCH) $* $(ARTIFACTS_BUCKET) $(ARTIFACTS_PATH)/$* $(ADDITIONAL_PAUSE_$(RELEASE_BRANCH)_FROM) $(LATEST)

.PHONY: build-ami-ubuntu-2004
build-ami-ubuntu-2004: MAKEFLAGS=
build-ami-ubuntu-2004: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ami
build-ami-ubuntu-2004: PACKER_TYPE_VAR_FILES=$(PACKER_AMI_VAR_FILES)
build-ami-ubuntu-2004: setup-ami-share deps-ami setup-packer-configs-ami
	PACKER_VAR_FILES="$(PACKER_VAR_FILES)" $(MAKE) -C $(IMAGE_BUILDER_DIR) validate-ami-ubuntu-2004

.PHONY: release-ami-ubuntu-2004
release-ami-ubuntu-2004: MAKEFLAGS=
release-ami-ubuntu-2004: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ami
release-ami-ubuntu-2004: PACKER_TYPE_VAR_FILES=$(PACKER_AMI_VAR_FILES)
release-ami-ubuntu-2004: export MANIFEST_OUTPUT=$(FULL_OUTPUT_DIR)/manifest.json
release-ami-ubuntu-2004: EXPORT_AMI_BUCKET?=$(ARTIFACTS_BUCKET)
release-ami-ubuntu-2004: AMI_S3_DST=$(EXPORT_AMI_BUCKET)/$(ARTIFACTS_UPLOAD_PATH)/ami
release-ami-ubuntu-2004: EXPORT_AMI_DST=$(AMI_S3_DST)/$(GIT_HASH)
release-ami-ubuntu-2004: LATEST_AMI_S3_URL=$(AMI_S3_DST)/$(LATEST)/ubuntu.raw
release-ami-ubuntu-2004: setup-ami-share deps-ami setup-packer-configs-ami
	PACKER_LOG=1 PACKER_VAR_FILES="$(PACKER_VAR_FILES)" $(MAKE) -C $(IMAGE_BUILDER_DIR) build-ami-ubuntu-2004
	build/export-ami-to-s3.sh $(RELEASE_BRANCH) $(MANIFEST_OUTPUT) raw $(EXPORT_AMI_DST) $(LATEST_AMI_S3_URL)

.PHONY: release-ova-ubuntu-2004
release-ova-ubuntu-2004: MAKEFLAGS=
release-ova-ubuntu-2004: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ova
release-ova-ubuntu-2004: PACKER_TYPE_VAR_FILES=$(PACKER_OVA_VAR_FILES)
release-ova-ubuntu-2004: deps-ova setup-vsphere setup-packer-configs-ova
	PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(ARTIFACTS_PATH)/ova/packer.log PACKER_VAR_FILES="$(PACKER_VAR_FILES)" \
		OVF_CUSTOM_PROPERTIES="$(FULL_OUTPUT_DIR)/config/ovf_custom_properties.json" \
		$(MAKE) -C $(IMAGE_BUILDER_DIR) build-node-ova-vsphere-ubuntu-2004

.PHONY: release-ova-bottlerocket
release-ova-bottlerocket: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ova
release-ova-bottlerocket: $(TUFTOOL_TARGET) $(BOTTLEROCKET_SETUP_TARGET)
	build/get_bottlerocket_artifacts.sh $(RELEASE_BRANCH) bottlerocket $(BOTTLEROCKET_DOWNLOAD_PATH) $(CARGO_HOME) $(PROJECT_PATH)/$(RELEASE_BRANCH) $(LATEST_TAG)

.PHONY: release-raw-ubuntu-2004
release-raw-ubuntu-2004: MAKEFLAGS=
release-raw-ubuntu-2004: FULL_OUTPUT_DIR=$(subst $(BASE_DIRECTORY),/home/ec2-user/$(shell basename $(BASE_DIRECTORY)),$(MAKE_ROOT)/$(OUTPUT_DIR))/raw
release-raw-ubuntu-2004: deps-raw setup-packer-configs-raw
	build/build_raw_image.sh $(BASE_DIRECTORY) $(PROJECT_PATH) $(IMAGE_BUILDER_DIR) "$(PACKER_VAR_FILES)" $(RAW_IMAGE_BUILD_AMI) $(RAW_IMAGE_BUILD_INSTANCE_TYPE) $(RAW_IMAGE_BUILD_KEY_NAME)

.PHONY: build-qemu-rhel-local
build-qemu-rhel-local: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/qemu
build-qemu-rhel-local: deps-qemu setup-packer-configs-qemu
	build/build_qemu_rhel_local.sh . $(BASE_IMAGE) $(IMAGE_BUILDER_DIR) "$(PACKER_VAR_FILES)"

.PHONY: validate-ubuntu-2004
validate-ubuntu-2004: check-env-validation $(GIT_PATCH_TARGET) setup-packer-configs-ova
	sh validate/validate.sh $(FULL_OUTPUT_DIR)/config/validate_goss_inline_vars.json

.PHONY: check-env-validation
check-env-validation:
	ifndef TEMPLATE
		$(error Environment var TEMPLATE not set. Example: TEMPLATE=<ubuntu.livecd> make ...)
	endif

.PHONY: s3-artifacts-%
s3-artifacts-%: ARTIFACTS_PATH_IMAGE_FORMAT=$(ARTIFACTS_PATH)/$*
s3-artifacts-%: $(S3_TARGET_PREREQUISITES)
	$(MAKE) -C $(MAKE_ROOT) s3-artifacts ARTIFACTS_PATH=$(ARTIFACTS_PATH_IMAGE_FORMAT) IMAGE_FORMAT=$*

.PHONY: upload-artifacts-%
upload-artifacts-%: ARTIFACTS_PATH_IMAGE_FORMAT=$(ARTIFACTS_PATH)/$*
upload-artifacts-%: ARTIFACTS_UPLOAD_PATH_IMAGE_FORMAT=$(ARTIFACTS_UPLOAD_PATH)/$*
upload-artifacts-%: s3-artifacts-%
	$(MAKE) -C $(MAKE_ROOT) upload-artifacts ARTIFACTS_PATH=$(ARTIFACTS_PATH_IMAGE_FORMAT) ARTIFACTS_UPLOAD_PATH=$(ARTIFACTS_UPLOAD_PATH_IMAGE_FORMAT) IMAGE_FORMAT=$*


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

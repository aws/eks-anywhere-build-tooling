BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
MAKE_ROOT=$(BASE_DIRECTORY)/projects/kubernetes-sigs/image-builder
GIT_TAG:=$(shell cat GIT_TAG)

REPO=image-builder
REPO_OWNER=kubernetes-sigs

IMAGE_BUILDER_DIR=$(REPO)/images/capi
RAW_IMAGE_BUILD_AMI?=ami-095413544ce52437d
RAW_IMAGE_BUILD_INSTANCE_TYPE?=t2.micro
RAW_IMAGE_BUILD_KEY_NAME?=raw-image-build

FULL_OUTPUT_DIR=
export BOTTLEROCKET_DOWNLOAD_PATH?=$(FULL_OUTPUT_DIR)/bottlerocket/downloads

RHSM_USERNAME?=
RHSM_PASSWORD?=
export BUILDER_ROOT=$(MAKE_ROOT)
export RHSM_USER=$(RHSM_USERNAME)
export RHSM_PASS=$(RHSM_PASSWORD)

VSPHERE_CONNECTION_DATA?={}
# Aws accounts to share built AMI with
DEV_ACCOUNTS?=

PACKER_AMI_VAR_FILES=$(MAKE_ROOT)/packer/ami/ami.json
PACKER_AMI_SHARE_FILE=$(MAKE_ROOT)/packer/ami/share-ami.json
PACKER_AMI_VAR_FILES+=$(PACKER_AMI_SHARE_FILE)

PACKER_OVA_CONF_FILE=$(MAKE_ROOT)/packer/ova/vsphere.json
PACKER_OVA_VAR_FILES=$(MAKE_ROOT)/packer/ova/vsphere-library.json
PACKER_OVA_VAR_FILES+=$(PACKER_OVA_CONF_FILE)

PACKER_TYPE_VAR_FILES?=
PACKER_VAR_FILES=$(FULL_OUTPUT_DIR)/config/kubernetes.json $(FULL_OUTPUT_DIR)/config/common.json $(FULL_OUTPUT_DIR)/config/cni.json $(FULL_OUTPUT_DIR)/config/additional_components.json $(PACKER_TYPE_VAR_FILES)

# Kubeadm hardcodes the version of the pause image it pulls, even tho containerd config
# overrides which image is actually used
# Setting up this mapping so we can build the additional image into the image
# The versions are dependent on the tag set for the Kubernetes version we are consuming vs what is hardcoded in the
# kubeadm constants for that specific version:
# Example for v1.22.4 on why we rely on the pause tag from v1.21.5, which is the previous version that we support:
# Kubernetes pause tag in v1.22.4: https://github.com/kubernetes/kubernetes/blob/v1.22.4/build/pause/Makefile#L20
# Kubernetes pause tag in v1.21.5: https://github.com/kubernetes/kubernetes/blob/v1.21.5/build/pause/Makefile#L20
# Kubeadm constants: https://github.com/kubernetes/kubernetes/blob/v1.22.4/cmd/kubeadm/app/constants/constants.go#L397
ADDITIONAL_PAUSE_1-21_FROM=1-20
ADDITIONAL_PAUSE_1-22_FROM=1-21
ADDITIONAL_PAUSE_1-23_FROM=1-22
ADDITIONAL_PAUSE_1-24_FROM=1-23
ADDITIONAL_PAUSE_1-25_FROM=1-24

HAS_RELEASE_BRANCHES=true
HAS_LICENSES=false
SIMPLE_CREATE_TARBALLS=false

# to support a no op attribution target
TARGETS_ALLOWED_WITH_NO_RELEASE_BRANCH=binaries checksums attribution release-ami-ubuntu-2004
BOTTLEROCKET_SETUP_TARGET=$(BOTTLEROCKET_DOWNLOAD_PATH)/bottlerocket-root-json-checksum

FINAL_OVA_PATH=$(ARTIFACTS_PATH)/ova/$(IMAGE_OS)/$(IMAGE_OS).ova
FINAL_RAW_IMAGE_PATH=$(ARTIFACTS_PATH)/raw/$(IMAGE_OS)/$(IMAGE_OS).gz
FINAL_CLOUDSTACK_IMAGE_PATH=$(ARTIFACTS_PATH)/cloudstack/$(IMAGE_OS)/$(IMAGE_OS).qcow2
FINAL_BOTTLEROCKET_OVA_PATH=$(ARTIFACTS_PATH)/ova/bottlerocket/bottlerocket.ova
FINAL_BOTTLEROCKET_RAW_PATH=$(ARTIFACTS_PATH)/raw/bottlerocket/bottlerocket.img.gz
FAKE_UBUNTU_OVA_PATH=fake-ubuntu.ova
FAKE_UBUNTU_RAW_PATH=fake-ubuntu.gz
BOTTLEROCKET_OVA_RELEASE_VERSION=$(shell yq e ".$(RELEASE_BRANCH).ova-release-version" $(MAKE_ROOT)/BOTTLEROCKET_RELEASES)
BOTTLEROCKET_RAW_RELEASE_VERSION=$(shell yq e ".$(RELEASE_BRANCH).raw-release-version" $(MAKE_ROOT)/BOTTLEROCKET_RELEASES)

IMAGE_OS?=ubuntu
IMAGE_FORMAT?=raw
BUILD_AMI_TARGETS=build-ami-ubuntu-2004
BUILD_OVA_TARGETS=setup-packer-configs-ova  $(FAKE_UBUNTU_OVA_PATH) $(FINAL_OVA_PATH)  upload-artifacts-ova
BUILD_RAW_TARGETS=release-image-build-on-metal-ubuntu $(FAKE_UBUNTU_RAW_PATH) $(FINAL_RAW_IMAGE_PATH) upload-artifacts-raw
ifneq ($(BOTTLEROCKET_RAW_RELEASE_VERSION),null)
	BUILD_RAW_TARGETS+=download-bottlerocket-raw $(FINAL_BOTTLEROCKET_RAW_PATH) upload-bottlerocket-raw
endif
ifneq ($(BOTTLEROCKET_OVA_RELEASE_VERSION),null)
	BUILD_OVA_TARGETS+=download-bottlerocket-ova $(FINAL_BOTTLEROCKET_OVA_PATH) upload-bottlerocket-ova
endif

BUILD_TARGETS=$(BUILD_RAW_TARGETS) $(BUILD_AMI_TARGETS) $(BUILD_OVA_TARGETS)
ifeq ($(IMAGE_FORMAT),ova)
	ifeq ($(IMAGE_OS),ubuntu)
		S3_TARGET_PREREQUISITES=$(FINAL_OVA_PATH)
		RELEASE_TARGETS=release-ova-$(IMAGE_OS) upload-artifacts-ova
	else ifeq ($(IMAGE_OS),bottlerocket)
		S3_TARGET_PREREQUISITES=$(FINAL_BOTTLEROCKET_OVA_PATH)
		RELEASE_TARGETS=download-bottlerocket-ova upload-bottlerocket-ova
	else ifeq ($(IMAGE_OS),redhat)
		S3_TARGET_PREREQUISITES=$(FINAL_OVA_PATH)
		RELEASE_TARGETS=release-ova-$(IMAGE_OS) upload-artifacts-ova
	else
		S3_TARGET_PREREQUISITES?=
		RELEASE_TARGETS=unsupported-release-target
	endif
else ifeq ($(IMAGE_FORMAT),raw)
	ifeq ($(IMAGE_OS),ubuntu)
		S3_TARGET_PREREQUISITES=$(FINAL_RAW_IMAGE_PATH)
		RELEASE_TARGETS=release-image-build-on-metal-$(IMAGE_OS) upload-artifacts-raw
	else ifeq ($(IMAGE_OS),bottlerocket)
		ifneq ($(RELEASE_BRANCH),1-20)
			S3_TARGET_PREREQUISITES=$(FINAL_BOTTLEROCKET_RAW_PATH)
			RELEASE_TARGETS=download-bottlerocket-raw upload-bottlerocket-raw
		else
			RELEASE_TARGETS=unsupported-release-target
		endif
	else ifeq ($(IMAGE_OS),redhat)
		S3_TARGET_PREREQUISITES=$(FINAL_RAW_IMAGE_PATH)
		RELEASE_TARGETS=release-image-build-on-metal-$(IMAGE_OS) upload-artifacts-raw
	else
		S3_TARGET_PREREQUISITES?=
		RELEASE_TARGETS=unsupported-release-target
	endif
else ifeq ($(IMAGE_FORMAT),cloudstack)
	ifeq ($(IMAGE_OS),redhat)
		S3_TARGET_PREREQUISITES=$(FINAL_CLOUDSTACK_IMAGE_PATH)
		RELEASE_TARGETS=release-image-build-on-metal-$(IMAGE_OS) upload-artifacts-cloudstack
	else
		S3_TARGET_PREREQUISITES?=
		RELEASE_TARGETS=unsupported-release-target
	endif
else ifeq ($(IMAGE_FORMAT),ami)
	ifeq ($(IMAGE_OS),ubuntu)
		RELEASE_TARGETS=release-ami-$(IMAGE_OS)
	else
		RELEASE_TARGETS=unsupported-release-target
	endif
else
	S3_TARGET_PREREQUISITES?=
	RELEASE_TARGETS=unsupported-release-target
endif

BUILDSPECS=projects/$(COMPONENT)/buildspecs/ova.yml projects/$(COMPONENT)/buildspecs/raw.yml
BUILDSPEC_DEPENDS_ON_OVERRIDE=none
BUILDSPEC_VARS_KEYS=IMAGE_OS RELEASE_BRANCH
BUILDSPEC_VARS_VALUES=SUPPORTED_OS_DISTROS SUPPORTED_K8S_VERSIONS
SUPPORTED_OS_DISTROS=bottlerocket

PROJECT_DEPENDENCIES=eksa/kubernetes-sigs/etcdadm eksa/kubernetes-sigs/cri-tools

include $(BASE_DIRECTORY)/Common.mk


export PATH:=$(MAKE_ROOT)/$(IMAGE_BUILDER_DIR)/.local/bin:$(PATH)
export GOVC_INSECURE?=true

# Since we do not build the ova in presubmit but want to validate upload-artifacts behavior
$(FAKE_UBUNTU_OVA_PATH):
	touch $@
	touch $(ARTIFACTS_PATH)/ova/ubuntu/packer.log

$(FAKE_UBUNTU_RAW_PATH):
	@mkdir -p $(@D)
	touch $@

# in the case of BR ova these targets are the same, guard against this to avoid the overwriting target warning
ifneq ($(FINAL_BOTTLEROCKET_OVA_PATH),$(FINAL_OVA_PATH))
$(FINAL_OVA_PATH):
	@mkdir -p $(@D)
	mv *.ova $@
endif

$(FINAL_RAW_IMAGE_PATH):
	@mkdir -p $(@D)
	mv *.gz $@

$(FINAL_CLOUDSTACK_IMAGE_PATH):
	@mkdir -p $(@D)
	mv *.qcow2 $@

$(FINAL_BOTTLEROCKET_OVA_PATH): FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/bottlerocket
$(FINAL_BOTTLEROCKET_OVA_PATH): BOTTLEROCKET_DOWNLOAD_PATH=$(FULL_OUTPUT_DIR)/bottlerocket/downloads
$(FINAL_BOTTLEROCKET_OVA_PATH):
	@mkdir -p $(@D)
	mv $(BOTTLEROCKET_DOWNLOAD_PATH)/ova/*.ova $@

$(FINAL_BOTTLEROCKET_RAW_PATH): FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/bottlerocket
$(FINAL_BOTTLEROCKET_RAW_PATH): BOTTLEROCKET_DOWNLOAD_PATH=$(FULL_OUTPUT_DIR)/bottlerocket/downloads
$(FINAL_BOTTLEROCKET_RAW_PATH):
	@mkdir -p $(@D)
	mv $(BOTTLEROCKET_DOWNLOAD_PATH)/raw/*.gz $@

.PHONY: setup-ami-share
setup-ami-share:
	jq -n '{"ami_users": "$(DEV_ACCOUNTS)", "ami_groups": "", "snapshot_users": "$(DEV_ACCOUNTS)", "snapshot_groups": ""}' > $(PACKER_AMI_SHARE_FILE)

$(BOTTLEROCKET_SETUP_TARGET): FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/bottlerocket
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
	build/setup_packer_configs.sh $(RELEASE_BRANCH) $* $(IMAGE_OS) $(ARTIFACTS_BUCKET) $(ARTIFACTS_PATH)/$*/$(IMAGE_OS) $(ADDITIONAL_PAUSE_$(RELEASE_BRANCH)_FROM) $(LATEST) $(IMAGE_BUILDER_DIR)

.PHONY: build-ami-ubuntu-2004
build-ami-ubuntu-2004: MAKEFLAGS=
build-ami-ubuntu-2004: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ami/$(IMAGE_OS)
build-ami-ubuntu-2004: PACKER_TYPE_VAR_FILES=$(PACKER_AMI_VAR_FILES)
build-ami-ubuntu-2004: setup-ami-share deps-ami setup-packer-configs-ami
	PACKER_VAR_FILES="$(PACKER_VAR_FILES)" $(MAKE) -C $(IMAGE_BUILDER_DIR) validate-ami-ubuntu-2004

.PHONY: release-ami-%
release-ami-%: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ami/$(IMAGE_OS)
release-ami-%: export MANIFEST_OUTPUT=$(FULL_OUTPUT_DIR)/manifest.json
release-ami-%: EXPORT_AMI_BUCKET?=$(ARTIFACTS_BUCKET)
release-ami-%: AMI_S3_DST=$(EXPORT_AMI_BUCKET)/$(ARTIFACTS_UPLOAD_PATH)/ami/$(IMAGE_OS)
release-ami-%: EXPORT_AMI_DST=$(AMI_S3_DST)/$(GIT_HASH)
release-ami-%: LATEST_AMI_S3_URL=$(AMI_S3_DST)/$(LATEST)/ubuntu.raw
release-ami-%: setup-ami-share
	build/build_image.sh $* $(RELEASE_BRANCH) $(IMAGE_FORMAT) $(ARTIFACTS_BUCKET) $(LATEST)
	build/export-ami-to-s3.sh $(RELEASE_BRANCH) $(MANIFEST_OUTPUT) raw $(EXPORT_AMI_DST) $(LATEST_AMI_S3_URL)

.PHONY: release-ova-%
release-ova-%:
	build/build_image.sh $* $(RELEASE_BRANCH) $(IMAGE_FORMAT) $(ARTIFACTS_BUCKET) $(LATEST)

.PHONY: local-build-ova-%
local-build-ova-%: MAKEFLAGS=
local-build-ova-%: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ova/$(IMAGE_OS)
local-build-ova-%: PACKER_TYPE_VAR_FILES=$(PACKER_OVA_VAR_FILES)
local-build-ova-%: deps-ova setup-packer-configs-ova
	PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(ARTIFACTS_PATH)/ova/$(IMAGE_OS)/packer.log PACKER_VAR_FILES="$(PACKER_VAR_FILES)" \
		OVF_CUSTOM_PROPERTIES="$(FULL_OUTPUT_DIR)/config/ovf_custom_properties.json" \
		$(MAKE) -C $(IMAGE_BUILDER_DIR) build-node-ova-vsphere-$*

.PHONY: download-bottlerocket-%
download-bottlerocket-%: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/bottlerocket
download-bottlerocket-%: BOTTLEROCKET_DOWNLOAD_PATH=$(FULL_OUTPUT_DIR)/bottlerocket/downloads
download-bottlerocket-%: $(BOTTLEROCKET_SETUP_TARGET)
	build/get_bottlerocket_artifacts.sh $(RELEASE_BRANCH) $* $(BOTTLEROCKET_DOWNLOAD_PATH) $(PROJECT_PATH)/$(RELEASE_BRANCH) $(LATEST_TAG)

.PHONY: upload-bottlerocket-%
upload-bottlerocket-%: IMAGE_OS=bottlerocket
upload-bottlerocket-%: ARTIFACTS_PATH_IMAGE_FORMAT=$(ARTIFACTS_PATH)/$*/bottlerocket
upload-bottlerocket-%: ARTIFACTS_UPLOAD_PATH_IMAGE_FORMAT=$(ARTIFACTS_UPLOAD_PATH)/$*/bottlerocket
upload-bottlerocket-%:
	build/upload_bottlerocket_artifacts.sh $(RELEASE_BRANCH) $* $(ARTIFACTS_PATH_IMAGE_FORMAT) $(ARTIFACTS_UPLOAD_PATH_IMAGE_FORMAT)

.PHONY: release-image-build-on-metal-%
release-image-build-on-metal-%: $(GIT_PATCH_TARGET)
	build/build_image_on_metal.sh $(BASE_DIRECTORY) $(PROJECT_PATH) $(RELEASE_BRANCH) $(RAW_IMAGE_BUILD_AMI) $(RAW_IMAGE_BUILD_INSTANCE_TYPE) $(RAW_IMAGE_BUILD_KEY_NAME) $* $(IMAGE_FORMAT) $(LATEST)

.PHONY: build-qemu-rhel-local
build-qemu-rhel-local: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/qemu/$(IMAGE_OS)
build-qemu-rhel-local: deps-qemu setup-packer-configs-qemu
	build/build_qemu_rhel_local.sh $(MAKE_ROOT) $(BASE_IMAGE) $(IMAGE_BUILDER_DIR) "$(PACKER_VAR_FILES)"

.PHONY: local-build-ami-%
local-build-ami-%: MAKEFLAGS=
local-build-ami-%: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/ami/$(IMAGE_OS)
local-build-ami-%: PACKER_TYPE_VAR_FILES=$(PACKER_AMI_VAR_FILES)
local-build-ami-%: deps-ami setup-packer-configs-ami
	PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(ARTIFACTS_PATH)/ami/$(IMAGE_OS)/packer.log PACKER_VAR_FILES="$(PACKER_VAR_FILES)" \
		$(MAKE) -C $(IMAGE_BUILDER_DIR) build-ami-$*

.PHONY: local-build-raw-%
local-build-raw-%: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/raw/$(IMAGE_OS)
local-build-raw-%: deps-raw setup-packer-configs-raw
	PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(ARTIFACTS_PATH)/raw/$(IMAGE_OS)/packer.log PACKER_VAR_FILES="$(PACKER_VAR_FILES)" \
		$(MAKE) -C $(IMAGE_BUILDER_DIR) build-raw-$*

.PHONY: local-build-qemu-%
local-build-qemu-%: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/qemu/$(IMAGE_OS)
local-build-qemu-%: deps-qemu setup-packer-configs-qemu
	PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(ARTIFACTS_PATH)/qemu/$(IMAGE_OS)/packer.log PACKER_VAR_FILES="$(PACKER_VAR_FILES)" \
		$(MAKE) -C $(IMAGE_BUILDER_DIR) build-qemu-$*

.PHONY: local-build-nutanix-ubuntu-2004
local-build-nutanix-ubuntu-2004: FULL_OUTPUT_DIR=$(MAKE_ROOT)/$(OUTPUT_DIR)/nutanix/$(IMAGE_OS)
local-build-nutanix-ubuntu-2004: deps-nutanix setup-packer-configs-nutanix
	PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(ARTIFACTS_PATH)/nutanix/$(IMAGE_OS)/packer.log PACKER_VAR_FILES="$(PACKER_VAR_FILES)" \
		$(MAKE) -C $(IMAGE_BUILDER_DIR) build-nutanix-ubuntu-2004

.PHONY: validate-ubuntu-2004
validate-ubuntu-2004: check-env-validation $(GIT_PATCH_TARGET) setup-packer-configs-ova
	sh validate/validate.sh $(FULL_OUTPUT_DIR)/config/validate_goss_inline_vars.json

.PHONY: check-env-validation
check-env-validation:
	ifndef TEMPLATE
		$(error Environment var TEMPLATE not set. Example: TEMPLATE=<ubuntu.livecd> make ...)
	endif

.PHONY: s3-artifacts-%
s3-artifacts-%: ARTIFACTS_PATH_IMAGE_FORMAT=$(ARTIFACTS_PATH)/$*/$(IMAGE_OS)
s3-artifacts-%: $(S3_TARGET_PREREQUISITES)
	$(MAKE) -C $(MAKE_ROOT) s3-artifacts ARTIFACTS_PATH=$(ARTIFACTS_PATH_IMAGE_FORMAT) IMAGE_FORMAT=$* IMAGE_OS=$(IMAGE_OS)

.PHONY: upload-artifacts-%
upload-artifacts-%: ARTIFACTS_PATH_IMAGE_FORMAT=$(ARTIFACTS_PATH)/$*/$(IMAGE_OS)
upload-artifacts-%: ARTIFACTS_UPLOAD_PATH_IMAGE_FORMAT=$(ARTIFACTS_UPLOAD_PATH)/$*/$(IMAGE_OS)
upload-artifacts-%: s3-artifacts-%
	$(MAKE) -C $(MAKE_ROOT) upload-artifacts ARTIFACTS_PATH=$(ARTIFACTS_PATH_IMAGE_FORMAT) ARTIFACTS_UPLOAD_PATH=$(ARTIFACTS_UPLOAD_PATH_IMAGE_FORMAT) IMAGE_FORMAT=$*


.PHONY: unsupported-release-target
unsupported-release-target:
	@echo "Unsupported combination of image format and image os chosen. Build will quit now..."

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

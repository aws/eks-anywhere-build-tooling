BASE_DIRECTORY:=$(abspath ../../../)
MAKE_ROOT=$(BASE_DIRECTORY)/projects/kubernetes-sigs/image-builder
GIT_TAG=$(shell cat GIT_TAG)

REPO=image-builder
REPO_OWNER=kubernetes-sigs

# no container images for this project
IMAGE_NAMES=

IMAGE_BUILDER_DIR=$(REPO)/images/capi
RAW_IMAGE_BUILD_AMI?=ami-095413544ce52437d
RAW_IMAGE_BUILD_INSTANCE_TYPE?=t3.micro
RAW_IMAGE_BUILD_KEY_NAME?=raw-image-build

# in the case of ubuntu ova builds we support either default bios or efi mode, for efi we put in a different output folder: $(ARTIFACTS_PATH)/ubuntu/2004/ova/efi/ubuntu.ova
# Note: raw ubuntu builds are always efi and are not in a seperate sub folder
FIRMWARE_PATH_SUFFIX=$(if $(and $(IMAGE_OS_FIRMWARE),$(filter $(IMAGE_OS)-$(IMAGE_FORMAT),$(ALLOWED_EFI_OS_FORMATS))),/efi,)

FULL_OUTPUT_DIR?=$(MAKE_ROOT)/$(OUTPUT_DIR)/$(IMAGE_FORMAT)/$(IMAGE_OS_DIR)
BOTTLEROCKET_DOWNLOAD_PATH?=$(FULL_OUTPUT_DIR)/downloads
ARTIFACTS_UPLOAD_PATH_IMAGE_FORMAT=$(ARTIFACTS_UPLOAD_PATH)/$(IMAGE_FORMAT)/$(IMAGE_OS_DIR)$(FIRMWARE_PATH_SUFFIX)
ARTIFACTS_PATH_IMAGE_FORMAT=$(ARTIFACTS_PATH)/$(IMAGE_OS_DIR)/$(IMAGE_FORMAT)$(FIRMWARE_PATH_SUFFIX)

RHSM_USERNAME?=
RHSM_PASSWORD?=
export BUILDER_ROOT=$(MAKE_ROOT)
export RHSM_USER=$(RHSM_USERNAME)
export RHSM_PASS=$(RHSM_PASSWORD)

REDHAT_CONFIG_TARGET=redhat-config.json
IF_REDHAT_CONFIG_TARGET=$(if $(filter redhat,$(IMAGE_OS)),$(REDHAT_CONFIG_TARGET),)

VSPHERE_CONNECTION_DATA?={}
NUTANIX_CONNECTION_DATA?={}
# Aws accounts to share built AMI with
DEV_ACCOUNTS?=

PACKER_AMI_VAR_FILES=$(MAKE_ROOT)/packer/ami/ami.json
PACKER_AMI_SHARE_FILE=$(MAKE_ROOT)/packer/ami/share-ami.json
PACKER_AMI_VAR_FILES+=$(PACKER_AMI_SHARE_FILE)

PACKER_OVA_CONF_FILE=$(MAKE_ROOT)/packer/ova/vsphere.json
PACKER_OVA_VAR_FILES=$(MAKE_ROOT)/packer/ova/vsphere-library.json
PACKER_OVA_VAR_FILES+=$(PACKER_OVA_CONF_FILE)

PACKER_RAW_VAR_FILES?=
PACKER_CLOUDSTACK_VAR_FILES?=

PACKER_NUTANIX_VAR_FILES?=$(MAKE_ROOT)/packer/nutanix/nutanix.json

PACKER_TYPE_VAR_FILES?=$(PACKER_$(call TO_UPPER,$(IMAGE_FORMAT))_VAR_FILES)
PACKER_ADDITIONAL_FILES_VAR_FILES?=
PACKER_VAR_FILES=$(FULL_OUTPUT_DIR)/config/kubernetes.json $(FULL_OUTPUT_DIR)/config/common.json $(FULL_OUTPUT_DIR)/config/cni.json $(FULL_OUTPUT_DIR)/config/additional_components.json $(PACKER_TYPE_VAR_FILES) $(PACKER_ADDITIONAL_FILES_VAR_FILES)

# Kubeadm hardcodes the version of the pause image it pulls, even tho containerd config
# overrides which image is actually used
# Setting up this mapping so we can build the additional image into the image
# The versions are dependent on the tag set for the Kubernetes version we are consuming vs what is hardcoded in the
# kubeadm constants for that specific version:
# Example for v1.22.4 on why we rely on the pause tag from v1.21.5, which is the previous version that we support:
# Kubernetes pause tag in v1.22.4: https://github.com/kubernetes/kubernetes/blob/v1.22.4/build/pause/Makefile#L20
# Kubernetes pause tag in v1.21.5: https://github.com/kubernetes/kubernetes/blob/v1.21.5/build/pause/Makefile#L20
# Kubeadm constants: https://github.com/kubernetes/kubernetes/blob/v1.22.4/cmd/kubeadm/app/constants/constants.go#L397
ADDITIONAL_PAUSE_1-23_FROM=1-22
ADDITIONAL_PAUSE_1-24_FROM=1-23
ADDITIONAL_PAUSE_1-25_FROM=1-24
ADDITIONAL_PAUSE_1-26_FROM=1-26
ADDITIONAL_PAUSE_1-27_FROM=1-27
ADDITIONAL_PAUSE_1-28_FROM=1-28
ADDITIONAL_PAUSE_1-29_FROM=1-28
ADDITIONAL_PAUSE_1-30_FROM=1-28
ADDITIONAL_PAUSE_1-31_FROM=1-31
ADDITIONAL_PAUSE_1-32_FROM=1-32
ADDITIONAL_PAUSE_1-33_FROM=1-33
ADDITIONAL_PAUSE_1-34_FROM=1-34

HAS_RELEASE_BRANCHES=true
HAS_LICENSES=false
SIMPLE_CREATE_TARBALLS=false

# if running locally and do not want to try and create ec2 instance
SKIP_METAL_INSTANCE_TEST?=false

# to support a no op attribution target
TARGETS_ALLOWED_WITH_NO_RELEASE_BRANCH=binaries checksums attribution update-dependency-versions validate-dependency-versions-raw validate-dependency-versions-cloudstack

IMAGE_OS_DIR?=$(IMAGE_OS)/$(IMAGE_OS_VERSION)
IMAGE_OS_WITH_VER?=$(IMAGE_OS)-$(IMAGE_OS_VERSION)

OVA_FORMAT_EXT=ova
RAW_FORMAT_EXT=gz
AMI_FORMAT_EXT=gz
CLOUDSTACK_FORMAT_EXT=qcow2
NUTANIX_FORMAT_EXT=img

# $1- image_format
# $2 - image_os
EXT_FROM_FORMAT=$(if $(and $(filter bottlerocket,$2),$(filter-out ova,$1)),img.,)$($(call TO_UPPER,$1)_FORMAT_EXT)

# $1 - target to convert
# upstream targets for ova are named slightly different
# we have tried to standardize on redhat vs rhel, but upstream uses rhel
# we are explcitly using cloudstack but upstream this is just a qemu build
# for raw ubuntu builds we always build the efi variant
# ex: validate-raw-ubuntu-2004, build-raw-ubuntu-2004
CONVERT_TO_UPSTREAM_TARGET=$(strip \
	$(eval _NEW_TARGET:=$1) \
	$(eval _NEW_TARGET:=$(subst validate-ova-,validate-node-ova-local-,$(_NEW_TARGET))) \
	$(eval _NEW_TARGET:=$(subst build-clone-ova-,build-node-ova-vsphere-clone-,$(_NEW_TARGET))) \
	$(eval _NEW_TARGET:=$(subst build-ova-,build-node-ova-vsphere-,$(_NEW_TARGET)))  \
	$(eval _NEW_TARGET:=$(subst -redhat-,-rhel-,$(_NEW_TARGET))) \
	$(eval _NEW_TARGET:=$(subst -cloudstack,-qemu,$(_NEW_TARGET))) \
	$(eval _NEW_TARGET:=$(if $(findstring raw-ubuntu,$(_NEW_TARGET)),$(_NEW_TARGET)-efi,$(_NEW_TARGET))) \
	$(eval _NEW_TARGET:=$(if $(findstring raw-rhel-9,$(_NEW_TARGET)),$(_NEW_TARGET)-efi,$(_NEW_TARGET))) \
	$(_NEW_TARGET) \
)

# ex: $(ARTIFACTS_PATH)/ubuntu/2004/ova/ubuntu.ova, $(ARTIFACTS_PATH)/bottlerocket/1/raw/bottlerocket.img.gz
FINAL_IMAGE_DIR=$(ARTIFACTS_PATH)/$(IMAGE_OS_DIR)/$(IMAGE_FORMAT)$(FIRMWARE_PATH_SUFFIX)
FINAL_IMAGE_FILE_NAME=$(IMAGE_OS).$(call EXT_FROM_FORMAT,$(IMAGE_FORMAT),$(IMAGE_OS))
FINAL_IMAGE_PATH=$(FINAL_IMAGE_DIR)/$(FINAL_IMAGE_FILE_NAME)

BOTTLEROCKET_AMI_RELEASE_VERSION=$(and $(RELEASE_BRANCH),$(shell yq e ".$(RELEASE_BRANCH).ami-release-version" $(MAKE_ROOT)/BOTTLEROCKET_RELEASES))
BOTTLEROCKET_OVA_RELEASE_VERSION=$(and $(RELEASE_BRANCH),$(shell yq e ".$(RELEASE_BRANCH).ova-release-version" $(MAKE_ROOT)/BOTTLEROCKET_RELEASES))
BOTTLEROCKET_RAW_RELEASE_VERSION=$(and $(RELEASE_BRANCH),$(shell yq e ".$(RELEASE_BRANCH).raw-release-version" $(MAKE_ROOT)/BOTTLEROCKET_RELEASES))

NOT_SUPPORTED_RELEASE_BRANCH_CONFIGURATION=$(if $(or $(filter-out bottlerocket,$(IMAGE_OS)),$(filter-out null,$(BOTTLEROCKET_$(call TO_UPPER,$(IMAGE_FORMAT))_RELEASE_VERSION))),false,true)
SKIPPED_K8S_VERSIONS=$(shell yq e 'with_entries(select(.value.$(IMAGE_FORMAT)-release-version == null)) | keys | join(" ")' BOTTLEROCKET_RELEASES)

IMAGE_OS?=ubuntu
IMAGE_OS_VERSION?=$(DEFAULT_$(call TO_UPPER,$(IMAGE_OS))_VERSION)
IMAGE_FORMAT?=ova
IMAGE_FORMATS_AVAIL=ami cloudstack nutanix ova raw
IMAGE_OS_FIRMWARE?=

###################### BUILD TARGETS ##############################
# build targets only run validation, they do not actually do an image build
# snow node images
BUILD_AMI_TARGETS=build-ami-ubuntu-2004 build-ami-ubuntu-2204

# $1 - image-format
# $2 - image-os-with-ver
# ex: ova, ubuntu-2204
FAKE_IMAGE_PATH_FROM_OS_WITH_VER=fake-$2-$1.$(call EXT_FROM_FORMAT,$(1),$(word 1,$(subst -, ,$2)))

# vsphere node images
# packer/ova/vsphere.json needs to exist for the validation, but can be empty
BUILD_OVA_TARGETS=packer/ova/vsphere.json $(foreach os,ubuntu-2004 ubuntu-2204 redhat-8,build-ova-$(os) $(call FAKE_IMAGE_PATH_FROM_OS_WITH_VER,ova,$(os)) upload-artifacts-ova-$(os))
BUILD_OVA_TARGETS+=build-ova-ubuntu-2004-efi fake-ubuntu-2004-ova-efi.ova upload-artifacts-ova-ubuntu-2004-efi build-ova-ubuntu-2204-efi fake-ubuntu-2204-ova-efi.ova upload-artifacts-ova-ubuntu-2204-efi
# vsphere clone node images
BUILD_OVA_CLONE_TARGETS=$(foreach os,ubuntu-2004 ubuntu-2204 redhat-8,build-clone-ova-$(os) $(call FAKE_IMAGE_PATH_FROM_OS_WITH_VER,ova,$(os)) upload-artifacts-clone-ova-$(os))
BUILD_OVA_CLONE_TARGETS+=build-clone-ova-ubuntu-2004-efi fake-ubuntu-2004-ova-efi.ova upload-artifacts-clone-ova-ubuntu-2004-efi build-clone-ova-ubuntu-2204-efi fake-ubuntu-2204-ova-efi.ova upload-artifacts-clone-ova-ubuntu-2204-efi build-clone-ova-ubuntu-2404-efi build-ova-ubuntu-2404-efi fake-ubuntu-2404-ova-efi.ova upload-artifacts-ova-ubuntu-2404-efi

# only running release-raw-ubuntu-2004 on presubmit here since we are just validating the instance creation and do not need to run for every OS
# tinkerbell/bare metal node images
BUILD_RAW_TARGETS=$(foreach os,ubuntu-2004 ubuntu-2204 ubuntu-2404 redhat-8 redhat-9,build-raw-$(os) $(call FAKE_IMAGE_PATH_FROM_OS_WITH_VER,raw,$(os)) upload-artifacts-raw-$(os)) $(if $(filter true,$(SKIP_METAL_INSTANCE_TEST)),,metal-instance-test)

# cloud stack node images
BUILD_CLOUDSTACK_TARGETS=$(foreach os, redhat-8 redhat-9,build-cloudstack-$(os) $(call FAKE_IMAGE_PATH_FROM_OS_WITH_VER,cloudstack,$(os)) upload-artifacts-cloudstack-$(os))

# nutanix node images
BUILD_NUTANIX_TARGETS=packer/nutanix/nutanix.json $(foreach os,ubuntu-2004 ubuntu-2204 ubuntu-2404 redhat-8 redhat-9,build-nutanix-$(os) $(call FAKE_IMAGE_PATH_FROM_OS_WITH_VER,nutanix,$(os)) upload-artifacts-nutanix-$(os))

ifneq ($(BOTTLEROCKET_AMI_RELEASE_VERSION),null)
	BUILD_AMI_TARGETS+=upload-bottlerocket-1-ami
endif
ifneq ($(BOTTLEROCKET_RAW_RELEASE_VERSION),null)
	BUILD_RAW_TARGETS+=upload-bottlerocket-1-raw
endif
ifneq ($(BOTTLEROCKET_OVA_RELEASE_VERSION),null)
	BUILD_OVA_TARGETS+=upload-bottlerocket-1-ova
endif

BUILD_TARGETS=$(if $(filter presubmit,$(JOB_TYPE)),update-dependency-versions,) $(BUILD_RAW_TARGETS) $(BUILD_AMI_TARGETS) $(BUILD_OVA_TARGETS) $(BUILD_CLOUDSTACK_TARGETS) $(BUILD_NUTANIX_TARGETS)
#####################################################################

###################### RELEASE TARGETS ##############################

# These targets do trigger real image builds or an error if the combination of OS to FORMAT is invalid
# The `release` target is called via codebuild with env vars (IMAGE_FORMAT, IMAGE_OS) to define what to actually call
# vs calling the targets specifically, this is so we can use on buildspec
# This blocks uses these env vars to determine what target to call

ALLOWED_OVA_OS=ubuntu bottlerocket redhat
ALLOWED_RAW_OS=ubuntu bottlerocket redhat
ALLOWED_CLOUDSTACK_OS=redhat
ALLOWED_AMI_OS=ubuntu bottlerocket
ALLOWED_NUTANIX_OS=ubuntu redhat

ALLOWED_UBUNTU_VERSIONS=2004 2204 2404
DEFAULT_UBUNTU_VERSION=2004

ALLOWED_REDHAT_VERSIONS=8 9
DEFAULT_REDHAT_VERSION=8

ALLOWED_BOTTLEROCKET_VERSIONS=1
DEFAULT_BOTTLEROCKET_VERSION=1

ALLOWED_EFI_OS_FORMATS=ubuntu-ova redhat-raw redhat-ova

# run image os/format validation before running release target(s)
RELEASE_TARGETS=validate-supported-image-all
# For bottlerocket the upload target has a prereq of the download target
# since in "build" we only run validation and fake the actual built image file
# the other OS builds needs to run the release- target before running upload-artifacts
ifeq ($(IMAGE_OS),bottlerocket)
	RELEASE_TARGETS+= upload-bottlerocket-1-$(IMAGE_FORMAT)
else
# for ami builds they are automatically uploaded during the ami export to s3 process, skip the upload target
	RELEASE_TARGETS+= release-$(IMAGE_FORMAT)-$(IMAGE_OS_WITH_VER) $(if $(filter ami,$(IMAGE_FORMAT)),,upload-artifacts-$(IMAGE_FORMAT)-$(IMAGE_OS_WITH_VER))
endif

#####################################################################

BUILDSPECS=projects/$(COMPONENT)/buildspecs/br.yml
# in the staging buildspec we only run the br "builds" which do not rely on etcdadm/crictl
BUILDSPEC_DEPENDS_ON_OVERRIDE=none
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_MEDIUM
SUPPORTED_OS_DISTROS=bottlerocket

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true

PROJECT_DEPENDENCIES=eksa/kubernetes-sigs/etcdadm eksa/kubernetes-sigs/cri-tools

DEBUG?=0

SKIP_ON_RELEASE_BRANCH=$(if $(filter bottlerocket,$(IMAGE_OS)),false,true)

# our file builds/downloads are technically intermediate targets which would get deleted when
# make finishes, this avoids that
ALL_FINAL_FILE_PATHS=\
	$(foreach format,$(IMAGE_FORMATS_AVAIL),\
		$(foreach os,$(ALLOWED_$(call TO_UPPER,$(format))_OS),\
			$(foreach ver,$(ALLOWED_$(call TO_UPPER,$(os))_VERSIONS),\
				$(foreach firmware,$(if $(filter $(os)-$(format),$(ALLOWED_EFI_OS_FORMATS)),/ /efi/,/),\
					$(ARTIFACTS_PATH)/$(os)/$(ver)/$(format)$(firmware)$(os).$(call EXT_FROM_FORMAT,$(format),$(os))\
				)\
			)\
		)\
	)

SPECIAL_TARGET_SECONDARY=$(strip $(ALL_FINAL_FILE_PATHS))

EXPECTED_FILES_PATH=expected-artifacts/expected_artifacts_$(IMAGE_FORMAT)$(if $(filter bottlerocket,$(IMAGE_OS)),_bottlerocket,)


include $(BASE_DIRECTORY)/Common.mk


export PATH:=$(MAKE_ROOT)/$(IMAGE_BUILDER_DIR)/.local/bin:$(PATH)

######################## DEPS #################################
.PHONY: setup-ami-share
setup-ami-share: | $$(ENABLE_LOGGING)
	@jq -n '{"ami_users": "$(DEV_ACCOUNTS)", "ami_groups": "", "snapshot_users": "$(DEV_ACCOUNTS)", "snapshot_groups": ""}' > $(PACKER_AMI_SHARE_FILE)

.PHONY: setup-packer-configs-%
setup-packer-configs-%: $(GIT_PATCH_TARGET) | ensure-jq ensure-yq $$(ENABLE_LOGGING)
	@build/setup_packer_configs.sh $(RELEASE_BRANCH) $(IMAGE_FORMAT) $(IMAGE_OS) $(ARTIFACTS_BUCKET) $(FINAL_IMAGE_DIR) $(ADDITIONAL_PAUSE_$(RELEASE_BRANCH)_FROM) $(LATEST) $(IMAGE_BUILDER_DIR) $(IMAGE_OS_DIR) $(IMAGE_OS_VERSION)

REDHAT_8_SEMVER=8.10
REDHAT_9_SEMVER=9.4
REDHAT_8_ISO_CHECKSUM=9b3c8e31bc2cdd2de9cf96abb3726347f5840ff3b176270647b3e66639af291b
REDHAT_9_ISO_CHECKSUM=398561d7b66f1a4bf23664f4aa8f2cfbb3641aa2f01a320068e86bd1fc0e9076
REDHAT_8_KMOD_RPM=kmod-megaraid_sas-07.727.03.00-1.el8_10.elrepo.x86_64.rpm
REDHAT_9_KMOD_RPM=kmod-megaraid_sas-07.727.03.00-1.el9_4.elrepo.x86_64.rpm

$(REDHAT_CONFIG_TARGET):
	@if [[ "$(IMAGE_FORMAT)" == "nutanix" ]]; then \
		jq --null-input \
			--arg rhel_username "$(RHSM_USERNAME)" \
			--arg rhel_password "$(RHSM_PASSWORD)" \
			--arg image_url "$$(aws s3 presign redhat-iso-pdx/$(REDHAT_$(IMAGE_OS_VERSION)_SEMVER)/rhel-$(REDHAT_$(IMAGE_OS_VERSION)_SEMVER)-x86_64-kvm.qcow2)" \
			'{"rhel_username": $$rhel_username, "rhel_password": $$rhel_password, "image_url": $$image_url}' > $@; \
	else \
		jq --null-input \
			--arg rhel_username "$(RHSM_USERNAME)" \
			--arg rhel_password "$(RHSM_PASSWORD)" \
			--arg iso_url "$$(aws s3 presign redhat-iso-pdx/$(REDHAT_$(IMAGE_OS_VERSION)_SEMVER)/rhel-$(REDHAT_$(IMAGE_OS_VERSION)_SEMVER)-x86_64-dvd.iso)" \
			--arg extra_rpms "$(if $(filter raw,$(IMAGE_FORMAT)),$$(aws s3 presign redhat-iso-pdx/$(REDHAT_$(IMAGE_OS_VERSION)_SEMVER)/rpms/$(REDHAT_$(IMAGE_OS_VERSION)_KMOD_RPM)),)" \
			--arg iso_checksum_type "sha256" \
			--arg iso_checksum "$(REDHAT_$(IMAGE_OS_VERSION)_ISO_CHECKSUM)" \
			'{"rhel_username": $$rhel_username, "rhel_password": $$rhel_password, "iso_url": $$iso_url, "iso_checksum_type": $$iso_checksum_type, "iso_checksum": $$iso_checksum, "extra_rpms": $$extra_rpms}' > $@; \
	fi
##############################################################

########################### FAKE TARGETS ###############################
# Since we do not build the ova in presubmit but want to validate upload-artifacts behavior
# ex: fake-ubuntu-2004-raw.gz,fake-ubuntu-2004-ova-efi.ova
fake-%: PACKER_LOG_PATH=$(ARTIFACTS_PATH)/$(subst -,/,$(basename $(*F)))/packer.log
fake-%: IMAGE_FORMAT=$(word 3,$(subst -, ,$(word 1,$(subst ., ,$*))))
fake-%: | $$(ENABLE_LOGGING)
	@touch $(PACKER_LOG_PATH)
	@if [[ "$(IMAGE_FORMAT)" == "ova" ]]; then \
		build/make_fake_image.sh $(@) $(FULL_OUTPUT_DIR)/config/kubernetes.json; \
	elif [[ "$(IMAGE_FORMAT)" != "nutanix" ]]; then \
		touch $(@); \
	fi

packer/ova/vsphere.json:
	echo "{}" > $@

packer/nutanix/nutanix.json:
	echo "{}" > $@

########################################################################

######################## BOTTLEROCKET #######################
$(ARTIFACTS_PATH)/bottlerocket/% upload-bottlerocket-%: IMAGE_OS=bottlerocket

#$ ex: $(ARTIFACTS_PATH)/bottlerocket/1/raw/bottlerocket.img.gz, $(ARTIFACTS_PATH)/bottlerocket/1/ova/bottlerocket.ova
$(ARTIFACTS_PATH)/bottlerocket/%: IMAGE_OS_VERSION=$(word 1,$(subst /, ,$(*D)))
$(ARTIFACTS_PATH)/bottlerocket/%: IMAGE_FORMAT=$(word 2,$(subst /, ,$(*D)))
$(ARTIFACTS_PATH)/bottlerocket/%: export BOTTLEROCKET_ROOT_JSON_PATH=$(BOTTLEROCKET_DOWNLOAD_PATH)/root.json
$(ARTIFACTS_PATH)/bottlerocket/%: | ensure-yq ensure-tuftool ensure-lz4
	@echo -e $(call TARGET_START_LOG)
	mkdir -p $(@D) $(BOTTLEROCKET_DOWNLOAD_PATH)
	# This configuration supports local installations and checksum validations of root.json file
	envsubst '$$BOTTLEROCKET_ROOT_JSON_PATH' \
		< $(MAKE_ROOT)/bottlerocket-root-json-checksum \
		> $(BOTTLEROCKET_DOWNLOAD_PATH)/bottlerocket-root-json-checksum
	curl https://cache.bottlerocket.aws/root.json -o $$BOTTLEROCKET_ROOT_JSON_PATH
	sha512sum -c $(BOTTLEROCKET_DOWNLOAD_PATH)/bottlerocket-root-json-checksum
	build/get_bottlerocket_artifacts.sh $(RELEASE_BRANCH) $(IMAGE_FORMAT) $(BOTTLEROCKET_DOWNLOAD_PATH) $(PROJECT_PATH)/$(RELEASE_BRANCH) $(LATEST_TAG)
	mv $(BOTTLEROCKET_DOWNLOAD_PATH)/$(IMAGE_FORMAT)/*$(suffix $(@F)) $@	
	echo $$(source $(BUILD_LIB)/common.sh && build::eksd_releases::get_release_yaml_url $(RELEASE_BRANCH)) > $(@D)/EKSD_MANIFEST_URL
	@echo -e $(call TARGET_END_LOG)

# ex: upload-bottlerocket-1-ova, upload-bottlerocket-1-raw
upload-bottlerocket-%: IMAGE_OS_VERSION=$(word 1,$(subst -, ,$*))
upload-bottlerocket-%: IMAGE_FORMAT=$(word 2,$(subst -, ,$*))
upload-bottlerocket-%: $$(FINAL_IMAGE_PATH) | $$(ENABLE_LOGGING)
	@build/upload_bottlerocket_artifacts.sh $(RELEASE_BRANCH) $(IMAGE_FORMAT) $(ARTIFACTS_PATH_IMAGE_FORMAT) $(ARTIFACTS_UPLOAD_PATH_IMAGE_FORMAT) $(IMAGE_OS_VERSION)
###############################################################

########################### FINAL TARGETS ##############################
# Moves packer produced file which is in the repo root to final location in artifacts_path
# $(ARTIFACTS_PATH)/ubuntu/2004/raw/ubuntu.gz
$(ARTIFACTS_PATH)/%: IMAGE_FORMAT=$(word 3,$(subst /, ,$*))
$(ARTIFACTS_PATH)/%:
	@echo -e $(call TARGET_START_LOG)
	@if [[ "$(IMAGE_FORMAT)" != "nutanix" ]]; then \
		mkdir -p $(@D); \
		mv *$(suffix $(@F)) $(@); \
	fi
	@echo -e $(call TARGET_END_LOG)

#######################################################################

######################## VALIDATE ###########################
# ex: build-ova-ubuntu-2004, build-raw-redhat-8, build-ova-ubuntu-2004-efi
build-%: MAKEFLAGS=
build-%: IMAGE_FORMAT=$(word 1,$(subst -, ,$*))
build-%: IMAGE_OS=$(word 2,$(subst -, ,$*))
build-%: IMAGE_OS_VERSION=$(word 3,$(subst -, ,$*))
build-%: IMAGE_OS_FIRMWARE=$(word 4,$(subst -, ,$*))
build-%: validate-supported-image-% setup-packer-configs-% $(if $(filter ami,$(IMAGE_FORMAT)),setup-ami-share,) | $$(ENABLE_LOGGING)
	@PACKER_VAR_FILES="$(PACKER_VAR_FILES)" $(MAKE) -C $(IMAGE_BUILDER_DIR) $(call CONVERT_TO_UPSTREAM_TARGET,validate-$(IMAGE_FORMAT)-$(IMAGE_OS_WITH_VER)$(if $(IMAGE_OS_FIRMWARE),-$(IMAGE_OS_FIRMWARE),))

# ex: build-clone-ova-ubuntu-2004, build-clone-ova-redhat-8, build-clone-ova-ubuntu-2004-efi
build-clone-%: MAKEFLAGS=
build-clone-%: IMAGE_FORMAT=$(word 1,$(subst -, ,$*))
build-clone-%: IMAGE_OS=$(word 2,$(subst -, ,$*))
build-clone-%: IMAGE_OS_VERSION=$(word 3,$(subst -, ,$*))
build-clone-%: IMAGE_OS_FIRMWARE=$(word 4,$(subst -, ,$*))
build-clone-%: validate-supported-image-% setup-packer-configs-% $(if $(filter ami,$(IMAGE_FORMAT)),setup-ami-share,) | $$(ENABLE_LOGGING)
	@PACKER_VAR_FILES="$(PACKER_VAR_FILES)" $(MAKE) -C $(IMAGE_BUILDER_DIR) $(call CONVERT_TO_UPSTREAM_TARGET,validate-$(IMAGE_FORMAT)-$(IMAGE_OS_WITH_VER)$(if $(IMAGE_OS_FIRMWARE),-$(IMAGE_OS_FIRMWARE),))

validate-ova-kubelet:
	@echo -e $(call TARGET_START_LOG)
	build/validate_ova_kubelet.sh $(ARTIFACTS_PATH)/$(FINAL_IMAGE_FILE_NAME) $(RELEASE_BRANCH) $(FULL_OUTPUT_DIR)/config/kubernetes.json
	@echo -e $(call TARGET_END_LOG)
#########################################################

######################## IMAGE BUILDER TARGETS #######################
# these targets end up calling the image-builder cli via the build_image.sh script
# this is to try and match customer usage of image-builder cli as close as possible during our CI builds
# the image-builder cli will end up calling the local-* targts

# These targets all end with <os>-<os-version>
# ex: release-ami-ubuntu-2004, release-ova-redhat-8
release-ami-% release-ova-% release-image-build-on-metal-%: IMAGE_OS=$(word 1,$(subst -, ,$*))
release-ami-% release-ova-% release-image-build-on-metal-%: IMAGE_OS_VERSION=$(word 2,$(subst -, ,$*))
release-ami-% release-ova-% release-image-build-on-metal-%: IMAGE_OS_FIRMWARE=$(word 3,$(subst -, ,$*))

release-ami-%: IMAGE_FORMAT=ami
release-ami-%: export MANIFEST_OUTPUT=$(FULL_OUTPUT_DIR)/manifest.json
release-ami-%: EXPORT_AMI_BUCKET?=$(ARTIFACTS_BUCKET)
release-ami-%: AMI_S3_DST=$(EXPORT_AMI_BUCKET)/$(ARTIFACTS_UPLOAD_PATH)/ami/$(IMAGE_OS_DIR)
release-ami-%: EXPORT_AMI_DST=$(AMI_S3_DST)/$(GIT_HASH)
release-ami-%: LATEST_AMI_S3_URL=$(AMI_S3_DST)/$(LATEST)/ubuntu.raw
release-ami-%: validate-supported-image-% setup-ami-share | ensure-jq
	@echo -e $(call TARGET_START_LOG)
	build/build_image.sh $(IMAGE_OS) $(IMAGE_OS_VERSION) $(RELEASE_BRANCH) $(IMAGE_FORMAT) $(ARTIFACTS_BUCKET) $(LATEST)
	build/export-ami-to-s3.sh $(RELEASE_BRANCH) $(MANIFEST_OUTPUT) raw $(EXPORT_AMI_DST) $(LATEST_AMI_S3_URL)
	@echo -e $(call TARGET_END_LOG)

release-ova-%: IMAGE_FORMAT=ova
release-ova-%: validate-supported-image-% $(IF_REDHAT_CONFIG_TARGET) | ensure-jq $$(ENABLE_LOGGING)
	@build/build_image.sh $(IMAGE_OS) $(IMAGE_OS_VERSION) $(RELEASE_BRANCH) $(IMAGE_FORMAT) $(ARTIFACTS_BUCKET) $(LATEST) $(IMAGE_OS_FIRMWARE)

release-clone-ova-%: IMAGE_FORMAT=ova
release-clone-ova-%: IMAGE_OS=$(word 1,$(subst -, ,$*))
release-clone-ova-%: IMAGE_OS_VERSION=$(word 2,$(subst -, ,$*))
release-clone-ova-%: IMAGE_OS_FIRMWARE=$(word 3,$(subst -, ,$*))
release-clone-ova-%: validate-supported-image-% $(IF_REDHAT_CONFIG_TARGET) | ensure-jq $$(ENABLE_LOGGING)
	@build/build_image.sh $(IMAGE_OS) $(IMAGE_OS_VERSION) $(RELEASE_BRANCH) $(IMAGE_FORMAT) $(ARTIFACTS_BUCKET) $(LATEST) $(IMAGE_OS_FIRMWARE)

release-raw-%: IMAGE_FORMAT=raw
release-raw-%: release-image-build-on-metal-%
	@echo Done building $@

release-cloudstack-%: IMAGE_FORMAT=cloudstack
release-cloudstack-%: release-image-build-on-metal-%
	@echo Done building $@

release-nutanix-%: IMAGE_FORMAT=nutanix
release-nutanix-%: validate-supported-image-% $(IF_REDHAT_CONFIG_TARGET) | ensure-jq $$(ENABLE_LOGGING)
	@build/build_image.sh $(IMAGE_OS) $(IMAGE_OS_VERSION) $(RELEASE_BRANCH) $(IMAGE_FORMAT) $(ARTIFACTS_BUCKET) $(LATEST) $(IMAGE_OS_FIRMWARE)

# used for cloudstack and raw
# clone the repo first since its scp'd to the temporary instance which runs kvm
release-image-build-on-metal-%: IMAGE_FORMAT?=raw
release-image-build-on-metal-%: validate-supported-image-% $(GIT_PATCH_TARGET) $(IF_REDHAT_CONFIG_TARGET) | ensure-jq $$(ENABLE_LOGGING)
	@build/build_image_on_metal.sh $(BASE_DIRECTORY) $(PROJECT_PATH) $(RELEASE_BRANCH) $(RAW_IMAGE_BUILD_AMI) $(RAW_IMAGE_BUILD_INSTANCE_TYPE) $(RAW_IMAGE_BUILD_KEY_NAME) $(IMAGE_OS) $(IMAGE_OS_VERSION) $(IMAGE_FORMAT) $(LATEST) $(FINAL_IMAGE_DIR) $(BRANCH_NAME)

#######################################################################

# These targets all end with <image format>-<os>-<os-version><optional -efi>
# ex: local-build-ami-ubuntu-2004, s3-artifacts-ova-ubuntu-2004, upload-artifacts-ova-ubuntu-2004, upload-artifacts-ova-ubuntu-2004-efi
# Special case for clone targets: local-clone-build-ova-redhat-9-efi
local-clone-build-%: IMAGE_FORMAT=$(word 1,$(subst -, ,$*))
local-clone-build-%: IMAGE_OS=$(word 2,$(subst -, ,$*))
local-clone-build-%: IMAGE_OS_VERSION=$(word 3,$(subst -, ,$*))
local-clone-build-%: IMAGE_OS_FIRMWARE=$(word 4,$(subst -, ,$*))

# Standard pattern for non-clone targets
local-build-% s3-artifacts-% upload-artifacts-%: IMAGE_FORMAT=$(word 1,$(subst -, ,$*))
local-build-% s3-artifacts-% upload-artifacts-%: IMAGE_OS=$(word 2,$(subst -, ,$*))
local-build-% s3-artifacts-% upload-artifacts-%: IMAGE_OS_VERSION=$(word 3,$(subst -, ,$*))
local-build-% s3-artifacts-% upload-artifacts-%: IMAGE_OS_FIRMWARE=$(word 4,$(subst -, ,$*))

######################## IMAGE BUILDER CLI TARGETS #######################
# OVF_CUSTOM_PROPERTIES is only for ova, but should be ok to passed since its ignored by other builds
# Special case for clone targets
local-clone-build-%: MAKEFLAGS=
local-clone-build-%: validate-supported-image-% validate-dependency-versions-$$(IMAGE_FORMAT) setup-packer-configs-%
	@echo -e $(call TARGET_START_LOG)
	if [[ "$(DEBUG)" == "1" ]]; then \
		touch $(FINAL_IMAGE_DIR)/packer.log; \
		tail -f $(FINAL_IMAGE_DIR)/packer.log & \
	fi
	DEBUG=$(DEBUG) PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(FINAL_IMAGE_DIR)/packer.log PACKER_VAR_FILES="$(PACKER_VAR_FILES)" \
		OVF_CUSTOM_PROPERTIES="$(FULL_OUTPUT_DIR)/config/ovf_custom_properties.json" \
		$(MAKE) -C $(IMAGE_BUILDER_DIR) $(call CONVERT_TO_UPSTREAM_TARGET,build-$(IMAGE_FORMAT)-clone-$(IMAGE_OS_WITH_VER)$(if $(IMAGE_OS_FIRMWARE),-$(IMAGE_OS_FIRMWARE),))
	@echo -e $(call TARGET_END_LOG)

# Standard pattern for non-clone targets
local-build-%: MAKEFLAGS=
local-build-%: validate-supported-image-% validate-dependency-versions-$$(IMAGE_FORMAT) setup-packer-configs-%
	@echo -e $(call TARGET_START_LOG)
	if [[ "$(DEBUG)" == "1" ]]; then \
		touch $(FINAL_IMAGE_DIR)/packer.log; \
		tail -f $(FINAL_IMAGE_DIR)/packer.log & \
	fi
	DEBUG=$(DEBUG) PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(FINAL_IMAGE_DIR)/packer.log PACKER_VAR_FILES="$(PACKER_VAR_FILES)" \
		OVF_CUSTOM_PROPERTIES="$(FULL_OUTPUT_DIR)/config/ovf_custom_properties.json" \
		$(MAKE) -C $(IMAGE_BUILDER_DIR) $(call CONVERT_TO_UPSTREAM_TARGET,build-$(IMAGE_FORMAT)-$(IMAGE_OS_WITH_VER)$(if $(IMAGE_OS_FIRMWARE),-$(IMAGE_OS_FIRMWARE),))
	@echo -e $(call TARGET_END_LOG)

validate-dependency-versions-%: IMAGE_FORMAT=$*
validate-dependency-versions-%:
	@$(MAKE) -C $(IMAGE_BUILDER_DIR) $(call CONVERT_TO_UPSTREAM_TARGET,deps-$(IMAGE_FORMAT))
	@build/validate_required_deps.sh $(IMAGE_FORMAT)

#######################################################################

######################## ARTIFACT TARGETS #############################
s3-artifacts-%: $$(FINAL_IMAGE_PATH) | $$(ENABLE_LOGGING)
	@$(MAKE) -C $(MAKE_ROOT) s3-artifacts ARTIFACTS_PATH=$(ARTIFACTS_PATH_IMAGE_FORMAT) IMAGE_FORMAT=$(IMAGE_FORMAT) IMAGE_OS=$(IMAGE_OS) IMAGE_OS_VERSION=$(IMAGE_OS_VERSION)

upload-artifacts-%: TARGETS=$(if $(filter ova,$(IMAGE_FORMAT)),validate-ova-kubelet,) upload-artifacts
upload-artifacts-%: s3-artifacts-% | $$(ENABLE_LOGGING)
	@$(MAKE) -C $(MAKE_ROOT) $(TARGETS) ARTIFACTS_PATH=$(ARTIFACTS_PATH_IMAGE_FORMAT) ARTIFACTS_UPLOAD_PATH=$(ARTIFACTS_UPLOAD_PATH_IMAGE_FORMAT) IMAGE_FORMAT=$(IMAGE_FORMAT) IMAGE_OS=$(IMAGE_OS) IMAGE_OS_VERSION=$(IMAGE_OS_VERSION)

#######################################################################

.PHONY: validate-supported-image-%
validate-supported-image-%:
	@if [[ "$(IMAGE_FORMAT)" == "" ]]; then \
		echo "IMAGE_FORMAT is required"; \
		exit 1; \
	elif [[ "$(IMAGE_OS)" == "" ]]; then \
		echo "IMAGE_OS is required"; \
		exit 1; \
	elif [[ "$(IMAGE_OS_VERSION)" == "" ]]; then \
		echo "IMAGE_OS_VERSION is required"; \
		exit 1; \
	elif [[ "$(ALLOWED_$(call TO_UPPER,$(IMAGE_FORMAT))_OS)" != *"$(IMAGE_OS)"* ]]; then \
		echo "Unsupported combination of image format ($(IMAGE_FORMAT)) and image os ($(IMAGE_OS)) chosen. Build will quit now..."; \
		exit 1; \
	elif [[ "$(ALLOWED_$(call TO_UPPER,$(IMAGE_OS))_VERSIONS)" != *"$(IMAGE_OS_VERSION)"* ]]; then \
		echo "Unsupported combination of image os ($(IMAGE_OS)) and image os version ($(IMAGE_OS_VERSION)) chosen. Build will quit now..."; \
		exit 1; \
	elif  [ -n "$(IMAGE_OS_FIRMWARE)" ] && [[ "$(IMAGE_OS_FIRMWARE)" != "efi" ]]; then \
		echo "efi is the only supported custom firmware"; \
		exit 1; \
	elif [ -n "$(IMAGE_OS_FIRMWARE)" ] && [[ "$(ALLOWED_EFI_OS_FORMATS)" != *"$(IMAGE_OS)-$(IMAGE_FORMAT)"* ]]; then \
		echo "EFI is only supported for Ubuntu/RedHat OVAs and Redhat RAW builds"; \
		exit 1; \
	fi

.PHONY: update-dependency-versions
update-dependency-versions: $(GIT_PATCH_TARGET) | $$(ENABLE_LOGGING)
	build/update_required_deps.sh; \
	if [ "$$(git status --porcelain -- REQUIRED_DEPENDENCY_VERSIONS.yaml | wc -l)" -gt 0 ]; then \
		echo "Error: Generated files, REQUIRED_DEPENDENCY_VERSIONS.yaml, do not match expected. Please run make update-dependency-versions to update"; \
		git diff -- REQUIRED_DEPENDENCY_VERSIONS.yaml; \
		exit 1; \
	fi


######################## LOCAL HELPERS ####################################

.PHONY: validate-ubuntu-%
validate-ubuntu-%: check-env-validation $(GIT_PATCH_TARGET) setup-packer-configs-ova
	sh validate/validate.sh $(FULL_OUTPUT_DIR)/config/validate_goss_inline_vars.json

.PHONY: check-env-validation
check-env-validation:
	ifndef TEMPLATE
		$(error Environment var TEMPLATE not set. Example: TEMPLATE=<ubuntu.livecd> make ...)
	endif
############################################################################

######################## TEMP BACK COMPAT ####################################

.PHONY: local-build-ova-rhel-8
local-build-ova-rhel-8: local-build-ova-redhat-8
	@echo Done building $@

.PHONY: local-build-raw-rhel-8
local-build-raw-rhel-8: local-build-raw-redhat-8
	@echo Done building $@

.PHONY: local-build-raw-rhel-9-efi
local-build-raw-rhel-9: local-build-raw-redhat-9-efi
	@echo Done building $@

.PHONY: local-build-qemu-rhel-8
local-build-qemu-rhel-8: local-build-cloudstack-redhat-8
	@echo Done building $@

.PHONY: local-build-raw-ubuntu-2004-efi
local-build-raw-ubuntu-2004-efi: local-build-raw-ubuntu-2004
	@echo Done building $@
############################################################################

.PHONY: clean-extra
clean-extra:
	@rm -rf fake-*

clean: clean-extra


.PHONY: metal-instance-test
metal-instance-test:
	@if [ "$(JOB_TYPE)" = "presubmit" ]; then \
		build/prow_build.sh release-raw-ubuntu-2004; \
	else \
		$(MAKE) release-raw-ubuntu-2004; \
	fi

	

## --------------------------------------
## Document dynamic targets
## --------------------------------------

build-ami-ubuntu-2004:
build-ami-ubuntu-2204:
release-ami-ubuntu-2004:
release-ami-ubuntu-2204:

build-nutanix-ubuntu-2004:
build-nutanix-ubuntu-2204:
build-nutanix-redhat-8:
build-nutanix-redhat-9:
release-nutanix-ubuntu-2004:
release-nutanix-ubuntu-2204:
release-nutanix-redhat-8:
release-nutanix-redhat-9:

build-ova-ubuntu-2004:
build-ova-ubuntu-2204:
build-ova-ubuntu-2404:
build-ova-ubuntu-2004-efi:
build-ova-ubuntu-2204-efi:
build-ova-ubuntu-2404-efi:
build-ova-redhat-8:
release-ova-ubuntu-2004:
release-ova-ubuntu-2204:
release-ova-ubuntu-2404:
release-ova-redhat-8:

build-raw-ubuntu-2004:
build-raw-ubuntu-2204:
build-raw-ubuntu-2404:
build-raw-redhat-8:
build-raw-redhat-9-efi:
release-raw-ubuntu-2004:
release-raw-ubuntu-2204:
release-raw-ubuntu-2404:
release-raw-redhat-8:
release-raw-redhat-9-efi:

build-cloudstack-redhat-8:
build-cloudstack-redhat-9:
release-cloudstack-redhat-8:
release-cloudstack-redhat-9:

upload-bottlerocket-1-ova:
upload-bottlerocket-1-raw:
upload-bottlerocket-1-ami:


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

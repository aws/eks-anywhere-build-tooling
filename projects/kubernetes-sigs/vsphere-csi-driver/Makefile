BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=vsphere-csi-driver
REPO_OWNER=kubernetes-sigs

BINARY_TARGET_FILES=vsphere-csi-driver vsphere-csi-syncer
SOURCE_PATTERNS=./cmd/vsphere-csi ./cmd/syncer

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CSI_DRIVER_IMAGE_COMPONENT=kubernetes-sigs/vsphere-csi-driver/csi/driver
CSI_SYNCER_IMAGE_COMPONENT=kubernetes-sigs/vsphere-csi-driver/csi/syncer

IMAGE_NAMES=csi-driver csi-syncer

OUTPUT_FILES_TARGET=$(OUTPUT_DIR)/files

# we are no longer build this project since the version of go is 
# no longer supported in our builder base
# this is being removed in the next release
EXCLUDE_FROM_STAGING_BUILDSPEC=true
EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true
BUILD_TARGETS=force-success
RELEASE_TARGETS=force-success
ATTRIBUTION_TARGETS=force-success

include $(BASE_DIRECTORY)/Common.mk


.PHONY: force-success
force-success:
	@echo "Skipping builds of vsphere-csi"

csi-driver/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-csi

$(call IMAGE_TARGETS_FOR_NAME, csi-syncer): $(OUTPUT_FILES_TARGET)

$(OUTPUT_FILES_TARGET):
	@mkdir -p $@
	cp $(REPO)/pkg/apis/cnsoperator/config/cnsregistervolume_crd.yaml \
		$(REPO)/pkg/apis/cnsoperator/config/cnsfileaccessconfig_crd.yaml \
		$(REPO)/pkg/internal/cnsoperator/config/cnsfilevolumeclient_crd.yaml $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

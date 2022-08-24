BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=vsphere-csi-driver
REPO_OWNER=kubernetes-sigs

BINARY_TARGET_FILES=vsphere-csi-driver vsphere-csi-syncer
SOURCE_PATTERNS=./cmd/vsphere-csi ./cmd/syncer
BINARY_TARGET_FILES_BUILD_ALONE=vsphere-csi-syncer
EXTRA_GO_LDFLAGS=-X sigs.k8s.io/vsphere-csi-driver/v2/pkg/csi/service.Version=$(GIT_TAG)

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CSI_DRIVER_IMAGE_COMPONENT=kubernetes-sigs/vsphere-csi-driver/csi/driver
CSI_SYNCER_IMAGE_COMPONENT=kubernetes-sigs/vsphere-csi-driver/csi/syncer

IMAGE_NAMES=csi-driver csi-syncer

include $(BASE_DIRECTORY)/Common.mk

$(OUTPUT_BIN_DIR)/linux-%/vsphere-csi-syncer: EXTRA_GO_LDFLAGS=-X sigs.k8s.io/vsphere-csi-driver/v2/pkg/syncer.Version=$(GIT_TAG)

csi-driver/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-csi


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

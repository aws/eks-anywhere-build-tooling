BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=openebs-k8s-provisioner
REPO_OWNER=openebs

BINARY_TARGET_FILES=openebs-provisioner snapshot-controller snapshot-pv-provisioner
SOURCE_PATTERNS=./cmd/openebs-provisioner ./cmd/snapshot-controller ./cmd/snapshot-pv-provisioner

OPENEBS_PROVISIONER_IMAGE_COMPONENT=openebs-k8s-provisioner/openebs-provisioner
SNAPSHOT_CONTROLLER_IMAGE_COMPONENT=openebs-k8s-provisioner/snapshot-controller
SNAPSHOT_PV_PROVISIONER_IMAGE_COMPONENT=openebs-k8s-provisioner/snapshot-pv-provisioner

IMAGE_NAMES=openebs-provisioner snapshot-controller snapshot-pv-provisioner
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################
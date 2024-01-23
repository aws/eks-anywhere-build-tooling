BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=cloud-provider-nutanix
REPO_OWNER=nutanix-cloud-native
COMPONENT=${REPO_OWNER}/${REPO}

BASE_IMAGE_NAME?=eks-distro-minimal-base-nonroot

BINARY_TARGET_FILES=controller
SOURCE_PATTERNS=.
EXTRA_GO_LDFLAGS=$(shell $(BUILD_LIB)/version.sh $(REPO) sigs.k8s.io/cloud-provider-nutanix/pkg/version)

IMAGE_COMPONENT=nutanix-cloud-native/cloud-provider-nutanix/controller

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

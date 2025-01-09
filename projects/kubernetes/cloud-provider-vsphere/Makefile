BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat ./$(RELEASE_BRANCH)/GIT_TAG)
GOLANG_VERSION=$(shell cat ./$(RELEASE_BRANCH)/GOLANG_VERSION)

REPO=cloud-provider-vsphere
REPO_OWNER=kubernetes

BASE_IMAGE_NAME?=eks-distro-minimal-base

FIX_LICENSES_VSPHERE_AUTOMATION_TARGET=$(REPO)/vendor/github.com/vmware/vsphere-automation-sdk-go/lib/LICENSE.txt

BINARY_TARGET_FILES=vsphere-cloud-controller-manager
SOURCE_PATTERNS=./cmd/vsphere-cloud-controller-manager

IMAGE_COMPONENT=kubernetes/cloud-provider-vsphere/cpi/manager

HAS_RELEASE_BRANCHES=true

# we need to set IMAGE_BUILD_ARGS here even though its the same as the default. 
# it is set in Common.mk on the images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L799)
# and the combine-images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L846)
# since combine-images has images as prereq target, the ?= does not really behavior as one might expect.
# the images target being the actual action, its version of the set takes prioirty and resets to empty
# setting it explicitly to empty here takes allows the combine-images override to take proirty
IMAGE_BUILD_ARGS=
DOCKERFILE_FOLDER=./docker/linux

BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_VSPHERE_AUTOMATION_TARGET)

$(FIX_LICENSES_VSPHERE_AUTOMATION_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
  # The vsphere-automation-sdk-go dependency github repo has a license however
  # it is broken up into three separate go modules. Since the license file does not live in the same
  # folder as the go.mod files it is not being included in the downloaded package. Manually
  # downloading from github and placing in each of the packages from vsphere-automation-sdk-go
  # under vendor to make go-licenses happy.  The license needs to be copied into each package
  # folder, otherwise go-licenses will group them all together as vsphere-automation-sdk-go
  # which would be wrong
	for package in lib runtime services ; do \
		wget --retry-connrefused -q https://raw.githubusercontent.com/vmware/vsphere-automation-sdk-go/master/LICENSE.txt -O \
			$(REPO)/vendor/github.com/vmware/vsphere-automation-sdk-go/$$package/LICENSE.txt; \
	done;


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

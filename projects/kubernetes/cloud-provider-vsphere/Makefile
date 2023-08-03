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
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

BUILDSPEC_VARS_KEYS=RELEASE_BRANCH
BUILDSPEC_VARS_VALUES=SUPPORTED_K8S_VERSIONS

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

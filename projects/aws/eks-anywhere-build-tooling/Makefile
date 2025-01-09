BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)

REPO=eks-anywhere-build-tooling
REPO_OWNER=aws

IMAGE_COMPONENT=eks-anywhere-cli-tools
IMAGE_NAMES=$(IMAGE_COMPONENT)

BASE_IMAGE_NAME?=eks-distro-minimal-base-docker-client

SIMPLE_CREATE_BINARIES=false
HAS_LICENSES=false

BUILD_TARGETS=local-images
RELEASE_TARGETS=images

PROJECT_DEPENDENCIES=eksa/fluxcd/flux2 eksa/kubernetes-sigs/cluster-api \
	eksa/kubernetes-sigs/kind eksa/replicatedhq/troubleshoot eksa/vmware/govmomi eksd/kubernetes/client eksa/helm/helm \
	eksa/apache/cloudstack-cloudmonkey

ORGANIZE_BINARIES_TARGETS=$(addsuffix /eks-a-tools,$(addprefix $(BINARY_DEPS_DIR)/linux-,amd64 arm64))

REPO_NO_CLONE=true

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true
# By default the kind combine-images build will be added to the dep list
# to avoid waiting for that to finish, just wait on the binaries
BUILDSPEC_DEPENDS_ON_OVERRIDE=fluxcd_flux2 kubernetes_sigs_cluster_api kubernetes_sigs_kind_binaries replicatedhq_troubleshoot vmware_govmomi helm_helm apache_cloudstack_cloudmonkey

include $(BASE_DIRECTORY)/Common.mk


$(call IMAGE_TARGETS_FOR_NAME, eks-anywhere-cli-tools): $(ORGANIZE_BINARIES_TARGETS)

$(ORGANIZE_BINARIES_TARGETS): $(HANDLE_DEPENDENCIES_TARGET)
	build/organize_binaries.sh $(lastword $(subst -, ,$(@D)))


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

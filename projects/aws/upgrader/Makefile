BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)

REPO=upgrader
REPO_OWNER=aws

GIT_TAG=v$(RELEASE_BRANCH)-$(shell yq e ".releases[] | select(.branch==\"${RELEASE_BRANCH}\").number" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)

HAS_RELEASE_BRANCHES=true

BUILDSPEC_VARS_KEYS=RELEASE_BRANCH
BUILDSPEC_VARS_VALUES=SUPPORTED_K8S_VERSIONS
BUILDSPEC_DEPENDS_ON_OVERRIDE=containerd_containerd_linux_amd64 containerd_containerd_linux_arm64 kubernetes_sigs_cri_tools
EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true

REPO_NO_CLONE=true
HAS_LICENSES=false

BASE_IMAGE_NAME=eks-distro-minimal-base-glibc
IMAGE_NAMES=upgrader

PROJECT_DEPENDENCIES=eksd/kubernetes/client eksd/kubernetes/server eksd/cni-plugins eksa/containerd/containerd eksa/kubernetes-sigs/cri-tools


include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

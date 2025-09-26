BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
CHART_TAG:=$(shell echo $(GIT_TAG) | sed 's/^v//')

REGISTRY=public.ecr.aws/eks/cilium
CILIUM_IMAGE=cilium
OPERATOR_IMAGE=operator-generic
CILIUM_CHART=cilium

REPO=cilium
REPO_OWNER=cilium

HAS_S3_ARTIFACTS=false
SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false
IMAGE_NAMES=

EXCLUDE_FROM_CHECKSUMS_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk


.PHONY: update-digests
update-digests: | ensure-skopeo
	build/update_digest.sh $(REGISTRY) $(CILIUM_IMAGE) $(GIT_TAG) images/cilium/IMAGE_DIGEST
	build/update_digest.sh $(REGISTRY) $(OPERATOR_IMAGE) $(GIT_TAG) images/operator-generic/IMAGE_DIGEST
	build/update_digest.sh $(REGISTRY) $(CILIUM_CHART) $(CHART_TAG) images/cilium-chart/IMAGE_DIGEST


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

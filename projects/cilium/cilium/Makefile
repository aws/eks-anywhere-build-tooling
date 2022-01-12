BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
HELM_REPO_VERSION?=$(shell cat HELM_REPO_VERSION)
CHART_TAG=$(shell echo $(GIT_TAG) | sed 's/^v//')

REGISTRY=public.ecr.aws/isovalent
CILIUM_IMAGE=cilium
OPERATOR_IMAGE=operator-generic
CILIUM_CHART=cilium

REPO=cilium
REPO_OWNER=cilium

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false

BUILD_TARGETS=upload-artifacts
RELEASE_TARGETS=upload-artifacts

include $(BASE_DIRECTORY)/Common.mk


s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests:
	build/create_manifests.sh $(GIT_TAG) $(HELM_REPO_VERSION) $(ARTIFACTS_PATH)

.PHONY: update-digests
update-digests:
	build/update_digest.sh $(REGISTRY) $(CILIUM_IMAGE) $(GIT_TAG) images/cilium/IMAGE_DIGEST
	build/update_digest.sh $(REGISTRY) $(OPERATOR_IMAGE) $(GIT_TAG) images/operator-generic/IMAGE_DIGEST
	build/update_digest.sh $(REGISTRY) $(CILIUM_CHART) $(CHART_TAG) images/cilium-chart/IMAGE_DIGEST

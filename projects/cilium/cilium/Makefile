BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
HELM_REPO_VERSION?=$(shell cat HELM_REPO_VERSION)

REPO=cilium
REPO_OWNER=cilium

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false

BUILD_TARGETS=s3-artifacts
RELEASE_TARGETS=s3-artifacts

include $(BASE_DIRECTORY)/Common.mk


s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests:
	build/create_manifests.sh $(GIT_TAG) $(HELM_REPO_VERSION) $(ARTIFACTS_PATH)

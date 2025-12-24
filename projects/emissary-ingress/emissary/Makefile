BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
APPMESH_VERSION:=$(shell cat APPMESH_VERSION)
IMAGE_BUILD_ARGS=APPMESH_VERSION COMPILER_IMAGE

REPO=emissary
REPO_OWNER=emissary-ingress

EXCLUDE_FROM_STAGING_BUILDSPEC=true

BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

SKIP_ON_RELEASE_BRANCH=true

BINARY_TARGET_FILES=apiext busyambassador capabilities_wrapper example-envoy-metrics-sink kat-client kat-server
SOURCE_PATTERNS=./cmd/apiext ./cmd/busyambassador ./cmd/capabilities_wrapper ./cmd/example-envoy-metrics-sink ./cmd/kat-client ./cmd/kat-server

EMISSARY_IMAGE_COMPONENT=emissary-ingress/emissary

IMAGE_NAMES?=emissary
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

BASE_IMAGE_NAME?=eks-distro-minimal-base-python
BASE_IMAGE_OS_VERSION=al2023
COMPILER_IMAGE_VERSION=3.11

HAS_HELM_CHART?=true
HELM_DIRECTORY=charts/emissary-ingress
HELM_CHART_NAMES=emissary/crds emissary-ingress/emissary
CHART_VERSION?=

AMBASSADOR_VERSION_TARGET=$(REPO)/python/ambassador.version

include $(BASE_DIRECTORY)/Common.mk


emissary-ingress/emissary/helm/%: PACKAGE_DEPENDENCIES=emissary-crds

emissary/crds/helm/%: HELM_DIRECTORY=charts/crds
# Originially emissary was created with `emissary/crds` as the HELM_DESTINATION_REPOSITORY for everything but the push
# which is `emissary-ingress/crds`
# this really only made a differnece in the helm_require.sh script which uses this value in the EksaPackageRequires file
# otherwise it was always `basename`, which was the same either way
# overriding for push to keep it consistent
emissary/crds/helm/build: HELM_DESTINATION_REPOSITORY=emissary/crds
emissary/crds/helm/push: HELM_DESTINATION_REPOSITORY=emissary-ingress/crds
emissary/crds/__helm__/create-ecr-repo: IMAGE_NAME=emissary-ingress/crds/__helm__

$(GO_MOD_DOWNLOAD_TARGETS): $(AMBASSADOR_VERSION_TARGET)

# do not pass our makeflags down to emissary make since there will be unncessary warnings
$(AMBASSADOR_VERSION_TARGET): MAKEFLAGS=
$(AMBASSADOR_VERSION_TARGET): export VERSION=v3.10.0
$(AMBASSADOR_VERSION_TARGET): export CHART_VERSION:=$(CHART_VERSION)
$(AMBASSADOR_VERSION_TARGET): export GOPATH=/var
$(AMBASSADOR_VERSION_TARGET): | $$(ENABLE_DOCKER)
	@source $(BUILD_LIB)/common.sh && build::common::use_go_version $(GOLANG_VERSION) && $(MAKE) -C $(REPO) python/ambassador.version


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

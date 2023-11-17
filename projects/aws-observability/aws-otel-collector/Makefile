BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
HELM_GIT_TAG=$(shell cat HELM_GIT_TAG)
 # Upstream images are used directly without re-building and re-tagging in build
 # tooling, so the IMAGE_TAG needs to be overwritten to match upstream ECR tags.
IMAGE_TAG=$(GIT_TAG)

REGISTRY=public.ecr.aws/aws-observability
COLLECTOR_IMAGE=aws-otel-collector

REPO_OWNER=aws-observability
REPO=aws-otel-collector

SIMPLE_CREATE_BINARIES=false
SIMPLE_CREATE_TARBALLS=false
HAS_LICENSES=false
IMAGE_NAMES=

BUILD_TARGETS=helm/build
RELEASE_TARGETS=helm/push
EXCLUDE_FROM_STAGING_BUILDSPEC=true
SKIP_ON_RELEASE_BRANCH=true

HAS_HELM_CHART=true

HELM_CHART_NAMES=adot/charts/aws-otel-collector

HELM_SOURCE_IMAGE_REPO=public.ecr.aws
HELM_SOURCE_OWNER=open-telemetry
HELM_SOURCE_REPOSITORY=opentelemetry-helm-charts
HELM_DIRECTORY=charts/opentelemetry-collector
HELM_IMAGE_LIST=aws-observability/aws-otel-collector
HELM_USE_UPSTREAM_IMAGE=true

include $(BASE_DIRECTORY)/Common.mk

.PHONY: update-digests
update-digests:
	source $(BUILD_LIB)/image_shasum.sh $(REGISTRY) $(COLLECTOR_IMAGE) $(GIT_TAG) > images/collector/IMAGE_DIGEST


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

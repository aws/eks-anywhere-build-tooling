BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO_OWNER=aws-observability
REPO=aws-otel-collector

BINARY_TARGET_FILES=aoc
SOURCE_PATTERNS=./cmd/awscollector

IMAGE_COMPONENT=aws-otel-collector/aws-otel-collector
IMAGE_NAMES=aws-otel-collector

HAS_HELM_CHART=true
HELM_GIT_TAG:=$(shell cat HELM_GIT_TAG)
HELM_SOURCE_OWNER=open-telemetry
HELM_SOURCE_REPOSITORY=opentelemetry-helm-charts
HELM_DIRECTORY=charts/opentelemetry-collector
HELM_DESTINATION_REPOSITORY=aws-otel-collector/aws-otel-collector-helm
HELM_IMAGE_LIST=aws-otel-collector/aws-otel-collector

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=emissary
REPO_OWNER=emissary-ingress

BINARY_TARGET_FILES=busyambassador capabilities_wrapper example-envoy-metrics-sink k8sregistryctl kat-client kat-server
SOURCE_PATTERNS=./cmd/busyambassador ./cmd/capabilities_wrapper ./cmd/example-envoy-metrics-sink ./cmd/k8sregistryctl ./cmd/kat-client ./cmd/kat-server
GO_MOD_PATHS=.

EMISSARY_IMAGE_COMPONENT=emissary-ingress/emissary

IMAGE_NAMES=emissary
DOCKERFILE_FOLDER=./docker/linux/emissary-image

CHART_VERSION=$(shell cat CHART_VERSION)

HELM_PULL=true
HELM_PULL_NAME=emissary-ingress
HELM_PULL_LOCATION=datawire/emissary-ingress

HAS_HELM_CHART=true
HELM_DIRECTORY=charts/emissary-ingress
HELM_DESTINATION_REPOSITORY=emissary-ingress/emissary
HELM_IMAGE_LIST=emissary-ingress/emissary
COPY_CRDS=$(BASE_DIRECTORY)/projects/emissary-ingress/emissary/helm/crds

include $(BASE_DIRECTORY)/Common.mk

# This will be replaced once the new eks-distro python images are pushed.
BASE_IMAGE_NAME?=public.ecr.aws/k1e6s8o8/eks-distro-minimal-base-python:3.9-2022-07-21-1658430881.2
emissary/images/%: public.ecr.aws/k1e6s8o8/eks-distro-minimal-base-python:3.9-2022-07-21-1658430881.2

$(GO_MOD_DOWNLOAD_TARGETS): $(OUTPUT_DIR)/$(REPO)

$(OUTPUT_DIR)/$(REPO):
	@export VERSION="v3.0.0" && export CHART_VERSION=$(CHART_VERSION) && cd emissary && GOPATH=/var && make python/ambassador.version

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
##########emissary-builder# END GENERATED ###########################
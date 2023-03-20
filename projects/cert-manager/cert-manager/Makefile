BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.19"

REPO=cert-manager
REPO_OWNER=cert-manager

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CERT_MANAGER_ACMESOLVER_IMAGE_COMPONENT=cert-manager/cert-manager-acmesolver
CERT_MANAGER_CAINJECTOR_IMAGE_COMPONENT=cert-manager/cert-manager-cainjector
CERT_MANAGER_CONTROLLER_IMAGE_COMPONENT=cert-manager/cert-manager-controller
CERT_MANAGER_WEBHOOK_IMAGE_COMPONENT=cert-manager/cert-manager-webhook
CERT_MANAGER_CTL_IMAGE_COMPONENT=cert-manager/cert-manager-ctl

IMAGE_NAMES=cert-manager-acmesolver cert-manager-cainjector cert-manager-controller cert-manager-webhook cert-manager-ctl

BINARY_TARGET_FILES=cert-manager-acmesolver cert-manager-cainjector cert-manager-controller cert-manager-webhook cert-manager-ctl
SOURCE_PATTERNS=./cmd/acmesolver ./cmd/cainjector ./cmd/controller ./cmd/webhook ./cmd/ctl

HAS_HELM_CHART=true
HELM_DIRECTORY=deploy/charts/cert-manager
HELM_IMAGE_LIST=cert-manager/cert-manager-cainjector cert-manager/cert-manager-controller cert-manager/cert-manager-webhook cert-manager/cert-manager-ctl

HAS_S3_ARTIFACTS=true

BUILDSPEC_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

include $(BASE_DIRECTORY)/Common.mk

$(call IMAGE_TARGETS_FOR_NAME, cert-manager-acmesolver): cert-manager-acmesolver-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, cert-manager-cainjector): cert-manager-cainjector-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, cert-manager-acmesolver): cert-manager-controller-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, cert-manager-webhook): cert-manager-webhook-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, cert-manager-ctl): cert-manager-ctl-useradd/images/export

cert-manager-acmesolver-useradd/images/export: IMAGE_USERADD_USER_NAME=acmesolver
cert-manager-cainjector-useradd/images/export: IMAGE_USERADD_USER_NAME=cainjector
cert-manager-controller-useradd/images/export: IMAGE_USERADD_USER_NAME=cert-manager
cert-manager-webhook-useradd/images/export: IMAGE_USERADD_USER_NAME=webhook
cert-manager-ctl-useradd/images/export: IMAGE_USERADD_USER_NAME=ctl

s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests:
	build/create_manifests.sh $(ARTIFACTS_PATH) $(BASE_DIRECTORY) $(GIT_TAG)

.PHONY: update-cert-manager-manifest
update-cert-manager-manifest:
	wget -q https://github.com/cert-manager/cert-manager/releases/download/$(GIT_TAG)/cert-manager.yaml -O manifests/cert-manager.yaml


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

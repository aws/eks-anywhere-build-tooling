BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=cert-manager
REPO_OWNER=jetstack

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CERT_MANAGER_ACMESOLVER_IMAGE_COMPONENT=jetstack/cert-manager-acmesolver
CERT_MANAGER_CAINJECTOR_IMAGE_COMPONENT=jetstack/cert-manager-cainjector
CERT_MANAGER_CONTROLLER_IMAGE_COMPONENT=jetstack/cert-manager-controller
CERT_MANAGER_WEBHOOK_IMAGE_COMPONENT=jetstack/cert-manager-webhook

IMAGE_NAMES=cert-manager-acmesolver cert-manager-cainjector cert-manager-controller cert-manager-webhook

BINARY_TARGET_FILES=cert-manager-acmesolver cert-manager-cainjector cert-manager-controller cert-manager-webhook
SOURCE_PATTERNS=./cmd/acmesolver ./cmd/cainjector ./cmd/controller ./cmd/webhook

HAS_HELM_CHART=true
HELM_DIRECTORY=deploy/charts/cert-manager

include $(BASE_DIRECTORY)/Common.mk

.PHONY: sign-file
sign-file:
	aws kms sign --key-id alias/signingKey --message fileb://helm/files.txt --signing-algorithm ECDSA_SHA_256 --message-type RAW

$(call IMAGE_TARGETS_FOR_NAME, cert-manager-acmesolver): cert-manager-acmesolver-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, cert-manager-cainjector): cert-manager-cainjector-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, cert-manager-acmesolver): cert-manager-controller-useradd/images/export

cert-manager-acmesolver-useradd/images/export: IMAGE_USERADD_USER_NAME=acmesolver
cert-manager-cainjector-useradd/images/export: IMAGE_USERADD_USER_NAME=cainjector
cert-manager-controller-useradd/images/export: IMAGE_USERADD_USER_NAME=cert-manager


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

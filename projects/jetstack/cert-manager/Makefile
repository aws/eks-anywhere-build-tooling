BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=cert-manager
REPO_OWNER=jetstack

LICENSE_PACKAGE_FILTER?=./cmd/acmesolver ./cmd/cainjector ./cmd/controller ./cmd/webhook

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CERT_MANAGER_ACMESOLVER_IMAGE_COMPONENT=jetstack/cert-manager-acmesolver
CERT_MANAGER_CAINJECTOR_IMAGE_COMPONENT=jetstack/cert-manager-cainjector
CERT_MANAGER_CONTROLLER_IMAGE_COMPONENT=jetstack/cert-manager-controller
CERT_MANAGER_WEBHOOK_IMAGE_COMPONENT=jetstack/cert-manager-webhook

BINARY_TARGET_FILES=cert-manager-acmesolver cert-manager-cainjector cert-manager-controller cert-manager-webhook
SOURCE_PATTERNS=./cmd/acmesolver ./cmd/cainjector ./cmd/controller ./cmd/webhook

include $(BASE_DIRECTORY)/Common.mk


.PHONY: local-images
local-images: cert-manager-acmesolver/images/amd64 cert-manager-cainjector/images/amd64 cert-manager-controller/images/amd64 cert-manager-webhook/images/amd64

.PHONY: images
images: cert-manager-acmesolver/images/push cert-manager-cainjector/images/push cert-manager-controller/images/push cert-manager-webhook/images/push

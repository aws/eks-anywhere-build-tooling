BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=hook
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

BOOTKIT_IMAGE_COMPONENT=tinkerbell/hook-bootkit
TINK_DOCKER_IMAGE_COMPONENT=tinkerbell/hook-docker
KERNEL_IMAGE_COMPONENT=tinkerbell/hook-kernel

IMAGE_NAMES=bootkit tink-docker kernel

BINARY_TARGET_FILES=bootkit tink-docker
SOURCE_PATTERNS=./ ./

GO_MOD_DOWNLOAD_TARGETS=$(REPO)/bootkit/eks-anywhere-go-mod-download $(REPO)/tink-docker/eks-anywhere-go-mod-download
GATHER_LICENSES_TARGETS=$(OUTPUT_DIR)/bootkit/attribution/go-license.csv $(OUTPUT_DIR)/tink-docker/attribution/go-license.csv
ATTRIBUTION_TARGETS=BOOTKIT_ATTRIBUTION.txt TINK_DOCKER_ATTRIBUTION.txt

include $(BASE_DIRECTORY)/Common.mk

$(OUTPUT_BIN_DIR)/%/bootkit $(REPO)/bootkit/eks-anywhere-go-mod-download \
	$(OUTPUT_DIR)/bootkit/attribution/go-license.csv: REPO_SUBPATH=./bootkit

$(OUTPUT_DIR)/bootkit/attribution/go-license.csv BOOTKIT_ATTRIBUTION.txt: LICENSES_OUTPUT_DIR=$(OUTPUT_DIR)/bootkit

$(OUTPUT_BIN_DIR)/%/tink-docker $(REPO)/tink-docker/eks-anywhere-go-mod-download \
	$(OUTPUT_DIR)/tink-docker/attribution/go-license.csv: REPO_SUBPATH=./tink-docker

$(OUTPUT_DIR)/tink-docker/attribution/go-license.csv TINK_DOCKER_ATTRIBUTION.txt: LICENSES_OUTPUT_DIR=$(OUTPUT_DIR)/tink-docker

kernel/images/%: DOCKERFILE_FOLDER=$(REPO)/kernel
kernel/images/%: IMAGE_CONTEXT_DIR=$(REPO)/kernel
kernel/images/%: IMAGE_BUILD_ARGS=KERNEL_VERSION KERNEL_SERIES EXTRA DEBUG
kernel/images/%: KERNEL_VERSION=5.10.11 
kernel/images/%: KERNEL_SERIES=5.10.x

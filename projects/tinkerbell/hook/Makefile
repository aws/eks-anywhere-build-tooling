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

CREATE_HOOK_FILES=$(REPO)/dist/vmlinuz-x86_64

HAS_S3_ARTIFACTS=true
SIMPLE_CREATE_TARBALLS=false

include $(BASE_DIRECTORY)/Common.mk

$(OUTPUT_BIN_DIR)/%/bootkit $(REPO)/bootkit/eks-anywhere-go-mod-download \
	$(OUTPUT_DIR)/bootkit/attribution/go-license.csv: REPO_SUBPATH=./bootkit

$(OUTPUT_DIR)/bootkit/attribution/go-license.csv BOOTKIT_ATTRIBUTION.txt: LICENSES_OUTPUT_DIR=$(OUTPUT_DIR)/bootkit

$(OUTPUT_BIN_DIR)/%/tink-docker $(REPO)/tink-docker/eks-anywhere-go-mod-download \
	$(OUTPUT_DIR)/tink-docker/attribution/go-license.csv: REPO_SUBPATH=./tink-docker

$(OUTPUT_DIR)/tink-docker/attribution/go-license.csv TINK_DOCKER_ATTRIBUTION.txt: LICENSES_OUTPUT_DIR=$(OUTPUT_DIR)/tink-docker

%/images/amd64: IMAGE_OUTPUT=push=true
%/images/amd64: IMAGE_OUTPUT_TYPE=image

kernel/images/%: BASE_IMAGE=quay.io/tinkerbell/hook-kernel:5.10.11

s3-artifacts: $(CREATE_HOOK_FILES)

$(CREATE_HOOK_FILES): tarballs
	sed -e 's/IMAGE_REPO/$(IMAGE_REPO)/g' linuxkit_config.yml > _linuxkit_config.yml
	LINUXKIT_CONFIG=$(BASE_DIRECTORY)/projects/tinkerbell/hook/_linuxkit_config.yml make dist-existing-images -C $(REPO)
	mkdir -p $(OUTPUT_DIR)/hook/$(GIT_TAG)
	cp $(REPO)/dist/* $(OUTPUT_DIR)/hook/$(GIT_TAG)/
	cp -rf $(OUTPUT_DIR)/hook $(ARTIFACTS_PATH)
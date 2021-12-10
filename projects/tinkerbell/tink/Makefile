BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=tink
REPO_OWNER=tinkerbell

LICENSE_PACKAGE_FILTER?=./cmd/tink-cli ./cmd/tink-server ./cmd/tink-worker

BASE_IMAGE_NAME?=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

TINK_CLI_IMAGE_COMPONENT=tinkerbell/tink/tink-cli
TINK_SERVER_IMAGE_COMPONENT=tinkerbell/tink/tink-server
TINK_WORKER_IMAGE_COMPONENT=tinkerbell/tink/tink-worker

IMAGE_NAMES=tink-cli tink-server tink-worker

BINARY_TARGET_FILES=tink-cli tink-server tink-worker
SOURCE_PATTERNS=./cmd/tink-cli ./cmd/tink-server ./cmd/tink-worker

VERSION?=$(shell git -C $(REPO) rev-parse --short HEAD)

# When we track a tag instead of commit, this would become $(GIT_TAG)-$(VERSION)
EXTRA_GO_LDFLAGS=-X main.version=$(VERSION)

FETCH_BINARIES_TARGETS=eksa/cloudflare/cfssl

include $(BASE_DIRECTORY)/Common.mk

$(call IMAGE_TARGETS_FOR_NAME, tink-server): $(call FULL_FETCH_BINARIES_TARGETS, $(FETCH_BINARIES_TARGETS))

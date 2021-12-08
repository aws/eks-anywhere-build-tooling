BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.14"

REPO=cri-tools
REPO_OWNER=kubernetes-sigs

BINARY_TARGET_FILES=crictl critest
SOURCE_PATTERNS=./cmd/crictl ./cmd/critest
EXTRA_GO_LDFLAGS=-X github.com/kubernetes-sigs/cri-tools/pkg/version.Version=$(GIT_TAG:v%=%)
EXTRA_GOBUILD_FLAGS=-tags selinux

HAS_S3_ARTIFACTS=true
IMAGE_NAMES=

include $(BASE_DIRECTORY)/Common.mk

$(OUTPUT_BIN_DIR)/linux-%/critest: GOBUILD_COMMAND=test
$(OUTPUT_BIN_DIR)/linux-%/critest: EXTRA_GOBUILD_FLAGS+=-c

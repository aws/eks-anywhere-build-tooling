BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.14"

REPO=cfssl
REPO_OWNER=cloudflare

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=cfssl cfssl-bundle cfssl-certinfo cfssl-newkey cfssl-scan cfssljson mkbundle multirootca
SOURCE_PATTERNS=./cmd/cfssl ./cmd/cfssl-bundle ./cmd/cfssl-certinfo ./cmd/cfssl-newkey ./cmd/cfssl-scan ./cmd/cfssljson ./cmd/mkbundle ./cmd/multirootca

EXTRA_GO_LDFLAGS=-s -w -X github.com/cloudflare/cfssl/cli/version.version=$(GIT_TAG)

HAS_S3_ARTIFACTS=true

RICE_BINARY_PLATFORMS=$(BUILDER_PLATFORM)

RICE_EMBED_TARGET=$(REPO)/cli/serve/rice-box.go

GIT_CLONE_DEPENDENCY=$(GIT_DEPS_DIR)/cfssl_trust/LICENSE

include $(BASE_DIRECTORY)/Common.mk

$(OUTPUT_BIN_DIR)/linux-%/rice: EXTRA_GO_LDFLAGS=

$(OUTPUT_BIN_DIR)/linux-amd64/cfssl $(OUTPUT_BIN_DIR)/linux-arm64/cfssl: $(RICE_EMBED_TARGET) $(GIT_CLONE_DEPENDENCY)

# make sure we build rice for the current os/arch so we can actually run it
$(eval $(foreach platform, $(RICE_BINARY_PLATFORMS), $(call BINARY_TARGET_BODY,$(platform),rice,./vendor/github.com/GeertJohan/go.rice/rice)))

$(RICE_EMBED_TARGET): $(OUTPUT_BIN_DIR)/$(subst /,-,$(RICE_BINARY_PLATFORMS))/rice
	$^ embed-go -i=./$(REPO)/cli/serve

$(GIT_CLONE_DEPENDENCY):
	git clone https://github.com/cloudflare/cfssl_trust.git $(GIT_DEPS_DIR)/cfssl_trust

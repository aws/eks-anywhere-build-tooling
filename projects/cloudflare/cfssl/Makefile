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

RICE_EMBED_TARGET=$(REPO)/cli/serve/rice-box.go

GIT_CLONE_DEPENDENCY=$(GIT_DEPS_DIR)/cfssl_trust/LICENSE

STATIC_FILES=$(shell git -C cfssl ls-tree -r -t --name-only HEAD ./cli/serve/static)

include $(BASE_DIRECTORY)/Common.mk


$(OUTPUT_BIN_DIR)/linux-%/rice: EXTRA_GO_LDFLAGS=

$(OUTPUT_BIN_DIR)/linux-amd64/cfssl $(OUTPUT_BIN_DIR)/linux-arm64/cfssl: $(RICE_EMBED_TARGET) $(GIT_CLONE_DEPENDENCY)

$(RICE_EMBED_TARGET): RICE=./$(OUTPUT_BIN_DIR)/$(subst /,-,$(BUILDER_PLATFORM))/rice
$(RICE_EMBED_TARGET):
	$(MAKE) $(RICE) BINARY_PLATFORMS=$(BUILDER_PLATFORM) BINARY_TARGET_FILES=rice SOURCE_PATTERNS=./vendor/github.com/GeertJohan/go.rice/rice
# We need to set the timestamp of the static files cfssl/cli/serve/*
# The timestamp is being used in the generated file from the embed-go command cfssl/cli/serve/rice-box.go
# When running git clone the timestamp gets set initially and we need this to be a static value to prevent checksum issues.
	$(foreach filename, $(STATIC_FILES), touch -a -m -t $(shell git -C $(REPO) log -1 --format="%cd" --date=format:%Y%m%d%H%M.%S -- $(filename)) ./$(REPO)/$(filename);)
	$(RICE) embed-go -i=./$(REPO)/cli/serve

$(GIT_CLONE_DEPENDENCY):
	git clone https://github.com/cloudflare/cfssl_trust.git $(GIT_DEPS_DIR)/cfssl_trust


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

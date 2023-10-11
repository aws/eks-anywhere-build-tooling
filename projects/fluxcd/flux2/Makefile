BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION:=$(shell cat GOLANG_VERSION)
REPO=flux2
REPO_OWNER=fluxcd
IMAGE_COMPONENT=fluxcd/flux-cli
IMAGE_NAMES=flux-cli

BASE_IMAGE_NAME=eks-distro-minimal-base

BINARY_TARGET_FILES=flux
SOURCE_PATTERNS=./cmd/flux
EXTRA_GO_LDFLAGS=-X main.VERSION=$(GIT_TAG:v%=%)

FLUX_MANIFESTS_TARGET=$(REPO)/cmd/flux/manifests

HAS_S3_ARTIFACTS=true

PROJECT_DEPENDENCIES=eksd/kubernetes/client

FIX_LICENSES_SHA1CD=$(REPO)/vendor/github.com/pjbgf/sha1cd/LICENSE

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS) $(OUTPUT_BIN_DIR)/linux-amd64/flux $(OUTPUT_BIN_DIR)/linux-arm64/flux: $(FLUX_MANIFESTS_TARGET) $(FIX_LICENSES_SHA1CD)

$(FLUX_MANIFESTS_TARGET): export PATH:=$(MAKE_ROOT)/$(OUTPUT_DIR):$(PATH)
$(FLUX_MANIFESTS_TARGET): $(KUSTOMIZE_TARGET)
	# EMBEDDED_MANIFESTS_TARGET from https://github.com/fluxcd/flux2/blob/main/Makefile#L2
	make -C $(REPO) cmd/flux/.manifests.done

$(FIX_LICENSES_SHA1CD): | $(GO_MOD_DOWNLOAD_TARGETS)
# pjbgf/sha1cd includes a LICENSE (apache2) in a later version than the currently depended on version
# Copying in manually for now until the dep is updated upstream
	wget -q --retry-connrefused -O $@ https://raw.githubusercontent.com/pjbgf/sha1cd/main/LICENSE

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

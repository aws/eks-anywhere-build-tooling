BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=govmomi
REPO_OWNER=vmware

BINARY_TARGET_FILES=govc
SOURCE_PATTERNS=./govc

PKG=github.com/vmware/govmomi/govc/flags
GIT_COMMIT=$(shell git -C $(REPO) rev-parse --short HEAD)
BUILD_DATE=$(shell git -C $(REPO) show -s --format=format:%cd --date=format:'%Y-%m-%dT%H:%M:%SZ' HEAD)
EXTRA_GO_LDFLAGS=-X $(PKG).BuildVersion=$(GIT_TAG) -X $(PKG).BuildCommit=$(GIT_COMMIT) -X $(PKG).BuildDate=$(BUILD_DATE)

HAS_S3_ARTIFACTS=true

IMAGE_NAMES=

include $(BASE_DIRECTORY)/Common.mk

build: validate-cli-version
release: validate-cli-version

.PHONY: validate-cli-version
validate-cli-version: CLI=./$(OUTPUT_BIN_DIR)/$(subst /,-,$(BUILDER_PLATFORM))/govc
validate-cli-version:
	$(MAKE) $(CLI) BINARY_PLATFORMS=$(BUILDER_PLATFORM)
	$(CLI) version
	@if [[ "$$($(CLI) version)" != "govc $(GIT_TAG:v%=%)" ]]; then \
		echo "Version set incorrectly on cli!"; \
		exit 1; \
	fi


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

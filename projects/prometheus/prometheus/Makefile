BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.18"

REPO=prometheus
REPO_OWNER=prometheus

BINARY_TARGET_FILES=prometheus promtool
SOURCE_PATTERNS=./cmd/prometheus ./cmd/promtool
# Add additional gobuild flag to embed the web UI in the go binary
EXTRA_GOBUILD_FLAGS=-tags builtinassets

DOCKERFILE_FOLDER=./docker/linux/prometheus

EXCLUDE_FROM_STAGING_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

# Before running GO_MOD_DOWNLOAD_TARGETS, we need to perform two preparation steps:
# (1) add the web UI component to the cloned repo through GIT_PATCH_TARGET;
# (2) copy files needed by the docker container to the OUTPUT_DIR.
$(GO_MOD_DOWNLOAD_TARGETS): $(OUTPUT_DIR)/$(REPO)

$(OUTPUT_DIR)/$(REPO): $(GIT_PATCH_TARGET)
	@mkdir -p $(OUTPUT_DIR)/documentation/examples
	@cp prometheus/documentation/examples/prometheus.yml $(OUTPUT_DIR)/documentation/examples/prometheus.yml

	@mkdir -p $(OUTPUT_DIR)/console_libraries
	@cp -r prometheus/console_libraries/ $(OUTPUT_DIR)/console_libraries/

	@mkdir -p $(OUTPUT_DIR)/consoles
	@cp -r prometheus/consoles/* $(OUTPUT_DIR)/consoles/

	@cp prometheus/npm_licenses.tar.bz2 $(OUTPUT_DIR)/npm_licenses.tar.bz2

# Fix licenses that have the incorrect module names.
# Module "github.com" should be "github.com/prometheus/prometheus"
$(ATTRIBUTION_TARGETS): fix-licenses

.phony: fix-licenses
fix-licenses: $(GATHER_LICENSES_TARGETS)
	build/fix_licenses.sh "github.com,https" "github.com\/prometheus\/prometheus,https" $(OUTPUT_DIR)/attribution/go-license.csv

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

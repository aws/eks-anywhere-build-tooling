BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=node_exporter
REPO_OWNER=prometheus

BINARY_TARGET_FILES=node_exporter
EXTRA_GOBUILD_FLAGS=-tags netgo,osusergo,static_build

NODE_EXPORTER_IMAGE_COMPONENT=prometheus/node-exporter

EXCLUDE_FROM_STAGING_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

# Fix licenses that have incorrect module names.
# Module "github.com" should be "github.com/prometheus/node_exporter"
$(ATTRIBUTION_TARGETS): fix-go-licenses

# Fix go.mod that has mismatched checksum from upstream.
$(GO_MOD_DOWNLOAD_TARGETS): fix-go-mod

.phony: fix-go-licenses
fix-go-licenses: $(GATHER_LICENSES_TARGETS)
	build/sed_replace.sh "github.com,https" "github.com\/prometheus\/node_exporter,https" $(OUTPUT_DIR)/attribution/go-license.csv

.phony: fix-go-mod
fix-go-mod:
	build/sed_replace.sh "github.com\/safchain\/ethtool v0.2.0" "github.com\/safchain\/ethtool v0.1.1-0.20220305102649-532e6d132c2e" node_exporter/go.mod

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

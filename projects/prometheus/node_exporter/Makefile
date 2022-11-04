BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=node_exporter
REPO_OWNER=prometheus

BINARY_TARGET_FILES=node_exporter
BASE_IMAGE_NAME=eks-distro-minimal-base
EXTRA_GOBUILD_FLAGS=-tags netgo,osusergo,static_build

NODE_EXPORTER_IMAGE_COMPONENT=prometheus/node-exporter

EXCLUDE_FROM_STAGING_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

# Fix licenses that have incorrect module names.
# Module "github.com" should be "github.com/prometheus/node_exporter"
$(ATTRIBUTION_TARGETS): fix-go-licenses

.phony: fix-go-licenses
fix-go-licenses: $(GATHER_LICENSES_TARGETS)
	build/sed_replace.sh "github.com,https" "github.com\/prometheus\/node_exporter,https" $(OUTPUT_DIR)/attribution/go-license.csv

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.18"

REPO=prometheus
REPO_OWNER=prometheus

BINARY_TARGET_FILES=prometheus promtool
SOURCE_PATTERNS=./cmd/prometheus ./cmd/promtool
# Add additional gobuild flag to embed the web UI in the go binary
EXTRA_GOBUILD_FLAGS=-tags builtinassets

EXCLUDE_FROM_STAGING_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

# Before running GO_MOD_DOWNLOAD_TARGETS, we need to perform two preparation steps:
# (1) run make assets cp-npm-licenses assets-compress to generate web UI related resources;
# (2) copy files needed by the docker container to the OUTPUT_DIR.
$(GO_MOD_DOWNLOAD_TARGETS): $(OUTPUT_DIR)/$(REPO)

$(OUTPUT_DIR)/$(REPO): assets cp-npm-licenses assets-compress

	mkdir -p $(OUTPUT_DIR)/extra/etc/prometheus
	cp prometheus/documentation/examples/prometheus.yml $(OUTPUT_DIR)/extra/etc/prometheus/prometheus.yml

	mkdir -p $(OUTPUT_DIR)/extra/usr/share/prometheus/console_libraries
	cp -r prometheus/console_libraries/* $(OUTPUT_DIR)/extra/usr/share/prometheus/console_libraries/

	mkdir -p $(OUTPUT_DIR)/extra/usr/share/prometheus/consoles
	cp -r prometheus/consoles/* $(OUTPUT_DIR)/extra/usr/share/prometheus/consoles/

	cp prometheus/npm_licenses.tar.bz2 $(OUTPUT_DIR)/npm_licenses.tar.bz2

# Before running ATTRIBUTION_TARGETS, we need to run fix-go-licenses first to correct
# module names. Module "github.com" should be "github.com/prometheus/prometheus".
$(ATTRIBUTION_TARGETS): fix-go-licenses

.PHONY: assets
assets: ui-install ui-build

.PHONY: assets-compress
assets-compress:
	@echo '>> compressing assets'
	build/compress_assets.sh

.PHONY: cp-npm-licenses
cp-npm-licenses:
	@echo ">> bundling npm licenses"
	build/cp_npm_licenses.sh $(OUTPUT_DIR)

.phony: fix-go-licenses
fix-go-licenses: $(GATHER_LICENSES_TARGETS)
	build/fix_go_licenses.sh "github.com,https" "github.com\/prometheus\/prometheus,https" $(OUTPUT_DIR)/attribution/go-license.csv

.PHONY: ui-install
ui-install:
	cd prometheus/web/ui && npm install

.PHONY: ui-build
ui-build:
	cd prometheus/web/ui && CI="" npm run build

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

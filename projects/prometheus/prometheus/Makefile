BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
HELM_GIT_TAG:=$(shell cat HELM_GIT_TAG)
GOLANG_VERSION?="1.18"

REPO=prometheus
REPO_OWNER=prometheus

BINARY_TARGET_FILES=prometheus promtool
BASE_IMAGE_NAME=eks-distro-minimal-base
SOURCE_PATTERNS=./cmd/prometheus ./cmd/promtool
# Add additional gobuild flag to embed the web UI in the go binary
EXTRA_GOBUILD_FLAGS=-tags netgo,builtinassets
WEB_UI_TARGETS=prometheus/web/ui/embed.go
EXTRA_DOCKER_FILE_TARGETS=$(OUTPUT_DIR)/extra/etc/prometheus/prometheus.yml

EXCLUDE_FROM_STAGING_BUILDSPEC=true

HAS_HELM_CHART=true
HELM_SOURCE_OWNER=prometheus-community
HELM_SOURCE_REPOSITORY=helm-charts
HELM_DIRECTORY=charts/prometheus
HELM_DESTINATION_REPOSITORY=prometheus/charts/prometheus
HELM_IMAGE_LIST=prometheus/prometheus prometheus/node-exporter

include $(BASE_DIRECTORY)/Common.mk

# Before running GO_MOD_DOWNLOAD_TARGETS, we need to generate web UI related 
# resources.
$(GO_MOD_DOWNLOAD_TARGETS): $(WEB_UI_TARGETS)

# Before running image relared targets, we need to copy files needed by the 
# docker container to the OUTPUT_DIR.
$(LOCAL_IMAGE_TARGETS): $(EXTRA_DOCKER_FILE_TARGETS)
$(IMAGE_TARGETS): $(EXTRA_DOCKER_FILE_TARGETS)

# Before running ATTRIBUTION_TARGETS, we need to run fix-go-licenses first to correct
# module names. Module "github.com" should be "github.com/prometheus/prometheus".
$(ATTRIBUTION_TARGETS): fix-go-licenses

$(WEB_UI_TARGETS):
	@echo '>> npm install and npm build'
	cd prometheus/web/ui && npm install && CI="" npm run build
	@echo ">> bundling npm licenses"
	build/cp_npm_licenses.sh
	@echo '>> compressing assets'
	build/compress_assets.sh

$(EXTRA_DOCKER_FILE_TARGETS):
	mkdir -p $(OUTPUT_DIR)/extra/etc/prometheus
	cp prometheus/documentation/examples/prometheus.yml $(OUTPUT_DIR)/extra/etc/prometheus/prometheus.yml

	mkdir -p $(OUTPUT_DIR)/extra/usr/share/prometheus/console_libraries
	cp -r prometheus/console_libraries/* $(OUTPUT_DIR)/extra/usr/share/prometheus/console_libraries/

	mkdir -p $(OUTPUT_DIR)/extra/usr/share/prometheus/consoles
	cp -r prometheus/consoles/* $(OUTPUT_DIR)/extra/usr/share/prometheus/consoles/

	cp prometheus/npm_licenses.tar.bz2 $(OUTPUT_DIR)/npm_licenses.tar.bz2

.phony: fix-go-licenses
fix-go-licenses: $(GATHER_LICENSES_TARGETS)
	build/fix_go_licenses.sh "github.com,https" "github.com\/prometheus\/prometheus,https" $(OUTPUT_DIR)/attribution/go-license.csv

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

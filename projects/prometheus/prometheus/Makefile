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
# (1) run make prometheus/web/ui/embed.go to generate web UI related resources;
# (2) copy files needed by the docker container to the OUTPUT_DIR.
$(GO_MOD_DOWNLOAD_TARGETS): $(OUTPUT_DIR)/$(REPO)

$(OUTPUT_DIR)/$(REPO): prometheus/web/ui/embed.go
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

prometheus/web/ui/embed.go:
	@echo '>> install required node version'
	if [ $(shell node -v | cut -d. -f1 | cut -dv -f2) -lt 16 ]; then \
		curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.34.0/install.sh | bash && \
    	source ~/.nvm/nvm.sh && nvm install 16; \
	fi
	@echo '>> npm install and npm build'
	cd prometheus/web/ui && source ~/.nvm/nvm.sh && nvm use 16 && npm install && CI="" npm run build
	@echo ">> bundling npm licenses"
	yum install -y bzip2 && build/cp_npm_licenses.sh
	@echo '>> compressing assets'
	build/compress_assets.sh

.phony: fix-go-licenses
fix-go-licenses: $(GATHER_LICENSES_TARGETS)
	build/fix_go_licenses.sh "github.com,https" "github.com\/prometheus\/prometheus,https" $(OUTPUT_DIR)/attribution/go-license.csv

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

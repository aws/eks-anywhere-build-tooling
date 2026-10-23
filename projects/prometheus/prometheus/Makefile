BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
HELM_GIT_TAG=$(shell cat HELM_GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
PROMETHEUS_NODE_EXPORTER_GIT_TAG=$(shell cat ../node_exporter/GIT_TAG)
PROMETHEUS_NODE_EXPORTER_IMAGE_TAG=$(PROMETHEUS_NODE_EXPORTER_GIT_TAG)-$(GIT_HASH)
REPO=prometheus
REPO_OWNER=prometheus

BINARY_TARGET_FILES=prometheus promtool
BASE_IMAGE_NAME=eks-distro-minimal-base
SOURCE_PATTERNS=./cmd/prometheus ./cmd/promtool
# Add additional gobuild flag to embed the web UI in the go binary
EXTRA_GOBUILD_FLAGS=-tags netgo,builtinassets
GIT_COMMIT=$(shell git -C $(REPO) rev-parse --short HEAD)
GIT_BRANCH=$(shell git -C $(REPO) rev-parse --abbrev-ref HEAD)
BUILD_USER=ec2-user
BUILD_DATE=$(shell git -C $(REPO) show -s --format=format:%cd --date=format:'%Y-%m-%dT%H:%M:%SZ' HEAD)
EXTRA_GO_LDFLAGS=-X github.com/prometheus/common/version.Version=$(GIT_TAG)	\
	-X github.com/prometheus/common/version.Revision=$(GIT_COMMIT)	\
	-X github.com/prometheus/common/version.Branch=$(GIT_BRANCH)	\
	-X github.com/prometheus/common/version.BuildUser=$(BUILD_USER)	\
	-X github.com/prometheus/common/version.BuildDate=$(BUILD_DATE)

WEB_UI_TARGETS=prometheus/web/ui/embed.go
EXTRA_DOCKER_FILE_TARGETS=$(OUTPUT_DIR)/extra/etc/prometheus/prometheus.yml

FIX_LICENSES_CRONEXPR_TARGET=$(REPO)/vendor/github.com/hashicorp/cronexpr/LICENSE

EXCLUDE_FROM_STAGING_BUILDSPEC=true
SKIP_ON_RELEASE_BRANCH=true
BUILDSPEC_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

HAS_HELM_CHART=true
HELM_SOURCE_OWNER=prometheus-community
HELM_SOURCE_REPOSITORY=helm-charts
HELM_DIRECTORY=charts/prometheus
HELM_IMAGE_LIST=prometheus/prometheus prometheus/node-exporter
HELM_IMAGE_TAG_LIST=$(IMAGE_TAG) $(PROMETHEUS_NODE_EXPORTER_IMAGE_TAG)

HELM_CHART_NAMES=prometheus/charts/prometheus

include $(BASE_DIRECTORY)/Common.mk

# Override to use AL2023 builder-base for Node.js v24 compatibility
# AL2023 has GLIBC 2.34 which supports Node.js v24 (requires GLIBC >= 2.28)
# Pinned in-repo so generated buildspecs don't change when the prow cluster's
# builder-base tag moves.
CURRENT_BUILDER_BASE_TAG := $(shell cat $(BASE_DIRECTORY)/BUILDER_BASE_AL2023_TAG_FILE)
CHECKSUMS_BUILDSPEC_IMAGE=$(BASE_IMAGE_REPO)/builder-base:$(CURRENT_BUILDER_BASE_TAG)

# Before running GO_MOD_DOWNLOAD_TARGETS, we need to generate web UI related
# resources.
$(GO_MOD_DOWNLOAD_TARGETS): $(WEB_UI_TARGETS)

# Before running image relared targets, we need to copy files needed by the 
# docker container to the OUTPUT_DIR.
$(LOCAL_IMAGE_TARGETS): $(EXTRA_DOCKER_FILE_TARGETS)
$(IMAGE_TARGETS): $(EXTRA_DOCKER_FILE_TARGETS)

$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_CRONEXPR_TARGET)

$(FIX_LICENSES_CRONEXPR_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# hashicorp/cronexpr is dual licensed as GPL and Apache, since the files are not called
# the standard license file name it is not found by go-licenses
	cp $(@D)/APLv2 $@

# Matches upstream CI's pnpm major (web/ui uses pnpm-lock.yaml v9).
PNPM_VERSION=12.14.0

$(WEB_UI_TARGETS): NPM_AVAIL=$(shell command -v npm &> /dev/null && echo "true" || echo "false")
$(WEB_UI_TARGETS): | $$(call ENABLE_DOCKER,$$(NPM_AVAIL))
	@echo '>> Checking Node.js version inside container'; \
	NODEJS_BIN_PATH=""; \
	if command -v node &> /dev/null; then \
		NODE_MAJOR=$$(node --version | cut -d'.' -f1 | sed 's/v//'); \
		if [ "$$NODE_MAJOR" -ge "24" ]; then \
			echo "✅ Node.js $$(node --version) from builder-base (>= 24)"; \
			NODEJS_BIN_PATH=$$(dirname $$(which node)); \
		else \
			echo "⚠️  Node.js $$(node --version) too old, installing v24.20.0..."; \
		fi; \
	else \
		echo "⚠️  Node.js not found, installing v24.20.0..."; \
	fi; \
	if [ -z "$$NODEJS_BIN_PATH" ]; then \
		NODEJS_DIR=$$(pwd)/.nodejs-24; \
		rm -rf $$NODEJS_DIR; \
		mkdir -p $$NODEJS_DIR; \
		ARCH=$$(uname -m); \
		if [ "$$ARCH" = "x86_64" ]; then \
			NODEJS_ARCH="x64"; \
			CHECKSUM="855d581f8a4eb1a8117e3426de25fe02770592febcfb31369aee1ffbfee9e8ec"; \
		elif [ "$$ARCH" = "aarch64" ]; then \
			NODEJS_ARCH="arm64"; \
			CHECKSUM="3515603e2487879a39bc75716f1a2affd027500c64ba50e845cf72cb33219013"; \
		else \
			echo "❌ Unsupported architecture: $$ARCH"; \
			exit 1; \
		fi; \
		NODEJS_TARBALL="node-v24.20.0-linux-$$NODEJS_ARCH.tar.gz"; \
		wget --no-verbose https://nodejs.org/dist/v24.20.0/$$NODEJS_TARBALL -O /tmp/$$NODEJS_TARBALL; \
		echo "$$CHECKSUM  /tmp/$$NODEJS_TARBALL" | sha256sum -c -; \
		tar -C $$NODEJS_DIR --strip-components=1 -xzf /tmp/$$NODEJS_TARBALL; \
		chmod -R +x $$NODEJS_DIR/bin/*; \
		rm -f /tmp/$$NODEJS_TARBALL; \
		echo "✅ Node.js v24.20.0 installed to $$NODEJS_DIR"; \
		NODEJS_BIN_PATH="$$NODEJS_DIR/bin"; \
	fi; \
	export PATH="$$NODEJS_BIN_PATH:$$PATH"; \
	echo ">> Node version: $$(node --version)"; \
	echo ">> npm version: $$(npm --version)"; \
	echo ">> PATH: $$PATH"; \
	echo '>> Installing pnpm $(PNPM_VERSION) (the web UI uses pnpm since Prometheus v3.13)'; \
	PNPM_DIR=$$(pwd)/.pnpm-$(PNPM_VERSION); \
	npm install --global --no-fund --no-audit --prefix "$$PNPM_DIR" pnpm@$(PNPM_VERSION); \
	export PATH="$$PNPM_DIR/bin:$$PATH"; \
	PNPM_STORE=$$(pwd)/.pnpm-store; \
	echo ">> pnpm version: $$(pnpm --version)"; \
	echo '>> Building Prometheus Web UI with Node.js v24'; \
	echo '>> Installing workspace dependencies (mantine-ui, module/*)'; \
	pnpm -C prometheus/web/ui install --frozen-lockfile --store-dir "$$PNPM_STORE"; \
	echo '>> Installing react-app dependencies (not in workspace)'; \
	pnpm -C prometheus/web/ui/react-app install --frozen-lockfile --store-dir "$$PNPM_STORE"; \
	echo '>> Building all UI components'; \
	CI="" pnpm -C prometheus/web/ui run build; \
	echo ">> bundling npm licenses"; \
	build/cp_npm_licenses.sh; \
	echo '>> compressing assets'; \
	build/compress_assets.sh

$(EXTRA_DOCKER_FILE_TARGETS):
	mkdir -p $(OUTPUT_DIR)/extra/etc/prometheus
	cp prometheus/documentation/examples/prometheus.yml $(OUTPUT_DIR)/extra/etc/prometheus/prometheus.yml

	# Example console templates and libraries removed in Prometheus v3.0
	# Console feature still exists - users can supply their own templates via --web.console.* flags
	# See: CHANGELOG.md #14807 - "Remove example files for the console feature"
	# mkdir -p $(OUTPUT_DIR)/extra/usr/share/prometheus/console_libraries
	# cp -r prometheus/console_libraries/* $(OUTPUT_DIR)/extra/usr/share/prometheus/console_libraries/
	# mkdir -p $(OUTPUT_DIR)/extra/usr/share/prometheus/consoles
	# cp -r prometheus/consoles/* $(OUTPUT_DIR)/extra/usr/share/prometheus/consoles/

	cp prometheus/npm_licenses.tar.bz2 $(OUTPUT_DIR)/npm_licenses.tar.bz2


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

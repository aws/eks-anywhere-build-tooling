BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)

REGISTRYVERSION=2.8.3-patch-redis
TRIVYVERSION=v0.59.1
TRIVYADAPTERVERSION=v0.32.3


define VERSIONS_FOR_PREPARE
REGISTRY_VERSION: $(REGISTRYVERSION)
TRIVY_VERSION: $(TRIVYVERSION)
TRIVY_ADAPTER_VERSION: $(TRIVYADAPTERVERSION)
endef

export VERSIONS_FOR_PREPARE

REPO=harbor
REPO_OWNER=goharbor

DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

HARBOR_DB_IMAGE_COMPONENT=harbor/harbor-db
HARBOR_PORTAL_IMAGE_COMPONENT=harbor/harbor-portal
HARBOR_CORE_IMAGE_COMPONENT=harbor/harbor-core
HARBOR_LOG_IMAGE_COMPONENT=harbor/harbor-log
HARBOR_NGINX_IMAGE_COMPONENT=harbor/harbor-nginx
HARBOR_JOBSERVICE_IMAGE_COMPONENT=harbor/harbor-jobservice
HARBOR_REGISTRY_IMAGE_COMPONENT=harbor/harbor-registry
HARBOR_REGISTRYCTL_IMAGE_COMPONENT=harbor/harbor-registryctl
HARBOR_REDIS_IMAGE_COMPONENT=harbor/harbor-redis
HARBOR_EXPORTER_IMAGE_COMPONENT=harbor/harbor-exporter
HARBOR_TRIVY_IMAGE_COMPONENT=harbor/harbor-trivy

PKG_PATH=github.com/goharbor/harbor/src/pkg
GITCOMMIT=$(shell git -C $(REPO) rev-parse --short=8 HEAD)
RELEASEVERSION=$(shell cat $(REPO)/VERSION)

IMAGE_NAMES?=harbor-db harbor-portal harbor-core harbor-log harbor-nginx harbor-jobservice harbor-registry harbor-registryctl harbor-redis harbor-exporter harbor-trivy

BINARY_TARGET_FILES=harbor-core harbor-jobservice harbor-registryctl harbor-migrate harbor-exporter
SOURCE_PATTERNS=./core ./jobservice ./registryctl ./cmd/standalone-db-migrator ./cmd/exporter
# harbor-core has additional ldflags
BINARY_TARGET_FILES_BUILD_ALONE=harbor-core
GO_MOD_PATHS=src src src src src src

TRIVY_VERSION=$(TRIVYVERSION)
IMAGE_BUILD_ARGS=TRIVY_VERSION

PROJECT_DEPENDENCIES=eksa/goharbor/distribution eksa/aquasecurity/trivy eksa/goharbor/harbor-scanner-trivy

FIX_LICENSES_XEIPUUV_TARGET=$(REPO)/vendor/github.com/xeipuuv/gojsonpointer/LICENSE.txt
FIX_LICENSES_VOLC_SDK_GOLANG_TARGET=$(REPO)/src/vendor/github.com/volcengine/volc-sdk-golang/service/sts/LICENSE

HAS_HELM_CHART?=true

HELM_GIT_TAG=$(shell cat HELM_GIT_TAG)
HELM_SOURCE_REPOSITORY=harbor-helm
HELM_IMAGE_LIST=harbor/harbor-core harbor/harbor-db harbor/harbor-log harbor/harbor-exporter harbor/harbor-jobservice harbor/harbor-nginx harbor/harbor-portal harbor/harbor-redis harbor/harbor-registry harbor/harbor-registryctl harbor/harbor-trivy

HELM_CHART_NAMES=harbor/$(HELM_SOURCE_REPOSITORY)

EXCLUDE_FROM_STAGING_BUILDSPEC=true
BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE


SKIP_ON_RELEASE_BRANCH=true

include $(BASE_DIRECTORY)/Common.mk


harbor-nginx/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-nginx
harbor-portal/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-nginx
harbor-redis/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-glibc

$(GO_MOD_DOWNLOAD_TARGETS): $(OUTPUT_DIR)/$(REPO)

$(OUTPUT_DIR)/$(REPO): $(GIT_PATCH_TARGET)
	@echo "$$VERSIONS_FOR_PREPARE" > $(PROJECT_ROOT)/$(REPO)/make/photon/prepare/versions
	@mkdir -p $@
	@cp $(PROJECT_ROOT)/$(REPO)/LICENSE $(PROJECT_ROOT)/$(REPO)/src/

	@mkdir -p $(OUTPUT_DIR)/harbor-core/harbor
	@cp harbor/make/photon/common/install_cert.sh harbor/make/photon/core/entrypoint.sh harbor/make/photon/prepare/versions $(OUTPUT_DIR)/harbor-core/harbor/
	@cp -r harbor/src/core/views $(OUTPUT_DIR)/harbor-core/harbor
	@cp -r harbor/make/migrations $(OUTPUT_DIR)/harbor-core/harbor
	@cp -r harbor/icons $(OUTPUT_DIR)/harbor-core/harbor

	@mkdir -p $(OUTPUT_DIR)/harbor-db/docker-entrypoint-initdb.d
	@cp harbor/make/photon/db/docker-entrypoint.sh harbor/make/photon/db/initdb.sh harbor/make/photon/db/upgrade.sh harbor/make/photon/db/docker-healthcheck.sh $(OUTPUT_DIR)/harbor-db/
	@cp harbor/make/photon/db/initial-registry.sql $(OUTPUT_DIR)/harbor-db/docker-entrypoint-initdb.d/

	@mkdir -p $(OUTPUT_DIR)/harbor-exporter/harbor
	@cp harbor/make/photon/exporter/entrypoint.sh harbor/make/photon/common/install_cert.sh $(OUTPUT_DIR)/harbor-exporter/harbor/

	@mkdir -p $(OUTPUT_DIR)/harbor-jobservice/harbor
	@cp harbor/make/photon/common/install_cert.sh harbor/make/photon/jobservice/entrypoint.sh $(OUTPUT_DIR)/harbor-jobservice/harbor/

	@mkdir -p $(OUTPUT_DIR)/harbor-log/etc/rsyslog.d
	@mkdir -p $(OUTPUT_DIR)/harbor-log/etc/cron.hourly
	@mkdir -p $(OUTPUT_DIR)/harbor-log/usr/local/bin
	@cp harbor/make/photon/log/rsyslog.conf $(OUTPUT_DIR)/harbor-log/etc/
	@cp harbor/make/photon/log/rsyslog_docker.conf $(OUTPUT_DIR)/harbor-log/etc/rsyslog.d/
	@cp harbor/make/photon/log/logrotate $(OUTPUT_DIR)/harbor-log/etc/cron.hourly/
	@cp harbor/make/photon/log/start.sh $(OUTPUT_DIR)/harbor-log/usr/local/bin/

	@mkdir -p $(OUTPUT_DIR)/harbor-portal/build_dir/dist
	@cp -r harbor/src/portal/* $(OUTPUT_DIR)/harbor-portal/build_dir/
	@cp harbor/api/v2.0/swagger.yaml $(OUTPUT_DIR)/harbor-portal/build_dir/swagger.yaml
	@cp harbor/LICENSE $(OUTPUT_DIR)/harbor-portal/build_dir/dist/LICENSE

	@mkdir -p $(OUTPUT_DIR)/harbor-redis/etc
	@cp harbor/make/photon/redis/redis.conf $(OUTPUT_DIR)/harbor-redis/etc/

	@mkdir -p $(OUTPUT_DIR)/harbor-registry/home/harbor
	@cp harbor/make/photon/common/install_cert.sh harbor/make/photon/registry/entrypoint.sh $(OUTPUT_DIR)/harbor-registry/home/harbor/

	@mkdir -p $(OUTPUT_DIR)/harbor-registryctl/home/harbor
	@cp harbor/make/photon/common/install_cert.sh harbor/make/photon/registryctl/start.sh $(OUTPUT_DIR)/harbor-registryctl/home/harbor/

	@mkdir -p $(OUTPUT_DIR)/harbor-trivy/home/scanner/bin
	@cp harbor/make/photon/common/install_cert.sh harbor/make/photon/trivy-adapter/entrypoint.sh $(OUTPUT_DIR)/harbor-trivy/home/scanner/

%/harbor-core: EXTRA_GO_LDFLAGS=-X $(PKG_PATH)/version.GitCommit=$(GITCOMMIT) -X $(PKG_PATH)/version.ReleaseVersion=$(RELEASEVERSION)

ifeq ($(JOB_TYPE),presubmit)
# space is very limited in presubmit jobs, the image builds can push the total used space over the limit
# after the binaries are built there are ~3GBs of go build/mod cache
$(call IMAGE_TARGETS_FOR_NAME, harbor-portal): clean-repo
endif

$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_XEIPUUV_TARGET) $(FIX_LICENSES_VOLC_SDK_GOLANG_TARGET)

$(FIX_LICENSES_VOLC_SDK_GOLANG_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
	wget -q --retry-connrefused https://raw.githubusercontent.com/volcengine/volc-sdk-golang/main/LICENSE -O $(REPO)/src/vendor/github.com/volcengine/volc-sdk-golang/service/sts/LICENSE;
	wget -q --retry-connrefused https://raw.githubusercontent.com/volcengine/volc-sdk-golang/main/LICENSE -O $(REPO)/src/vendor/github.com/volcengine/volc-sdk-golang/base/LICENSE;

$(FIX_LICENSES_XEIPUUV_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The xeipuuv dependency github repos all have licenses however they all do not have go.mod files
# checked in to the repo. Hence we need to manually download licenses from Github for each of them 
# and place them in the respective folders under vendor directory so that they is available for 
# go-licenses to pick up	
	for package in gojsonpointer gojsonreference gojsonschema ; do \
		dest=$(REPO)/vendor/github.com/xeipuuv/$$package/LICENSE.txt; \
		mkdir -p $$(dirname $$dest); \
		wget --retry-connrefused -q https://raw.githubusercontent.com/xeipuuv/$$package/master/LICENSE-APACHE-2.0.txt -O \
			$$dest; \
	done;


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

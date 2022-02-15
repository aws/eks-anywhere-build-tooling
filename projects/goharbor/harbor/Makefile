BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"
REGISTRYVERSION=2.7.1-patch-2819-2553-redis
NOTARYVERSION=v0.6.1
TRIVYVERSION=v0.20.1
TRIVYADAPTERVERSION=v0.24.0
CHARTMUSEUM_SRC_TAG=v0.13.1
CHARTMUSEUMVERSION=$(CHARTMUSEUM_SRC_TAG)-redis

define VERSIONS_FOR_PREPARE
REGISTRY_VERSION: $(REGISTRYVERSION)
NOTARY_VERSION: $(NOTARYVERSION)
TRIVY_VERSION: $(TRIVYVERSION)
TRIVY_ADAPTER_VERSION: $(TRIVYADAPTERVERSION)
CHARTMUSEUM_VERSION: $(CHARTMUSEUMVERSION)
endef

export VERSIONS_FOR_PREPARE

REPO=harbor
REPO_OWNER=goharbor
REPO_SUBPATH=./src

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

PKG_PATH=github.com/goharbor/harbor/src/pkg
GITCOMMIT=$(shell git rev-parse --short=8 HEAD)
RELEASEVERSION=$(shell cat VERSION)

IMAGE_NAMES=harbor-db harbor-portal harbor-core harbor-log harbor-nginx harbor-jobservice harbor-registry harbor-registryctl harbor-redis harbor-exporter

BINARY_TARGET_FILES=harbor-core harbor-jobservice harbor-registryctl harbor-migrate-patch harbor-migrate harbor-exporter
SOURCE_PATTERNS=./core ./jobservice ./registryctl ./cmd/migrate-patch ./cmd/standalone-db-migrator ./cmd/exporter


PRECOMPILE_PREPARE=harbor_prepare

HARBOR_BASE_NAMESPACE=goharbor
IMAGE_BUILD_ARGS=HARBOR_BASE_NAMESPACE

FETCH_BINARIES_TARGETS=eksa/distribution/distribution

include $(BASE_DIRECTORY)/Common.mk

$(call IMAGE_TARGETS_FOR_NAME, harbor-db): harbor-db-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-portal): harbor-portal-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-core): harbor-core-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-log): harbor-log-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-nginx): harbor-nginx-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-jobservice): harbor-jobservice-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-registry): harbor-registry-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-registryctl): harbor-registryctl-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-redis): harbor-redis-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-exporter): harbor-exporter-useradd/images/export

IMAGE_USERADD_USER_ID=10000
IMAGE_USERADD_USER_NAME=harbor

COPY_FILE_LIST=harbor/make/photon/db/docker-entrypoint.sh harbor/make/photon/db/initdb.sh harbor/make/photon/db/upgrade.sh harbor/make/photon/db/docker-healthcheck.sh harbor/make/photon/db/initial-notaryserver.sql harbor/make/photon/db/initial-notarysigner.sql harbor/make/photon/db/initial-registry.sql harbor/make/photon/core/entrypoint.sh harbor/src/core/views harbor/make/migrations harbor/icons harbor/make/photon/prepare/versions harbor/api/v2.0/legacy_swagger.yaml harbor/api/v2.0/swagger.yaml harbor/api/swagger.yaml harbor/src/portal harbor/LICENSE harbor/make/photon/log/rsyslog.conf harbor/make/photon/log/rsyslog_docker.conf harbor/make/photon/log/logrotate harbor/make/photon/log/start.sh harbor/make/photon/common/install_cert.sh harbor/make/photon/jobservice/entrypoint.sh harbor/make/photon/registry/entrypoint.sh harbor/make/photon/registryctl/start.sh harbor/make/photon/redis/docker-healthcheck harbor/make/photon/redis/redis.conf harbor/make/photon/exporter/entrypoint.sh

harbor-db-useradd/images/export: IMAGE_USERADD_USER_NAME=postgres
harbor-db-useradd/images/export: IMAGE_USERADD_USER_ID=999
harbor-portal-useradd/images/export: IMAGE_USERADD_USER_NAME=nginx 
harbor-log-useradd/images/export: IMAGE_USERADD_USER_NAME=syslog
harbor-nginx-useradd/images/export: IMAGE_USERADD_USER_NAME=nginx
harbor-redis-useradd/images/export: IMAGE_USERADD_USER_NAME=redis

harbor_prepare:
	@echo "$$VERSIONS_FOR_PREPARE" > $(PROJECT_ROOT)/$(REPO)/make/photon/prepare/versions
	@make -C $(REPO) gen_apis
	@mkdir -p $(OUTPUT_DIR)
	@cp -r --parents $(COPY_FILE_LIST) $(OUTPUT_DIR)

PKG_PATH=github.com/goharbor/harbor/src/pkg
GITCOMMIT=$(shell git -C $(REPO) rev-parse --short=8 HEAD)
RELEASEVERSION=$(shell cat $(REPO)/VERSION)
%/harbor-core: EXTRA_GO_LDFLAGS=-X $(PKG_PATH)/version.GitCommit=$(GITCOMMIT) -X $(PKG_PATH)/version.ReleaseVersion=$(RELEASEVERSION)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

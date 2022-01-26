BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"
REGISTRYVERSION=v2.7.1-patch-2819-2553-redis
NOTARYVERSION=v0.7.0
TRIVYVERSION=v0.20.1
TRIVYADAPTERVERSION=v0.24.0
CHARTMUSEUM_SRC_TAG=v0.13.1
CHARTMUSEUMVERSION=$(CHARTMUSEUM_SRC_TAG)-redis
REGISTRY_SRC_TAG=v2.7.1

define VERSIONS_FOR_PREPARE
REGISTRY_VERSION: $(REGISTRYVERSION)
NOTARY_VERSION: $(NOTARYVERSION)
TRIVY_VERSION: $(TRIVYVERSION)
TRIVY_ADAPTER_VERSION: $(TRIVYADAPTERVERSION)
CHARTMUSEUM_VERSION: $(CHARTMUSEUMVERSION)
endef

REPO=harbor
REPO_OWNER=goharbor
REPO_SUBPATH=./src

BASE_IMAGE_NAME?=eks-distro-base
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

IMAGE_NAMES=harbor-db harbor-portal harbor-core harbor-log harbor-nginx harbor-jobservice harbor-registry harbor-registryctl harbor-redis harbor-exporter

BINARY_TARGET_FILES=harbor-core harbor-jobservice harbor-registry harbor-registryctl harbor-exporter
SOURCE_PATTERNS=./core ./jobservice ./registry ./registryctl ./cmd/exporter
EXTRA_GOBUILD_FLAGS=-mod=vendor

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

harbor-db-useradd/images/export: IMAGE_USERADD_USER_NAME=postgres IMAGE_USERADD_USER_ID=999
harbor-portal-useradd/images/export: IMAGE_USERADD_USER_NAME=nginx IMAGE_USERADD_USER_ID=10000
harbor-core-useradd/images/export: IMAGE_USERADD_USER_NAME=harbor IMAGE_USERADD_USER_ID=10000
harbor-log-useradd/images/export: IMAGE_USERADD_USER_NAME=syslog IMAGE_USERADD_USER_ID=10000
harbor-nginx-useradd/images/export: IMAGE_USERADD_USER_NAME=nginx IMAGE_USERADD_USER_ID=10000
harbor-jobservice-useradd/images/export: IMAGE_USERADD_USER_NAME=harbor IMAGE_USERADD_USER_ID=10000
harbor-registry-useradd/images/export: IMAGE_USERADD_USER_NAME=harbor IMAGE_USERADD_USER_ID=10000
harbor-registryctl-useradd/images/export: IMAGE_USERADD_USER_NAME=harbor IMAGE_USERADD_USER_ID=10000
harbor-redis-useradd/images/export: IMAGE_USERADD_USER_NAME=redis
harbor-exporter-useradd/images/export: IMAGE_USERADD_USER_NAME=harbor IMAGE_USERADD_USER_ID=10000

PRECOMPILE_PREPARE=harbor_prepare

DOCKERFILEPATH=$(PROJECT_ROOT)/$(REPO)/make/photon
DOCKERFILEPATH_REG=$(DOCKERFILEPATH)/registry
DOCKERFILEPATH_CHART_SERVER=$(DOCKERFILEPATH)/chartserver
MAKE_PREPARE_PATH=$(DOCKERFILEPATH)/prepare

harbor_base_namespace=goharbor
IMAGE_BUILD_ARGS=harbor_base_namespace

harbor_prepare: version_prepare swagger_prepare registry_prepare

export VERSIONS_FOR_PREPARE
version_prepare:
	@echo "$$VERSIONS_FOR_PREPARE" > $(MAKE_PREPARE_PATH)/versions

swagger_prepare:
	make -C $(REPO) gen_apis

registry_prepare:
	cd $(PROJECT_ROOT)/$(REPO)/$(REPO_SUBPATH) \
	&& git clone -b $(REGISTRY_SRC_TAG) https://github.com/distribution/distribution.git registry && cd registry \
	&& wget https://github.com/distribution/distribution/pull/2879.patch && git apply 2879.patch \
	&& wget https://github.com/distribution/distribution/pull/3370.patch && git apply 3370.patch \
	&& git apply $(DOCKERFILEPATH_REG)/redis.patch \
	&& cp ./cmd/registry/config-example.yml $(DOCKERFILEPATH_REG)/config.yml

include $(BASE_DIRECTORY)/Common.mk

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)

REPO=harbor-helm
REPO_OWNER=goharbor

HELM_DESTINATION_REPOSITORY=harbor/$(REPO)

HAS_HELM_CHART=true

HELM_ADDITIONAL_KEY_VALUES=\
  NGINX_IMAGE_DIGEST=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} harbor/harbor-nginx $(LATEST)) \
  PORTAL_IMAGE_DIGEST=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} harbor/harbor-portal $(LATEST)) \
  CORE_IMAGE_DIGEST=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} harbor/harbor-core $(LATEST)) \
  JOBSERVICE_IMAGE_DIGEST=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} harbor/harbor-jobservice $(LATEST)) \
  REGISTRY_IMAGE_DIGEST=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} harbor/harbor-registry $(LATEST)) \
  REGISTRYCTL_IMAGE_DIGEST=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} harbor/harbor-registryctl $(LATEST)) \
  CHARTMUSEUM_IMAGE_DIGEST="" \
  TRIVY_ADAPTER_IMAGE_DIGEST="" \
  NOTARY_SERVER_IMAGE_DIGEST="" \
  NOTARY_SIGNER_IMAGE_DIGEST="" \
  DATABASE_IMAGE_DIGEST=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} harbor/harbor-db $(LATEST)) \
  REDIS_IMAGE_DIGEST=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} harbor/harbor-redis $(LATEST)) \
  EXPORTER_IMAGE_DIGEST=$(shell $(BASE_DIRECTORY)/build/lib/image_shasum.sh ${IMAGE_REPO} harbor/harbor-exporter $(LATEST))

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

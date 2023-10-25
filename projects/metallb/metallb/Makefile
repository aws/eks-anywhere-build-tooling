BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=metallb
REPO_OWNER=metallb

BINARY_TARGET_FILES=controller speaker
SOURCE_PATTERNS=go.universe.tf/metallb/controller go.universe.tf/metallb/speaker


CONTROLLER_IMAGE_COMPONENT=metallb/controller
SPEAKER_IMAGE_COMPONENT=metallb/speaker

IMAGE_NAMES=controller speaker
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)


HAS_HELM_CHART=true
HELM_DIRECTORY=charts/metallb
HELM_DESTINATION_REPOSITORY=metallb/metallb
HELM_IMAGE_LIST=metallb/controller metallb/speaker

EXCLUDE_FROM_STAGING_BUILDSPEC=true
SKIP_ON_RELEASE_BRANCH=true

include $(BASE_DIRECTORY)/Common.mk

# Build helm chart
.PHONY: helm/build
helm/build: $(LICENSES_TARGETS_FOR_PREREQ)
helm/build: $(if $(filter true,$(REPO_NO_CLONE)),,$(HELM_GIT_CHECKOUT_TARGET))
helm/build: $(if $(wildcard $(PROJECT_ROOT)/helm/patches),$(HELM_GIT_PATCH_TARGET),)
	# Metallb requires 2 charts. Consider this a temporary work-around.
	$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) metallb/crds charts/crds $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_require.sh $(HELM_SOURCE_IMAGE_REPO) metallb/crds  $(OUTPUT_DIR) $(IMAGE_TAG) $(HELM_TAG) $(PROJECT_ROOT) $(LATEST) $(HELM_USE_UPSTREAM_IMAGE) "$(PROJECT_DEPENDENCIES)" "$(FORCE_JSON_SCHEMA_FILE)" $(HELM_IMAGE_LIST)
	$(BUILD_LIB)/helm_replace.sh metallb/crds $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) metallb/crds
	$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) $(HELM_DESTINATION_REPOSITORY) $(HELM_DIRECTORY) $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_require.sh $(HELM_SOURCE_IMAGE_REPO) $(HELM_DESTINATION_REPOSITORY) $(OUTPUT_DIR) $(IMAGE_TAG) $(HELM_TAG) $(PROJECT_ROOT) $(LATEST) $(HELM_USE_UPSTREAM_IMAGE) metallb-crds "$(FORCE_JSON_SCHEMA_FILE)" $(HELM_IMAGE_LIST)
	$(BUILD_LIB)/helm_replace.sh $(HELM_DESTINATION_REPOSITORY) $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) $(HELM_DESTINATION_REPOSITORY)

# Build helm chart and push to registry defined in IMAGE_REPO.
.PHONY: helm/push
helm/push: helm/build
	$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) $(HELM_DESTINATION_REPOSITORY) $(HELM_TAG) $(GIT_TAG) $(OUTPUT_DIR) $(LATEST)
	$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) metallb/crds $(HELM_TAG) $(GIT_TAG) $(OUTPUT_DIR) $(LATEST)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

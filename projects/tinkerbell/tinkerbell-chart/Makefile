BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell yq eval '.version' chart/Chart.yaml)

REPO_OWNER=tinkerbell
REPO=tinkerbell-chart

REPO_NO_CLONE=true

HELM_DIRECTORY=chart
HELM_SOURCE_REPOSITORY=.
HELM_DESTINATION_REPOSITORY=$(REPO_OWNER)/$(REPO)
IMAGE_TAG=$(GIT_TAG)
REQUIRES_CONFIG_FILE=false

BUILD_TARGETS=build-chart
RELEASE_TARGETS=publish-chart

include $(BASE_DIRECTORY)/Common.mk

.PHONY: build-chart
build-chart:
	$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) $(HELM_DESTINATION_REPOSITORY) $(HELM_DIRECTORY) $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) $(HELM_DESTINATION_REPOSITORY)

.PHONY: publish-chart
publish-chart: build-chart
	$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) $(HELM_DESTINATION_REPOSITORY) $(IMAGE_TAG) $(OUTPUT_DIR)

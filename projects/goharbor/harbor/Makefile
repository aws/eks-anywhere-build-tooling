BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=harbor
REPO_OWNER=goharbor
REPO_SUBPATH=./src

DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

FETCH_BINARIES_TARGETS=eksa/helm/chartmuseum

IMAGE_NAMES=harbor-chartserver

include $(BASE_DIRECTORY)/Common.mk

HARBOR_CHARTSERVER_IMAGE_COMPONENT=harbor/harbor-chartserver

$(call IMAGE_TARGETS_FOR_NAME, harbor-chartserver): harbor-chartserver-useradd/images/export

harbor-chartserver-useradd/images/export: scripts_prepare IMAGE_USERADD_USER_NAME=chart IMAGE_USERADD_USER_ID=10000

PRECOMPILE_PREPARE=swagger_prepare

HARBOR_BASE_NAMESPACE=goharbor
IMAGE_BUILD_ARGS=HARBOR_BASE_NAMESPACE

swagger_prepare:
	make -C $(REPO) gen_apis

scripts_prepare:
	@mkdir -p $(OUTPUT_DIR)
	@cp -r $(PROJECT_ROOT)/$(REPO)/make $(OUTPUT_DIR)

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

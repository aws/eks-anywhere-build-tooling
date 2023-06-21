REPO_OWNER=aws
REPO=eks-anywhere-packages
BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
PROJECT_COMMIT_HASH=$(shell git -C $(REPO) rev-parse --short HEAD)
TIME=$(shell git -C $(REPO) show -s --format=%ct HEAD)
GIT_TAG?=$(shell cat GIT_TAG)
GIT_HASH?=$(shell git -C $(BASE_DIRECTORY) rev-parse HEAD)
IMAGE_TAG=$(RELEASE_BRANCH)-$(GIT_TAG)-$(GIT_HASH)
HELM_TAG=$(RELEASE_BRANCH)-$(GIT_TAG)-$(GIT_HASH)
GOLANG_VERSION?="1.19"
GO_MOD_NAME=github.com/aws/eks-anywhere-packages/
GO_VAR_PREFIX=pkg/config.
GO_VARS=version=${GIT_TAG} gitHash=${PROJECT_COMMIT_HASH} buildTime=${TIME}
EXTRA_GO_LDFLAGS=$(foreach v,$(GO_VARS),-X ${GO_MOD_NAME}${GO_VAR_PREFIX}${v})
PROJECT_DEPENDENCIES=eksd/kubernetes/client eksa/aws/rolesanywhere-credential-helper eksa/kubernetes/cloud-provider-aws

HAS_RELEASE_BRANCHES=true
BUILDSPEC_VARS_KEYS=RELEASE_BRANCH
BUILDSPEC_VARS_VALUES=SUPPORTED_K8S_VERSIONS

# This project is a bit of odd ball when it comes to using the standard makefile targets
# - the image builds are release branch aware but the binary is the same across all versions

# force binaries to go to non-release branch bin folder
BINARIES_ARE_RELEASE_BRANCHED=false

# do not look for checksums in release branch folder, instead use project root
PROJECT_ROOT=$(MAKE_ROOT)

DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)
GO_MOD_PATHS=. ecrtokenrefresher credentialproviderpackage

IMAGE_COMPONENT=eks-anywhere-packages
ECRTOKENREFRESHER_IMAGE_COMPONENT=ecr-token-refresher
CREDENTIALPROVIDER_IMAGE_COMPONENT=credential-provider-package
IMAGE_NAMES=eks-anywhere-packages ecrtokenrefresher credentialprovider

BINARY_TARGET_FILES=package-manager ecrtokenrefresher credentialprovider
BASE_IMAGE_NAME=eks-distro-minimal-base
SOURCE_PATTERNS=. ./cmd/ecr-token-refresher ./cmd/aws-credential-provider

FIX_LICENSES_ECRTOKENREFRESHER_TARGET=$(REPO)/ecrtokenrefresher/LICENSE
FIX_LICENSES_CREDENTIALPROVIDERPACKAGE_TARGET=$(REPO)/credentialproviderpackage/LICENSE

HAS_HELM_CHART=true
HELM_DIRECTORY=charts/eks-anywhere-packages
HELM_IMAGE_LIST=eks-anywhere-packages ecr-token-refresher
CRED_HELM_IMAGE_LIST=credential-provider-package

BUILDSPEC_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

include $(BASE_DIRECTORY)/Common.mk

$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_ECRTOKENREFRESHER_TARGET) $(FIX_LICENSES_CREDENTIALPROVIDERPACKAGE_TARGET)

# Copy the License to each directory that has a unique go.mod path,
# and convert "license" to "LICENSE" for licenses under vendor folders.
$(FIX_LICENSES_ECRTOKENREFRESHER_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
	cp $(REPO)/LICENSE $(REPO)/ecrtokenrefresher/LICENSE
	mv $(REPO)/ecrtokenrefresher/vendor/github.com/josharian/intern/license.md $(REPO)/ecrtokenrefresher/vendor/github.com/josharian/intern/LICENSE.md

$(FIX_LICENSES_CREDENTIALPROVIDERPACKAGE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
	cp $(REPO)/LICENSE $(REPO)/credentialproviderpackage/LICENSE

# Build helm chart
.PHONY: helm/build
helm/build: $(LICENSES_TARGETS_FOR_PREREQ)
helm/build: $(if $(filter true,$(REPO_NO_CLONE)),,$(HELM_GIT_CHECKOUT_TARGET))
helm/build: $(if $(wildcard $(PROJECT_ROOT)/helm/patches),$(HELM_GIT_PATCH_TARGET),)
	# Pkg controller requires multiple charts. Consider this a temporary work-around (turned too permanent).
	$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) eks-anywhere-packages-migrations charts/eks-anywhere-packages-migrations $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_require.sh $(HELM_SOURCE_IMAGE_REPO) eks-anywhere-packages-migrations $(OUTPUT_DIR) $(IMAGE_TAG) $(HELM_TAG) $(PROJECT_ROOT) $(LATEST) $(HELM_USE_UPSTREAM_IMAGE) $(HELM_IMAGE_LIST)
	$(BUILD_LIB)/helm_replace.sh eks-anywhere-packages-migrations $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) eks-anywhere-packages-migrations

	$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) eks-anywhere-packages-crds charts/eks-anywhere-packages-crds $(OUTPUT_DIR)
	PACKAGE_DEPENDENCIES=eks-anywhere-packages-migrations $(BUILD_LIB)/helm_require.sh $(HELM_SOURCE_IMAGE_REPO) eks-anywhere-packages-crds $(OUTPUT_DIR) $(IMAGE_TAG) $(HELM_TAG) $(PROJECT_ROOT) $(LATEST) $(HELM_USE_UPSTREAM_IMAGE) $(HELM_IMAGE_LIST)
	$(BUILD_LIB)/helm_replace.sh eks-anywhere-packages-crds $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) eks-anywhere-packages-crds

	$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) $(HELM_DESTINATION_REPOSITORY) $(HELM_DIRECTORY) $(OUTPUT_DIR)
	PACKAGE_DEPENDENCIES=eks-anywhere-packages-crds $(BUILD_LIB)/helm_require.sh $(HELM_SOURCE_IMAGE_REPO) $(HELM_DESTINATION_REPOSITORY) $(OUTPUT_DIR) $(IMAGE_TAG) $(HELM_TAG) $(PROJECT_ROOT) $(LATEST) $(HELM_USE_UPSTREAM_IMAGE) $(HELM_IMAGE_LIST)
	$(BUILD_LIB)/helm_replace.sh $(HELM_DESTINATION_REPOSITORY) $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) $(HELM_DESTINATION_REPOSITORY)
	
	# Credential Provider Package lives in pkg controller repo. Workaround for its sepearate chart
	$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) credential-provider-package credentialproviderpackage/charts/credential-provider-package $(OUTPUT_DIR)
	FORCE_JSON_SCHEMA_FILE=$(PROJECT_ROOT)/helm/credential-provider-package.schema.json $(BUILD_LIB)/helm_require.sh $(HELM_SOURCE_IMAGE_REPO) credential-provider-package $(OUTPUT_DIR) $(IMAGE_TAG) $(HELM_TAG) $(PROJECT_ROOT) $(LATEST) $(HELM_USE_UPSTREAM_IMAGE) $(CRED_HELM_IMAGE_LIST)
	$(BUILD_LIB)/helm_replace.sh credential-provider-package $(OUTPUT_DIR)
	$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) credential-provider-package

# Build helm chart and push to registry defined in IMAGE_REPO.
.PHONY: helm/push
helm/push: helm/build
	$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) $(HELM_DESTINATION_REPOSITORY) $(HELM_TAG) $(GIT_TAG) $(OUTPUT_DIR) $(LATEST)
	$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) eks-anywhere-packages-crds $(HELM_TAG) $(GIT_TAG) $(OUTPUT_DIR) $(LATEST)
	$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) eks-anywhere-packages-migrations $(HELM_TAG) $(GIT_TAG) $(OUTPUT_DIR) $(LATEST)
	$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) credential-provider-package $(HELM_TAG) $(GIT_TAG) $(OUTPUT_DIR) $(LATEST)


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

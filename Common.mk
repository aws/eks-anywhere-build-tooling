# Disable built-in rules and variables
MAKEFLAGS+=--no-builtin-rules --warn-undefined-variables
.SHELLFLAGS:=-eu -o pipefail -c
.SUFFIXES:

RELEASE_ENVIRONMENT?=development
GIT_HASH?=$(shell git rev-parse HEAD)

COMPONENT=$(REPO_OWNER)/$(REPO)
MAKE_ROOT=$(BASE_DIRECTORY)/projects/$(COMPONENT)
BUILD_LIB=${BASE_DIRECTORY}/build/lib
OUTPUT_BIN_DIR?=$(OUTPUT_DIR)/bin/$(REPO)

#################### AWS ###########################
AWS_REGION?=us-west-2
AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
ARTIFACT_BUCKET?=my-s3-bucket
IMAGE_REPO?=$(if $(AWS_ACCOUNT_ID),$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com,localhost:5000)
####################################################

#################### LATEST TAG ####################
BRANCH_NAME?=main
LATEST_TAG?=latest
ifneq ("$(BRANCH_NAME)","main")
	LATEST_TAG=$(BRANCH_NAME)
endif
####################################################

#################### CODEBUILD #####################
ifdef CODEBUILD_SRC_DIR
	ARTIFACTS_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	CLONE_URL=https://git-codecommit.$(AWS_REGION).amazonaws.com/v1/repos/$(REPO_OWNER).$(REPO)
else
	ARTIFACTS_PATH?=$(MAKE_ROOT)/_output/tar
	CLONE_URL=https://github.com/$(COMPONENT).git	
endif
####################################################

#################### GIT ###########################
GIT_CHECKOUT_TARGET?=$(REPO)/eks-anywhere-checkout-$(GIT_TAG)
GIT_PATCH_TARGET?=$(REPO)/eks-anywhere-patched
####################################################

#################### RELEASE BRANCHES ##############
HAS_RELEASE_BRANCHES?=false
RELEASE_BRANCH?=
SUPPORTED_K8S_VERSIONS=$(shell yq e 'keys | .[]' $(BASE_DIRECTORY)/projects/kubernetes-sigs/image-builder/BOTTLEROCKET_OVA_RELEASES)
ifneq ($(RELEASE_BRANCH),)
	RELEASE_BRANCH_SUFFIX=/$(RELEASE_BRANCH)

	ARTIFACTS_PATH:=$(ARTIFACTS_PATH)/$(RELEASE_BRANCH_SUFFIX)
	PROJECT_ROOT?=$(MAKE_ROOT)$(RELEASE_BRANCH_SUFFIX)
	OUTPUT_DIR?=_output$(RELEASE_BRANCH_SUFFIX)

	# include release branch info in latest tag
	LATEST_TAG:=$(GIT_TAG)-$(LATEST_TAG)
else ifneq ($(and $(filter true,$(HAS_RELEASE_BRANCHES)), \
	$(filter-out build release upload-artifacts release-upload clean,$(MAKECMDGOALS))),)
	# if project has release branches and not calling one of the above targets
$(error When running targets for this project other than `build` or `release` a `RELEASE_BRANCH` is required)
else ifeq ($(HAS_RELEASE_BRANCHES),true)
	# project has release branches and one was not specified, trigger target for all
	BUILD_TARGETS=build/release-branches/all
	RELEASE_TARGETS=release/release-branches/all
	RELEASE_UPLOAD_TARGETS=release-upload/release-branches/all

	# avoid warnings when trying to read GIT_TAG file which wont exist when no release_branch is given
	GIT_TAG=non-existent
else
	PROJECT_ROOT?=$(MAKE_ROOT)
	OUTPUT_DIR?=_output
endif

####################################################

#################### BASE IMAGES ###################
BASE_IMAGE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-base
BASE_IMAGE_TAG_FILE?=$(BASE_DIRECTORY)/$(shell echo $(BASE_IMAGE_NAME) | tr '[:lower:]' '[:upper:]' | tr '-' '_')_TAG_FILE
BASE_IMAGE_TAG?=$(shell cat $(BASE_IMAGE_TAG_FILE))
BASE_IMAGE?=$(BASE_IMAGE_REPO)/$(BASE_IMAGE_NAME):$(BASE_IMAGE_TAG)
BUILDER_IMAGE?=$(BASE_IMAGE_REPO)/$(BASE_IMAGE_NAME)-builder:$(BASE_IMAGE_TAG)
####################################################

#################### IMAGES ########################
IMAGE_COMPONENT?=$(COMPONENT)
IMAGE_DESCRIPTION?=$(COMPONENT)
IMAGE_OUTPUT_DIR?=/tmp
IMAGE_CONTEXT_DIR?=.
IMAGE_OUTPUT_NAME?=$(IMAGE_NAME)
IMAGE_BUILD_ARGS?=
DOCKERFILE_FOLDER?=./docker/linux

# This tag is overwritten in the prow job to point to the upstream git tag and this repo's commit hash
IMAGE_TAG?=$(GIT_TAG)-$(GIT_HASH)

# For projects with multiple containers this is defined to override the default
# ex: CLUSTER_API_CONTROLLER_IMAGE_COMPONENT
IMAGE_COMPONENT_VARIABLE=$(shell echo '$(IMAGE_NAME)' | tr '[:lower:]' '[:upper:]' | tr '-' '_' )_IMAGE_COMPONENT
IMAGE=$(IMAGE_REPO)/$(if $(value $(IMAGE_COMPONENT_VARIABLE)),$(value $(IMAGE_COMPONENT_VARIABLE)),$(IMAGE_COMPONENT)):$(IMAGE_TAG)
LATEST_IMAGE=$(IMAGE:$(lastword $(subst :, ,$(IMAGE)))=$(LATEST_TAG))
####################################################

#################### BINARIES ######################
BINARY_PLATFORMS?=linux/amd64 linux/arm64
SIMPLE_CREATE_BINARIES?=true

BINARY_TARGETS?=$(call BINARY_TARGETS_FROM_FILES_PLATFORMS, $(BINARY_TARGET_FILES), $(BINARY_PLATFORMS))
BINARY_TARGET_FILES?=
SOURCE_PATTERNS?=.

#### BUILD FLAGS ####
EXTRA_GO_LDFLAGS?=
GO_LDFLAGS=-s -w -buildid= -extldflags -static $(EXTRA_GO_LDFLAGS)
GOBUILD_COMMAND?=build
EXTRA_GOBUILD_FLAGS?=
######################

#### HELPERS ########
# https://riptutorial.com/makefile/example/23643/zipping-lists
# Used to generate binary targets based on BINARY_TARGET_FILES
list-rem = $(wordlist 2,$(words $1),$1)
pairmap = $(and $(strip $2),$(strip $3),$(call \
    $1,$(firstword $2),$(firstword $3)) $(call \
    pairmap,$1,$(call list-rem,$2),$(call list-rem,$3)))
######################

####################################################

#################### LICENSES ######################
LICENSE_PACKAGE_FILTER?=
REPO_SUBPATH?=

ATTRIBUTION_TARGET=$(wildcard ATTRIBUTION.txt) $(wildcard $(RELEASE_BRANCH)/ATTRIBUTION.txt)
GATHER_LICENSES_TARGET=$(OUTPUT_DIR)/attribution/go-license.csv
####################################################

#################### TARBALLS ######################

# TODO: currently all projects push artifact tars to s3, change this to only those that are actually consumed
# etcdadm,kind,cri-tools
HAS_S3_ARTIFACTS?=false

SIMPLE_CREATE_TARBALLS?=true
TAR_FILE_PREFIX?=$(REPO)
FAKE_ARM_IMAGES_FOR_VALIDATION?=false
####################################################

#################### OTHER #########################
KUSTOMIZE_TARGET=$(OUTPUT_DIR)/kustomize
####################################################

#################### TARGETS FOR OVERRIDING ########
BUILD_TARGETS?=validate-checksums local-images generate-attribution $(if $(filter true,$(HAS_S3_ARTIFACTS)),s3-artifacts,)
RELEASE_TARGETS?=validate-checksums images $(if $(filter true,$(HAS_S3_ARTIFACTS)),s3-artifacts,)
RELEASE_UPLOAD_TARGETS?=release $(if $(filter true,$(HAS_S3_ARTIFACTS)),upload-artifacts,)
####################################################

define BUILDCTL
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=$(IMAGE_PLATFORMS) \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--opt build-arg:RELEASE_BRANCH=$(RELEASE_BRANCH) \
		$(foreach BUILD_ARG,$(IMAGE_BUILD_ARGS),--opt build-arg:$(BUILD_ARG)=$($(BUILD_ARG))) \
		--progress plain \
		--local dockerfile=$(DOCKERFILE_FOLDER) \
		--local context=$(IMAGE_CONTEXT_DIR) \
		--output type=$(IMAGE_OUTPUT_TYPE),oci-mediatypes=true,\"name=$(IMAGE),$(LATEST_IMAGE)\",$(IMAGE_OUTPUT)
endef 

define WRITE_LOCAL_IMAGE_TAG
	echo $(IMAGE_TAG) > $(IMAGE_OUTPUT_DIR)/$(IMAGE_OUTPUT_NAME).docker_tag
	echo $(IMAGE) > $(IMAGE_OUTPUT_DIR)/$(IMAGE_OUTPUT_NAME).docker_image_name	
endef

define IMAGE_TARGETS_FOR_NAME
	$(addsuffix /images/push, $(1)) $(addsuffix /images/amd64, $(1)) $(addsuffix /images/arm64, $(1))
endef

define BINARY_TARGETS_FROM_FILES_PLATFORMS
	$(foreach platform, $(2), $(foreach target, $(1), $(OUTPUT_BIN_DIR)/$(subst /,-,$(platform))/$(target)))
endef

define BINARY_TARGET_BODY_ALL_PLATFORMS
	$(eval $(foreach platform, $(BINARY_PLATFORMS), $(call BINARY_TARGET_BODY,$(platform),$(1),$(2))))
endef

define BINARY_TARGET_BODY
	$(OUTPUT_BIN_DIR)/$(subst /,-,$(1))/$(2): | $(if $(wildcard patches),$(GIT_PATCH_TARGET),) $(GIT_CHECKOUT_TARGET)
		$(BASE_DIRECTORY)/build/lib/simple_create_binaries.sh $$(MAKE_ROOT) \
			$$(MAKE_ROOT)/$(OUTPUT_BIN_DIR)/$(subst /,-,$(1))/$(2) $$(REPO) $$(GOLANG_VERSION) $$(GIT_TAG) $(1) $(3) \
			"$$(GOBUILD_COMMAND)" "$$(EXTRA_GOBUILD_FLAGS)" "$$(GO_LDFLAGS)" $$(REPO_SUBPATH)

endef

## --------------------------------------
## Help
## --------------------------------------
##@ Helpers
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-35s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Source repo + binary Targets

$(REPO):
	git clone $(CLONE_URL) $(REPO)

$(GIT_CHECKOUT_TARGET): | $(REPO)
	(cd $(REPO) && $(BASE_DIRECTORY)/build/lib/wait_for_tag.sh $(GIT_TAG))
	git -C $(REPO) checkout -f $(GIT_TAG)
	touch $@

$(GIT_PATCH_TARGET): $(GIT_CHECKOUT_TARGET)
	git -C $(REPO) config user.email prow@amazonaws.com
	git -C $(REPO) config user.name "Prow Bot"
	git -C $(REPO) am $(MAKE_ROOT)/patches/*
	@touch $@


ifeq ($(SIMPLE_CREATE_BINARIES),true)
$(call pairmap,BINARY_TARGET_BODY_ALL_PLATFORMS,$(BINARY_TARGET_FILES),$(SOURCE_PATTERNS))
endif

.PHONY: binaries
binaries: ## Build binaries by calling build/lib/simple_create_binaries.sh unless SIMPLE_CREATE_BINARIES=false, then calls build/create_binaries.sh from the project root.
binaries: $(BINARY_TARGETS)

$(KUSTOMIZE_TARGET):
	@mkdir -p $(OUTPUT_DIR)
	curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash -s -- $(OUTPUT_DIR)


## File/Folder Targets

$(OUTPUT_DIR)/images/%:
	@mkdir -p $(@D)

$(OUTPUT_DIR)/ATTRIBUTION.txt:
	@cp $(ATTRIBUTION_TARGET) $(OUTPUT_DIR)


##@ License Targets

## Gather licenses for project based on dependencies in REPO.
$(GATHER_LICENSES_TARGET): $(BINARY_TARGETS)
	$(BASE_DIRECTORY)/build/lib/gather_licenses.sh $(REPO) $(MAKE_ROOT)/$(OUTPUT_DIR) "$(LICENSE_PACKAGE_FILTER)" $(REPO_SUBPATH)

$(ATTRIBUTION_TARGET): $(GATHER_LICENSES_TARGET)
	$(BASE_DIRECTORY)/build/lib/create_attribution.sh $(MAKE_ROOT) $(GOLANG_VERSION) $(MAKE_ROOT)/$(OUTPUT_DIR) ATTRIBUTION.txt $(RELEASE_BRANCH)

.PHONY: gather-licenses
gather-licenses: ## Helper to call $(GATHER_LICENSES_TARGET) which gathers all licenses
gather-licenses: $(GATHER_LICENSES_TARGET)

.PHONY: generate-attribution
generate-attribution: ## Generates attribution from licenses gathered during `gather-licenses`.
generate-attribution: $(ATTRIBUTION_TARGET)

##@ Tarball Targets

.PHONY: tarballs
tarballs: ## Create tarballs by calling build/lib/simple_create_tarballs.sh unless SIMPLE_CREATE_TARBALLS=false, then calls build/create_tarballs.sh from project directory
tarballs: $(GATHER_LICENSES_TARGET) $(OUTPUT_DIR)/ATTRIBUTION.txt
ifeq ($(SIMPLE_CREATE_TARBALLS),true)
	$(BASE_DIRECTORY)/build/lib/simple_create_tarballs.sh $(TAR_FILE_PREFIX) $(MAKE_ROOT)/$(OUTPUT_DIR) $(MAKE_ROOT)/$(OUTPUT_BIN_DIR) $(GIT_TAG) "$(BINARY_PLATFORMS)" $(ARTIFACTS_PATH) $(GIT_HASH)
else
	build/create_tarballs.sh $(REPO) $(GIT_TAG) $(RELEASE_BRANCH)
endif

.PHONY: upload-artifacts
upload-artifacts:
	$(BASE_DIRECTORY)/build/lib/upload_artifacts.sh $(ARTIFACTS_PATH) $(ARTIFACTS_BUCKET) $(PROJECT_PATH) $(CODEBUILD_BUILD_NUMBER) $(GIT_HASH) $(LATEST_TAG)

.PHONY: s3-artifacts
s3-artifacts: tarballs
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)

##@ Checksum Targets
	
.PHONY: checksums
checksums: ## Update checksums file based on currently built binaries.
checksums: $(BINARY_TARGETS)
	$(BASE_DIRECTORY)/build/lib/update_checksums.sh $(MAKE_ROOT) $(PROJECT_ROOT) $(MAKE_ROOT)/$(OUTPUT_BIN_DIR)

.PHONY: validate-checksums
validate-checksums: ## Validate checksums of currently built binaries against checksums file.
validate-checksums: $(BINARY_TARGETS)
	$(BASE_DIRECTORY)/build/lib/validate_checksums.sh $(MAKE_ROOT) $(PROJECT_ROOT) $(MAKE_ROOT)/$(OUTPUT_BIN_DIR)

## Image Targets


#	IMAGE_NAME is dynamically set based on target prefix. \
#	BASE_IMAGE BUILDER_IMAGE RELEASE_BRANCH are automatically passed as build-arg(s) to buildctl. args: \
#	DOCKERFILE_FOLDER: folder containing dockerfile, defaults ./docker/linux \
#	IMAGE_BUILD_ARGS:  additional build-args passed to buildctl, set to name of variable defined in makefile \
#	IMAGE_CONTEXT_DIR: context directory for buildctl, default: .

.PHONY: %/images/push %/images/amd64 %/images/arm64
%/images/push %/images/amd64 %/images/arm64: IMAGE_NAME=$*

%/images/push: ## Build image using buildkit for all platforms, by default pushes to registry defined in IMAGE_REPO.
%/images/push: IMAGE_PLATFORMS?=linux/amd64,linux/arm64
%/images/push: IMAGE_OUTPUT_TYPE?=image
%/images/push: IMAGE_OUTPUT?=push=true
%/images/push: $(GATHER_LICENSES_TARGET) $(OUTPUT_DIR)/ATTRIBUTION.txt
	$(BUILDCTL)

.PHONY: helm/build
helm/build: ## Build helm chart
	$(BUILD_LIB)/helm_build.sh $(IMAGE_COMPONENT) $(IMAGE_TAG) $(IMAGE_DESCRIPTION)

.PHONY: helm/push
helm/push: ## Build helm chart and push to registry defined in IMAGE_REPO.
	$(BUILD_LIB)/helm_build.sh $(IMAGE_COMPONENT) $(IMAGE_TAG) $(IMAGE_DESCRIPTION) $(IMAGE_REPO)

%/images/amd64: ## Build image using buildkit only builds linux/amd64 and saves to local tar.
%/images/amd64: IMAGE_PLATFORMS?=linux/amd64

%/images/arm64: ## Build image using buildkit only builds linux/arm64 and saves to local tar.
%/images/arm64: IMAGE_PLATFORMS?=linux/arm64

%/images/amd64 %/images/arm64: IMAGE_OUTPUT_TYPE?=oci
%/images/amd64 %/images/arm64: IMAGE_OUTPUT?=dest=$(IMAGE_OUTPUT_DIR)/$(IMAGE_OUTPUT_NAME).tar

%/images/amd64: $(GATHER_LICENSES_TARGET) $(OUTPUT_DIR)/ATTRIBUTION.txt
	@mkdir -p $(IMAGE_OUTPUT_DIR)
	$(BUILDCTL)
	$(WRITE_LOCAL_IMAGE_TAG)

%/images/arm64: $(GATHER_LICENSES_TARGET) $(OUTPUT_DIR)/ATTRIBUTION.txt
	@mkdir -p $(IMAGE_OUTPUT_DIR)
	$(BUILDCTL)
	$(WRITE_LOCAL_IMAGE_TAG)

##@ Build Targets

.PHONY: build
build: ## Called via prow presubmit, calls `binaries gather-licenses clean-repo local-images generate-attribution checksums` by default
build: FAKE_ARM_IMAGES_FOR_VALIDATION=true
build: $(BUILD_TARGETS)

.PHONY: release
release: ## Called via prow postsubmit + release jobs, calls `binaries gather-licenses clean-repo images` by default
release: $(RELEASE_TARGETS)

.PHONY: release-upload
release-upload: $(RELEASE_UPLOAD_TARGETS)

.PHONY: %/release-branches/all
%/release-branches/all:
	@for version in $(SUPPORTED_K8S_VERSIONS) ; do \
		$(MAKE) $* RELEASE_BRANCH=$$version; \
	done;

##@ Clean Targets

.PHONY: clean-repo
clean-repo: ## Removes source directory
clean-repo:
	@rm -rf $(REPO)	

.PHONY: clean
clean: ## Removes source and _output directory
clean: clean-repo
	@rm -rf _output	

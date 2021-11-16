BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
BUILD_LIB=$(BASE_DIRECTORY)/build/lib

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
PROJECT_PATH?=projects/kubernetes-sigs/image-builder
PROJECT_DIRECTORY=$(BASE_DIRECTORY)/$(PROJECT_PATH)
OUTPUT_DIRECTORY=$(PROJECT_DIRECTORY)/_output
BOTTLEROCKET_DOWNLOAD_PATH?=$(PROJECT_DIRECTORY)/bottlerocket/downloads
CARGO_HOME=$(OUTPUT_DIRECTORY)/cargo
ARTIFACTS_BUCKET?=my-s3-bucket
GIT_TAG?=$(shell cat GIT_TAG)
BRANCH_NAME?=main
LATEST_TAG?=latest
ifneq ("$(BRANCH_NAME)","main")
	LATEST_TAG=$(BRANCH_NAME)
endif

REPO=image-builder
COMPONENT=kubernetes-sigs/$(REPO)
ifdef CODEBUILD_SRC_DIR
	PROJECT_PATH=projects/kubernetes-sigs/image-builder
	OVA_PATH=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/image-builder/images/capi/output
	CLONE_URL?=https://git-codecommit.$(AWS_REGION).amazonaws.com/v1/repos/kubernetes-sigs.image-builder
else
	OVA_PATH=image-builder/images/capi/output
	CLONE_URL?=https://github.com/$(COMPONENT).git
endif
IMAGE_BUILDER_DIR=$(REPO)/images/capi

VSPHERE_CONNECTION_DATA?={}
# Aws accounts to share built AMI with
DEV_ACCOUNTS?=

ifeq ("$(CLONE_URL)","")
	$(error No clone url was provided.)
endif

ifeq ("$(REPO)","")
	$(error No repository name was provided.)
endif

ifeq ("$(GIT_TAG)","")
	$(error No git tag was provided.)
endif

PACKER_AMI_VAR_FILES=$(PROJECT_DIRECTORY)/packer/ami/packer.json
PACKER_AMI_SHARE_FILE=$(PROJECT_DIRECTORY)/packer/ami/share-ami.json
PACKER_AMI_VAR_FILES+=$(PACKER_AMI_SHARE_FILE)

PACKER_OVA_CONF_FILE=$(PROJECT_DIRECTORY)/packer/ova/vsphere.json
PACKER_OVA_VAR_FILES=$(PACKER_OVA_CONF_FILE)
PACKER_OVA_VAR_FILES+=$(PROJECT_DIRECTORY)/packer/ova/vsphere-library.json

SUPPORTED_K8S_VERSIONS=$(shell yq e 'keys | .[]' $(PROJECT_DIRECTORY)/BOTTLEROCKET_OVA_RELEASES)
BUILD_TARGETS=$(addprefix build-ami-ubuntu-2004-, $(SUPPORTED_K8S_VERSIONS))
BUILD_TARGETS+=$(addprefix release-ova-bottlerocket-, $(SUPPORTED_K8S_VERSIONS))
RELEASE_TARGETS=$(addprefix release-ami-ubuntu-2004-, $(SUPPORTED_K8S_VERSIONS))
RELEASE_TARGETS+=$(addprefix release-ova-ubuntu-2004-, $(SUPPORTED_K8S_VERSIONS))
RELEASE_TARGETS+=$(addprefix release-ova-bottlerocket-, $(SUPPORTED_K8S_VERSIONS))

# Kubeadm hardcodes the version of the pause image it pulls, even tho containerd config
# overrides which image is actually used
# Setting up this mapping so we can build the additional image into the image
ADDITIONAL_PAUSE_1-18_FROM=1-18
ADDITIONAL_PAUSE_1-19_FROM=1-18
ADDITIONAL_PAUSE_1-20_FROM=1-19
ADDITIONAL_PAUSE_1-21_FROM=1-20

.PHONY: setup-ami-share
setup-ami-share:
	jq -n '{"ami_users": "$(DEV_ACCOUNTS)", "snapshot_users": "$(DEV_ACCOUNTS)"}' > $(PACKER_AMI_SHARE_FILE)

.PHONY: setup-image-builder
setup-image-builder:
	build/setup_image_builder.sh $(REPO) $(CLONE_URL) $(GIT_TAG)

.PHONY: setup-vsphere
setup-vsphere:
	echo $(VSPHERE_CONNECTION_DATA) > $(PACKER_OVA_CONF_FILE)

.PHONY: setup-bottlerocket
setup-bottlerocket:
	build/setup_bottlerocket.sh $(BOTTLEROCKET_DOWNLOAD_PATH) $(CARGO_HOME)

.PHONY: deps-ami
deps-ami: setup-image-builder
	make -C $(IMAGE_BUILDER_DIR) deps-ami

.PHONY: deps-ova
deps-ova: setup-image-builder
	make -C $(IMAGE_BUILDER_DIR) deps-ova

## --------------------------------------
## Help
## --------------------------------------
##@ Helpers
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-35s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


.PHONY: upload-artifacts-%
upload-artifacts-%:
	build/upload_artifacts.sh $(OVA_PATH) $(ARTIFACTS_BUCKET) $(PROJECT_PATH)/$* $(CODEBUILD_BUILD_NUMBER) $(CODEBUILD_RESOLVED_SOURCE_VERSION) $(LATEST_TAG)

.PHONY: setup-packer-configs-%
setup-packer-configs-%:
	build/setup_packer_configs.sh $* $(ARTIFACTS_BUCKET) $(OVA_PATH) $(ADDITIONAL_PAUSE_$*_FROM)

.PHONY: build-ami-ubuntu-2004-%
build-ami-ubuntu-2004-%: setup-ami-share deps-ami setup-packer-configs-%
	PACKER_VAR_FILES="$(OUTPUT_DIRECTORY)/$*/config/kubernetes.json $(PACKER_AMI_VAR_FILES) $(OUTPUT_DIRECTORY)/$*/config/common.json $(OUTPUT_DIRECTORY)/$*/config/cni.json $(OUTPUT_DIRECTORY)/$*/config/additional_components.json" \
		make -C $(IMAGE_BUILDER_DIR) validate-ami-ubuntu-2004

.PHONY: build
build: ## Validate configs for all versions of kubernetes
build: $(BUILD_TARGETS)

.PHONY: release-ami-ubuntu-2004-%
release-ami-ubuntu-2004-%: setup-ami-share deps-ami setup-packer-configs-%
	PACKER_VAR_FILES="$(OUTPUT_DIRECTORY)/$*/config/kubernetes.json $(PACKER_AMI_VAR_FILES) $(OUTPUT_DIRECTORY)/$*/config/common.json $(OUTPUT_DIRECTORY)/$*/config/cni.json $(OUTPUT_DIRECTORY)/$*/config/additional_components.json" \
		make -C $(IMAGE_BUILDER_DIR) build-ami-ubuntu-2004

.PHONY: release-ova-ubuntu-2004-%
release-ova-ubuntu-2004-%: deps-ova setup-vsphere setup-packer-configs-%
	PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_LOG_PATH=$(OVA_PATH)/packer.log PACKER_VAR_FILES="$(OUTPUT_DIRECTORY)/$*/config/kubernetes.json $(PACKER_OVA_VAR_FILES) $(OUTPUT_DIRECTORY)/$*/config/common.json $(OUTPUT_DIRECTORY)/$*/config/cni.json $(OUTPUT_DIRECTORY)/$*/config/additional_components.json" \
		OVF_CUSTOM_PROPERTIES="$(OUTPUT_DIRECTORY)/$*/config/ovf_custom_properties.json" \
		make -C $(IMAGE_BUILDER_DIR) build-node-ova-vsphere-ubuntu-2004
	build/get_artifacts.sh $(IMAGE_BUILDER_DIR) ubuntu

.PHONY: release-upload-ova-ubuntu-2004-%
release-upload-ova-ubuntu-2004-%: 
	$(MAKE) release-ova-ubuntu-2004-$*
	$(MAKE) upload-artifacts-$*

.PHONY: release-ova-bottlerocket-%
release-ova-bottlerocket-%:
release-ova-bottlerocket-%: setup-bottlerocket
	build/upload_bottlerocket_artifacts.sh $* bottlerocket $(BOTTLEROCKET_DOWNLOAD_PATH) $(CARGO_HOME) $(PROJECT_PATH)/$* $(LATEST_TAG)

.PHONY: release
release: ## Build and share amis and push ovas to s3 for all kubernetes versions
release: $(RELEASE_TARGETS)

.PHONY: validate-ubuntu-2004-%
validate-ubuntu-2004-%: check-env-validation setup-image-builder setup-packer-configs-%
	sh validate/validate.sh $(OUTPUT_DIRECTORY)/$*/config/validate_goss_inline_vars.json

.PHONY: check-env-validation
check-env-validation:
	ifndef TEMPLATE
		$(error Environment var TEMPLATE not set. Example: TEMPLATE=<ubuntu.livecd> make ...)
	endif

## --------------------------------------
## Document dynamic build targets
## --------------------------------------
##@ Builds
build-ami-ubuntu-2004-1-18: ## Validates configs for ubuntu 2004 based ami with Kubernetes 1-18
build-ami-ubuntu-2004-1-19: ## Validates configs for ubuntu 2004 based ami with Kubernetes 1-19
build-ami-ubuntu-2004-1-20: ## Validates configs for ubuntu 2004 based ami with Kubernetes 1-20
build-ami-ubuntu-2004-1-21: ## Validates configs for ubuntu 2004 based ami with Kubernetes 1-21

##@ Releases Ami
release-ami-ubuntu-2004-1-18: ## Builds and shares ubuntu 2004 based ami with Kubernetes 1-18
release-ami-ubuntu-2004-1-19: ## Builds and shares ubuntu 2004 based ami with Kubernetes 1-19
release-ami-ubuntu-2004-1-20: ## Builds and shares ubuntu 2004 based ami with Kubernetes 1-20
release-ami-ubuntu-2004-1-21: ## Builds and shares ubuntu 2004 based ami with Kubernetes 1-21

##@ Builds OVA for release
release-ova-ubuntu-2004-1-18: ## Builds ubuntu 2004 based ova with Kubernetes 1-18
release-ova-ubuntu-2004-1-19: ## Builds ubuntu 2004 based ova with Kubernetes 1-19
release-ova-ubuntu-2004-1-20: ## Builds ubuntu 2004 based ova with Kubernetes 1-20
release-ova-ubuntu-2004-1-21: ## Builds ubuntu 2004 based ova with Kubernetes 1-21

##@ Builds and uploads OVA
release-upload-ova-ubuntu-2004-1-18: ## Builds and pushes to s3 ubuntu 2004 based ova with Kubernetes 1-18
release-upload-ova-ubuntu-2004-1-19: ## Builds and pushes to s3 ubuntu 2004 based ova with Kubernetes 1-19
release-upload-ova-ubuntu-2004-1-20: ## Builds and pushes to s3 ubuntu 2004 based ova with Kubernetes 1-20
release-upload-ova-ubuntu-2004-1-21: ## Builds and pushes to s3 ubuntu 2004 based ova with Kubernetes 1-21
release-ova-bottlerocket-1-20: ## Pushes to s3 Bottlerocket based ova with Kubernetes 1-20
release-ova-bottlerocket-1-21: ## Pushes to s3 Bottlerocket based ova with Kubernetes 1-21
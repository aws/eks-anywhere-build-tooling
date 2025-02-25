BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)

REPO_OWNER=aws
REPO=eks-anywhere-packages

GO_MOD_NAME=github.com/aws/eks-anywhere-packages/
GO_VAR_PREFIX=pkg/config.
PROJECT_COMMIT_HASH=$(shell git -C $(REPO) rev-parse --short HEAD)
TIME=$(shell git -C $(REPO) show -s --format=%ct HEAD)
GO_VARS=version=${GIT_TAG} gitHash=${PROJECT_COMMIT_HASH} buildTime=${TIME}
EXTRA_GO_LDFLAGS=$(foreach v,$(GO_VARS),-X ${GO_MOD_NAME}${GO_VAR_PREFIX}${v})

PROJECT_DEPENDENCIES=eksd/kubernetes/client eksa/aws/rolesanywhere-credential-helper eksa/kubernetes/cloud-provider-aws/1-28

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
HELM_DIRECTORY?=charts/$(HELM_CHART_NAME)
HELM_IMAGE_LIST=eks-anywhere-packages ecr-token-refresher
CRED_HELM_IMAGE_LIST=credential-provider-package

HELM_CHART_NAMES=eks-anywhere-packages-migrations eks-anywhere-packages-crds eks-anywhere-packages credential-provider-package

BUILDSPEC_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

# building multiple helm charts requires overrides
eks-anywhere-packages-crds/helm/%: PACKAGE_DEPENDENCIES=eks-anywhere-packages-migrations

eks-anywhere-packages/helm/%: PACKAGE_DEPENDENCIES=eks-anywhere-packages-crds

credential-provider-package/helm/%: HELM_IMAGE_LIST=$(CRED_HELM_IMAGE_LIST)
credential-provider-package/helm/%: HELM_DIRECTORY=credentialproviderpackage/charts/credential-provider-package
credential-provider-package/helm/%: FORCE_JSON_SCHEMA_FILE=$(PROJECT_ROOT)/helm/credential-provider-package.schema.json

$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_ECRTOKENREFRESHER_TARGET) $(FIX_LICENSES_CREDENTIALPROVIDERPACKAGE_TARGET)

# Copy the License to each directory that has a unique go.mod path,
# and convert "license" to "LICENSE" for licenses under vendor folders.
$(FIX_LICENSES_ECRTOKENREFRESHER_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
	cp $(REPO)/LICENSE $(REPO)/ecrtokenrefresher/LICENSE
	mv $(REPO)/ecrtokenrefresher/vendor/github.com/josharian/intern/license.md $(REPO)/ecrtokenrefresher/vendor/github.com/josharian/intern/LICENSE.md

$(FIX_LICENSES_CREDENTIALPROVIDERPACKAGE_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
	cp $(REPO)/LICENSE $(REPO)/credentialproviderpackage/LICENSE



########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

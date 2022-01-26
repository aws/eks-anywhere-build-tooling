BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.14"

NOTARYMIGRATEVERSION=v4.11.0

REPO=notary
REPO_OWNER=notaryproject

BASE_IMAGE_NAME?=eks-distro-base

HARBOR_NOTARY_SERVER_IMAGE_COMPONENT=harbor/harbor-notary-server
HARBOR_NOTARY_SIGNER_IMAGE_COMPONENT=harbor/harbor-notary-signer

BINARY_TARGET_FILES=harbor-notary-server harbor-notary-signer harbor-migrate
SOURCE_PATTERNS=./cmd/notary-server ./cmd/notary-signer ./migrate/cli

IMAGE_NAMES=harbor-notary-server harbor-notary-signer

$(call IMAGE_TARGETS_FOR_NAME, harbor-notary-server): harbor-notary-server-useradd/images/export
$(call IMAGE_TARGETS_FOR_NAME, harbor-notary-signer): harbor-notary-signer-useradd/images/export

harbor-notary-server-useradd/images/export: IMAGE_USERADD_USER_NAME=notary IMAGE_USERADD_USER_ID=10000
harbor-notary-signer-useradd/images/export: IMAGE_USERADD_USER_NAME=notary IMAGE_USERADD_USER_ID=10000

PRECOMPILE_PREPARE=migration_prepare

harbor_base_namespace=goharbor
IMAGE_BUILD_ARGS=harbor_base_namespace

COMPONENT=goharbor/harbor-notary
CLONE_URL=https://github.com/$(REPO_OWNER)/$(REPO).git

# TODO: migrate dependency needs to be pulled out as it doesn't compile along with notary
migration_prepare: 
	cd $(PROJECT_ROOT)/$(REPO) \
	&& git clone -b $(NOTARYMIGRATEVERSION) https://github.com/golang-migrate/migrate migrate

include $(BASE_DIRECTORY)/Common.mk

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################
BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat ./$(RELEASE_BRANCH)/GIT_TAG)
GOLANG_VERSION=$(shell cat ./$(RELEASE_BRANCH)/GOLANG_VERSION)

REPO=cloud-provider-aws
REPO_OWNER=kubernetes

BINARY_TARGET_FILES=ecr-credential-provider
SOURCE_PATTERNS=./cmd/ecr-credential-provider
EXTRA_GO_LDFLAGS=-X k8s.io/component-base/version.gitVersion=$(shell cat ./$(RELEASE_BRANCH)/GIT_TAG)

HAS_S3_ARTIFACTS=true
HAS_RELEASE_BRANCHES=true
BUILDSPEC_VARS_KEYS=RELEASE_BRANCH
SUPPORTED_K8S_VERSIONS=1-34
BUILDSPEC_VARS_VALUES=SUPPORTED_K8S_VERSIONS

IMAGE_NAMES=


include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.18"

REPO=rolesanywhere-credential-helper
REPO_OWNER=aws

BINARY_TARGET_FILES=aws_signing_helper
SOURCE_PATTERNS=./cmd/aws_signing_helper
EXTRA_GO_LDFLAGS=-X k8s.io/component-base/version.gitVersion=$(GIT_TAG)

HAS_S3_ARTIFACTS=true

IMAGE_NAMES=


include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

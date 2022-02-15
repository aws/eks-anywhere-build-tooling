BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=distribution
REPO_OWNER=distribution

BINARY_TARGET_FILES=registry
SOURCE_PATTERNS=./cmd/registry

VERSION=$(shell git describe --match 'v[0-9]*' --dirty='.m' --always)
REVISION=$(shell git rev-parse HEAD)$(shell if ! git diff --no-ext-diff --quiet --exit-code; then echo .m; fi)
PKG=github.com/docker/distribution

EXTRA_GO_LDFLAGS=-X $(PKG)/version.Version=$(VERSION) -X $(PKG)/version.Revision=$(REVISION) -X $(PKG)/version.Package=$(PKG)
EXTRA_GOBUILD_FLAGS=-tags include_oss,include_gcs

HAS_S3_ARTIFACTS=true
IMAGE_NAMES=

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=distribution
REPO_OWNER=goharbor

BINARY_TARGET_FILES=registry
SOURCE_PATTERNS=./cmd/registry

VERSION=$(shell git -C $(REPO) describe --match 'v[0-9]*' --dirty='.m' --always)
REVISION=$(shell git -C $(REPO) rev-parse HEAD)$(shell if ! git -C $(REPO) diff --no-ext-diff --quiet --exit-code; then echo .m; fi)
PKG=github.com/docker/distribution

EXTRA_GO_LDFLAGS=-X $(PKG)/version.Version=$(VERSION) -X $(PKG)/version.Revision=$(REVISION) -X $(PKG)/version.Package=$(PKG)
EXTRA_GOBUILD_FLAGS=-tags include_oss,include_gcs

HAS_S3_ARTIFACTS=true
EXCLUDE_FROM_STAGING_BUILDSPEC=true
SKIP_ON_RELEASE_BRANCH=true
IMAGE_NAMES=

FIX_LICENSES_REDIS_TARGET=$(REPO)/vendor/github.com/garyburd/redigo/LICENSE.txt

include $(BASE_DIRECTORY)/Common.mk

$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_REDIS_TARGET)

$(FIX_LICENSES_REDIS_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The garyburd/redigo dependency github repo has a license however redis and internal are subfolders 
# without their own licenses. Hence we need to manually download parent license from Github for each of them 
# and place them in the respective folders under vendor directory so that they is available for 
# go-licenses to pick up	
	for package in redis internal ; do \
		dest=$(REPO)/vendor/github.com/garyburd/redigo/$$package/LICENSE.txt; \
		mkdir -p $$(dirname $$dest); \
		wget -q --retry-connrefused https://raw.githubusercontent.com/garyburd/redigo/master/LICENSE -O \
				$$dest; \
	done;


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

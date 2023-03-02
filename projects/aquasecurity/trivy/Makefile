BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.19"

REPO=trivy
REPO_OWNER=aquasecurity

BINARY_TARGET_FILES=trivy
SOURCE_PATTERNS=./cmd/trivy

VERSION=$(shell git -C $(REPO) describe --tags)

EXTRA_GO_LDFLAGS=-X main.version=$(VERSION)

HAS_S3_ARTIFACTS=true
IMAGE_NAMES=

# modernc.org/mathutil's license is actually mismatched
# the current version of go-license finds it to match BSD-2 where its
# really BSD-3.  The library go-liceses uses to determine the type
# has a newer version which more accurately matches between the two BSD types
# the confidence to BSD-2 is > .85 but under the .9 default whereas
# the newer version matching BSD-3 ends up being > .99 confidence
# We could remove this in the future when go-licenses upgrades
# in the meantime we have confirmed the license type and that it is acceptable
LICENSE_THRESHOLD=.85

EXCLUDE_FROM_STAGING_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk

$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_SHA1CD_TARGET)

$(FIX_LICENSES_SHA1CD_TARGET):
# The cgo dependency github repondoesn't have license however its parent repo does.
# Hence we need to manually download licenses from parent folder and place it in its
# respective folder under vendor directory so that it is available for go-licenses to pick up	
	for package in cgo ; do \
		dest=$(REPO)/vendor/github.com/pjbgf/sha1cd/$$package/LICENSE.txt; \
		mkdir -p $$(dirname $$dest); \
		wget --retry-connrefused -q https://raw.githubusercontent.com/pjbgf/sha1cd/main/LICENSE -O \
			$$dest; \
	done;


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

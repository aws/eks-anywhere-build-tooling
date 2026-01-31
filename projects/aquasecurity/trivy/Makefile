BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=trivy
REPO_OWNER=aquasecurity

BINARY_TARGET_FILES=trivy
SOURCE_PATTERNS=./cmd/trivy

VERSION=$(shell git -C $(REPO) describe --tags)

EXTRA_GO_LDFLAGS=-X main.version=$(VERSION)

# Trivy packages need jsonv2 in v0.68.2
# May remove this in the future when jsonv2 is not in GOEXPERIMENT mode
# Link for the fix: https://github.com/aquasecurity/trivy/discussions/9646
# Link for jsonv2: https://go.dev/blog/jsonv2-exp
export GOEXPERIMENT=jsonv2

HAS_S3_ARTIFACTS=true
IMAGE_NAMES=

BUILDSPEC_COMPUTE_TYPE=BUILD_GENERAL1_LARGE
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
SKIP_ON_RELEASE_BRANCH=true

include $(BASE_DIRECTORY)/Common.mk

FIX_LICENSES_CSAF_TARGET=$(REPO)/vendor/github.com/gocsaf/csaf/v3/LICENSE.txt

$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_CSAF_TARGET)

$(FIX_LICENSES_CSAF_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The csaf-poc dependency github repos all have licenses however they all do not have go.mod files
# checked in to the repo. Hence we need to manually download licenses from Github for each of them 
# and place them in the respective folders under vendor directory so that they is available for 
# go-licenses to pick up	
	wget -q --retry-connrefused https://raw.githubusercontent.com/gocsaf/csaf/main/LICENSES/Apache-2.0.txt -O \
		$(REPO)/vendor/github.com/gocsaf/csaf/v3/LICENSE.txt; \


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

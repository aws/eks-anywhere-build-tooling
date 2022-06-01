BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=trivy
REPO_OWNER=aquasecurity

BINARY_TARGET_FILES=trivy
SOURCE_PATTERNS=./cmd/trivy

VERSION=$(shell git -C $(REPO) describe --tags)

EXTRA_GO_LDFLAGS=-X main.version=$(VERSION)

HAS_S3_ARTIFACTS=true
IMAGE_NAMES=

FIX_LICENSES_MODERNC_TARGET=$(REPO)/vendor/modernc.org/LICENSE.txt

include $(BASE_DIRECTORY)/Common.mk

$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_MODERNC_TARGET)

$(FIX_LICENSES_MODERNC_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
	for package in cc/v3 ccgo/v3 libc mathutil memory opt sqlite strutil token ; do \
		dest=$(REPO)/vendor/modernc.org/$$package/LICENSE.txt; \
		cp $(REPO)/vendor/modernc.org/$$package/LICENSE $$dest; \
	done;

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

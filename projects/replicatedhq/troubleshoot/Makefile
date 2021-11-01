BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.14"

REPO=troubleshoot
REPO_OWNER=replicatedhq

LICENSE_PACKAGE_FILTER?=./cmd/troubleshoot

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/support-bundle
FIX_LICENSES_TARGET=$(REPO)/vendor/github.com/tj/go-spin/LICENSE

include $(BASE_DIRECTORY)/Common.mk


build: s3-artifacts
release: upload-artifacts

$(GATHER_LICENSES_TARGET): $(FIX_LICENSES_TARGET)

.PHONY: s3-artifacts
s3-artifacts: tarballs
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)

$(FIX_LICENSES_TARGET):
# The tj/go-spin dependency github repo has a license file however it does not have a go.mod file
# checked in to the repo. Hence we need to manually download the license from Github and place
# it in the respective folder under vendor directory so that they is available for go-licenses
# to pick up
	wget https://raw.githubusercontent.com/tj/go-spin/master/LICENSE -O $@

.PHONY: local-images
local-images: #noop

.PHONY: images
images: #noop
BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=node_exporter
REPO_OWNER=prometheus

BINARY_TARGET_FILES=node_exporter
BASE_IMAGE_NAME=eks-distro-minimal-base
EXTRA_GOBUILD_FLAGS=-tags netgo,osusergo,static_build

NODE_EXPORTER_IMAGE_COMPONENT=prometheus/node-exporter

FIX_LICENSES_DENNWC_TARGETS=$(REPO)/vendor/github.com/dennwc/btrfs/LICENSE $(REPO)/vendor/github.com/dennwc/ioctl/LICENSE

EXCLUDE_FROM_STAGING_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): | $(FIX_LICENSES_DENNWC_TARGETS)

$(FIX_LICENSES_DENNWC_TARGETS): | $(GO_MOD_DOWNLOAD_TARGETS)
  # The modules from the github org dennwc are properly licensed but for some
  # unknown reason are not included in the module download via go mod down
  # Manually downloading from github and placing in each of the packages
  # under vendor to make go-licenses happy
	for package in btrfs ioctl ; do \
		wget --retry-connrefused -q https://raw.githubusercontent.com/dennwc/$$package/master/LICENSE -O \
			$(REPO)/vendor/github.com/dennwc/$$package/LICENSE; \
	done;


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.14"

REPO=troubleshoot
REPO_OWNER=replicatedhq

BINARY_TARGET_FILES=support-bundle
SOURCE_PATTERNS=./cmd/troubleshoot
EXTRA_GO_LDFLAGS=-X github.com/replicatedhq/troubleshoot/pkg/version.version=$(GIT_TAG)	\
	-X github.com/replicatedhq/troubleshoot/pkg/version.gitSHA=$(shell git -C $(REPO) rev-list -n 1 $(GIT_TAG))
EXTRA_GOBUILD_FLAGS=-installsuffix netgo -tags netgo,containers_image_ostree_stub,exclude_graphdriver_devicemapper,exclude_graphdriver_btrfs,containers_image_openpgp
FIX_LICENSES_TARGET=$(REPO)/vendor/github.com/tj/go-spin/LICENSE

HAS_S3_ARTIFACTS=true

IMAGE_NAMES=

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_TARGET)

$(FIX_LICENSES_TARGET):
# The tj/go-spin dependency github repo has a license file however it does not have a go.mod file
# checked in to the repo. Hence we need to manually download the license from Github and place
# it in the respective folder under vendor directory so that they is available for go-licenses
# to pick up
	wget https://raw.githubusercontent.com/tj/go-spin/master/LICENSE -O $@

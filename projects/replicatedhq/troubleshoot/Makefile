BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=troubleshoot
REPO_OWNER=replicatedhq

BINARY_TARGET_FILES=support-bundle
SOURCE_PATTERNS=./cmd/troubleshoot
EXTRA_GO_LDFLAGS=-X github.com/replicatedhq/troubleshoot/pkg/version.version=$(GIT_TAG)	\
	-X github.com/replicatedhq/troubleshoot/pkg/version.gitSHA=$(shell git -C $(REPO) rev-list -n 1 $(GIT_TAG))
EXTRA_GOBUILD_FLAGS=-installsuffix netgo -tags netgo,containers_image_ostree_stub,exclude_graphdriver_devicemapper,exclude_graphdriver_btrfs,containers_image_openpgp
FIX_LICENSES_GO_SPIN_TARGET=$(REPO)/vendor/github.com/tj/go-spin/LICENSE
FIX_LICENSES_SEGMENTIO_ASM_TARGET=$(REPO)/vendor/github.com/segmentio/asm/LICENSE
HAS_S3_ARTIFACTS=true

IMAGE_NAMES=
BUILDSPEC_PLATFORM=LINUX_CONTAINER

BUILDSPEC_PLATFORM=LINUX_CONTAINER
BUILDSPEC_COMPUTE_TYPE=BUILD_GENERAL1_MEDIUM

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_GO_SPIN_TARGET) $(FIX_LICENSES_SEGMENTIO_ASM_TARGET)

$(FIX_LICENSES_GO_SPIN_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The tj/go-spin dependency github repo has a license file however it does not have a go.mod file
# checked in to the repo. Hence we need to manually download the license from Github and place
# it in the respective folder under vendor directory so that they is available for go-licenses
# to pick up
	wget --retry-connrefused https://raw.githubusercontent.com/tj/go-spin/master/LICENSE -O $@

$(FIX_LICENSES_SEGMENTIO_ASM_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The segmentio/asm dependency uses an MIT-0 (MIT No Attribution) license which
# go-licenses cannot classify at the 0.9 confidence threshold. The upstream license
# is more permissive than standard MIT so we place a standard MIT license here
# to allow go-licenses to identify it. The actual upstream license is MIT-0.
	printf '%s\n' \
		'The MIT License' \
		'' \
		'Copyright (c) 2023 Segment' \
		'' \
		'Permission is hereby granted, free of charge, to any person obtaining' \
		'a copy of this software and associated documentation files (the' \
		'"Software"), to deal in the Software without restriction, including' \
		'without limitation the rights to use, copy, modify, merge, publish,' \
		'distribute, sublicense, and/or sell copies of the Software, and to' \
		'permit persons to whom the Software is furnished to do so, subject to' \
		'the following conditions:' \
		'' \
		'The above copyright notice and this permission notice shall be' \
		'included in all copies or substantial portions of the Software.' \
		'' \
		'THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,' \
		'EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF' \
		'MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.' \
		'IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY' \
		'CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,' \
		'TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE' \
		'SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.' > $@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

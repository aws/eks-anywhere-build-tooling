BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=hegel
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=hegel
SOURCE_PATTERNS=./cmd/hegel

FIX_LICENSES_STRPTIME_TARGET=$(REPO)/vendor/github.com/pbnjay/strptime/LICENSE.txt

IMAGE_USERADD_USER_NAME=tinkerbell

include $(BASE_DIRECTORY)/Common.mk

$(call IMAGE_TARGETS_FOR_NAME, hegel): hegel-useradd/images/export

$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_STRPTIME_TARGET)

$(FIX_LICENSES_STRPTIME_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# One of the project dependencies github.com/pbnjay/strptime does not have a license file.
# The project owner had mentioned a MIT License and has it embedded in github.com/pbnjay/strptime/strptime.go.
# We copied the license out to STRPTIME_LICENSE.txt and moved it into vendor during build.
	cp STRPTIME_LICENSE.txt $(REPO)/vendor/github.com/pbnjay/strptime/LICENSE.txt


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

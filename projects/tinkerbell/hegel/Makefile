BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=hegel
REPO_OWNER=tinkerbell

LICENSE_PACKAGE_FILTER?=./cmd/hegel

BASE_IMAGE_NAME?=eks-distro-base

BINARY_TARGET_FILES=hegel
SOURCE_PATTERNS=./cmd/hegel

FIX_LICENSES_STRPTIME_TARGET=$(REPO)/vendor/github.com/pbnjay/strptime/LICENSE.txt

include $(BASE_DIRECTORY)/Common.mk

$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_STRPTIME_TARGET)

$(FIX_LICENSES_STRPTIME_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# One of the project dependencies github.com/pbnjay/strptime does not have a license file.
# The project owner had mentioned a MIT License and has it embeded in strptime.go file.
# We copied the license out to STRPTIME_LICENSE.txt and moved it into vendor during build.
	cp STRPTIME_LICENSE.txt $(REPO)/vendor/github.com/pbnjay/strptime/LICENSE.txt

.PHONY: local-images
local-images: hegel/images/amd64

.PHONY: images
images: hegel/images/push

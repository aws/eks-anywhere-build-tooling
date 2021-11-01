BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=notification-controller
REPO_OWNER=fluxcd

LICENSE_PACKAGE_FILTER?=.

BINARY_TARGET_FILES=notification-controller
EXTRA_GO_LDFLAGS=-extldflags -static

FIX_LICENSES_API_LICENSE_TARGET=$(REPO)/vendor/github.com/fluxcd/notification-controller/api/LICENSE
FIX_LICENSES_K0KUBUN_TARGET=$(REPO)/vendor/github.com/k0kubun/pp/LICENSE.txt

BASE_IMAGE_NAME=eks-distro-minimal-base

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGET): $(FIX_LICENSES_K0KUBUN_TARGET) $(FIX_LICENSES_API_LICENSE_TARGET)

.PHONY: local-images
local-images: notification-controller/images/amd64

.PHONY: images
images: notification-controller/images/push

$(FIX_LICENSES_API_LICENSE_TARGET): | $(GIT_CHECKOUT_TARGET)
# Internal go.mod under /api directory
	cp $(REPO)/LICENSE $@

$(FIX_LICENSES_K0KUBUN_TARGET): | $(GIT_CHECKOUT_TARGET)
# The notification-controller project consumes an older version (v2.3.0) of the k0kubun/pp module 
# which does not have a LICENSE file committed to the repo. Hence we are fetching the license file 
# from the master branch so that it is available in the vendor directory for go-licenses to pick up
	wget -q -O $@  https://raw.githubusercontent.com/k0kubun/pp/master/LICENSE.txt

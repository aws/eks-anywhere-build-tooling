BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION:=$(shell cat GOLANG_VERSION)
REPO=boots
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=boots
SOURCE_PATTERNS=./cmd/boots

GITREV?=$(shell git -C $(REPO) rev-parse --short HEAD)

EXTRA_GO_LDFLAGS=-X main.GitRev=${GITREV}

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

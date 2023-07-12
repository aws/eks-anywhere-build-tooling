# avoid finding the local fake git repo by looking one level up when looking for base_directory
BASE_DIRECTORY:=$(shell git -C .. rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION:=$(shell cat GOLANG_VERSION)

REPO=image-builder
REPO_OWNER=aws

BINARY_TARGET_FILES=image-builder
GO_MOD_PATHS=..
SIMPLE_CREATE_BINARIES=true
HAS_S3_ARTIFACTS=true
IMAGE_NAMES=

REPO_NO_CLONE=true
GIT_CHECKOUT_TARGET=main.go
FAKE_GIT_REPO_TARGET=.git

EXTRA_GO_LDFLAGS=-X github.com/aws/eks-anywhere-build-tooling/image-builder/cmd.version=${GIT_TAG}

# Go 1.18 automatically embeds vcs information if there is a .git folder in the current path
# this makes builds not reproducible due a sync timestamp 
EXTRA_GOBUILD_FLAGS=-buildvcs=false

include $(BASE_DIRECTORY)/Common.mk


$(REPO):
	@mkdir $@

$(GATHER_LICENSES_TARGETS): | $(FAKE_GIT_REPO_TARGET)

$(FAKE_GIT_REPO_TARGET):
	@git init
	@git remote add origin https://github.com/aws/eks-anywhere-build-tooling.git

.PHONY: clean-extra
clean-extra:
	@rm -rf vendor .git $(GO_MOD_DOWNLOAD_TARGETS) $(REPO)

clean: clean-extra

unit-test:
	go test ./...

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

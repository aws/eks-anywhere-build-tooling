# avoid finding the local fake git repo by looking one level up when looking for base_directory
BASE_DIRECTORY:=$(shell git -C .. rev-parse --show-toplevel)
GOLANG_VERSION?=1.18
GIT_TAG:=0.1.0-dev
BUILD_TOOLING_GIT_TAG=$(shell cat BUILD_TOOLING_GIT_TAG)

REPO=image-builder
REPO_OWNER=aws

BINARY_TARGET_FILES=image-builder
GO_MOD_PATHS=..
SIMPLE_CREATE_BINARIES=true

REPO_NO_CLONE=true
GIT_CHECKOUT_TARGET=main.go
FAKE_GIT_REPO_TARGET=.git

include $(BASE_DIRECTORY)/Common.mk

EXTRA_GO_LDFLAGS=-X github.com/aws/eks-anywhere-build-tooling/image-builder/cmd.version=${GIT_TAG} -X github.com/aws/eks-anywhere-build-tooling/image-builder/builder.buildToolingTag=${BUILD_TOOLING_GIT_TAG}

$(REPO):
	@mkdir $@

$(GATHER_LICENSES_TARGETS): | $(FAKE_GIT_REPO_TARGET)

$(ATTRIBUTION_TARGETS): GIT_TAG

GIT_TAG:
	echo $(GIT_TAG) > $@

$(FAKE_GIT_REPO_TARGET):
	@git init
	@git remote add origin https://github.com/aws/eks-anywhere-build-tooling.git

.PHONY: clean-extra
clean-extra:
	@rm -rf vendor .git GIT_TAG $(GO_MOD_DOWNLOAD_TARGETS) $(REPO)

clean: clean-extra


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

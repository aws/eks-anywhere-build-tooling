# avoid finding the local fake git repo by looking one level up when looking for base_directory
BASE_DIRECTORY=$(shell git -C .. rev-parse --show-toplevel)
GOLANG_VERSION?=1.16

REPO=bottlerocket-bootstrap
REPO_OWNER=aws

BASE_IMAGE_NAME?=eks-distro-base
IMAGE_COMPONENT=bottlerocket-bootstrap

BINARY_TARGET_FILES=bottlerocket-bootstrap
SOURCE_PATTERNS=.
REPO_SUBPATH=..

GIT_TAG=v$(RELEASE_BRANCH)-$(shell yq e ".$(RELEASE_BRANCH).number" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)

HAS_RELEASE_BRANCHES=true

# This project is a bit of odd ball when it comes to using the standard makefile targets
# - the image builds are release branch aware but the binary is the same across all versions
# - there is no clone of source since the source for the bootstrap exist in this repo already
# - go-licenses supports only HTTPS URLs so a temporary git repo is created in this project root
#   to add the https origin
# - REPO target is overridden to create empty folder

# force binarys to go to non-release branch bin folder
BINARIES_ARE_RELEASE_BRANCHED=false

# do not look for checksums in release branch folder, instead use project root
PROJECT_ROOT=$(MAKE_ROOT)

# fake git checkout
REPO_NO_CLONE=true
GIT_CHECKOUT_TARGET=main.go
FAKE_GIT_REPO_TARGET=.git

FETCH_BINARIES_TARGETS=eksd/kubernetes/client eksd/kubernetes/server eksa/kubernetes-sigs/etcdadm

include $(BASE_DIRECTORY)/Common.mk

$(REPO):
	@mkdir $@

$(call IMAGE_TARGETS_FOR_NAME, bottlerocket-bootstrap): $(call FULL_FETCH_BINARIES_TARGETS, $(FETCH_BINARIES_TARGETS))

$(GATHER_LICENSES_TARGETS): | $(FAKE_GIT_REPO_TARGET)

$(ATTRIBUTION_TARGETS): GIT_TAG

$(FAKE_GIT_REPO_TARGET):
	@git init
	@git remote add origin https://github.com/aws/eks-anywhere-build-tooling.git 

GIT_TAG:
	echo 0.0.0 > $@

.PHONY: clean-extra
clean-extra:
	@rm -rf vendor .git GIT_TAG

clean: clean-extra

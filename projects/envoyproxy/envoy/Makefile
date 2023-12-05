BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)

ENVOY_IMAGE=public.ecr.aws/appmesh/aws-appmesh-envoy:$(GIT_TAG)
IMAGE_BUILD_ARGS=ENVOY_IMAGE

REPO=envoy
REPO_OWNER=envoyproxy

BASE_IMAGE_NAME?=eks-distro-minimal-base-glibc

REPO_NO_CLONE=true
SIMPLE_CREATE_BINARIES=false
EXCLUDE_FROM_UPGRADE_BUILDSPEC=true

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

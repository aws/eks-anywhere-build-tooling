# Disable built-in rules and variables
MAKEFLAGS+=--no-builtin-rules --warn-undefined-variables
.SUFFIXES:
.SECONDEXPANSION:

RELEASE_ENVIRONMENT?=development

GIT_HASH=$(shell git -C $(BASE_DIRECTORY) rev-parse HEAD)
COMPONENT?=$(REPO_OWNER)/$(REPO)
MAKE_ROOT=$(BASE_DIRECTORY)/projects/$(COMPONENT)
PROJECT_PATH?=$(subst $(BASE_DIRECTORY)/,,$(MAKE_ROOT))
BUILD_LIB=$(BASE_DIRECTORY)/build/lib
OUTPUT_BIN_DIR?=$(OUTPUT_DIR)/bin/$(REPO)

SHELL_TRACE?=false
TRACE_SHELL=$(BUILD_LIB)/make_shell.sh trace true
LOGGING_SHELL=$(BUILD_LIB)/make_shell.sh log $(SHELL_TRACE)
DOCKER_SHELL=$(BUILD_LIB)/make_shell.sh docker $(SHELL_TRACE)
NOOP_SHELL=true
DEFAULT_SHELL=$(if $(filter true,$(SHELL_TRACE)),$(TRACE_SHELL),bash)
SHELL=$(DEFAULT_SHELL)
.SHELLFLAGS:=-eu -o pipefail -c
#################### AWS ###########################
AWS_REGION?=us-west-2
AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
ARTIFACTS_BUCKET?=s3://my-s3-bucket
IMAGE_REPO?=$(if $(AWS_ACCOUNT_ID),$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com,localhost:5000)
####################################################

#################### LATEST TAG ####################
# ensure local execution uses the 'main' or 'release-X' branch bundle
# similiar to https://github.com/aws/eks-anywhere/blob/main/Makefile#L32
# codebuild var
PULL_BASE_REF?=
BRANCH_NAME?=main
ifneq ($(PULL_BASE_REF),) # PULL_BASE_REF originates from prow
	BRANCH_NAME=$(PULL_BASE_REF)
endif

LATEST=latest
ifneq ($(BRANCH_NAME),main)
	LATEST=$(BRANCH_NAME)
endif

SKIP_ON_RELEASE_BRANCH?=false

# for some projects like the BR image build, we do not always have
# the artifacts avialable upstrea when adding a new kube version
NOT_SUPPORTED_RELEASE_BRANCH_CONFIGURATION?=false
####################################################

#################### CODEBUILD #####################
CODEBUILD_CI?=false
CI?=false
JOB_TYPE?=
INCLUDE_OUTPUT_IN_PROW_ARTIFACTS?=false
# prow artifacts location env var
ARTIFACTS?=
CODEBUILD_BUILD_IMAGE?=
CLONE_URL?=$(call GET_CLONE_URL,$(REPO_OWNER),$(REPO))
#HELM_CLONE_URL=$(call GET_CLONE_URL,$(HELM_SOURCE_OWNER),$(HELM_SOURCE_REPOSITORY))
HELM_CLONE_URL=https://github.com/$(HELM_SOURCE_OWNER)/$(HELM_SOURCE_REPOSITORY).git
ifeq ($(CODEBUILD_CI),true)
	ARTIFACTS_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
	UPLOAD_DRY_RUN=false
	BUILD_IDENTIFIER=$(CODEBUILD_BUILD_NUMBER)
else
	ARTIFACTS_PATH?=$(MAKE_ROOT)/_output/tar
	UPLOAD_DRY_RUN=$(if $(findstring postsubmit,$(JOB_TYPE)),false,true)
	ifeq ($(CI),true)
		BUILD_IDENTIFIER=$(PROW_JOB_ID)
	else
		BUILD_IDENTIFIER=$(shell date "+%F-%s")
	endif
endif
EXCLUDE_FROM_STAGING_BUILDSPEC?=false
EXCLUDE_FROM_CHECKSUMS_BUILDSPEC?=false
EXCLUDE_FROM_UPGRADE_BUILDSPEC?=false
DO_NOT_EXCLUDE_FROM_BUILDSPEC=false
BUILDSPECS?=buildspec.yml
CHECKSUMS_BUILDSPECS?=buildspecs/checksums-buildspec.yml
UPGRADE_BUILDSPECS?=buildspecs/upgrade-buildspec.yml
BUILDSPEC_VARS_KEYS?=
BUILDSPEC_VARS_VALUES?=
BUILDSPEC_PLATFORM?=ARM_CONTAINER
BUILDSPEC_COMPUTE_TYPE?=BUILD_GENERAL1_SMALL
BUILDSPECS_FOR_COMBINE_IMAGES=buildspec.yml buildspecs/combine-images.yml
BUILDSPEC_1_VARS_KEYS?=$(if $(findstring $(BUILDSPECS_FOR_COMBINE_IMAGES),$(BUILDSPECS)),IMAGE_PLATFORMS,)
BUILDSPEC_1_VARS_VALUES?=$(if $(findstring $(BUILDSPECS_FOR_COMBINE_IMAGES),$(BUILDSPECS)),IMAGE_PLATFORMS,)
BUILDSPEC_2_DEPENDS_ON_OVERRIDE?=$(if $(filter buildspecs/combine-images.yml,$(word 2,$(BUILDSPECS))),BUILDSPEC_1,)
####################################################

#################### GIT ###########################
GIT_CHECKOUT_TARGET?=$(REPO)/eks-anywhere-checkout-$(subst /,-,$(GIT_TAG))
GIT_PATCH_TARGET?=$(REPO)/eks-anywhere-patched
REPO_NO_CLONE?=false
PATCHES_DIR=$(or $(wildcard $(PROJECT_ROOT)/patches),$(wildcard $(MAKE_ROOT)/patches))
HELM_PATCHES_DIR=$(or $(wildcard $(PROJECT_ROOT)/helm/patches),$(wildcard $(MAKE_ROOT)/helm/patches))
REPO_SPARSE_CHECKOUT?=
####################################################

#################### RELEASE BRANCHES ##############
HAS_RELEASE_BRANCHES?=false
RELEASE_BRANCH?=
SUPPORTED_K8S_VERSIONS?=$(shell cat $(BASE_DIRECTORY)/release/SUPPORTED_RELEASE_BRANCHES)
# Comma-separated list of Kubernetes versions to skip building artifacts for
SKIPPED_K8S_VERSIONS?=
BINARIES_ARE_RELEASE_BRANCHED?=true
IS_RELEASE_BRANCH_BUILD=$(filter true,$(HAS_RELEASE_BRANCHES))
UNRELEASE_BRANCH_BINARY_TARGETS=patch-repo binaries attribution checksums validate-checksums
IS_UNRELEASE_BRANCH_TARGET=$(and $(filter false,$(BINARIES_ARE_RELEASE_BRANCHED)),$(filter $(UNRELEASE_BRANCH_BINARY_TARGETS) $(foreach target,$(UNRELEASE_BRANCH_BINARY_TARGETS),run-$(target)-in-docker run-in-docker/$(target)),$(MAKECMDGOALS)))
TARGETS_ALLOWED_WITH_NO_RELEASE_BRANCH?=
TARGETS_ALLOWED_WITH_NO_RELEASE_BRANCH+=build release clean clean-extra clean-go-cache help start-docker-builder stop-docker-builder create-ecr-repos all-attributions all-checksums all-attributions-checksums update-patch-numbers check-for-release-branch-skip run-buildkit-and-registry $(if $(filter false, $(HAS_LICENSES)),attribution,) $(if $(filter true, $(HAS_HELM_CHART)),,helm/push)
MAKECMDGOALS_WITHOUT_VAR_VALUE=$(foreach t,$(MAKECMDGOALS),$(if $(findstring var-value-,$(t)),,$(t)))
ifneq ($(and $(IS_RELEASE_BRANCH_BUILD),$(or $(RELEASE_BRANCH),$(IS_UNRELEASE_BRANCH_TARGET))),)
	RELEASE_BRANCH_SUFFIX=$(if $(filter true,$(BINARIES_ARE_RELEASE_BRANCHED)),/$(RELEASE_BRANCH),)

	ARTIFACTS_PATH:=$(ARTIFACTS_PATH)$(RELEASE_BRANCH_SUFFIX)
	OUTPUT_DIR?=_output$(RELEASE_BRANCH_SUFFIX)
	PROJECT_ROOT?=$(MAKE_ROOT)$(RELEASE_BRANCH_SUFFIX)
	ARTIFACTS_UPLOAD_PATH?=$(PROJECT_PATH)$(RELEASE_BRANCH_SUFFIX)

	# Deps are always released branched
	BINARY_DEPS_DIR?=_output/$(RELEASE_BRANCH)/dependencies

	# include release branch info in latest tag
	LATEST_TAG?=$(GIT_TAG)-$(LATEST)
else ifneq ($(and $(IS_RELEASE_BRANCH_BUILD), $(filter-out $(TARGETS_ALLOWED_WITH_NO_RELEASE_BRANCH),$(MAKECMDGOALS_WITHOUT_VAR_VALUE))),)
	# if project has release branches and not calling one of the above targets
$(error When running targets for this project other than `$(TARGETS_ALLOWED_WITH_NO_RELEASE_BRANCH)` a `RELEASE_BRANCH` is required)
else ifneq ($(IS_RELEASE_BRANCH_BUILD),)
	# project has release branches and one was not specified, trigger target for all
	# if BUILD_TARGETS or RELEASE_TARGETS are set via an env and we change them here
	# it will change the env var for the sub shell/make calls which is not what we want
	BUILD_TARGETS_OVERRIDE=build/release-branches/all
	RELEASE_TARGETS_OVERRIDE=release/release-branches/all

	# avoid warnings when trying to read GIT_TAG file which wont exist when no release_branch is given
	GIT_TAG=non-existent
	OUTPUT_DIR=non-existent
else
	PROJECT_ROOT?=$(MAKE_ROOT)
	ARTIFACTS_UPLOAD_PATH?=$(PROJECT_PATH)
	OUTPUT_DIR?=_output
	LATEST_TAG?=$(LATEST)
endif

####################################################

#################### BASE IMAGES ###################
BASE_IMAGE_REPO?=public.ecr.aws/eks-distro-build-tooling
BASE_IMAGE_NAME?=eks-distro-base
BASE_IMAGE_OS_VERSION?=al2
COMPILER_IMAGE_VERSION?=
BASE_IMAGE_TAG_FILE?=$(BASE_DIRECTORY)/$(call TO_UPPER,$(BASE_IMAGE_NAME))_$(if $(COMPILER_IMAGE_VERSION),$(COMPILER_IMAGE_VERSION)_,)$(if $(filter-out al2,$(BASE_IMAGE_OS_VERSION)),$(call TO_UPPER,$(BASE_IMAGE_OS_VERSION))_,)TAG_FILE
BASE_IMAGE_TAG?=$(shell cat $(BASE_IMAGE_TAG_FILE))
BASE_IMAGE?=$(BASE_IMAGE_REPO)/$(BASE_IMAGE_NAME):$(BASE_IMAGE_TAG)
BUILDER_IMAGE?=$(BASE_IMAGE_REPO)/$(BASE_IMAGE_NAME)-builder:$(BASE_IMAGE_TAG)
COMPILER_IMAGE?=$(BASE_IMAGE_REPO)/$(BASE_IMAGE_NAME:eks-distro-minimal-base-%=%):$(BASE_IMAGE_TAG)
EKS_DISTRO_BASE_IMAGE=$(BASE_IMAGE_REPO)/eks-distro-base:$(shell cat $(BASE_DIRECTORY)/EKS_DISTRO_BASE_TAG_FILE)
####################################################

#################### IMAGES ########################
IMAGE_COMPONENT?=$(COMPONENT)
IMAGE_OUTPUT_DIR?=/tmp
IMAGE_OUTPUT_NAME?=$(IMAGE_NAME)
IMAGE_TARGET?=

IMAGE_NAMES?=$(REPO)

# This tag is overwritten in the prow job to point to the upstream git tag and this repo's commit hash
IMAGE_TAG?=$(GIT_TAG)-$(GIT_HASH)
IMAGE_TAG_SUFFIX?=
# For projects with multiple containers this is defined to override the default
# ex: CLUSTER_API_CONTROLLER_IMAGE_COMPONENT
IMAGE_COMPONENT_VARIABLE=$(call TO_UPPER,$(IMAGE_NAME))_IMAGE_COMPONENT
IMAGE_REPO_COMPONENT=$(call IF_OVERRIDE_VARIABLE,$(IMAGE_COMPONENT_VARIABLE),$(IMAGE_COMPONENT))
IMAGE=$(IMAGE_REPO)/$(IMAGE_REPO_COMPONENT):$(IMAGE_TAG)$(IMAGE_TAG_SUFFIX)
LATEST_IMAGE=$(IMAGE:$(lastword $(subst :, ,$(IMAGE)))=$(LATEST_TAG))$(IMAGE_TAG_SUFFIX)

IMAGE_USERADD_USER_ID?=1000
IMAGE_USERADD_USER_NAME?=

# When building from outside the build account, such as prow, this is set to enable pulling the cache from the build acct
ADDITIONAL_IMAGE_CACHE_REPOS?=
# Branch builds should look at the current branch latest image for cache as well as main branch latest for cache to cover the cases
# where its the first build from a new release branch
IMAGE_CACHE_TAGS=$(foreach tag,$(LATEST_TAG) $(if $(filter latest,$(LATEST_TAG)),,latest),$(tag)$(IMAGE_TAG_SUFFIX) $(tag)$(IMAGE_TAG_SUFFIX)-cache)
IMAGE_IMPORT_CACHE?=$(foreach repo,$(strip $(IMAGE_REPO) $(ADDITIONAL_IMAGE_CACHE_REPOS)),$(foreach tag,$(IMAGE_CACHE_TAGS),type=registry,ref=$(repo)/$(IMAGE_REPO_COMPONENT):$(tag)))
IMAGE_EXPORT_CACHE?=--export-cache type=registry,mode=max,image-manifest=true,oci-mediatypes=true,ref=$(LATEST_IMAGE)-cache

BUILD_OCI_TARS?=false

LOCAL_IMAGE_TARGETS=$(foreach image,$(IMAGE_NAMES),$(image)/images/$(BUILDER_PLATFORM_ARCH))
IMAGE_TARGETS=$(foreach image,$(IMAGE_NAMES),$(if $(filter true,$(BUILD_OCI_TARS)),$(call IMAGE_TARGETS_FOR_NAME,$(image)),$(image)/images/push))

# intentionally not setting a default verison, we do not want projects depending on a default
GOLANG_VERSION?=
# If running in the builder base on prow or codebuild, grab the current tag to be used when building with cgo or docker
CURRENT_BUILDER_BASE_TAG=$(or \
	$(and $(wildcard /config/BUILDER_BASE_TAG_FILE),$(shell cat /config/BUILDER_BASE_TAG_FILE))\
	,$(shell curl -s https://raw.githubusercontent.com/aws/eks-anywhere-prow-jobs/main/BUILDER_BASE_TAG_FILE))
CURRENT_BUILDER_BASE_IMAGE=$(if $(CODEBUILD_BUILD_IMAGE),$(CODEBUILD_BUILD_IMAGE),$(BASE_IMAGE_REPO)/builder-base:$(CURRENT_BUILDER_BASE_TAG))
GOLANG_GCC_BUILDER_IMAGE=$(BASE_IMAGE_REPO)/golang:$(shell cat $(BASE_DIRECTORY)/EKS_DISTRO_MINIMAL_BASE_GOLANG_COMPILER_$(GOLANG_VERSION)_GCC_TAG_FILE)

# in CODEBUILD always use buildctl
BUILDCTL_AVAILABLE=$(or $(filter true,$(IS_ON_BUILDER_BASE)),$(shell command -v buildctl &> /dev/null && buildctl debug workers &> /dev/null && echo "true" || echo "false"))
BUILDX_AVAILABLE=$(shell command -v docker &> /dev/null && docker info &> /dev/null && docker buildx inspect &> /dev/null && echo "true" || echo "false")
DOCKER_AVAILABLE=$(shell command -v docker &> /dev/null && docker info &> /dev/null && echo "true" || echo "false")
####################################################

#################### HELM ##########################
HAS_HELM_CHART?=false
HELM_SOURCE_OWNER?=$(REPO_OWNER)
HELM_SOURCE_REPOSITORY?=$(REPO)
HELM_SOURCE_IMAGE_REPO?=$(IMAGE_REPO)
HELM_GIT_TAG?=$(GIT_TAG)
HELM_TAG?=$(GIT_TAG)-$(GIT_HASH)
HELM_USE_UPSTREAM_IMAGE?=false
# HELM_DIRECTORY must be a relative path from project root to the directory that contains a chart
HELM_DIRECTORY?=.
HELM_DESTINATION_REPOSITORY?=$(HELM_CHART_NAME)
HELM_CHART_FOLDER?=.
HELM_IMAGE_LIST?=$(IMAGE_COMPONENT)
HELM_IMAGE_TAG_LIST?=$(foreach _,$(HELM_IMAGE_LIST),$(IMAGE_TAG))
HELM_GIT_CHECKOUT_TARGET?=$(HELM_SOURCE_REPOSITORY)/eks-anywhere-checkout-$(subst /,-,$(HELM_GIT_TAG))
HELM_GIT_PATCH_TARGET?=$(HELM_SOURCE_REPOSITORY)/eks-anywhere-helm-patched
BUILD_HELM_DEPENDENCIES?=false
PACKAGE_DEPENDENCIES?=
FORCE_JSON_SCHEMA_FILE?=
HELM_CHART_NAMES?=$(IMAGE_COMPONENT)

# *theory* for the enable_logging preqreq when targets are prereqs of other targets
# for some reason the logging_shell does not get set for the targets which are acting as prereqs
# so if you call helm/build it will only show the logging around helm/build, not all
# by also defining all the actual targets (removing the wildcard) it seems to help make
# figure something out so that the logging shell is properly set
# this is only an issue in the newer version of Make running in AL2/23 vs on mac
# $1 - helm action (copy|require|replace|build|push)
FULL_CHART_TARGETS=$(addsuffix /helm/$1,$(HELM_CHART_NAMES))
####################################################

#### HELPERS ########
# https://riptutorial.com/makefile/example/23643/zipping-lists
# Used to generate binary targets based on BINARY_TARGET_FILES
list-rem = $(wordlist 2,$(words $1),$1)

pairmap = $(and $(strip $2),$(strip $3),$(call \
    $1,$(firstword $2),$(firstword $3)) $(call \
    pairmap,$1,$(call list-rem,$2),$(call list-rem,$3)))

trimap = $(and $(strip $2),$(strip $3),$(strip $4),$(call \
    $1,$(firstword $2),$(firstword $3),$(firstword $4)) $(call \
    trimap,$1,$(call list-rem,$2),$(call list-rem,$3),$(call list-rem,$4)))

_pos = $(if $(filter $1,$2),$(call _pos,$1,\
       $(call list-rem,$2),x $3),$3)
pos = $(words $(call _pos,$1,$2,))

# TODO: this exist in the gmsl, https://gmsl.sourceforge.io/
# look into introducting gmsl for things like this
# this function gets called a few dozen times and the alternative of using shell with tr takes
# noticeablely longer
TO_UPPER = $(subst a,A,$(subst b,B,$(subst c,C,$(subst d,D,$(subst e,E,$(subst \
	f,F,$(subst g,G,$(subst h,H,$(subst i,I,$(subst j,J,$(subst k,K,$(subst l,L,$(subst \
	m,M,$(subst n,N,$(subst o,O,$(subst p,P,$(subst q,Q,$(subst r,R,$(subst s,S,$(subst \
	t,T,$(subst u,U,$(subst v,V,$(subst w,W,$(subst x,X,$(subst y,Y,$(subst z,Z,$(subst -,_,$(1))))))))))))))))))))))))))))

TO_LOWER = $(subst A,a,$(subst B,b,$(subst C,c,$(subst D,d,$(subst E,e,$(subst \
	F,f,$(subst G,g,$(subst H,h,$(subst I,i,$(subst J,j,$(subst K,k,$(subst L,l,$(subst \
	M,m,$(subst N,n,$(subst O,o,$(subst P,p,$(subst Q,q,$(subst R,r,$(subst S,s,$(subst \
	T,t,$(subst U,u,$(subst V,v,$(subst W,w,$(subst X,x,$(subst Y,y,$(subst Z,z,$(subst _,-,$(1))))))))))))))))))))))))))))

# $1 - variable name to resolve and cache
CACHE_RESULT=$(if $(filter undefined,$(origin _cached-$1)),$(eval _cached-$1:=1)$(eval _cache-$1:=$($1)),)$(_cache-$1)

# $1 - variable name
CACHE_VARIABLE=$(eval _old-$(1)=$(value $(1)))$(eval $(1)=$$(call CACHE_RESULT,_old-$(1)))

# $1 - potential override variable name
# $2 - value if variable not set
# returns value of override var if one is set, otherwise returns $(2)
# intentionally no tab/space since it would come out in the result of calling this func
IF_OVERRIDE_VARIABLE=$(if $(filter undefined,$(origin $1)),$(2),$(value $(1)))

# $1 - image name
IMAGE_TARGETS_FOR_NAME=$(addsuffix /images/push, $(1)) $(addsuffix /images/amd64, $(1)) $(addsuffix /images/arm64, $(1))

# $1 - binary file name
FULL_FETCH_BINARIES_TARGETS=$(foreach platform,$(BINARY_PLATFORMS),$(addprefix $(BINARY_DEPS_DIR)/$(subst /,-,$(platform))/, $(1)))

# $1 - targets
# $2 - platforms
BINARY_TARGETS_FROM_FILES_PLATFORMS=$(foreach platform, $(2), $(foreach target, $(1), \
		$(OUTPUT_BIN_DIR)/$(subst /,-,$(platform))/$(if $(findstring windows,$(platform)),$(target).exe,$(target))))

# This "function" is used to construt the git clone URL for projects.
# Indenting the block results in the URL getting prefixed with a
# space, hence no indentation below.
# $1 - repo owner
# $2 - repo
GET_CLONE_URL=$(shell source $(BUILD_LIB)/common.sh && build::common::get_clone_url $(1) $(2) $(AWS_REGION) $(CODEBUILD_CI))

# $1 - binary file name
# $2 - go mod path for binary
# returns full target path for given binary + go mod path
# if the go mod path is `.` then do not prefix attribution dir, otherwise use binary name
LICENSE_TARGET_FROM_BINARY_GO_MOD=$(call LICENSE_OUTPUT_FROM_BINARY_GO_MOD,$(1),$(2))attribution/go-license.csv

# $1 - binary file name
# $2 - go mod path for binary
# return $1 if the go mod path is not the first, unless there is an override var for the binary
ATTRIBUTION_PREFIX_FROM_BINARY_GO_MOD=$(or \
	$(call IF_OVERRIDE_VARIABLE,$(call TO_UPPER,$(1))_ATTRIBUTION_OVERRIDE,), \
	$(if $(strip $(filter-out $(word 1,$(GO_MOD_PATHS)),$(2))),$(1),))

# $1 - binary file name
# $2 - go mod path for binary
# returns full path to create attribution/licenses directory
LICENSE_OUTPUT_FROM_BINARY_GO_MOD=$(LICENSES_OUTPUT_DIR)/$(call ADD_TRAILING_CHAR,$(call ATTRIBUTION_PREFIX_FROM_BINARY_GO_MOD,$(1),$(2)),/)

# $1 - binary file name
# $2 - go mod path for binary
# returns attribution target for given binary + go mod path
ATTRIBUTION_TARGET_FROM_BINARY_GO_MOD=$(if $(and $(IS_RELEASE_BRANCH_BUILD),$(filter \
	true,$(BINARIES_ARE_RELEASE_BRANCHED))),$(RELEASE_BRANCH)/,)$(call ADD_TRAILING_CHAR,$(call TO_UPPER,$(call ATTRIBUTION_PREFIX_FROM_BINARY_GO_MOD,$(1),$(2))),_)ATTRIBUTION.txt

# $1 - go mod path
GO_MOD_DOWNLOAD_TARGET_FROM_GO_MOD_PATH=$(REPO)/$(if $(filter-out .,$(1)),$(1)/,)eks-anywhere-go-mod-download

# $1 - binary file name
GO_MOD_TARGET_FOR_BINARY_VAR_NAME= \
	GO_MOD_TARGET_FOR_BINARY_$(call TO_UPPER,$(call IF_OVERRIDE_VARIABLE,$(call TO_UPPER,$(1))_ATTRIBUTION_OVERRIDE,$(1)))

# $1 - value
# $2 - char
# if value is non empty, add trailing $2
# intentionally no tab/space since it would come out in the result of calling this func
ADD_TRAILING_CHAR=$(if $(1),$(1)$(2),)

# check if pass variable has length of 1
IS_ONE_WORD=$(if $(filter 1,$(words $(1))),true,false)

SED_CMD=$(shell source $(BUILD_LIB)/common.sh && build::find::gnu_variant_on_mac sed)

# creating a space character by using a empty var seperated by space
_EMPTY_VAR=
SPACE=$(_EMPTY_VAR) $(_EMPTY_VAR)
####################################################

#################### BINARIES ######################
# if the pattern ends in the same as a previous pattern, binary must be built seperately
# if the go mod path has changed from the main, must be built seperately
# if binary is already in the BINARY_TARGET_FILES_BUILD_ALONE list do not add, but properly add source pattern and go mod
# $1 - binary file name
# $2 - source pattern
# $3 - go mod path for binary
setup_build_alone_vs_together = \
	$(eval type:=$(if $(or \
			$(call IF_OVERRIDE_VARIABLE,_UNIQ_PATTERN_$(notdir $(2)),), \
			$(filter-out $(word 1,$(GO_MOD_PATHS)),$(3)), \
			$(filter $(1),$(BINARY_TARGET_FILES_BUILD_ALONE))) \
		,ALONE,TOGETHER)) \
	$(if $(filter $(1),$(BINARY_TARGET_FILES_BUILD_ALONE)),,$(eval BINARY_TARGET_FILES_BUILD_$(type)+=$(1))) \
	$(eval SOURCE_PATTERNS_BUILD_$(type)+=$(2)) \
	$(eval GO_MOD_PATHS_BUILD_$(type)+=$(3)) \
	$(eval _UNIQ_PATTERN_$(notdir $(2)):=1)

# Setup vars UNIQ_GO_MOD_PATHS UNIQ_GO_MOD_TARGET_FILES
# which will store the mapping of uniq go_mod paths to first target file for repsective go mod
# $1 - binary file name
# $2 - source pattern
# $3 - go mod path for binary
setup_uniq_go_mod_license_filters = \
	$(if $(call IF_OVERRIDE_VARIABLE,GO_MOD_$(subst /,_,$(3))_LICENSE_PACKAGE_FILTER,),, \
			$(eval UNIQ_GO_MOD_PATHS+=$(3)) \
			$(eval UNIQ_GO_MOD_TARGET_FILES+=$(1))) \
			$(eval $(call GO_MOD_TARGET_FOR_BINARY_VAR_NAME,$(1))=$(3)) \
	$(eval GO_MOD_$(subst /,_,$(3))_LICENSE_PACKAGE_FILTER+=$(call IF_OVERRIDE_VARIABLE,LICENSE_PACKAGE_FILTER,$(2)))

BINARY_PLATFORMS?=linux/amd64 linux/arm64
SIMPLE_CREATE_BINARIES?=true

BINARY_TARGETS?=$(call BINARY_TARGETS_FROM_FILES_PLATFORMS, $(BINARY_TARGET_FILES), $(BINARY_PLATFORMS))
BINARY_TARGET_FILES?=
SOURCE_PATTERNS?=$(foreach _,$(BINARY_TARGET_FILES),.)
GO_MOD_PATHS?=$(foreach _,$(BINARY_TARGET_FILES),.)

# There may not any that need building alone, defining empty vars in case not set from above
BINARY_TARGET_FILES_BUILD_ALONE?=
SOURCE_PATTERNS_BUILD_ALONE?=
GO_MOD_PATHS_BUILD_ALONE?=
UNIQ_GO_MOD_PATHS?=
$(call trimap,setup_build_alone_vs_together,$(BINARY_TARGET_FILES),$(SOURCE_PATTERNS),$(GO_MOD_PATHS))
$(call trimap,setup_uniq_go_mod_license_filters,$(BINARY_TARGET_FILES),$(SOURCE_PATTERNS),$(GO_MOD_PATHS))

GO_MOD_DOWNLOAD_TARGETS?=$(foreach path, $(UNIQ_GO_MOD_PATHS), $(call GO_MOD_DOWNLOAD_TARGET_FROM_GO_MOD_PATH,$(path)))

VENDOR_UPDATE_SCRIPT?=
#### CGO ############
CGO_CREATE_BINARIES?=false
IS_ON_BUILDER_BASE=$(if $(wildcard /buildkit.sh),true,false)
BUILDER_PLATFORM_OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
BUILDER_PLATFORM_ARCH=$(if $(filter x86_64,$(shell uname -m)),amd64,arm64)
BUILDER_PLATFORM=$(BUILDER_PLATFORM_OS)/$(BUILDER_PLATFORM_ARCH)
NEEDS_CGO_BUILDER=$(and $(if $(filter true,$(CGO_CREATE_BINARIES)),true,),$(if $(filter true,$(IS_ON_BUILDER_BASE)),,true))
USE_DOCKER_FOR_CGO_BUILD?=false
GO_MOD_CACHE=$(shell if source $(BUILD_LIB)/common.sh && build::common::use_go_version $(GOLANG_VERSION) > /dev/null 2>&1 && command -v go &> /dev/null; then go env GOMODCACHE; else echo $${HOME}/.cache/go/pkg/mod; fi)
GO_BUILD_CACHE=$(shell if source $(BUILD_LIB)/common.sh && build::common::use_go_version $(GOLANG_VERSION) > /dev/null 2>&1 && command -v go &> /dev/null; then go env GOCACHE; else echo $${HOME}/.cache/go-build; fi)
GO_MODS_VENDORED?=false
DOCKER_PLATFORM?=
######################

#### BUILD FLAGS ####
CGO_ENABLED=$(if $(filter true,$(CGO_CREATE_BINARIES)),1,0)
GO_LDFLAGS?=$(if $(filter true,$(CGO_CREATE_BINARIES)),-s -w -buildid= $(EXTRA_GO_LDFLAGS),-s -w -buildid= -extldflags -static $(EXTRA_GO_LDFLAGS))
CGO_BUILD_ID=-Wl,--build-id=none
CGO_LDFLAGS?=$(if $(filter true,$(CGO_CREATE_BINARIES)),$(CGO_BUILD_ID),)
EXTRA_GOBUILD_FLAGS?=$(if $(filter true,$(CGO_CREATE_BINARIES)),-gcflags=-trimpath=$(MAKE_ROOT) -asmflags=-trimpath=$(MAKE_ROOT),)
CGO_CFLAGS_ALLOW?=

EXTRA_GO_LDFLAGS?=
GOBUILD_COMMAND?=build
######################

############### BINARIES DEPS ######################
BINARY_DEPS_DIR?=$(OUTPUT_DIR)/dependencies
PROJECT_DEPENDENCIES?=
HANDLE_DEPENDENCIES_TARGET?=handle-dependencies

# Based on PROJECT_DEPENDENCIES, generate fetch binaries targets, only projects with s3 artifacts will be fetched
PROJECT_DEPENDENCIES_TARGETS=$(foreach dep,$(PROJECT_DEPENDENCIES), \
	$(eval project_path_parts:=$(subst /, ,$(dep))) \
	$(eval project_path:=$(BASE_DIRECTORY)/projects/$(word 2,$(project_path_parts))/$(word 3,$(project_path_parts))) \
	$(eval release_branch:=$(or $(word 4,$(project_path_parts)),$(RELEASE_BRANCH))) \
	$(if $(or $(findstring eksd,$(dep)), \
		$(and \
			$(if $(wildcard $(project_path)),true,$(error Non-existent dependency: $(dep))), \
			$(filter true,$(shell $(MAKE) -C $(project_path) var-value-HAS_S3_ARTIFACTS RELEASE_BRANCH=$(release_branch))) \
		)),$(call FULL_FETCH_BINARIES_TARGETS,$(dep)),))
####################################################

#################### LICENSES ######################
HAS_LICENSES?=true
ATTRIBUTION_TARGETS?=$(call pairmap,ATTRIBUTION_TARGET_FROM_BINARY_GO_MOD,$(BINARY_TARGET_FILES),$(GO_MOD_PATHS))
GATHER_LICENSES_TARGETS?=$(call pairmap,LICENSE_TARGET_FROM_BINARY_GO_MOD,$(BINARY_TARGET_FILES),$(GO_MOD_PATHS))
LICENSES_OUTPUT_DIR?=$(OUTPUT_DIR)
LICENSES_TARGETS_FOR_PREREQ?=$(if $(filter true,$(HAS_LICENSES)),$(GATHER_LICENSES_TARGETS) \
	$(foreach target,$(ATTRIBUTION_TARGETS),_output/$(target)),)
# .9 is the default if nothing is passed to go-licenses
# allow override on a per project basis for super specific cases
LICENSE_THRESHOLD?=.9
####################################################

#################### TARBALLS ######################
HAS_S3_ARTIFACTS?=false

SIMPLE_CREATE_TARBALLS?=true
TAR_FILE_PREFIX?=$(REPO)
FAKE_ARM_BINARIES_FOR_VALIDATION?=$(if $(filter linux/arm64,$(BINARY_PLATFORMS)),false,true)
FAKE_AMD_BINARIES_FOR_VALIDATION?=$(if $(filter linux/amd64,$(BINARY_PLATFORMS)),false,true)
FAKE_ARM_IMAGES_FOR_VALIDATION?=false
IMAGE_FORMAT?=
IMAGE_OS?=
UPLOAD_DO_NOT_DELETE?=false
UPLOAD_CREATE_PUBLIC_ACL?=true
EXPECTED_FILES_PATH?=expected_artifacts
####################################################

#################### OTHER #########################
KUSTOMIZE_VERSION=5.4.3
KUSTOMIZE_TARGET=$(OUTPUT_DIR)/kustomize
GIT_DEPS_DIR?=$(OUTPUT_DIR)/gitdependencies
SPECIAL_TARGET_SECONDARY+=$(strip $(PROJECT_DEPENDENCIES_TARGETS) $(GO_MOD_DOWNLOAD_TARGETS))
SKIP_CHECKSUM_VALIDATION?=false
IN_DOCKER_TARGETS=all-attributions all-attributions-checksums all-checksums attribution attribution-checksums binaries checksums clean clean-go-cache validate-checksums $(GO_MOD_DOWNLOAD_TARGETS) $(BINARY_TARGETS) $(GATHER_LICENSES_TARGETS)
PRUNE_BUILDCTL?=false
GITHUB_TOKEN?=

# if this is set we are running in the context of a run-<>-in-docker target
DOCKER_RUN_BASE_DIRECTORY?=
DEPENDENCY_TOOLS=buildctl helm jq lz4 skopeo tuftool yq

# set to true if do not care about checksum match, default to false to limit confusion
# around checksum mismatch
FORCE_BUILD_ON_HOST?=false

# $1 - should run on host
MAYBE_RUN_IN_DOCKER?=$(if \
	$(or \
		$(filter true,$(IS_ON_BUILDER_BASE)), \
		$(filter true,$(FORCE_BUILD_ON_HOST)), \
		$(filter true,$2), \
		$(filter false,$(DOCKER_AVAILABLE)) \
	),,true)

# this can be used as a normal macro, $(ENABLE_DOCKER), or as a func with 1 param, $(call ENABLE_DOCKER)
# $1 - should run on host (optional)
ENABLE_DOCKER=$(eval $(call ENABLE_DOCKER_BODY,$@,$(if $(filter undefined,$(origin 1)),false,$(value 1)),))
# $1 - container platform to use
ENABLE_DOCKER_PLATFORM=$(eval $(call ENABLE_DOCKER_BODY,$@,false,$(if $(filter undefined,$(origin 1)),,$(value 1))))

# $1 - target name
# $2 - should run on host
# $3 - container platform to use (for cgo)
# SHELL is overriden with the run_in_docker.sh $(RUN_IN_DOCKER_BODY) args passed at the beginning
# in the make_shell.sh, if the action is docker, then it will wait until its the actual target
# signified by LOGGING_TARGET being set and then run that script instead of the actual target script
# since we are manipulating the SHELL here we have to make sure any variables which are going to be used in the
# context of MAYBE_RUN_IN_DOCKER cannot themselves call out to a shell, we can force eval to avoid
define ENABLE_DOCKER_BODY
$(eval _TARGET:=$1)
$(eval _RUN_ON_HOST:=$2)
$(eval _DOCKER_PLATFORM:=$3)
$(eval _USE_DOCKER:=$(call MAYBE_RUN_IN_DOCKER,$(_RUN_ON_HOST)))
$(_TARGET): export LOGGING_TARGET=$(_TARGET)
$(_TARGET): export RUN_IN_DOCKER_ARGS=$$(RUN_IN_DOCKER_ARGS_BODY)
$(_TARGET): MAKE_TARGET:=$(_TARGET)
$(_TARGET): DOCKER_PLATFORM:=$(_DOCKER_PLATFORM)
$(_TARGET): USE_DOCKER:=$(_USE_DOCKER)
$(_TARGET): SHELL=$(if $(_USE_DOCKER),$(DOCKER_SHELL),$(LOGGING_SHELL))
endef

define RUN_IN_DOCKER_ARGS_BODY
$(COMPONENT) $(MAKE_TARGET) $(IMAGE_REPO) "$(RELEASE_BRANCH)" "$(ARTIFACTS_BUCKET)" "$(BASE_DIRECTORY)" "$(GO_MOD_CACHE)" "$(BUILDER_PLATFORM_ARCH)" true "$(CURRENT_BUILDER_BASE_TAG)" "$(DOCKER_PLATFORM)"
endef

####################################################

#################### LOGGING #######################
DATE_CMD=TZ=utc $(shell source $(BUILD_LIB)/common.sh && build::find::gnu_variant_on_mac date)
DATE_NANO=$(shell if [ "$$(uname -s)" = "Linux" ] || command -v gdate &> /dev/null; then echo %3N; fi)
TARGET_START_LOG?=$(eval _START_TIME:=$(shell $(DATE_CMD) +%s.$(DATE_NANO)))\\n------------------- $(shell $(DATE_CMD) +"%Y-%m-%dT%H:%M:%S.$(DATE_NANO)%z") Starting target=$@ -------------------
TARGET_END_LOG?="------------------- `$(DATE_CMD) +'%Y-%m-%dT%H:%M:%S.$(DATE_NANO)%z'` Finished target=$@ duration=`echo $$($(DATE_CMD) +%s.$(DATE_NANO)) - $(_START_TIME) | bc` seconds -------------------\\n"

ENABLE_LOGGING=$(eval $(call ENABLE_LOGGING_BODY,$@))
# $1 - target name
# target is exported so that the script can tell the difference between
# a recipe and a $(shell) call, which we do not want to log around
# this style of enable logging only works for single line recipes
define ENABLE_LOGGING_BODY
$(eval _TARGET:=$1)
$(_TARGET): export LOGGING_TARGET=$(_TARGET)
$(_TARGET): SHELL=$$(LOGGING_SHELL)
endef
####################################################

#################### TARGETS FOR OVERRIDING ########
BUILD_TARGETS?=github-rate-limit-pre validate-checksums attribution $(if $(IMAGE_NAMES),local-images,) $(if $(filter true,$(HAS_HELM_CHART)),helm/build,) $(if $(filter true,$(HAS_S3_ARTIFACTS)),upload-artifacts,) attribution-pr github-rate-limit-post
RELEASE_TARGETS?=validate-checksums $(if $(IMAGE_NAMES),images,) $(if $(filter true,$(HAS_HELM_CHART)),helm/push,) $(if $(filter true,$(HAS_S3_ARTIFACTS)),upload-artifacts,)
BUILD_TARGETS_OVERRIDE?=
RELEASE_TARGETS_OVERRIDE?=
####################################################

# convert commonly used, usually shell call, variables to lazily resolved cached variables
CACHE_VARS=AWS_ACCOUNT_ID BUILD_IDENTIFIER BUILDER_PLATFORM_ARCH BUILDER_PLATFORM_OS DATE_CMD DATE_NANO \
	GIT_HASH GIT_TAG GO_BUILD_CACHE GO_MOD_CACHE GOLANG_VERSION HELM_GIT_TAG IS_ON_BUILDER_BASE \
	PROJECT_DEPENDENCIES_TARGETS SUPPORTED_K8S_VERSIONS BUILDCTL_AVAILABLE BUILDX_AVAILABLE DOCKER_AVAILABLE CURRENT_BUILDER_BASE_TAG
$(foreach v,$(strip $(CACHE_VARS)),$(call CACHE_VARIABLE,$(v)))

define BUILDCTL
	$(BUILD_LIB)/buildkit.sh \
		build \
		--frontend dockerfile.v0 \
		--opt platform=$(IMAGE_PLATFORMS) \
		--opt build-arg:BASE_IMAGE=$(BASE_IMAGE) \
		--opt build-arg:BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--opt build-arg:RELEASE_BRANCH=$(RELEASE_BRANCH) \
		$(foreach BUILD_ARG,$(IMAGE_BUILD_ARGS),--opt build-arg:$(BUILD_ARG)=$($(BUILD_ARG))) \
		--progress plain \
		--local dockerfile=$(DOCKERFILE_FOLDER) \
		--local context=$(IMAGE_CONTEXT_DIR) \
		$(if $(filter push=true,$(IMAGE_OUTPUT)),$(IMAGE_EXPORT_CACHE)) \
		$(foreach IMPORT_CACHE,$(IMAGE_IMPORT_CACHE),--import-cache $(IMPORT_CACHE)) \
		--opt target=$(IMAGE_TARGET) \
		--output type=$(IMAGE_OUTPUT_TYPE),oci-mediatypes=true,\"name=$(IMAGE),$(LATEST_IMAGE)\",$(IMAGE_OUTPUT)
endef 

define WRITE_LOCAL_IMAGE_TAG
	echo $(IMAGE_TAG) > $(IMAGE_OUTPUT_DIR)/$(IMAGE_OUTPUT_NAME).docker_tag
	echo $(IMAGE) > $(IMAGE_OUTPUT_DIR)/$(IMAGE_OUTPUT_NAME).docker_image_name	
endef

# Do not binary deps + go mod download file as intermediate files
ifneq ($(SPECIAL_TARGET_SECONDARY),)
.SECONDARY: $(SPECIAL_TARGET_SECONDARY)
endif

#### Source repo + binary Targets
ifneq ($(REPO_NO_CLONE),true)
$(REPO): | ensure-bash-version
	@echo -e $(call TARGET_START_LOG)
ifneq ($(REPO_SPARSE_CHECKOUT),)
	source $(BUILD_LIB)/common.sh && retry git clone --quiet --depth 1 --filter=blob:none --sparse -b $(GIT_TAG) $(CLONE_URL) $(REPO)
	git -C $(REPO) sparse-checkout set $(REPO_SPARSE_CHECKOUT) --cone --skip-checks
else
	source $(BUILD_LIB)/common.sh && retry git clone --quiet $(CLONE_URL) $(REPO)
endif
	@echo -e $(call TARGET_END_LOG)
endif

$(GIT_CHECKOUT_TARGET): | $(REPO)
	@echo -e $(call TARGET_START_LOG)
	@rm -f $(REPO)/eks-anywhere-*
	(cd $(REPO) && $(BASE_DIRECTORY)/build/lib/wait_for_tag.sh $(GIT_TAG))
	git -C $(REPO) checkout --quiet -f $(GIT_TAG)
	@touch $@
	@echo -e $(call TARGET_END_LOG)

$(GIT_PATCH_TARGET): $(GIT_CHECKOUT_TARGET)
	@echo -e $(call TARGET_START_LOG)
	git -C $(REPO) config user.email prow@amazonaws.com
	git -C $(REPO) config user.name "Prow Bot"
	if [ -n "$(PATCHES_DIR)" ]; then git -C $(REPO) am --committer-date-is-author-date $(PATCHES_DIR)/*; fi
	@touch $@
	@echo -e $(call TARGET_END_LOG)

ifneq ($(PATCHES_DIR),)
update-patch-numbers:
	$(SED_CMD) -i -E "s|PATCH (.*)/[0-9]+|PATCH \1/$(shell ls -1 $(PATCHES_DIR) | wc -l | tr -d ' ')|" $(PATCHES_DIR)/*
endif

## GO mod download targets
$(REPO)/%ks-anywhere-go-mod-download: REPO_SUBPATH=$(if $(filter e,$*),,$(*:%/e=%))
$(REPO)/%ks-anywhere-go-mod-download: $(if $(PATCHES_DIR),$(GIT_PATCH_TARGET),$(GIT_CHECKOUT_TARGET)) | $$(ENABLE_DOCKER)
	@if [[ "$(GO_MODS_VENDORED)" == "false" ]]; then $(BASE_DIRECTORY)/build/lib/go_mod_download.sh $(MAKE_ROOT) $(REPO) $(GIT_TAG) $(GOLANG_VERSION) "$(REPO_SUBPATH)"; fi; \
	touch $@

ifneq ($(REPO),$(HELM_SOURCE_REPOSITORY))
$(HELM_SOURCE_REPOSITORY):
	@echo -e $(call TARGET_START_LOG)
	source $(BUILD_LIB)/common.sh && retry git clone --quiet $(HELM_CLONE_URL) $(HELM_SOURCE_REPOSITORY)
	@echo -e $(call TARGET_END_LOG)
endif

ifneq ($(GIT_TAG),$(HELM_GIT_TAG))
$(HELM_GIT_CHECKOUT_TARGET): | $(HELM_SOURCE_REPOSITORY)
	@echo -e $(call TARGET_START_LOG)
	@echo rm -f $(HELM_SOURCE_REPOSITORY)/eks-anywhere-*
	(cd $(HELM_SOURCE_REPOSITORY) && $(BASE_DIRECTORY)/build/lib/wait_for_tag.sh $(HELM_GIT_TAG))
	git -C $(HELM_SOURCE_REPOSITORY) checkout --quiet -f $(HELM_GIT_TAG)
	touch $@
	@echo -e $(call TARGET_END_LOG)
endif

$(HELM_GIT_PATCH_TARGET): $(HELM_GIT_CHECKOUT_TARGET)
	@echo -e $(call TARGET_START_LOG)
	git -C $(HELM_SOURCE_REPOSITORY) config user.email prow@amazonaws.com
	git -C $(HELM_SOURCE_REPOSITORY) config user.name "Prow Bot"
	if [ -n "$(HELM_PATCHES_DIR)" ]; then git -C $(HELM_SOURCE_REPOSITORY) am --committer-date-is-author-date $(HELM_PATCHES_DIR)/*; fi
	@touch $@
	@echo -e $(call TARGET_END_LOG)

ifeq ($(SIMPLE_CREATE_BINARIES),true)
# GO_MOD_TARGET_FOR_BINARY_<binary> variables are created earlier in the makefile when determining which binaries can be built together vs alone
# if target is included in BINARY_TARGET_FILES_BUILD_TOGETHER list, use SOURCE_PATTERNS_BUILD_TOGETHER, otherewise use source pattern at the same index as binary_target in binary_target_files
$(OUTPUT_BIN_DIR)/%: PLATFORM=$(subst -,/,$(*D))
$(OUTPUT_BIN_DIR)/%: BINARY_TARGET=$(@F:%.exe=%)
$(OUTPUT_BIN_DIR)/%: SOURCE_PATTERN=$(if $(filter $(BINARY_TARGET),$(BINARY_TARGET_FILES_BUILD_TOGETHER)),$(SOURCE_PATTERNS_BUILD_TOGETHER),$(word $(call pos,$(BINARY_TARGET),$(BINARY_TARGET_FILES)),$(SOURCE_PATTERNS)))
$(OUTPUT_BIN_DIR)/%: OUTPUT_PATH=$(if $(and $(if $(filter false,$(call IS_ONE_WORD,$(BINARY_TARGET_FILES_BUILD_TOGETHER))),$(filter $(BINARY_TARGET),$(BINARY_TARGET_FILES_BUILD_TOGETHER)))),$(@D)/,$@)
$(OUTPUT_BIN_DIR)/%: GO_MOD_PATH=$($(call GO_MOD_TARGET_FOR_BINARY_VAR_NAME,$(BINARY_TARGET)))
$(OUTPUT_BIN_DIR)/%: GO_MOD_DOWNLOAD_TARGET=$(call GO_MOD_DOWNLOAD_TARGET_FROM_GO_MOD_PATH,$(GO_MOD_PATH))
# if cgo build, set platform when running the docker container
$(OUTPUT_BIN_DIR)/%: SET_PLATFORM=$(if $(filter true,$(NEEDS_CGO_BUILDER)),$(PLATFORM),)
$(OUTPUT_BIN_DIR)/%: $$(GO_MOD_DOWNLOAD_TARGET) | $$(call ENABLE_DOCKER_PLATFORM,$$(SET_PLATFORM))
	@$(BASE_DIRECTORY)/build/lib/simple_create_binaries.sh $(MAKE_ROOT) $(MAKE_ROOT)/$(OUTPUT_PATH) $(REPO) $(GOLANG_VERSION) $(PLATFORM) "$(SOURCE_PATTERN)" "$(GOBUILD_COMMAND)" "$(EXTRA_GOBUILD_FLAGS)" "$(GO_LDFLAGS)" $(CGO_ENABLED) "$(CGO_LDFLAGS)" "$(CGO_CFLAGS_ALLOW)" "$(GO_MOD_PATH)" "$(BINARY_TARGET_FILES_BUILD_TOGETHER)"
endif

.PHONY: binaries
binaries: $(BINARY_TARGETS)

$(KUSTOMIZE_TARGET):
	@echo -e $(call TARGET_START_LOG)
ifeq ($(GITHUB_TOKEN),)
	$(warning No GITHUB_TOKEN set, may get rate limited while trying to install kustomize)
endif
	@mkdir -p $(OUTPUT_DIR)
	source $(BUILD_LIB)/common.sh && retry curl -o /tmp/install_kustomize.sh -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
	chmod +x /tmp/install_kustomize.sh
	source $(BUILD_LIB)/common.sh && retry /tmp/install_kustomize.sh $(KUSTOMIZE_VERSION) $(OUTPUT_DIR)
	@echo -e $(call TARGET_END_LOG)

.PHONY: clone-repo
clone-repo: $(REPO)

.PHONY: checkout-repo
checkout-repo: $(if $(filter true,$(REPO_NO_CLONE)),,$(if $(PATCHES_DIR),$(GIT_PATCH_TARGET),$(GIT_CHECKOUT_TARGET)))

.PHONY: checkout-helm-repo
checkout-helm-repo: $(if $(filter true,$(REPO_NO_CLONE)),,$(if $(HELM_PATCHES_DIR),$(HELM_GIT_PATCH_TARGET),$(HELM_GIT_CHECKOUT_TARGET)))

.PHONY: patch-repo
patch-repo: checkout-repo

## File/Folder Targets

$(OUTPUT_DIR)/images/%:
	@mkdir -p $(@D)

$(OUTPUT_DIR)/%TTRIBUTION.txt: SOURCE_FILE=$(@:_output/%=%) # we want to keep the release branch part which is in the OUTPUT var, hardcoding _output
$(OUTPUT_DIR)/%TTRIBUTION.txt:
	@mkdir -p $(OUTPUT_DIR)
	@cp $(SOURCE_FILE) $(OUTPUT_DIR)


## License Targets
# if there is only one go mod path then licenses are gathered to _output, `%` will equal `a`
# multiple go mod paths are in use and licenses are gathered and stored in sub folders, `%` will equal `<binary>/a`
# GO_MOD_TARGET_FOR_BINARY_<binary> variables are created earlier in the makefile when determining which binaries can be built together vs alone
$(OUTPUT_DIR)/%ttribution/go-license.csv: BINARY_TARGET=$(if $(filter .,$(*D)),,$(*D))
$(OUTPUT_DIR)/%ttribution/go-license.csv: GO_MOD_PATH=$(if $(BINARY_TARGET),$(GO_MOD_TARGET_FOR_BINARY_$(call TO_UPPER,$(BINARY_TARGET))),$(word 1,$(UNIQ_GO_MOD_PATHS)))
$(OUTPUT_DIR)/%ttribution/go-license.csv: LICENSE_PACKAGE_FILTER=$(GO_MOD_$(subst /,_,$(GO_MOD_PATH))_LICENSE_PACKAGE_FILTER)
$(OUTPUT_DIR)/%ttribution/go-license.csv: GO_MOD_DOWNLOAD_TARGET=$(call GO_MOD_DOWNLOAD_TARGET_FROM_GO_MOD_PATH,$(GO_MOD_PATH)) 
$(OUTPUT_DIR)/%ttribution/go-license.csv: $$(GO_MOD_DOWNLOAD_TARGET) | ensure-jq $$(ENABLE_DOCKER)
	@$(BASE_DIRECTORY)/build/lib/gather_licenses.sh $(REPO) $(MAKE_ROOT)/$(OUTPUT_DIR)/$(BINARY_TARGET) "$(LICENSE_PACKAGE_FILTER)" $(GO_MOD_PATH) $(GOLANG_VERSION) $(LICENSE_THRESHOLD) $(CGO_ENABLED)

.PHONY: gather-licenses
gather-licenses: $(GATHER_LICENSES_TARGETS)

## Attribution Targets
# if there is only one go mod path so only one attribution is created, the file will be named ATTRIBUTION.txt and licenses will be stored in _output, `%` will equal `A`
# if multiple attributions are being generated, the file will be <binary>_ATTRIBUTION.txt and licenses will be stored in _output/<binary>, `%` will equal `<BINARY>_A`
%TTRIBUTION.txt: LICENSE_OUTPUT_PATH=$(OUTPUT_DIR)$(if $(filter A,$(*F)),,/$(call TO_LOWER,$(*F:%_A=%)))
%TTRIBUTION.txt: GENERATE_ATTR_AVAIL=$(shell command -v generate-attribution &> /dev/null && echo "true" || echo "false")
%TTRIBUTION.txt: LICENSE_TARGET=$(LICENSE_OUTPUT_PATH)/attribution/go-license.csv
%TTRIBUTION.txt: $$(LICENSE_TARGET) | $$(call ENABLE_DOCKER,$$(GENERATE_ATTR_AVAIL))
	@$(BASE_DIRECTORY)/build/lib/create_attribution.sh $(MAKE_ROOT) $(GOLANG_VERSION) $(MAKE_ROOT)/$(LICENSE_OUTPUT_PATH) $(@F) $(RELEASE_BRANCH)

.PHONY: attribution
attribution: $(and $(filter true,$(HAS_LICENSES)),$(ATTRIBUTION_TARGETS))

.PHONY: attribution-pr
attribution-pr: attribution | $$(ENABLE_LOGGING)
	@$(BASE_DIRECTORY)/build/update-attribution-files/create_pr.sh

.PHONY: all-attributions
all-attributions:
	$(BASE_DIRECTORY)/build/update-attribution-files/make_attribution.sh projects/$(COMPONENT) attribution

#### Tarball Targets

.PHONY: tarballs
tarballs: $(LICENSES_TARGETS_FOR_PREREQ) | $$(ENABLE_LOGGING)
ifeq ($(SIMPLE_CREATE_TARBALLS),true)
	@$(BUILD_LIB)/simple_create_tarballs.sh $(TAR_FILE_PREFIX) $(MAKE_ROOT)/$(OUTPUT_DIR) $(MAKE_ROOT)/$(OUTPUT_BIN_DIR) $(GIT_TAG) "$(BINARY_PLATFORMS)" $(ARTIFACTS_PATH) $(GIT_HASH)
endif

.PHONY: upload-artifacts
upload-artifacts: s3-artifacts upload-output-to-prow-artifacts-s3-artifacts | $$(ENABLE_LOGGING)
	@$(BUILD_LIB)/upload_artifacts.sh $(ARTIFACTS_PATH) $(ARTIFACTS_BUCKET) $(ARTIFACTS_UPLOAD_PATH) $(BUILD_IDENTIFIER) $(GIT_HASH) $(LATEST) $(UPLOAD_DRY_RUN) $(UPLOAD_DO_NOT_DELETE) $(UPLOAD_CREATE_PUBLIC_ACL)

.PHONY: s3-artifacts
s3-artifacts: tarballs
	@echo -e $(call TARGET_START_LOG)
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG) $(FAKE_ARM_BINARIES_FOR_VALIDATION) $(FAKE_AMD_BINARIES_FOR_VALIDATION) $(MAKE_ROOT)/$(EXPECTED_FILES_PATH) $(IMAGE_OS)
	@echo -e $(call TARGET_END_LOG)

.PHONY: upload-output-to-prow-artifacts-%
upload-output-to-prow-artifacts-%:
	@if [[ "$(JOB_TYPE)" == "presubmit" ]] && [[ "$(INCLUDE_OUTPUT_IN_PROW_ARTIFACTS)" == "true" ]]; then \
		cp -rf $(OUTPUT_DIR) $(ARTIFACTS); \
	fi

### Checksum Targets

.PHONY: checksums
checksums: $(BINARY_TARGETS) | $$(ENABLE_LOGGING)
ifneq ($(strip $(BINARY_TARGETS)),)
	@$(BASE_DIRECTORY)/build/lib/update_checksums.sh $(MAKE_ROOT) $(PROJECT_ROOT) $(MAKE_ROOT)/$(OUTPUT_BIN_DIR)
endif

.PHONY: validate-checksums
validate-checksums: $(BINARY_TARGETS) upload-output-to-prow-artifacts-validate-checksums | $$(ENABLE_LOGGING)
ifneq ($(and $(strip $(BINARY_TARGETS)), $(filter false, $(SKIP_CHECKSUM_VALIDATION))),)
	@$(BASE_DIRECTORY)/build/lib/validate_checksums.sh $(MAKE_ROOT) $(PROJECT_ROOT) $(MAKE_ROOT)/$(OUTPUT_BIN_DIR) $(FAKE_ARM_BINARIES_FOR_VALIDATION) $(FAKE_AMD_BINARIES_FOR_VALIDATION)
endif

.PHONY: attribution-checksums
attribution-checksums: attribution checksums

.PHONY: all-checksums
all-checksums:
	$(BASE_DIRECTORY)/build/update-attribution-files/make_attribution.sh projects/$(COMPONENT) checksums

.PHONY: all-attributions-checksums
all-attributions-checksums:
	$(BASE_DIRECTORY)/build/update-attribution-files/make_attribution.sh projects/$(COMPONENT) "attribution checksums"

#### Image Helpers

ifneq ($(IMAGE_NAMES),)
.PHONY: local-images images
local-images: clean-job-caches $(LOCAL_IMAGE_TARGETS)
images: $(IMAGE_TARGETS)
endif

.PHONY: clean-job-caches
# space is very limited in presubmit jobs, the image builds can push the total used space over the limit.
# go-build cache and pkg mod cache handled by target above
# prune is handled by buildkit.sh
clean-job-caches: $(and $(findstring presubmit,$(JOB_TYPE)),$(filter true,$(PRUNE_BUILDCTL)),clean-go-cache)

.PHONY: %/images/push %/images/amd64 %/images/arm64
%/images/push %/images/amd64 %/images/arm64: IMAGE_NAME=$*
%/images/push %/images/amd64 %/images/arm64: DOCKERFILE_FOLDER?=./docker/linux
%/images/push %/images/amd64 %/images/arm64: IMAGE_CONTEXT_DIR?=.
%/images/push %/images/amd64 %/images/arm64: IMAGE_BUILD_ARGS?=
# if there is neither buildctl or buildx, use ensure-buildctl to show the user that error
%/images/push %/images/amd64 %/images/arm64 %-useradd/images/export: ENSURE_PREREQ=$(if $(or $(filter true,$(BUILDCTL_AVAILABLE)),$(filter false,$(BUILDX_AVAILABLE))),ensure-buildkitd-host,)
%/images/push %/images/amd64 %/images/arm64 %-useradd/images/export: export USE_BUILDX=$(if $(filter ensure-buildkitd-host,$(ENSURE_PREREQ)),false,true)

# Build image using buildkit for all platforms, by default pushes to registry defined in IMAGE_REPO.
%/images/push: IMAGE_PLATFORMS?=linux/amd64,linux/arm64
%/images/push: IMAGE_OUTPUT_TYPE?=image
%/images/push: IMAGE_OUTPUT?=push=true

# Build image using buildkit only builds linux/amd64 oci and saves to local tar.
%/images/amd64: IMAGE_PLATFORMS?=linux/amd64

# Build image using buildkit only builds linux/arm64 oci and saves to local tar.
%/images/arm64: IMAGE_PLATFORMS?=linux/arm64

%/images/amd64 %/images/arm64: IMAGE_OUTPUT_TYPE?=oci
%/images/amd64 %/images/arm64: IMAGE_OUTPUT?=dest=$(IMAGE_OUTPUT_DIR)/$(IMAGE_OUTPUT_NAME).tar

%/images/push: $(BINARY_TARGETS) $(LICENSES_TARGETS_FOR_PREREQ) $(HANDLE_DEPENDENCIES_TARGET) | $$(ENSURE_PREREQ) $$(ENABLE_LOGGING)
	@$(BUILDCTL)

%/images/amd64: $(BINARY_TARGETS) $(LICENSES_TARGETS_FOR_PREREQ) $(HANDLE_DEPENDENCIES_TARGET) | $$(ENSURE_PREREQ)
	@echo -e $(call TARGET_START_LOG)
	@mkdir -p $(IMAGE_OUTPUT_DIR)
	$(BUILDCTL)
	$(WRITE_LOCAL_IMAGE_TAG)
	@echo -e $(call TARGET_END_LOG)

%/images/arm64: $(BINARY_TARGETS) $(LICENSES_TARGETS_FOR_PREREQ) $(HANDLE_DEPENDENCIES_TARGET) | $$(ENSURE_PREREQ)
	@echo -e $(call TARGET_START_LOG)
	@mkdir -p $(IMAGE_OUTPUT_DIR)
	$(BUILDCTL)
	$(WRITE_LOCAL_IMAGE_TAG)
	@echo -e $(call TARGET_END_LOG)

## Useradd targets
%-useradd/images/export: IMAGE_OUTPUT_TYPE=local
%-useradd/images/export: IMAGE_OUTPUT_DIR=$(OUTPUT_DIR)/files/$*
%-useradd/images/export: IMAGE_OUTPUT?=dest=$(IMAGE_OUTPUT_DIR)
%-useradd/images/export: IMAGE_BUILD_ARGS=IMAGE_USERADD_USER_ID IMAGE_USERADD_USER_NAME
%-useradd/images/export: DOCKERFILE_FOLDER=$(BUILD_LIB)/docker/linux/useradd
%-useradd/images/export: IMAGE_PLATFORMS=linux/$(BUILDER_PLATFORM_ARCH)
%-useradd/images/export: | $$(ENSURE_PREREQ) $$(ENABLE_LOGGING)
	@mkdir -p $(IMAGE_OUTPUT_DIR) && $(BUILDCTL)

PHONY: combine-images
combine-images: IMAGE_BUILD_ARGS=IMAGE
combine-images: DOCKERFILE_FOLDER=$(BUILD_LIB)/docker/linux/combine
combine-images: IMAGE_EXPORT_CACHE=--export-cache type=inline
combine-images: IMAGE_TARGET=
combine-images: IMAGE_CONTEXT_DIR=.
combine-images: images

## Helm Targets
.PHONY: helm/pull 
helm/pull: | $$(ENABLE_LOGGING)
	@$(BUILD_LIB)/helm_pull.sh $(HELM_PULL_LOCATION) $(HELM_REPO_URL) $(HELM_PULL_NAME) $(REPO) $(HELM_DIRECTORY) $(CHART_VERSION) $(COPY_CRDS)

.PHONY: %/helm/copy %/helm/require %/helm/replace %/helm/build %/helm/push
%/helm/copy %/helm/require %/helm/replace %/helm/build %/helm/push: HELM_CHART_NAME=$*

$(call FULL_CHART_TARGETS,copy) : %/helm/copy: checkout-repo checkout-helm-repo $(LICENSES_TARGETS_FOR_PREREQ) | ensure-helm ensure-skopeo $$(ENABLE_LOGGING)
	@$(BUILD_LIB)/helm_copy.sh $(HELM_SOURCE_REPOSITORY) $(HELM_DESTINATION_REPOSITORY) $(HELM_DIRECTORY) $(OUTPUT_DIR)

$(call FULL_CHART_TARGETS,require) : %/helm/require: %/helm/copy | $$(ENABLE_LOGGING)
	@$(BUILD_LIB)/helm_require.sh $(HELM_SOURCE_IMAGE_REPO) $(HELM_DESTINATION_REPOSITORY) $(OUTPUT_DIR) $(IMAGE_TAG) $(HELM_TAG) $(PROJECT_ROOT) $(LATEST) $(HELM_USE_UPSTREAM_IMAGE) "$(PACKAGE_DEPENDENCIES)" "$(FORCE_JSON_SCHEMA_FILE)" "$(HELM_IMAGE_LIST)" "$(HELM_IMAGE_TAG_LIST)"

$(call FULL_CHART_TARGETS,replace) : %/helm/replace: %/helm/require | $$(ENABLE_LOGGING)
	@$(BUILD_LIB)/helm_replace.sh $(HELM_DESTINATION_REPOSITORY) $(HELM_CHART_FOLDER) $(OUTPUT_DIR)

$(call FULL_CHART_TARGETS,build) : %/helm/build: %/helm/replace | $$(ENABLE_LOGGING)
	@$(BUILD_LIB)/helm_build.sh $(OUTPUT_DIR) $(HELM_DESTINATION_REPOSITORY) $(HELM_CHART_FOLDER) $(BUILD_HELM_DEPENDENCIES)

$(call FULL_CHART_TARGETS,push) : %/helm/push: %/helm/build | $$(ENABLE_LOGGING)
	@$(BUILD_LIB)/helm_push.sh $(IMAGE_REPO) $(HELM_DESTINATION_REPOSITORY) $(HELM_CHART_FOLDER) $(HELM_TAG) $(GIT_TAG) $(OUTPUT_DIR) $(LATEST)

# Build helm chart
.PHONY: helm/build
helm/build: $(foreach chart,$(HELM_CHART_NAMES),$(chart)/helm/build)

# Build helm chart and push to registry defined in IMAGE_REPO.
.PHONY: helm/push
helm/push: $(if $(filter true,$(HAS_HELM_CHART)),$(foreach chart,$(HELM_CHART_NAMES),$(chart)/helm/push),)

#@ Fetch Binary Targets
.PHONY: handle-dependencies 
handle-dependencies: # Download and extract TARs for each dependency listed in PROJECT_DEPENDENCIES
handle-dependencies: $(PROJECT_DEPENDENCIES_TARGETS)

$(BINARY_DEPS_DIR)/linux-%: | $$(ENABLE_LOGGING)
	@$(BUILD_LIB)/fetch_binaries.sh $(BINARY_DEPS_DIR) $* $(ARTIFACTS_BUCKET) $(LATEST) $(RELEASE_BRANCH)

## Build Targets
.PHONY: build
build: $(or $(BUILD_TARGETS_OVERRIDE),$(BUILD_TARGETS))

.PHONY: release
release: $(or $(RELEASE_TARGETS_OVERRIDE),$(RELEASE_TARGETS))

# Iterate over release branch versions, avoiding branches explicitly marked as skipped
.PHONY: %/release-branches/all
%/release-branches/all:
	@set -e; \
	for version in $(SUPPORTED_K8S_VERSIONS) ; do \
	    if ! [[ "$(SKIPPED_K8S_VERSIONS)" =~ $$version  ]]; then \
			$(MAKE) $* $(if $(filter true,$(BINARIES_ARE_RELEASE_BRANCHED)),clean-output,) RELEASE_BRANCH=$$version IMAGE_PLATFORMS=$(call IF_OVERRIDE_VARIABLE,IMAGE_PLATFORMS,); \
		fi \
	done;

###  Clean Targets

# When go downloads pkg to the module cache, GOPATH/pkg/mod, it removes the write permissions
# prevent accident modifications since files/checksums are tightly controlled
# adding the perms necessary to perform the delete
# When building go bins using mods which have been downloaded by go mod download/vendor which will exist in the go_mod_cache
# there is additional checksum (?) information that is not preserved in the vendor directory within the project folder
# This additional information gets written out into the resulting binary. If we did not run go mod vendor, which we do 
# for all project builds, we could get checksum mismatches on the final binaries due to sometimes having the mod previously
# downloaded in the go_mod_cahe.  Running go mod vendor always ensures that the go mod has always been downloaded
# to the go_mod_cache directory. If we clear the go_mod_cache we need to delete the go_mod_download sentinel file
# so the next time we run build go mods will be redownloaded
.PHONY: clean-go-cache
clean-go-cache: | $$(ENABLE_LOGGING)
	@if [ -n "$(GOLANG_VERSION)" ]; then \
		chmod -fR 777 $(GO_MOD_CACHE) &> /dev/null || :; \
		$(foreach folder,$(GO_MOD_CACHE) $(GO_BUILD_CACHE),$(if $(wildcard $(folder)),du -hs $(folder) && rm -rf $(folder);,)) \
		$(foreach file,$(GO_MOD_DOWNLOAD_TARGETS),$(if $(wildcard $(file)),rm -f $(file);,)) \
	fi

.PHONY: clean-repo
clean-repo:
	@rm -rf $(REPO)	$(HELM_SOURCE_REPOSITORY)

# intentionally not using OUTPUT_DIR variable to ensure we
# delete the entire output folder for release branch'd projects
.PHONY: clean-output
clean-output:
	@if [ -d _output ]; then \
		du -hs _output; \
		rm -rf _output; \
	fi

.PHONY: clean
clean: $(if $(filter true,$(REPO_NO_CLONE)),,clean-repo) clean-output

## --------------------------------------
## Help
## --------------------------------------
#@  Helpers
.PHONY: help
help: # Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% \/a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-55s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 4) } ' $(MAKEFILE_LIST)

.PHONY: help-list
help-list: 
	@awk 'BEGIN {FS = ":.*#";} /^[$$()% \/a-zA-Z0-9_-]+:.*?#/ { printf "%s: ##%s\n", $$1, $$2 } /^#@/ { printf "\n##@%s\n", substr($$0, 4) } ' $(MAKEFILE_LIST)

.PHONY: add-generated-help-block
add-generated-help-block: # Add or update generated help block to document project make file and support shell auto completion
add-generated-help-block:
	$(BUILD_LIB)/generate_help_body.sh $(MAKE_ROOT) "$(BINARY_TARGET_FILES)" "$(BINARY_PLATFORMS)" "${BINARY_TARGETS}" \
		$(REPO) $(if $(PATCHES_DIR),true,false) "$(LOCAL_IMAGE_TARGETS)" "$(IMAGE_TARGETS)" "$(BUILD_TARGETS)" "$(RELEASE_TARGETS)" \
		"$(HAS_S3_ARTIFACTS)" "$(HAS_LICENSES)" "$(REPO_NO_CLONE)" "$(PROJECT_DEPENDENCIES_TARGETS)" \
		"$(HAS_HELM_CHART)" "$(IN_DOCKER_TARGETS)"

## --------------------------------------
## Update Helpers
## --------------------------------------
#@ Update Helpers

.PHONY: start-docker-builder
start-docker-builder: # Start long lived builder base docker container
	$(BUILD_LIB)/run_target_docker.sh $(COMPONENT) var-value-PULL_BASE_REF $(IMAGE_REPO) "$(RELEASE_BRANCH)" "$(ARTIFACTS_BUCKET)" "$(BASE_DIRECTORY)" "$(GO_MOD_CACHE)" "$(BUILDER_PLATFORM_ARCH)" false "$(CURRENT_BUILDER_BASE_TAG)"

.PHONY: stop-docker-builder
stop-docker-builder: # Clean up builder base docker container
	$(MAKE) -C $(BASE_DIRECTORY) stop-docker-builder

.PHONY: run-buildkit-and-registry
run-buildkit-and-registry: # Run buildkitd and a local docker registry as containers
	$(MAKE) -C $(BASE_DIRECTORY) run-buildkit-and-registry

.PHONY: stop-buildkit-and-registry
stop-buildkit-and-registry: # Stop the buildkitd and a local docker registry containers
	$(MAKE) -C $(BASE_DIRECTORY) stop-buildkit-and-registry

.PHONY: generate
generate: | ensure-locale # Update UPSTREAM_PROJECTS.yaml
	$(BUILD_LIB)/generate_projects_list.sh $(BASE_DIRECTORY)

.PHONY: update-go-mods
update-go-mods: # Update locally checked-in go sum to assist in vuln scanning
update-go-mods: DEST_PATH=$(if $(IS_RELEASE_BRANCH_BUILD),$(RELEASE_BRANCH)/$$gomod,$$gomod)
update-go-mods: checkout-repo
	for gomod in $(GO_MOD_PATHS); do \
		mkdir -p $(DEST_PATH); \
		cp $(REPO)/$$gomod/go.{mod,sum} $(DEST_PATH); \
	done

.PHONY: update-vendor-for-dep-patch
update-vendor-for-dep-patch: # After bumping dep in go.mod file, uses generic vendor update script or one provided from upstream project
update-vendor-for-dep-patch: checkout-repo
	$(BUILD_LIB)/update_vendor.sh $(PROJECT_ROOT) $(REPO) $(GIT_TAG) $(GOLANG_VERSION) $(VENDOR_UPDATE_SCRIPT)

.PHONY: patch-for-dep-update
patch-for-dep-update: # After bumping dep in go.mod file and updating vendor, generates patch
patch-for-dep-update: checkout-repo
	$(BUILD_LIB)/patch_for_dep_update.sh $(REPO) $(GIT_TAG) $(PROJECT_ROOT)/patches

.PHONY: %/create-ecr-repo
%/create-ecr-repo: IMAGE_NAME=$*
%/create-ecr-repo: | $$(ENABLE_LOGGING)
	@cmd=( ecr ); \
	if [[ "${IMAGE_REPO}" =~ ^public\.ecr\.aws/ ]]; then \
		cmd=( ecr-public --region us-east-1 ); \
	fi; \
	repo=$(IMAGE_REPO_COMPONENT); \
	if [[ "$(IMAGE_NAME)" = *"__helm__"* ]]; then \
		repo="$(IMAGE_NAME:%/__helm__=%)"; \
	fi; \
	if ! aws $${cmd[*]} describe-repositories --repository-name "$$repo" > /dev/null 2>&1; then \
		aws $${cmd[*]} create-repository --repository-name "$$repo"; \
	fi;

.PHONY: create-ecr-repos
create-ecr-repos: # Create repos in ECR for project images for local testing
create-ecr-repos: $(foreach image,$(IMAGE_NAMES),$(image)/create-ecr-repo) $(if $(filter true,$(HAS_HELM_CHART)),$(foreach chart,$(HELM_CHART_NAMES),$(chart)/__helm__/create-ecr-repo),)

.PHONY: var-value-%
var-value-%:
	@echo $($*)

.PHONY: check-for-supported-release-branch
check-for-supported-release-branch:
	@if [ "$(NOT_SUPPORTED_RELEASE_BRANCH_CONFIGURATION)" == "true" ]; then \
		echo "Not a supported version to build"; \
		exit 1;	\
	elif [ -d $(MAKE_ROOT)/$(RELEASE_BRANCH) ]; then \
		echo "Supported version to build"; \
		exit 0; \
	elif { [ "false" == "$(BINARIES_ARE_RELEASE_BRANCHED)" ] || [ -z "$(BINARY_TARGET_FILES)" ]; } && \
		{ [ "true" == "$$(yq e ".releases[] | select(.branch==\"$(RELEASE_BRANCH)\") | has(\"branch\")" $(BASE_DIRECTORY)/EKSD_LATEST_RELEASES)" ] && grep $(RELEASE_BRANCH) $(BASE_DIRECTORY)/release/SUPPORTED_RELEASE_BRANCHES &> /dev/null; }; then \
		echo "Supported version to build"; \
		exit 0; \
	else \
		echo "Not a supported version to build"; \
		exit 1; \
	fi	

.PHONY: check-for-release-branch-skip
check-for-release-branch-skip:
	@if [ "$(BRANCH_NAME)" != "main" ] && [ "$(SKIP_ON_RELEASE_BRANCH)" = "true" ]; then \
		echo "Skipping build on release branch"; \
		exit 1; \
	fi

.PHONY: github-rate-limit-%
github-rate-limit-%:
	@if [[ -n "$(GITHUB_TOKEN)" ]] && [[  "presubmit" == "$(JOB_TYPE)" ]]; then \
		echo "Current Github rate limits:"; \
		GH_PAGER='' gh api rate_limit; \
	fi

# Locale settings impact file ordering in ls or shell file expansion. The file order is used to
# generate files that are subsequently validated by the CI. If local environments use different 
# locales to the CI we get unexpected failures that are tricky to debug without knowledge of 
# locales so we'll explicitly warn here.
# In a AL2 container image (like builder base), LANG will be empty which is equivalent to posix
# In a AL2 (or other distro) full instance the LANG will be en-us.UTF-8 which produces different sorts
# On Mac, LANG will be en-us.UTF-8 but has a fix applied to sort to avoid the difference
.PHONY: ensure-locale
ensure-locale:
	@if [ "Linux" = "$$(uname -s)" ]; then \
		LOCALE=$$(locale | grep LANG | cut -d= -f2 | tr -d '"' | tr '[:upper:]' '[:lower:]'); \
		if [[ "c.utf-8 posix" != *"$${LOCALE:-posix}"* ]]; then \
			echo WARNING: Environment locale set to $$LOCALE. On Linux systems this may create \
				non-deterministic behavior when running generation recipes. If the CI fails validation try \
				exporting LANG=C.UTF-8 to generate files instead.; \
		fi; \
	fi

.PHONY: ensure-/
ensure/%: CMD=$*
ensure/%:
	@if ! command -v "$(CMD)" &> /dev/null; then \
		echo "'$(CMD)' is required for this target, please install."; \
		exit 1; \
	fi

# needs to be defined explictly to avoid messing up usage as prereqs for targets with their own %/stems
.PHONY: $(foreach tool,$(DEPENDENCY_TOOLS), ensure-$(tool))
$(foreach tool,$(DEPENDENCY_TOOLS),$(eval ensure-$(tool): ensure/$(tool)))

.PHONY: ensure-docker
ensure-docker: ensure/docker
	@if ! docker info > /dev/null 2>&1 ; then \
		echo "Please ensure docker is running to make this target"; \
		exit 1; \
	fi

# in code build we use the /buildkit.sh to launch buildkitd when we need it
# in that case skip this check
.PHONY: ensure-buildkitd-host
ensure-buildkitd-host: | ensure-buildctl
	@if [ "true" = "$(IS_ON_BUILDER_BASE)" ]; then \
		exit 0; \
	elif [ -z "$${BUILDKIT_HOST:-}" ] && [ ! -S /run/buildkit/buildkitd.sock ]; then \
		echo "Please set the 'BUILDKIT_HOST' environment variable."; \
		echo "If you want to run buildkitd via a docker container use"; \
		echo "export BUILDKIT_HOST=docker-container://buildkitd && make run-buildkit-and-registry"; \
		exit 1; \
	elif ! buildctl debug workers > /dev/null 2>&1; then \
		echo "buildkit does not appear to be running."; \
		echo "If you want to run buildkitd via a docker container use"; \
		echo "make run-buildkit-and-registry"; \
		exit 1; \
	fi

.PHONY: ensure-bash-version
ensure-bash-version:
	@if (($${BASH_VERSINFO[0]}<4)) || ( (($${BASH_VERSINFO[0]}==4)) && (($${BASH_VERSINFO[1]}<2)) ); then \
    	echo "Bash version 4.2 or newer is required."; \
		if [ "$$(uname)" = 'Darwin' ]; then \
			echo "Install with 'brew install bash' on Mac OS X."; \
		fi; \
    	exit 1; \
  	fi

## --------------------------------------
## Docker Helpers
## --------------------------------------
# since these targets will likely be file paths it has to be run-in-docker/ vs run-in-docker- otherwise make
# will not properly match the stem. this requires a change to how we used to name these targets
.PHONY: run-in-docker/%
run-in-docker/%: MAKE_TARGET=$*
# on AL2/AL23 the $$(ENABLE_LOGGING) approach does not appear to work, but it does on Mac
# on Mac the explicit vars here do not work. doing both seems to work in both cases
run-in-docker/%: export LOGGING_TARGET=$@
run-in-docker/%: SHELL=$(LOGGING_SHELL)
run-in-docker/%: | ensure-docker $$(ENABLE_LOGGING)
	@$(BUILD_LIB)/run_target_docker.sh $(RUN_IN_DOCKER_ARGS_BODY)

# backcompat old style run-<>-in-docker style targets which work for anything that does not have a / in the target
.PHONY: run-%-in-docker
run-%-in-docker: run-in-docker/%
	@echo "Please switch to 'run-in-docker/$*' style targets"

# make sure by default all targets use the 
# if we do not have this as a catch all target then the first target
# which sets it to 0 will affect all the prereq targets for that target
%: SHELL=$(DEFAULT_SHELL)
# reset to default in case we are running a submake or some other target which has set it
%: export LOGGING_TARGET=
%: export RUN_IN_DOCKER_ARGS=

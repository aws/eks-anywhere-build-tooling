BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=linuxkit
REPO_OWNER=linuxkit

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=linuxkit rngd sysfs sysctl init rc.init service
GO_MOD_PATHS=src/cmd/linuxkit pkg/rngd pkg/sysfs pkg/sysctl pkg/init pkg/init pkg/init
SOURCE_PATTERNS=. ./cmd/rngd . . ./cmd/init ./cmd/rc.init ./cmd/service


# this is a native build which we can not cross compile
# ensure we are building to a directory based on the current
# host platform to avoid ever creating a different arch'd
# binary in the wrong folder
BINARY_PLATFORMS?=linux/$(BUILDER_PLATFORM_ARCH)
IMAGE_PLATFORMS?=linux/$(BUILDER_PLATFORM_ARCH)
USERMODE_HELPER_TARGETS=$(foreach platform,$(IMAGE_PLATFORMS),$(OUTPUT_BIN_DIR)/$(subst /,-,$(platform))/usermode-helper)

IMAGE_NAMES=init ca-certificates firmware rngd sysctl sysfs modprobe dhcpcd openntpd getty

# we are building these using the upstream dockerfiles for now
# the plan is to convert these to our standard build process
# We are using this pattern of setting the dockerfile_folder instead of local to the target
# because this project uses the combine-images target which overrides the dockerfile_folder
# to the standard combine-images dockerfile
# If we setup our override in this Makefile at the target it does not allow the combine-images
# to override it
IMAGES_USING_UPSTREAM_DOCKERFILES=ca-certificates dhcpcd firmware getty init modprobe openntpd
DOCKERFILE_FOLDER=$(if $(filter $(IMAGE_NAME),$(IMAGES_USING_UPSTREAM_DOCKERFILES)),$(REPO)/pkg/$(IMAGE_NAME),./docker/linux/$(IMAGE_NAME))

INIT_IMAGE_COMPONENT=linuxkit/init
CA_CERTIFICATES_IMAGE_COMPONENT=linuxkit/ca-certificates
FIRMWARE_IMAGE_COMPONENT=linuxkit/firmware
RNGD_IMAGE_COMPONENT=linuxkit/rngd
SYSCTL_IMAGE_COMPONENT=linuxkit/sysctl
SYSFS_IMAGE_COMPONENT=linuxkit/sysfs
MODPROBE_IMAGE_COMPONENT=linuxkit/modprobe
DHCPCD_IMAGE_COMPONENT=linuxkit/dhcpcd
OPENNTPD_IMAGE_COMPONENT=linuxkit/openntpd
GETTY_IMAGE_COMPONENT=linuxkit/getty

# when using yq to conver to json 0666 is being parsed as a number and converted to 438
# replace it back to 0600
# https://github.com/mikefarah/yq/issues/1634
MOBY_CONFIG='$(shell yq '.config' $(REPO)/pkg/$(IMAGE_NAME)/build.yml | yq -o=json -I=0 | sed 's/438/"0666"/g')'

# we need to set IMAGE_BUILD_ARGS here even though its the same as the default. 
# it is set in Common.mk on the images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L799)
# and the combine-images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L846)
# since combine-images has images as prereq target, the ?= does not really behavior as one might expect.
# the images target being the actual action, its version of the set takes prioirty and resets to empty
# setting it explicitly to empty here takes allows the combine-images override to take proirty
IMAGE_BUILD_ARGS=
VERSION=$(shell git -C $(REPO) describe --tags)
GIT_COMMIT=$(shell git -C $(REPO) rev-parse HEAD)

HAS_S3_ARTIFACTS=true

FIX_LICENSES_SUB_MOD_TARGETS=$(REPO)/src/cmd/linuxkit/LICENSE $(REPO)/pkg/rngd/LICENSE $(REPO)/pkg/sysctl/LICENSE $(REPO)/pkg/sysfs/LICENSE $(REPO)/pkg/init/LICENSE
FIX_LICENSES_HYPERKIT_TARGE=$(REPO)/src/cmd/linuxkit/vendor/github.com/moby/hyperkit/go/README.md

BUILDSPECS=$(BUILDSPECS_FOR_COMBINE_IMAGES)
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true

include $(BASE_DIRECTORY)/Common.mk

# build usermode-helper with binary targets
validate-checksums: $(USERMODE_HELPER_TARGETS)

dhcpcd/images/% getty/images/% modprobe/images/% openntpd/images/% rngd/images/% sysctl/images/% sysfs/images/%: IMAGE_BUILD_ARGS+=MOBY_CONFIG

$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_SUB_MOD_TARGETS) $(FIX_LICENSES_HYPERKIT_TARGE)

$(OUTPUT_BIN_DIR)/linux-%/linuxkit: EXTRA_GO_LDFLAGS=-X github.com/linuxkit/linuxkit/src/cmd/linuxkit/version.GitCommit=$(GIT_COMMIT) -X github.com/linuxkit/linuxkit/src/cmd/linuxkit/version.Version=$(VERSION)

$(OUTPUT_BIN_DIR)/linux-%/rngd: CGO_CREATE_BINARIES=true
$(OUTPUT_BIN_DIR)/linux-%/rngd: CGO_CFLAGS_ALLOW=(-mrdrnd|-mrdseed)
$(OUTPUT_BIN_DIR)/linux-%/rngd: EXTRA_GO_LDFLAGS=-extldflags -static
$(OUTPUT_BIN_DIR)/linux-%/rngd: EXTRA_GOBUILD_FLAGS=-tags netgo,osusergo,static_build
$(OUTPUT_DIR)/rngd/attribution/go-license.csv: CGO_CREATE_BINARIES=true

ca-certificates/images/% dhcpcd/images/% firmware/images/% getty/images/% modprobe/images/% openntpd/images/%: IMAGE_CONTEXT_DIR=$(REPO)/pkg/$(IMAGE_NAME)

$(OUTPUT_BIN_DIR)/%/usermode-helper: MAKEFLAGS=
$(OUTPUT_BIN_DIR)/%/usermode-helper: PLATFORM=$(subst -,/,$*)
$(OUTPUT_BIN_DIR)/%/usermode-helper: $(GIT_PATCH_TARGET) | $$(call ENABLE_DOCKER_PLATFORM,$$(PLATFORM))
	@mkdir -p $(MAKE_ROOT)/$(@D); \
	LDFLAGS=-static CFLAGS=-Werror $(MAKE) -C $(REPO)/pkg/init usermode-helper; \
	mv $(REPO)/pkg/init/usermode-helper $@

$(FIX_LICENSES_SUB_MOD_TARGETS): $(GIT_PATCH_TARGET)
# the src/cmd/linuxkit dir does not have the license file, copy from repo root
	cp $(REPO)/LICENSE $@

$(FIX_LICENSES_HYPERKIT_TARGE): | $(GO_MOD_DOWNLOAD_TARGETS)
# the moby/hyperkit dep in linuxkit is quite old, however the version required
# did have the license in the README.md upstream. Copy readme from that commit
	wget -q --retry-connrefused https://raw.githubusercontent.com/moby/hyperkit/d65b09c1c28a2bfb6a976c86ecd885d2ee4c71d3/README.md -O \
		$@


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
REPO=linuxkit
REPO_OWNER=linuxkit

BASE_IMAGE_NAME?=eks-distro-minimal-base

BINARY_TARGET_FILES=linuxkit
GO_MOD_PATHS=src/cmd/linuxkit

IMAGE_NAMES=init ca-certificates firmware rngd sysctl sysfs modprobe dhcpcd openntpd getty
# we are building these using the upstream dockerfiles for now
# the plan is to convert these to our standard build process
DOCKERFILE_FOLDER=$(REPO)/pkg/$(IMAGE_NAME)
IMAGE_CONTEXT_DIR=$(REPO)/pkg/$(IMAGE_NAME)

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
MOBY_CONFIG='$(shell set -x && yq '.config' $(REPO)/pkg/$(IMAGE_NAME)/build.yml | yq -o=json -I=0 | sed 's/438/"0666"/g')'

# we need to set IMAGE_BUILD_ARGS here even though its the same as the default. 
# it is set in Common.mk on the images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L799)
# and the combine-images target (https://github.com/aws/eks-anywhere-build-tooling/blob/8b6d6d66974e881b22e3c9c8ea29adc26f7df5fd/Common.mk#L846)
# since combine-images has images as prereq target, the ?= does not really behavior as one might expect.
# the images target being the actual action, its version of the set takes prioirty and resets to empty
# setting it explicitly to empty here takes allows the combine-images override to take proirty
IMAGE_BUILD_ARGS=

VERSION=$(shell git -C $(REPO) describe --tags)
GIT_COMMIT=$(shell git -C $(REPO) rev-parse HEAD)
EXTRA_GO_LDFLAGS=-X github.com/linuxkit/linuxkit/src/cmd/linuxkit/version.GitCommit=$(GIT_COMMIT) -X github.com/linuxkit/linuxkit/src/cmd/linuxkit/version.Version=$(VERSION)

HAS_S3_ARTIFACTS=true

FIX_LICENSES_LINUXKIT_TARGET=$(REPO)/src/cmd/linuxkit/LICENSE
FIX_LICENSES_HYPERKIT_TARGE=$(REPO)/src/cmd/linuxkit/vendor/github.com/moby/hyperkit/go/README.md

BUILDSPECS=buildspec.yml buildspecs/combine-images.yml
BUILDSPEC_1_COMPUTE_TYPE=BUILD_GENERAL1_LARGE
BUILDSPEC_1_VARS_KEYS=IMAGE_PLATFORMS
BUILDSPEC_1_VARS_VALUES=IMAGE_PLATFORMS
BUILDSPEC_1_ARCH_TYPES=LINUX_CONTAINER ARM_CONTAINER
BUILDSPEC_2_DEPENDS_ON_OVERRIDE=linuxkit_linuxkit_linux_amd64 linuxkit_linuxkit_linux_arm64

# Since we build the arm and amd binaries on difference instances in codebuild
# we do not want to delete missing files when s3 sync-ing from local to the bucket
# since this would clobber the other build's artifacts
UPLOAD_DO_NOT_DELETE=true

include $(BASE_DIRECTORY)/Common.mk

dhcpcd/images/% getty/images/% modprobe/images/% openntpd/images/% rngd/images/% sysctl/images/% sysfs/images/%: IMAGE_BUILD_ARGS+=MOBY_CONFIG

$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_LINUXKIT_TARGET) $(FIX_LICENSES_HYPERKIT_TARGE)

$(FIX_LICENSES_LINUXKIT_TARGET): $(REPO)
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

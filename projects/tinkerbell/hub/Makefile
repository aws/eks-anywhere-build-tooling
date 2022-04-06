BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.15"

REPO=hub
REPO_OWNER=tinkerbell

BASE_IMAGE_NAME?=eks-distro-minimal-base

CEXEC_IMAGE_COMPONENT=tinkerbell/hub/cexec
KEXEC_IMAGE_COMPONENT=tinkerbell/hub/kexec
IMAGE2DISK_IMAGE_COMPONENT=tinkerbell/hub/image2disk
OCI2DISK_IMAGE_COMPONENT=tinkerbell/hub/oci2disk
WRITEFILE_IMAGE_COMPONENT=tinkerbell/hub/writefile
REBOOT_IMAGE_COMPONENT=tinkerbell/hub/reboot

BINARY_TARGET_FILES=cexec kexec image2disk oci2disk writefile
SOURCE_PATTERNS=. . . . .
IMAGE_NAMES=$(BINARY_TARGET_FILES) reboot

DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

IMAGE_BUILD_ARGS=ATTRIBUTION_FILE
ATTRIBUTION_FILE=$(call TO_UPPER,$(IMAGE_NAME))_ATTRIBUTION.txt

ATTRIBUTION_TARGETS=$(foreach binary,$(BINARY_TARGET_FILES),$(call TO_UPPER,$(binary))_ATTRIBUTION.txt)
GATHER_LICENSES_TARGETS=$(foreach binary,$(BINARY_TARGET_FILES),$(OUTPUT_DIR)/$(binary)/attribution/go-license.csv)
GO_MOD_DOWNLOAD_TARGETS=$(foreach binary,$(BINARY_TARGET_FILES),$(REPO)/actions/$(binary)/v1/eks-anywhere-go-mod-download)


define HUB_ACTION_TARGETS
    $(OUTPUT_BIN_DIR)/linux-%/$(1) $(REPO)/actions/$(1)/v1/eks-anywhere-go-mod-download $(OUTPUT_DIR)/$(1)/attribution/go-license.csv: REPO_SUBPATH=actions/$(1)/v1
    $(OUTPUT_DIR)/$(1)/attribution/go-license.csv $(2)_ATTRIBUTION.txt: LICENSE_PACKAGE_FILTER=.
    $(OUTPUT_DIR)/$(1)/attribution/go-license.csv $(2)_ATTRIBUTION.txt: LICENSES_OUTPUT_DIR=$(OUTPUT_DIR)/$(1)

endef

include $(BASE_DIRECTORY)/Common.mk


$(eval $(foreach binary,$(BINARY_TARGET_FILES),$(call HUB_ACTION_TARGETS,$(binary),$(call TO_UPPER,$(binary)))))

# reboot action image requires touch.
# We are using eks-distro-minimal-base-glibc as the base to install touch
reboot/images/%: BASE_IMAGE_NAME=eks-distro-minimal-base-glibc

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################

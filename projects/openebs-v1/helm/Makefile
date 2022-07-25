BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.17"

REPO=charts
REPO_OWNER=openebs

HAS_HELM_CHART=true
HELM_DIRECTORY=./charts/openebs
HELM_DESTINATION_REPOSITORY=openebs/openebs
HELM_IMAGE_LIST=openebs/cvc-operator openebs/pool-manager openebs/volume-manager openebs/jiva openebs/m-exporter openebs/admission-server openebs/cstor-pool-mgmt openebs/cstor-volume-mgmt openebs/maya-apiserver openebs/openebs-provisioner openebs/snapshot-controller openebs/snapshot-pv-provisioner

include $(BASE_DIRECTORY)/Common.mk


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################
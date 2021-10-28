BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=cluster-api-provider-aws
REPO_OWNER=kubernetes-sigs

LICENSE_PACKAGE_FILTER?=./cmd/clusterawsadm sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks sigs.k8s.io/cluster-api-provider-aws/controlplane/eks ./

BINARY_TARGET=$(OUTPUT_BIN_DIR)/linux-amd64/clusterawsadm

BASE_IMAGE_NAME=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CLUSTER_API_AWS_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api-provider-aws/cluster-api-aws-controller
EKS_BOOTSTRAP_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api-provider-aws/eks-bootstrap-controller
EKS_CONTROL_PLANE_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api-provider-aws/eks-control-plane-controller

include $(BASE_DIRECTORY)/Common.mk


build: s3-artifacts
release: upload-artifacts

.PHONY: local-images
local-images: cluster-api-aws-controller/images/amd64 eks-bootstrap-controller/images/amd64 eks-control-plane-controller/images/amd64

.PHONY: images
images: cluster-api-aws-controller/images/push eks-bootstrap-controller/images/push eks-control-plane-controller/images/push

.PHONY: s3-artifacts
s3-artifacts: tarballs create-manifests
	$(BUILD_LIB)/create_release_checksums.sh $(ARTIFACTS_PATH)
	$(BUILD_LIB)/validate_artifacts.sh $(MAKE_ROOT) $(ARTIFACTS_PATH) $(GIT_TAG)

.PHONY: create-manifests
create-manifests:
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

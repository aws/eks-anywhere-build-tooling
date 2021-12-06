BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
GIT_TAG?=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.13"

REPO=cluster-api-provider-aws
REPO_OWNER=kubernetes-sigs

BINARY_TARGET_FILES=manager eks-bootstrap-manager eks-controlplane-manager clusterawsadm
SOURCE_PATTERNS=. sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks sigs.k8s.io/cluster-api-provider-aws/controlplane/eks ./cmd/clusterawsadm

BASE_IMAGE_NAME=eks-distro-minimal-base
DOCKERFILE_FOLDER=./docker/linux/$(IMAGE_NAME)

CLUSTER_API_AWS_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api-provider-aws/cluster-api-aws-controller
EKS_BOOTSTRAP_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api-provider-aws/eks-bootstrap-controller
EKS_CONTROL_PLANE_CONTROLLER_IMAGE_COMPONENT=kubernetes-sigs/cluster-api-provider-aws/eks-control-plane-controller

IMAGE_NAMES= cluster-api-aws-controller eks-bootstrap-controller eks-control-plane-controller

HAS_S3_ARTIFACTS=true

include $(BASE_DIRECTORY)/Common.mk


s3-artifacts: create-manifests

.PHONY: create-manifests
create-manifests: tarballs
	build/create_manifests.sh $(REPO) $(OUTPUT_DIR) $(ARTIFACTS_PATH) $(GIT_TAG) $(IMAGE_REPO) $(IMAGE_TAG)

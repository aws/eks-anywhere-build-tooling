ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-docker-client
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

WORKDIR /

ARG EKS_A_TOOL_BINARY_DIR=/eks-a-tools/binary
ARG EKS_A_TOOL_LICENSE_DIR=/eks-a-tools/licenses

COPY ./_output/$TARGETOS-$TARGETARCH/eks-a-tools/binary $EKS_A_TOOL_BINARY_DIR
COPY ./_output/$TARGETOS-$TARGETARCH/eks-a-tools/licenses $EKS_A_TOOL_LICENSE_DIR

ENV PATH="${EKS_A_TOOL_BINARY_DIR}:${PATH}"

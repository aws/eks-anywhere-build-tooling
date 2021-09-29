ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

RUN yum install tar -y && \
    yum clean all && \
    rm -rf /var/cache/yum

WORKDIR /

ARG EKS_A_DIAGNOSTIC_COLLECTOR_BINARY_DIR=/eks-a-diagnostic-collector-tools/binary
ARG EKS_A_DIAGNOSTIC_COLLECTOR_TOOL_LICENSE_DIR=/eks-a-diagnostic-collector-tools/licenses

COPY ./_output/eks-a-diagnostic-collector-tools/binary $EKS_A_DIAGNOSTIC_COLLECTOR_BINARY_DIR
COPY ./_output/eks-a-diagnostic-collector-tools/licenses $EKS_A_DIAGNOSTIC_COLLECTOR_TOOL_LICENSE_DIR
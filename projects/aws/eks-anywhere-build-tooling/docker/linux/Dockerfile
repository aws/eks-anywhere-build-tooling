ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-docker-client
ARG BUILDER_IMAGE
FROM $BASE_IMAGE
FROM $BUILDER_IMAGE as builder

# sleep is needed to reuse the container during cli runs instead of starting new containers each time
# but it comes from the coreutils package which pulls in more than we need
# manually installing sleep from the rpm
# sleep only depends on glibc so there are no additional deps needed
RUN set -x && \
    yumdownloader --destdir=/tmp/downloads coreutils && \
    cd /newroot && \
    rpm2cpio /tmp/downloads/coreutils*.rpm | cpio -idv ./usr/bin/sleep

FROM $BASE_IMAGE
ARG TARGETARCH
ARG TARGETOS

WORKDIR /

ARG EKS_A_TOOL_BINARY_DIR=/eks-a-tools/binary
ARG EKS_A_TOOL_LICENSE_DIR=/eks-a-tools/licenses

COPY --from=builder /newroot/usr/bin/sleep /usr/bin/
COPY ./_output/dependencies/$TARGETOS-$TARGETARCH/eks-a-tools/binary $EKS_A_TOOL_BINARY_DIR
COPY ./_output/dependencies/$TARGETOS-$TARGETARCH/eks-a-tools/licenses $EKS_A_TOOL_LICENSE_DIR

ENV PATH="${EKS_A_TOOL_BINARY_DIR}:${PATH}"

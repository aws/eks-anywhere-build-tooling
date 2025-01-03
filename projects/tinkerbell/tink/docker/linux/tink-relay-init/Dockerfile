ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nsenter
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE as builder

WORKDIR /newroot

RUN set -x && \
    install_binary /usr/sbin/awk /usr/bin/nmap /usr/bin/head /usr/sbin/ip /usr/bin/echo /usr/bin/env /usr/bin/getopt && \
    install_rpm bash && \
    cleanup "deps"

ARG TARGETARCH
ARG TARGETOS

FROM $BASE_IMAGE

COPY --from=builder /newroot /

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-glibc
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE as builder

WORKDIR /newroot

RUN set -x && \
    enable_extra docker && \
    install_binary /usr/bin/docker-init && \
    cleanup "deps"

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/isc-projects/dhcp/dhcrelay /usr/sbin
COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/isc-projects/dhcp/LICENSES /DHCP_LICENSES

COPY --from=builder /newroot /

ENTRYPOINT ["/usr/bin/docker-init", "--", "dhcrelay", "-d"]
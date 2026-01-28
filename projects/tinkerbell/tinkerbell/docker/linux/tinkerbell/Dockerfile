# Tinkerbell mono-repo server image
# Requires ipmitool for bmclib BMC operations

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-glibc
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE as ipmitool-builder

WORKDIR /newroot

# Install ipmitool and its dependencies for bmclib
RUN set -x && \
    clean_install "systemd" true true && \
    clean_install "ipmitool openssl-libs ncurses-libs readline" && \
    remove_package "bash coreutils gawk info grep sed shadow-utils systemd-sysv" && \
    remove_package "systemd" true && \
    cleanup "glibc"

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY --from=ipmitool-builder /newroot /

COPY _output/bin/tinkerbell/$TARGETOS-$TARGETARCH/tinkerbell /tinkerbell
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT ["/tinkerbell"]

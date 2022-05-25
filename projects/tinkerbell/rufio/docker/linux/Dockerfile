ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-glibc
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE as ipmitool-builder

WORKDIR /

RUN yum install -y \
    autoconf \
    automake \
    git \
    libtool \
    make \
    readline-devel \
    openssl-devel

COPY ipmitool /ipmitool
RUN cd ipmitool && \
    ./bootstrap && \
    ./configure \
        --prefix=/usr/local \
        --enable-ipmievd \
        --enable-ipmishell \
        --enable-intf-lan \
        --enable-intf-lanplus \
        --enable-intf-open && \
    make && \
    make DESTDIR=/newroot install

WORKDIR /newroot

RUN set -x && \
    clean_install "openssl-libs ncurses-libs readline" && \
    cleanup "glibc"

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY --from=ipmitool-builder /newroot /

COPY _output/bin/rufio/$TARGETOS-$TARGETARCH/manager /usr/bin/manager
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 65532:65532

ENTRYPOINT ["/manager"]
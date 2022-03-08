ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
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

COPY _output/files/pbnj /
COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/grpc-ecosystem/grpc-health-probe /usr/local/bin
COPY _output/bin/pbnj/$TARGETOS-$TARGETARCH/pbnj /usr/bin/pbnj
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENV GIN_MODE release
USER pbnj
EXPOSE 50051 9090 8080

ENTRYPOINT ["/usr/bin/pbnj"]
CMD ["server"]

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE as ipmitool-builder

ARG IPMITOOL_CLONE_URL
ARG IPMITOOL_GIT_TAG

WORKDIR /

RUN yum install -y \
    autoconf \
    automake \
    git \
    libtool \
    make \
    readline-devel

RUN git clone $IPMITOOL_CLONE_URL && \
    cd ipmitool && \
    git checkout $IPMITOOL_COMMIT && \
    ./bootstrap && \
    ./configure \
        --prefix=/usr/local \
        --enable-ipmievd \
        --enable-ipmishell \
        --enable-intf-lan \
        --enable-intf-lanplus \
        --enable-intf-open && \
    make && \
    make install

WORKDIR /newroot

RUN set -x && \
    clean_install "openssl-libs ncurses-libs readline" true true && \
    cleanup "glibc"

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/grpc-ecosystem/grpc-health-probe /usr/local/bin
COPY _output/bin/pbnj/$TARGETOS-$TARGETARCH/pbnj /usr/bin/pbnj
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

COPY --from=ipmitool-builder /usr/local/bin/ipmitool /usr/local/bin
COPY --from=ipmitool-builder /ipmitool/COPYING /LICENSES/github.com/ipmitool/ipmitool/COPYING

COPY --from=ipmitool-builder /newroot /

ENV GIN_MODE release
USER pbnj
EXPOSE 50051 9090 8080

ENTRYPOINT ["/usr/bin/pbnj"]
CMD ["server"]

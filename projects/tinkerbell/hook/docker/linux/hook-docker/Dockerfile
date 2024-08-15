ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as docker-builder

ARG TARGETARCH

WORKDIR /newroot

RUN set -x && \
    enable_extra docker && \
    install_rpm containerd docker e2fsprogs runc pigz procps-ng  && \
    if [ $TARGETARCH = "amd64" ]; then BUSYBOX_ARCH="x86_64"; else BUSYBOX_ARCH="armv81"; fi && \
    curl https://busybox.net/downloads/binaries/1.31.0-defconfig-multiarch-musl/busybox-$BUSYBOX_ARCH -o /newroot/usr/bin/busybox && \
    chmod +x /newroot/usr/bin/busybox && \
    ln -sf /usr/bin/busybox /newroot/usr/sbin/reboot && \
    ln -sf /usr/bin/docker-init /newroot/usr/local/bin/docker-init && \
    ln -sf /usr/bin/dockerd /newroot/usr/local/bin/dockerd && \
    cleanup "tink-docker"


FROM docker-builder as dind-builder

RUN set -x && \
    install_rpm bash && \
    install_binary /usr/bin/cp /usr/bin/rm


FROM $BASE_IMAGE as hook-docker

ARG TARGETARCH
ARG TARGETOS

WORKDIR /

COPY --from=docker-builder /newroot /

COPY _output/bin/hook/$TARGETOS-$TARGETARCH/hook-docker /usr/bin/tink-docker
COPY _output/hook-docker/LICENSES /LICENSES
COPY HOOK_DOCKER_ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT [ "/usr/bin/tink-docker" ]

FROM hook-docker as hook-dind

COPY --from=dind-builder /newroot /

VOLUME /var/lib/docker

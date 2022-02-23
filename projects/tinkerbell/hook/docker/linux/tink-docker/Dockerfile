ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as docker-builder

WORKDIR /newroot

RUN set -x && \
    clean_install "systemd" true true && \
    clean_install "docker procps e2fsprogs" && \
    remove_package "bash coreutils gawk info sed shadow-utils grep" && \
    remove_package "systemd" true && \
    cleanup "tink-docker" && \
    ln -sf /usr/bin/docker-init /newroot/usr/local/bin/docker-init && \
    ln -sf /usr/bin/dockerd /newroot/usr/local/bin/dockerd

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

WORKDIR /

COPY --from=docker-builder /newroot /

COPY _output/bin/hook/$TARGETOS-$TARGETARCH/tink-docker /usr/bin/tink-docker
COPY _output/tink-docker/LICENSES /LICENSES
COPY TINK_DOCKER_ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT [ "/usr/bin/tink-docker" ]

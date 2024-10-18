ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-glibc
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as builder

RUN set -x && \
    mkdir -p /home/_ntp && \    
    chmod 750 /home/_ntp

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS
ARG MOBY_CONFIG

LABEL "org.mobyproject.config"="${MOBY_CONFIG}"

COPY --from=builder /home/_ntp/ /home/_ntp/
COPY _output/files/openntpd /
COPY _output/bin/linuxkit/$TARGETOS-$TARGETARCH/openntpd/ /
COPY linuxkit/pkg/openntpd/etc/ /etc/

CMD ["/usr/sbin/ntpd", "-d"]

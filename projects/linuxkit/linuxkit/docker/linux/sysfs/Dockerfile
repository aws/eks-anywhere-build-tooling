FROM scratch

ARG TARGETARCH
ARG TARGETOS

ARG MOBY_CONFIG
LABEL "org.mobyproject.config"="${MOBY_CONFIG}"

COPY _output/bin/linuxkit/$TARGETOS-$TARGETARCH/sysfs /usr/bin/sysfs
COPY SYSFS_ATTRIBUTION.txt /

COPY linuxkit/pkg/sysfs/etc/ /etc/

CMD ["/usr/bin/sysfs"]

FROM scratch

ARG TARGETARCH
ARG TARGETOS

ARG MOBY_CONFIG
LABEL "org.mobyproject.config"="${MOBY_CONFIG}"

COPY _output/bin/linuxkit/$TARGETOS-$TARGETARCH/sysctl /usr/bin/sysctl
COPY SYSCTL_ATTRIBUTION.txt /

COPY linuxkit/pkg/sysctl/etc/ /etc/

CMD ["/usr/bin/sysctl"]

FROM scratch

ARG TARGETARCH
ARG TARGETOS

ARG MOBY_CONFIG
LABEL "org.mobyproject.config"="${MOBY_CONFIG}"

COPY _output/bin/linuxkit/$TARGETOS-$TARGETARCH/rngd /sbin/rngd
COPY RNGD_ATTRIBUTION.txt /

CMD ["/sbin/rngd"]

ARG BASE_IMAGE
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/node_exporter/$TARGETOS-$TARGETARCH/node_exporter /bin/node_exporter
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

EXPOSE      9100
USER        1000
ENTRYPOINT  [ "/bin/node_exporter" ]

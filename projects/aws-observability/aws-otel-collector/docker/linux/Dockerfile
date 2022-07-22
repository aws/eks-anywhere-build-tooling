ARG BASE_IMAGE
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/aws-otel-collector/$TARGETOS-$TARGETARCH/aoc /opt/aws-observability/aws-otel-collector

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT ["/opt/aws-observability/aws-otel-collector"]

ARG BASE_IMAGE
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/cstor-operators/$TARGETOS-$TARGETARCH/cspc-operator /opt/cstor-operators/cspc-operator

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000

ENTRYPOINT ["/opt/cstor-operators/cspc-operator"]
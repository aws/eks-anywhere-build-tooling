ARG BASE_IMAGE
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/cstor-operators/$TARGETOS-$TARGETARCH/cvc-operator /opt/cstor-operators/cvc-operator

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000

ENTRYPOINT ["/opt/cstor-operators/cvc-operator"]
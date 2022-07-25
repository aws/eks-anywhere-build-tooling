ARG BASE_IMAGE
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/maya/$TARGETOS-$TARGETARCH/cstor-volume-mgmt /opt/maya/cstor-volume-mgmt

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000

ENTRYPOINT ["/opt/maya/cstor-volume-mgmt"]
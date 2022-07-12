ARG BASE_IMAGE
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/maya/$TARGETOS-$TARGETARCH/cstor-pool-mgmt /opt/maya/cstor-pool-mgmt

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000

ENTRYPOINT ["/opt/maya/cstor-pool-mgmt"]
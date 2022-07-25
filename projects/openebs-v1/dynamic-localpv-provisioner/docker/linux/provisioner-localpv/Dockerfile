ARG BASE_IMAGE
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/dynamic-localpv-provisioner/$TARGETOS-$TARGETARCH/provisioner-localpv /opt/dynamic-localpv-provisioner/provisioner-localpv

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT ["/opt/dynamic-localpv-provisioner/provisioner-localpv"]
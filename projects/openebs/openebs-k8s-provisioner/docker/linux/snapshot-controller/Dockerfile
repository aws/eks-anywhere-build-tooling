ARG BASE_IMAGE
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/openebs-k8s-provisioner/$TARGETOS-$TARGETARCH/snapshot-controller /opt/openebs-k8s-provisioner/snapshot-controller

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000

ENTRYPOINT ["/opt/openebs-k8s-provisioner/snapshot-controller"]
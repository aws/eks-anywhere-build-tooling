ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

ARG ATTRIBUTION_FILE

COPY _output/bin/hub/$TARGETOS-$TARGETARCH/oci2disk /usr/bin/oci2disk
COPY _output/oci2disk/LICENSES /LICENSES
COPY $ATTRIBUTION_FILE /ATTRIBUTION.txt

ENTRYPOINT ["/usr/bin/oci2disk"]

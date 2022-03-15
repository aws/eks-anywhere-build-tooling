ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

ARG ATTRIBUTION_FILE

COPY _output/bin/hub/$TARGETOS-$TARGETARCH/image2disk /usr/bin/image2disk
COPY _output/image2disk/LICENSES /LICENSES
COPY $ATTRIBUTION_FILE /ATTRIBUTION.txt

ENTRYPOINT ["/usr/bin/image2disk"]

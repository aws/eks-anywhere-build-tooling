ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/tink/$TARGETOS-$TARGETARCH/tink-worker /usr/bin/tink-worker
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT [ "/usr/bin/tink-worker" ]

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/hook/$TARGETOS-$TARGETARCH/bootkit /usr/bin/bootkit
COPY _output/bootkit/LICENSES /LICENSES
COPY BOOTKIT_ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT [ "/usr/bin/bootkit" ]

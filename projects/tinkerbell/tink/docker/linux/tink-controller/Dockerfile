ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nonroot
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

WORKDIR /

COPY _output/bin/tink/$TARGETOS-$TARGETARCH/tink-controller /usr/bin/tink-controller
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

EXPOSE 42113
EXPOSE 42114

ENTRYPOINT ["/usr/bin/tink-controller"]

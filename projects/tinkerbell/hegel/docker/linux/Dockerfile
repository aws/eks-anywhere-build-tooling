ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/files/hegel /
COPY _output/bin/hegel/$TARGETOS-$TARGETARCH/hegel /usr/bin/hegel
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER tinkerbell

EXPOSE 50060
EXPOSE 50061

ENTRYPOINT ["/usr/bin/hegel"]

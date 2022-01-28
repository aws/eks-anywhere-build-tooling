ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/boots/$TARGETOS-$TARGETARCH/boots /usr/bin/boots
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

EXPOSE 67 69 80

ENTRYPOINT ["/usr/bin/boots"]

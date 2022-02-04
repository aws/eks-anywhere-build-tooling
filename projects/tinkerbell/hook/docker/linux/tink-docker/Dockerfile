ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/hook/$TARGETOS-$TARGETARCH/tink-docker /usr/bin/tink-docker
COPY _output/tink-docker/LICENSES /LICENSES
COPY TINK_DOCKER_ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT [ "/usr/bin/tink-docker" ]

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/eks-anywhere-packages/$TARGETOS-$TARGETARCH/ecrtokenrefresher /ecr-refresh
COPY _output/ecrtokenrefresher/LICENSES /LICENSES
COPY ECRTOKENREFRESHER_ATTRIBUTION.txt /ATTRIBUTION.txt

USER 65532
ENTRYPOINT ["/ecr-refresh"]

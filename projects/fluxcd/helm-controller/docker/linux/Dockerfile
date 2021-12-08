ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/files/helm-controller /
COPY _output/bin/helm-controller/$TARGETOS-$TARGETARCH/helm-controller /usr/local/bin/helm-controller
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER controller

ENTRYPOINT ["helm-controller"]

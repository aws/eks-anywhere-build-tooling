ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/files/cert-manager-controller /
COPY _output/bin/cert-manager/$TARGETOS-$TARGETARCH/cert-manager-controller /usr/bin/cert-manager
COPY _output/cert-manager-controller/LICENSES /LICENSES
COPY CERT_MANAGER_CONTROLLER_ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000
ENTRYPOINT ["/usr/bin/cert-manager"]

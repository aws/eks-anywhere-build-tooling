ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/files/cert-manager-acmesolver /
COPY _output/bin/cert-manager/$TARGETOS-$TARGETARCH/cert-manager-acmesolver /usr/bin/acmesolver
COPY _output/cert-manager-acmesolver/LICENSES /LICENSES
COPY CERT_MANAGER_ACMESOLVER_ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000
ENTRYPOINT ["/usr/bin/acmesolver"]

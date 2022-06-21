ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/files/cert-manager-cainjector /
COPY _output/bin/cert-manager/$TARGETOS-$TARGETARCH/cert-manager-cainjector /usr/bin/cainjector
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000
ENTRYPOINT ["/usr/bin/cainjector"]

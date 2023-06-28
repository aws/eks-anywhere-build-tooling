ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/files/cert-manager-ctl /
COPY _output/bin/cert-manager/$TARGETOS-$TARGETARCH/cert-manager-ctl /usr/bin/ctl
COPY _output/cert-manager-ctl/LICENSES /LICENSES
COPY CERT_MANAGER_WEBHOOK_ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000
ENTRYPOINT ["/usr/bin/ctl"]

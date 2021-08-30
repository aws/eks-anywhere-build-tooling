ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

COPY _output/bin/cert-manager/linux-amd64/cert-manager-webhook /usr/bin/webhook
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT ["/usr/bin/webhook"]

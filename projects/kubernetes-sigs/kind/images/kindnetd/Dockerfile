ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-iptables
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY --chown=root:root _output/bin/kind/$TARGETOS-$TARGETARCH/kindnetd /bin/kindnetd
COPY _output/kindnetd/LICENSES /KINDNETD_LICENSES
COPY KINDNETD_ATTRIBUTION.txt /KINDNETD_ATTRIBUTION.txt

CMD ["/bin/kindnetd"]

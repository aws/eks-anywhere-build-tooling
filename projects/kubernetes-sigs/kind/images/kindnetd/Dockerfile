ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-iptables
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY --chown=root:root _output/bin/kind/$TARGETOS-$TARGETARCH/kindnetd /bin/kindnetd
COPY _output/kindnetd/LICENSES /KINDNETD_LICENSES
COPY KINDNETD_ATTRIBUTION.txt /KINDNETD_ATTRIBUTION.txt

# Opt into using the iptables-wrapper script to determine iptables mode
RUN ["update-alternatives", "--set", "iptables", "/usr/sbin/iptables-wrapper"]

CMD ["/bin/kindnetd"]

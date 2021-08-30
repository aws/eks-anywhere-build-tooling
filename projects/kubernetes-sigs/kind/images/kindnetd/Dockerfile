ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-iptables
FROM $BASE_IMAGE

COPY --chown=root:root _output/bin/kind/linux-amd64/kindnetd /bin/kindnetd
COPY _output/kindnetd/LICENSES /KINDNETD_LICENSES
COPY KINDNETD_ATTRIBUTION.txt /KINDNETD_ATTRIBUTION.txt

CMD ["/bin/kindnetd"]

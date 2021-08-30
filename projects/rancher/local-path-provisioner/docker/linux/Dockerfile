ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

COPY _output/bin/local-path-provisioner/linux-amd64/local-path-provisioner-linux-amd64 /usr/local/bin/local-path-provisioner
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

CMD ["/usr/local/bin/local-path-provisioner"]

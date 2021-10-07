ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/local-path-provisioner/$TARGETOS-$TARGETARCH/local-path-provisioner /usr/local/bin/local-path-provisioner
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

CMD ["/usr/local/bin/local-path-provisioner"]

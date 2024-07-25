ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/boots/$TARGETOS-$TARGETARCH/smee /usr/bin/boots
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt
COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/tinkerbell/ipxedust/LICENSES /IPXE_LICENSES
COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/tinkerbell/ipxedust/ATTRIBUTION.txt /IPXE_ATTRIBUTION.txt

EXPOSE 67 69 80

ENTRYPOINT ["/usr/bin/boots"]

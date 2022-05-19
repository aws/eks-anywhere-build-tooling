ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

ARG ATTRIBUTION_FILE

COPY _output/bin/hub/$TARGETOS-$TARGETARCH/writefile /usr/bin/writefile
COPY _output/writefile/LICENSES /LICENSES
COPY $ATTRIBUTION_FILE /ATTRIBUTION.txt

COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/torvalds/linux/tools/bootconfig /usr/bin/bootconfig
COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/torvalds/linux/ATTRIBUTION.txt /BOOTCONFIG_ATTRIBUTION.txt
COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/torvalds/linux/LICENSES /BOOTCONFIG_LICENSES

ENTRYPOINT ["/usr/bin/writefile"]

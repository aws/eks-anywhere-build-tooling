ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nonroot
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

WORKDIR /

COPY _output/bin/cluster-api-provider-vsphere/$TARGETOS-$TARGETARCH/manager /manager
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 65532
ENTRYPOINT ["/manager"]

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/files/ /config

COPY _output/bin/vsphere-csi-driver/$TARGETOS-$TARGETARCH/vsphere-csi-syncer /bin/vsphere-csi-syncer
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT ["/bin/vsphere-csi-syncer"]

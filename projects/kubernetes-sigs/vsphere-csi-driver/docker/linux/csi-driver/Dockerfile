ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-csi
FROM $BASE_IMAGE

COPY _output/bin/vsphere-csi-driver/linux-amd64/vsphere-csi-driver /bin/vsphere-csi-driver
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT ["/bin/vsphere-csi-driver"]

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG RELEASE_BRANCH
ARG TARGETARCH
ARG TARGETOS

COPY _output/$RELEASE_BRANCH/bin/cloud-provider-vsphere/$TARGETOS-$TARGETARCH/vsphere-cloud-controller-manager /bin/vsphere-cloud-controller-manager
COPY _output/$RELEASE_BRANCH/LICENSES /LICENSES
COPY $RELEASE_BRANCH/ATTRIBUTION.txt /ATTRIBUTION.txt

USER nobody

ENTRYPOINT ["/bin/vsphere-cloud-controller-manager"]

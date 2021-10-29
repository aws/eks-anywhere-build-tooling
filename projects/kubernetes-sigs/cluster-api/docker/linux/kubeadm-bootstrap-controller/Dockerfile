ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nonroot
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

WORKDIR /

COPY _output/bin/cluster-api/$TARGETOS-$TARGETARCH/kubeadm-bootstrap-manager /manager
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT ["/manager"]

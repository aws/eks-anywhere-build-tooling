ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG RELEASE_BRANCH
ARG TARGETARCH
ARG TARGETOS

COPY _output/$RELEASE_BRANCH/bin/autoscaler/$TARGETOS-$TARGETARCH/cluster-autoscaler /cluster-autoscaler
COPY _output/$RELEASE_BRANCH/LICENSES /LICENSES
COPY $RELEASE_BRANCH/ATTRIBUTION.txt /ATTRIBUTION.txt

USER nobody

CMD ["/cluster-autoscaler"]

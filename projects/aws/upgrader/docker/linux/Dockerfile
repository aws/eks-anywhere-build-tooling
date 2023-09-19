ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-glibc
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as builder

RUN set -x && \
    install_binary /usr/bin/nsenter \
        /usr/bin/cp && \
    cleanup "deps"

FROM $BASE_IMAGE

ARG RELEASE_BRANCH
ARG TARGETARCH
ARG TARGETOS

COPY --from=builder /newroot /

COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksa/containerd/containerd /eksa-upgrades/binaries/containerd
COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksa/kubernetes-sigs/cri-tools /eksa-upgrades/binaries/containerd/usr/local/bin
COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksd/cni-plugins /eksa-upgrades/binaries/cni-plugins/opt/cni/bin
COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksd/kubernetes/client/bin /eksa-upgrades/binaries/kubernetes/usr/bin
COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksd/kubernetes/server/bin/kubelet \
    _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksd/kubernetes/server/bin/kubeadm \
    /eksa-upgrades/binaries/kubernetes/usr/bin/

# TODO: ATTR/Licenses

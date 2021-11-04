ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG RELEASE_BRANCH
ARG TARGETARCH
ARG TARGETOS

COPY _output/$TARGETOS-$TARGETARCH/kubernetes/$RELEASE_BRANCH/kubernetes/client/bin/kubectl \
     _output/$TARGETOS-$TARGETARCH/kubernetes/$RELEASE_BRANCH/kubernetes/server/bin/kubeadm \
     _output/$TARGETOS-$TARGETARCH/etcdadm \
     /opt/bin/
COPY _output/$TARGETOS-$TARGETARCH/kubernetes/$RELEASE_BRANCH/kubernetes/ATTRIBUTION.txt /KUBERNETES_ATTRIBUTION.txt
COPY _output/$TARGETOS-$TARGETARCH/kubernetes/$RELEASE_BRANCH/kubernetes/LICENSES /KUBERNETES_LICENSES
COPY _output/$TARGETOS-$TARGETARCH/etcdadm/ATTRIBUTION.txt /ETCDADM_ATTRIBUTION.txt
COPY _output/$TARGETOS-$TARGETARCH/etcdadm/LICENSES /ETCDADM_LICENSES

RUN ulimit -n 1024 && \
    mkdir -p /opt/bin && \
    chmod +x /opt/bin/kube{adm,ctl} && \
    chmod +x /opt/bin/etcdadm

COPY _output/bin/bottlerocket-bootstrap/$TARGETOS-$TARGETARCH/bottlerocket-bootstrap /bottlerocket-bootstrap
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

CMD ["/bottlerocket-bootstrap"]

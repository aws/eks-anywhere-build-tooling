ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG RELEASE_BRANCH

COPY _output/$RELEASE_BRANCH/kubernetes/kubernetes/client/bin/kubectl \
     _output/$RELEASE_BRANCH/kubernetes/kubernetes/server/bin/kubeadm \
     _output/$RELEASE_BRANCH/etcdadm \
     /opt/bin/
COPY _output/$RELEASE_BRANCH/kubernetes/kubernetes/ATTRIBUTION.txt /KUBERNETES_ATTRIBUTION.txt
COPY _output/$RELEASE_BRANCH/kubernetes/kubernetes/LICENSES /KUBERNETES_LICENSES
COPY _output/$RELEASE_BRANCH/etcdadm/ATTRIBUTION.txt /ETCDADM_ATTRIBUTION.txt
COPY _output/$RELEASE_BRANCH/etcdadm/LICENSES /ETCDADM_LICENSES
COPY scripts/kubeadm-* /usr/sbin/

RUN ulimit -n 1024 && \
    mkdir -p /opt/bin && \
    chmod +x /opt/bin/kube{adm,ctl} && \
    chmod +x /opt/bin/etcdadm && \
    chmod +x /usr/sbin/kubeadm-*

COPY bottlerocket-bootstrap /bottlerocket-bootstrap
RUN chmod +x /bottlerocket-bootstrap

CMD ["/bottlerocket-bootstrap"]

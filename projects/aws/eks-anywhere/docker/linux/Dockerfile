ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

RUN yum install tar -y && \
    yum clean all && \
    rm -rf /var/cache/yum

COPY kubernetes/kubernetes/client/bin/kubectl /usr/local/bin/kubectl
COPY kubernetes/kubernetes/ATTRIBUTION.txt /KUBERNETES_ATTRIBUTION.txt
COPY kubernetes/kubernetes/LICENSES /KUBERNETES_LICENSES
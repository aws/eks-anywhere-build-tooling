ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

RUN yum install tar -y && \
    yum clean all && \
    rm -rf /var/cache/yum

COPY kubernetes/$TARGETOS-$TARGETARCH/kubernetes/client/bin/kubectl /usr/local/bin/kubectl
COPY kubernetes/$TARGETOS-$TARGETARCH/kubernetes/ATTRIBUTION.txt /KUBERNETES_ATTRIBUTION.txt
COPY kubernetes/$TARGETOS-$TARGETARCH/kubernetes/LICENSES /KUBERNETES_LICENSES

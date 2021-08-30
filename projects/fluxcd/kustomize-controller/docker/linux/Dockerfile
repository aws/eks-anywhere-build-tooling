ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-git
ARG BUILDER_IMAGE

FROM $BASE_IMAGE as base

FROM $BUILDER_IMAGE as builder

COPY --from=base /etc/group /etc/group
COPY --from=base /etc/passwd /etc/passwd
RUN yum install -y shadow-utils && \
    groupadd -r controller && \
    useradd -r -g controller controller


FROM $BASE_IMAGE

COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd

COPY _output/files/kubeconfig /home/controller/.kube/config

COPY _output/bin/kustomize-controller/linux-amd64/kustomize-controller /usr/local/bin/kustomize-controller
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

COPY kubernetes/kubernetes/client/bin/kubectl /usr/local/bin/kubectl
COPY kubernetes/kubernetes/ATTRIBUTION.txt /KUBERNETES_ATTRIBUTION.txt
COPY kubernetes/kubernetes/LICENSES /KUBERNETES_LICENSES

USER controller
ENV GNUPGHOME=/tmp

ENTRYPOINT ["kustomize-controller"]

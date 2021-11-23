ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-git
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/files/kubeconfig /home/controller/.kube/config
COPY _output/files/kustomize-controller /
COPY _output/bin/kustomize-controller/$TARGETOS-$TARGETARCH/kustomize-controller /usr/local/bin/kustomize-controller
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

COPY kubernetes/$TARGETOS-$TARGETARCH/kubernetes/client/bin/kubectl /usr/local/bin/kubectl
COPY kubernetes/$TARGETOS-$TARGETARCH/kubernetes/ATTRIBUTION.txt /KUBERNETES_ATTRIBUTION.txt
COPY kubernetes/$TARGETOS-$TARGETARCH/kubernetes/LICENSES /KUBERNETES_LICENSES

USER controller
ENV GNUPGHOME=/tmp

ENTRYPOINT ["kustomize-controller"]

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-git
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/kustomize-controller/$TARGETOS-$TARGETARCH/kustomize-controller /usr/local/bin/kustomize-controller
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 65534
ENV GNUPGHOME=/tmp

ENTRYPOINT ["kustomize-controller"]

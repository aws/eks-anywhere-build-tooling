ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS
ARG RELEASE_BRANCH

COPY _output/bin/eks-anywhere-packages/$TARGETOS-$TARGETARCH/package-manager /package-manager
COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksd/kubernetes/client/bin/kubectl /usr/local/bin/kubectl
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 65532
ENTRYPOINT ["/package-manager"]
CMD ["server"]

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-docker-client
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

WORKDIR /

COPY _output/bin/cluster-api/$TARGETOS-$TARGETARCH/cluster-api-provider-docker-manager /manager
COPY _output/capd/LICENSES /CAPD_LICENSES
COPY CAPD_ATTRIBUTION.txt /CAPD_ATTRIBUTION.txt

COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksd/kubernetes/client/bin/kubectl /kubectl
COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksd/kubernetes/ATTRIBUTION.txt /KUBERNETES_ATTRIBUTION.txt
COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksd/kubernetes/LICENSES /KUBERNETES_LICENSES

# NOTE: CAPD can't use non-root because docker requires access to the docker socket
USER root
ENTRYPOINT ["/manager"]

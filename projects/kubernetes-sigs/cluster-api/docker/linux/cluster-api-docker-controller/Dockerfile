ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-docker-client
FROM $BASE_IMAGE

WORKDIR /

COPY _output/bin/cluster-api/linux-amd64/cluster-api-provider-docker-manager /manager
COPY _output/capd/LICENSES /CAPD_LICENSES
COPY CAPD_ATTRIBUTION.txt /CAPD_ATTRIBUTION.txt

COPY kubernetes/kubernetes/client/bin/kubectl /kubectl
COPY kubernetes/kubernetes/ATTRIBUTION.txt /KUBERNETES_ATTRIBUTION.txt
COPY kubernetes/kubernetes/LICENSES /KUBERNETES_LICENSES

# NOTE: CAPD can't use non-root because docker requires access to the docker socket
USER root
ENTRYPOINT ["/manager"]

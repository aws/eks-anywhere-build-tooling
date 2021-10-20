ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nonroot
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/kube-rbac-proxy/$TARGETOS-$TARGETARCH/kube-rbac-proxy /usr/local/bin/kube-rbac-proxy
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/kube-rbac-proxy"]

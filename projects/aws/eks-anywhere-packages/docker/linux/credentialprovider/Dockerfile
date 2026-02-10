ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/aws/rolesanywhere-credential-helper/aws_signing_helper /eksa-binaries/
COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/kubernetes/cloud-provider-aws/1-35/ecr-credential-provider /eksa-binaries/v1/ecr-credential-provider
COPY _output/bin/eks-anywhere-packages/$TARGETOS-$TARGETARCH/credentialprovider /
COPY _output/LICENSES /LICENSES
COPY CREDENTIALPROVIDER_ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT ["/credentialprovider"]

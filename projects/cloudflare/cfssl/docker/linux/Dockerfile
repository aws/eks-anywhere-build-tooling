ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/gitdependencies/cfssl_trust /etc/cfssl
COPY _output/bin/cfssl/$TARGETOS-$TARGETARCH/* /usr/bin/
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

EXPOSE 8888

ENTRYPOINT ["cfssl"]
CMD ["--help"]

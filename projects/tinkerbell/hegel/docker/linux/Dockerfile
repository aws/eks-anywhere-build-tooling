ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

RUN yum install update upgrade -y shadow-utils ca-certificates

RUN useradd -u 1000 tinkerbell
USER tinkerbell

COPY _output/bin/hegel/$TARGETOS-$TARGETARCH/hegel /usr/bin/hegel
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

EXPOSE 50060
EXPOSE 50061

ENTRYPOINT ["/usr/bin/hegel"]

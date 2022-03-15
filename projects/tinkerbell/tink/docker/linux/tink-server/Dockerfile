ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE as wget-builder

WORKDIR /newroot

RUN set -x && \
    clean_install "wget" && \
    remove_package "bash coreutils gawk info sed grep" && \
    cleanup "wget"

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

WORKDIR /

COPY --from=wget-builder /newroot /

COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/cloudflare/cfssl /usr/local/bin
COPY _output/bin/tink/$TARGETOS-$TARGETARCH/tink-server /usr/bin/tink-server
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

EXPOSE 42113
EXPOSE 42114

ENTRYPOINT ["/usr/bin/tink-server"]

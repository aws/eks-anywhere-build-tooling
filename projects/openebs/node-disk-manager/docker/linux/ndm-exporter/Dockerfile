ARG BASE_IMAGE

FROM public.ecr.aws/eks-distro-build-tooling/builder-base:latest as builder
ADD node-disk-manager node-disk-manager
WORKDIR node-disk-manager
RUN go build -o ndm-exporter ./cmd/ndm-exporter

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

# COPY _output/bin/node-disk-manager/$TARGETOS-$TARGETARCH/ndm-exporter /opt/node-disk-manager/ndm-exporter

COPY --from=builder /node-disk-manager/ndm-exporter /opt/node-disk-manager/ndm-exporter

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000

ENTRYPOINT ["/opt/node-disk-manager/ndm-exporter"]
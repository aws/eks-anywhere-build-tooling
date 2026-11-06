ARG BASE_IMAGE
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as builder

# The chart's init containers use sh and cp to populate shared volumes.
# Keep RPM scriptlet setup bounded on hosts with large descriptor limits.
RUN set -x && \
    ulimit -n 1024 && \
    install_rpm bash && \
    install_binary /usr/bin/cp && \
    cleanup "frr-k8s"

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY --from=builder /newroot /
COPY _output/bin/frr-k8s/$TARGETOS-$TARGETARCH/frr-k8s /frr-k8s
COPY _output/bin/frr-k8s/$TARGETOS-$TARGETARCH/frr-metrics /frr-metrics
COPY _output/bin/frr-k8s/$TARGETOS-$TARGETARCH/frr-status /frr-status
COPY _output/bin/frr-k8s/$TARGETOS-$TARGETARCH/statuscleaner /statuscleaner
COPY --chmod=0755 frr-k8s/cmd/reloader/frr-reloader.sh /frr-reloader.sh
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

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

ARG GIT_HASH
COPY ATTRIBUTION.txt _output/generate_runtime_attribution.sh /source-notices/

# Resolve exact-version source URLs from the actual RPMs, including the base.
RUN set -eux && \
    test -n "$GIT_HASH" && \
    bash /source-notices/generate_runtime_attribution.sh /newroot /source-notices/RUNTIME_ATTRIBUTION.txt && \
    cat /source-notices/ATTRIBUTION.txt /source-notices/RUNTIME_ATTRIBUTION.txt > /newroot/ATTRIBUTION.txt && \
    printf '\nBuild recipes: https://github.com/aws/eks-anywhere-build-tooling/tree/%s/projects/metallb/frr-k8s\n' "$GIT_HASH" >> /newroot/ATTRIBUTION.txt && \
    for directory in licenses doc; do \
        if [ -d "/newroot/usr/share/$directory" ]; then \
            mkdir -p "/newroot/LICENSES/amazon-linux/$directory"; \
            cp -a "/newroot/usr/share/$directory/." "/newroot/LICENSES/amazon-linux/$directory/"; \
        fi; \
    done

FROM scratch AS runtime-attribution
COPY --from=builder /source-notices/RUNTIME_ATTRIBUTION.txt /RUNTIME_ATTRIBUTION.txt

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

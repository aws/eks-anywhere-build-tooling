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
COPY ATTRIBUTION.txt RUNTIME_ATTRIBUTION.txt /source-notices/

# Source notices must cover the actual RPM versions, including the base layer.
RUN set -eux && \
    test -n "$GIT_HASH" && \
    rpm --root /newroot -qa --qf '%{LICENSE}\t%{SOURCERPM}\n' > /tmp/runtime-rpms && \
    for source in $(awk -F '\t' '$1 ~ /GPL|SISSL|Sleepycat/ {print $2}' /tmp/runtime-rpms | sort -u); do \
        grep -F "/$source" /source-notices/RUNTIME_ATTRIBUTION.txt > /dev/null || \
            { echo "Missing source notice for $source" >&2; exit 1; }; \
    done && \
    cat /source-notices/ATTRIBUTION.txt /source-notices/RUNTIME_ATTRIBUTION.txt > /newroot/ATTRIBUTION.txt && \
    printf '\nBuild recipes: https://github.com/aws/eks-anywhere-build-tooling/tree/%s/projects/metallb/frr-k8s\n' "$GIT_HASH" >> /newroot/ATTRIBUTION.txt && \
    for directory in licenses doc; do \
        if [ -d "/newroot/usr/share/$directory" ]; then \
            mkdir -p "/newroot/LICENSES/amazon-linux/$directory"; \
            cp -a "/newroot/usr/share/$directory/." "/newroot/LICENSES/amazon-linux/$directory/"; \
        fi; \
    done

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

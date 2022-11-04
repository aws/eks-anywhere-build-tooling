ARG BASE_IMAGE
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as builder

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/prometheus/$TARGETOS-$TARGETARCH/prometheus /newroot/bin/prometheus
COPY _output/bin/prometheus/$TARGETOS-$TARGETARCH/promtool /newroot/bin/promtool
COPY _output/npm_licenses.tar.bz2 ATTRIBUTION.txt /newroot
COPY _output/LICENSES /newroot/LICENSES
COPY _output/extra/ /newroot

RUN ln -s /usr/share/prometheus/console_libraries /usr/share/prometheus/consoles/ \
    /newroot/etc/prometheus/ && mkdir -p /newroot/prometheus \
    && chown -R 1000:1000 /newroot/etc/prometheus /newroot/prometheus

FROM $BASE_IMAGE

COPY --from=builder /newroot /

WORKDIR /prometheus

USER       1000
EXPOSE     9090
VOLUME     [ "/prometheus" ]
ENTRYPOINT [ "/bin/prometheus" ]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml", \
             "--storage.tsdb.path=/prometheus", \
             "--web.console.libraries=/usr/share/prometheus/console_libraries", \
             "--web.console.templates=/usr/share/prometheus/consoles" ]

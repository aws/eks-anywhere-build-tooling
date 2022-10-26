ARG BASE_IMAGE
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/prometheus/$TARGETOS-$TARGETARCH/prometheus /bin/prometheus
COPY _output/bin/prometheus/$TARGETOS-$TARGETARCH/promtool /bin/promtool
COPY _output/documentation/examples/prometheus.yml   /etc/prometheus/prometheus.yml
COPY _output/console_libraries/                      /usr/share/prometheus/console_libraries/
COPY _output/consoles/                               /usr/share/prometheus/consoles/
COPY _output/npm_licenses.tar.bz2                    /npm_licenses.tar.bz2
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

WORKDIR /prometheus
RUN ln -s /usr/share/prometheus/console_libraries /usr/share/prometheus/consoles/ /etc/prometheus/ && \
    chown -R 1000:1000 /etc/prometheus /prometheus

USER 1000
EXPOSE     9090
VOLUME     [ "/prometheus" ]
ENTRYPOINT ["/bin/prometheus"]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml", \
             "--storage.tsdb.path=/prometheus", \
             "--web.console.libraries=/usr/share/prometheus/console_libraries", \
             "--web.console.templates=/usr/share/prometheus/consoles" ]

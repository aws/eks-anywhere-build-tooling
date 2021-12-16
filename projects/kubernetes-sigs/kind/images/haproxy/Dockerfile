ARG BASE_IMAGE
FROM $BASE_IMAGE

# add haproxy user 
COPY _output/files/haproxy/etc /etc

# upstream minimal config
COPY --chown=haproxy:haproxy _output/files/haproxy/usr /usr

USER haproxy

# below roughly matches the standard haproxy image
STOPSIGNAL SIGUSR1
ENTRYPOINT ["haproxy", "-sf", "7", "-W", "-db", "-f", "/usr/local/etc/haproxy/haproxy.cfg"]

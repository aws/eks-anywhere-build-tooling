ARG BASE_IMAGE
FROM $BASE_IMAGE

# upstream minimal config
COPY --chown=haproxy:haproxy _output/files/haproxy/usr /usr

# below roughly matches the standard haproxy image
STOPSIGNAL SIGUSR1
ENTRYPOINT ["haproxy", "-sf", "7", "-W", "-db", "-f", "/usr/local/etc/haproxy/haproxy.cfg"]

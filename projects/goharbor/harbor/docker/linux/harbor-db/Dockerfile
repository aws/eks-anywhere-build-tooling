ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ENV PGDATA /var/lib/postgresql/data

RUN amazon-linux-extras enable postgresql13 \
    && yum install -y shadow-utils gzip postgresql postgresql-server findutils bc >> /dev/null \
    && userdel postgres \
    && groupadd -f -r postgres --gid=999 \
    && useradd -m -r -g postgres --uid=999 postgres \
    && mkdir -p /docker-entrypoint-initdb.d \
    && mkdir -p /run/postgresql \
    && chown -R postgres:postgres /run/postgresql \
    && chmod 2777 /run/postgresql \
    && mkdir -p "$PGDATA" && chown -R postgres:postgres "$PGDATA" && chmod 777 "$PGDATA" \
    && sed -i "s|#listen_addresses = 'localhost'.*|listen_addresses = '*'|g" /usr/share/pgsql/postgresql.conf.sample \
    && sed -i "s|#unix_socket_directories = '/tmp'.*|unix_socket_directories = '/run/postgresql'|g" /usr/share/pgsql/postgresql.conf.sample \
    && sed -i "s|logging_collector = on.*|logging_collector = off|g" /usr/share/pgsql/postgresql.conf.sample \
    && yum clean all \
    && yum erase -y toybox && yum install -y util-linux net-tools

VOLUME /var/lib/postgresql/data

COPY --chown=postgres:postgres _output/make/photon/db/docker-entrypoint.sh /docker-entrypoint.sh
COPY --chown=postgres:postgres _output/make/photon/db/initdb.sh /initdb.sh
COPY --chown=postgres:postgres _output/make/photon/db/upgrade.sh /upgrade.sh
COPY --chown=postgres:postgres _output/make/photon/db/docker-healthcheck.sh /docker-healthcheck.sh
COPY --chown=postgres:postgres _output/make/photon/db/initial-notaryserver.sql /docker-entrypoint-initdb.d/
COPY --chown=postgres:postgres _output/make/photon/db/initial-notarysigner.sql /docker-entrypoint-initdb.d/
COPY --chown=postgres:postgres _output/make/photon/db/initial-registry.sql /docker-entrypoint-initdb.d/
# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

RUN chmod u+x /docker-entrypoint.sh /docker-healthcheck.sh

ENTRYPOINT ["/docker-entrypoint.sh", "96", "13"]
HEALTHCHECK CMD ["/docker-healthcheck.sh"]

USER postgres

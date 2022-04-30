ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ENV PGDATA /var/lib/postgresql/data

COPY _output/harbor-db/ _output/LICENSES ATTRIBUTION.txt /

RUN amazon-linux-extras enable postgresql13 \
    && yum install -y shadow-utils gzip postgresql postgresql-server findutils bc util-linux net-tools >> /dev/null \
    && userdel postgres \
    && groupadd -f -r postgres --gid=999 \
    && useradd -m -r -g postgres --uid=999 postgres \
    && mkdir -p /run/postgresql \
    && chown -R postgres:postgres /run/postgresql /docker-entrypoint-initdb.d /docker-entrypoint.sh /initdb.sh /upgrade.sh /docker-healthcheck.sh \
    && chmod 2777 /run/postgresql \
    && mkdir -p "$PGDATA" && chown -R postgres:postgres "$PGDATA" && chmod 777 "$PGDATA" \
    && sed -i "s|#listen_addresses = 'localhost'.*|listen_addresses = '*'|g" /usr/share/pgsql/postgresql.conf.sample \
    && sed -i "s|#unix_socket_directories = '/tmp'.*|unix_socket_directories = '/run/postgresql'|g" /usr/share/pgsql/postgresql.conf.sample \
    && yum erase -y toyboxs \
    && yum clean all \
    && chmod u+x /docker-entrypoint.sh /docker-healthcheck.sh

VOLUME /var/lib/postgresql/data

ENTRYPOINT ["/docker-entrypoint.sh", "96", "13"]
HEALTHCHECK CMD ["/docker-healthcheck.sh"]

USER postgres
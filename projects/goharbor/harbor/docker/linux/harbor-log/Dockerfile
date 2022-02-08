ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

RUN yum install -y cronie rsyslog logrotate shadow-utils tar gzip sudo >> /dev/null \
    && mkdir /var/spool/rsyslog \
    && groupadd -f -r -g 10000 syslog && useradd --no-log-init -r -g 10000 -u 10000 syslog \
    && yum clean all \
    && chage -M 99999 root

COPY --chown=10000:10000 _output/make/photon/log/rsyslog.conf /etc/rsyslog.conf
COPY --chown=10000:10000 _output/make/photon/log/rsyslog_docker.conf /etc/rsyslog.d/
RUN rm /etc/cron.daily/logrotate
COPY _output/make/photon/log/logrotate /etc/cron.hourly/
COPY _output/make/photon/log/start.sh /usr/local/bin/
# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

RUN chmod +x /usr/local/bin/start.sh /etc/rsyslog.d/ && \
    chown -R 10000:10000 /run /var/lib/logrotate/

HEALTHCHECK CMD netstat -ltun|grep 10514

VOLUME /var/log/docker/ /run/ /etc/logrotate.d/

CMD /usr/local/bin/start.sh
ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

COPY _output/harbor-log/ _output/LICENSES ATTRIBUTION.txt /

RUN yum install -y cronie rsyslog logrotate shadow-utils tar gzip sudo >> /dev/null \
    && mkdir /var/spool/rsyslog \
    && groupadd -f -r -g 10000 syslog && useradd --no-log-init -r -g 10000 -u 10000 syslog \
    && yum clean all \
    && chage -M 99999 root \
    && rm /etc/cron.daily/logrotate \
    && chmod +x /usr/local/bin/start.sh /etc/rsyslog.d/ \
    && chown -R 10000:10000 /run /var/lib/logrotate/ /etc/rsyslog.conf /etc/rsyslog.d/

HEALTHCHECK CMD netstat -ltun|grep 10514

VOLUME /var/log/docker/ /run/ /etc/logrotate.d/

CMD /usr/local/bin/start.sh
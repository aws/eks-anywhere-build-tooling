ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

RUN amazon-linux-extras enable nginx1 \
    && yum install -y nginx shadow-utils >> /dev/null \
    && yum clean all \
    && userdel nginx \
    && groupadd -f -r -g 10000 nginx && useradd --no-log-init -r -g 10000 -u 10000 nginx \
    && chmod 755 /var/log/nginx \
    && ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log \
    && yum erase -y shadow-utils

# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

VOLUME /var/cache/nginx /var/log/nginx /run

STOPSIGNAL SIGQUIT

HEALTHCHECK CMD curl --fail -s http://localhost:8080 || exit 1

USER nginx

CMD ["nginx", "-g", "daemon off;"]
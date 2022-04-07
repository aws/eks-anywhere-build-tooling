ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nginx
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

RUN pwd
RUN ls -la

ADD conf/index.html /usr/share/nginx/index.template
ADD conf/default.conf /etc/nginx/conf.d/

ADD conf/hello.sh /docker-entrypoint.d/99-hello.sh
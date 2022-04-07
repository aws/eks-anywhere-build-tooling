ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nginx
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

ADD ./hello-eks-anywhere/conf/index.html /usr/share/nginx/index.template
ADD ./hello-eks-anywhere/conf/default.conf /etc/nginx/conf.d/

ADD ./hello-eks-anywhere/conf/hello.sh /docker-entrypoint.d/99-hello.sh
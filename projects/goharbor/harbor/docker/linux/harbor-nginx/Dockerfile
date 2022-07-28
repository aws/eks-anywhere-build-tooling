ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nginx
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE as builder

RUN chmod 755 /newroot/var/log/nginx \
    && ln -sf /dev/stdout /newroot/var/log/nginx/access.log \
    && ln -sf /dev/stderr /newroot/var/log/nginx/error.log

FROM $BASE_IMAGE

COPY --from=builder /newroot /

VOLUME /var/cache/nginx /var/log/nginx /run

STOPSIGNAL SIGQUIT

HEALTHCHECK CMD curl --fail -s http://localhost:8080 || exit 1

USER nginx

CMD ["nginx", "-g", "daemon off;"]
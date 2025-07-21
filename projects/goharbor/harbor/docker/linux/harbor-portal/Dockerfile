ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nginx

FROM public.ecr.aws/docker/library/node:16.18.0 as nodeportal

WORKDIR /build_dir

COPY _output/harbor-portal/ /

ENV NPM_CONFIG_REGISTRY=https://registry.npmjs.org

RUN sed -i 's|http://deb.debian.org/debian|http://archive.debian.org/debian|g' /etc/apt/sources.list \
    && sed -i 's|http://deb.debian.org/debian-security|http://archive.debian.org/debian-security|g' /etc/apt/sources.list \
    && sed -i '/buster-updates/d' /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends python-yaml \
    && npm install --unsafe-perm \ 
    && npm run generate-build-timestamp \
    && npm --prefix /build_dir/app-swagger-ui install --unsafe-perm \
    && npm --prefix /build_dir/app-swagger-ui run build \
    && node --max_old_space_size=2048 'node_modules/@angular/cli/bin/ng' build --configuration production \
    && python -c 'import sys, yaml, json; y=yaml.load(sys.stdin.read()); print json.dumps(y)' < swagger.yaml > dist/swagger.json \
    && cp swagger.yaml dist \
    && mkdir -p /newroot/var/log/nginx /newroot/usr/share/nginx/html \
    && chmod 755 /newroot/var/log/nginx \
    && ln -sf /dev/stdout /newroot/var/log/nginx/access.log \
    && ln -sf /dev/stderr /newroot/var/log/nginx/error.log \
    && mv /build_dir/package*.json /newroot/usr/share/nginx/ \
    && mv /build_dir/dist/* /newroot/usr/share/nginx/html/ \
    && mv /build_dir/app-swagger-ui/dist/* /newroot/usr/share/nginx/html/ 

FROM $BASE_IMAGE

COPY --from=nodeportal /newroot /

VOLUME /var/cache/nginx /var/log/nginx /run

STOPSIGNAL SIGQUIT

HEALTHCHECK CMD curl --fail -s http://localhost:8080 || curl -k --fail -s https://localhost:8443 || exit 1
USER nginx
CMD ["nginx", "-g", "daemon off;"]

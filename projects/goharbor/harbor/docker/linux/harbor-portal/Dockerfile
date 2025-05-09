ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nginx

FROM public.ecr.aws/docker/library/node:16.18.0 as nodeportal

WORKDIR /build_dir

COPY _output/harbor-portal/ /

ENV NPM_CONFIG_REGISTRY=https://registry.npmjs.org

RUN apt-get update \
    && apt-get install -y --no-install-recommends python-yaml \
    # Clean apt cache to free up space
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && npm install --unsafe-perm \ 
    && npm run generate-build-timestamp \
    && npm --prefix /build_dir/app-swagger-ui install --unsafe-perm \
    && npm --prefix /build_dir/app-swagger-ui run build \
    # Clean npm cache to free up space
    && npm cache clean --force \
    && node --max_old_space_size=2048 'node_modules/@angular/cli/bin/ng' build --configuration production \
    # Remove node_modules after build to free up space
    && rm -rf node_modules \
    && rm -rf /build_dir/app-swagger-ui/node_modules \
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
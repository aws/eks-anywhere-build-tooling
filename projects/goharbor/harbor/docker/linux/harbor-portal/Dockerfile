ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nginx

FROM public.ecr.aws/docker/library/node:22.22.3 as nodeportal

WORKDIR /build_dir

COPY _output/harbor-portal/ /

ENV NPM_CONFIG_REGISTRY=https://registry.npmjs.org
ENV NODE_OPTIONS=--max_old_space_size=4096

RUN npm ci --unsafe-perm --ignore-scripts --no-optional --no-audit --no-fund --loglevel=error \
    && npm run postinstall \
    && npm run generate-build-timestamp \
    && node --max_old_space_size=2048 'node_modules/@angular/cli/bin/ng' build --configuration production \
    && node -e "const yaml = require('js-yaml'); const fs = require('fs'); const swagger = yaml.load(fs.readFileSync('swagger.yaml', 'utf8')); fs.writeFileSync('dist/swagger.json', JSON.stringify(swagger));" \
    && cp swagger.yaml dist \
    && npm --prefix /build_dir/app-swagger-ui ci --unsafe-perm --no-optional --no-audit --no-fund --loglevel=error \
    && npm --prefix /build_dir/app-swagger-ui run build \
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

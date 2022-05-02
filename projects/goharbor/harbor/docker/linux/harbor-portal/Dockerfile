ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nginx

FROM node:16.10.0 as nodeportal

WORKDIR /build_dir

COPY _output/harbor-portal/ /

ENV NPM_CONFIG_REGISTRY=https://registry.npmjs.org

RUN apt-get update \
    && apt-get install -y --no-install-recommends python-yaml \
    && npm install --unsafe-perm \ 
    && npm run generate-build-timestamp \
    && node --max_old_space_size=2048 'node_modules/@angular/cli/bin/ng' build --configuration production \
    && python -c 'import sys, yaml, json; y=yaml.load(sys.stdin.read()); print json.dumps(y)' < swagger.yaml > dist/swagger.json \
    && python -c 'import sys, yaml, json; y=yaml.load(sys.stdin.read()); print json.dumps(y)' < swagger2.yaml > dist/swagger2.json \
    && python -c 'import sys, yaml, json; y=yaml.load(sys.stdin.read()); print json.dumps(y)' < swagger3.yaml > dist/swagger3.json \
    && cp swagger.yaml dist


FROM $BASE_IMAGE

COPY --from=nodeportal /build_dir/dist /build_dir/package*.json /usr/share/nginx/html/

RUN chmod 755 /var/log/nginx \
    && ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log \
    && mv /usr/share/nginx/html/package*.json /usr/share/nginx/

VOLUME /var/cache/nginx /var/log/nginx /run

STOPSIGNAL SIGQUIT

HEALTHCHECK CMD curl --fail -s http://localhost:8080 || curl -k --fail -s https://localhost:8443 || exit 1
USER nginx
CMD ["nginx", "-g", "daemon off;"]
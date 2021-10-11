FROM public.ecr.aws/nginx/nginx:1.21-alpine

WORKDIR /

ADD conf/index.template /usr/share/nginx/index.template
ADD conf/default.conf /etc/nginx/conf.d/
ADD conf/test.sh /docker-entrypoint.d/99-test.sh

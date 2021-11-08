#!/bin/sh

export foo="$(cat /config/foo | tr -d '\n')"
export bar="$(cat /secrets/bar | tr -d '\n')"
export version="$(cat /IMAGE_TAG | tr -d '\n')"
export pvcsize="$(df --output=size /pvc | tail -1)"
date=$(date +"%Y%m%d%H%M%S")
date_value=$(date +"%Y-%m-%d %H:%M:%S")
echo "${date_value}" >/pvc/$date
export history="$(cat /pvc/* | tr -d '\n' ',')"

mkdir -p /usr/share/nginx/txt/
echo "⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢

Thank you for using

███████╗██╗  ██╗███████╗
██╔════╝██║ ██╔╝██╔════╝
█████╗  █████╔╝ ███████╗
██╔══╝  ██╔═██╗ ╚════██║
███████╗██║  ██╗███████║
╚══════╝╚═╝  ╚═╝╚══════╝

 █████╗ ███╗   ██╗██╗   ██╗██╗    ██╗██╗  ██╗███████╗██████╗ ███████╗
██╔══██╗████╗  ██║╚██╗ ██╔╝██║    ██║██║  ██║██╔════╝██╔══██╗██╔════╝
███████║██╔██╗ ██║ ╚████╔╝ ██║ █╗ ██║███████║█████╗  ██████╔╝█████╗
██╔══██║██║╚██╗██║  ╚██╔╝  ██║███╗██║██╔══██║██╔══╝  ██╔══██╗██╔══╝
██║  ██║██║ ╚████║   ██║   ╚███╔███╔╝██║  ██║███████╗██║  ██║███████╗
╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝    ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝

You have successfully deployed the eks-anywhere-test pod $POD_NAME

For more information check out
https://anywhere.eks.amazonaws.com

config value foo: ${foo}
secret value bar: ${bar}
version: ${version}
pvc size 1K blocks: ${pvcsize}
history: ${history}

⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢" \
	| tee /usr/share/nginx/txt/index.html
cat /usr/share/nginx/index.template | envsubst > /usr/share/nginx/html/index.html

echo '{"podname":"${POD_NAME}","nodename":"$NODE_NAME","foo":"$foo","bar":"$bar","pvcsize":"$pvcsize","version":"$version"}' \
    | envsubst > /usr/share/nginx/txt/index.json
ln -s /usr/share/nginx/txt/index.json /usr/share/nginx/html/index.json

export -n foo bar version pvcsize history

exec "$@"

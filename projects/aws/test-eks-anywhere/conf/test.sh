#!/bin/sh

export foo="$(cat /config/foo | tr -d '\n')"
export bar="$(cat /secrets/bar | tr -d '\n')"

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

You have successfully deployed the hello-eks-a pod $POD_NAME

For more information check out
https://anywhere.eks.amazonaws.com

config value foo: ${foo}
secret value bar: ${bar}

⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢" \
	| tee /usr/share/nginx/txt/index.html
cat /usr/share/nginx/index.template | envsubst > /usr/share/nginx/html/index.html

echo '{"podname":"${POD_NAME}","nodename":"$NODE_NAME","foo":"$foo","bar":"$bar"}' \
    | envsubst > /usr/share/nginx/txt/index.json
ln -s /usr/share/nginx/txt/index.json /usr/share/nginx/html/index.json

export -n foo bar

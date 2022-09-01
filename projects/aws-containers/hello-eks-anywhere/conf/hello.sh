#!/bin/sh

export version="$(cat /IMAGE_TAG | tr -d '\n')"

echo "⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢

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

${TITLE}
${SUBTITLE}
version: ${version}

⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢⬡⬢" \
	| tee /usr/share/nginx/txt/index.html
cat /usr/share/nginx/index.template | envsubst > /usr/share/nginx/html/index.html

echo '{"podname":"${POD_NAME}","nodename":"$NODE_NAME","title":"$TITLE","subtitle":"$SUBTITLE","version":"$version"}' \
    | envsubst > /usr/share/nginx/txt/index.json
ln -s /usr/share/nginx/txt/index.json /usr/share/nginx/html/index.json

exec "$@"

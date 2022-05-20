#!/bin/bash

set -euxo pipefail

source /etc/eks/logging.sh

SCRIPT_LOG=/var/log/eks-bootstrap.log
touch $SCRIPT_LOG

# save stdout and stderr to file descriptors 3 and 4,
# then redirect them to "$SCRIPT_LOG"
# restore stdout and stderr at bottom of script
exec 3>&1 4>&2 >>$SCRIPT_LOG 2>&1

KUBE_VIP_IMAGE=$1
VIP=$2

# if it's control plane node to join, generate the manifest after `kubeadm join` command complete successfully
if zgrep -q "kubeadm join --config /run/kubeadm/kubeadm-join-config.yaml" /var/lib/cloud/instance/user-data.txt && grep -q success /run/cluster-api/bootstrap-success.complete ; then
  log::info "Joining as control plane node, not the first control plane node to join"
  /etc/eks/generate-kube-vip-manifest.sh "$KUBE_VIP_IMAGE" "$VIP"
fi

# restore stdout and stderr
exec 1>&3 2>&4

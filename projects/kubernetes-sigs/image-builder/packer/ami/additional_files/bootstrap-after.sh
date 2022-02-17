#!/bin/bash

VIP=$1

# if it's control plane node to join, generate the manifest after `kubeadm join` command complete successfully
if grep -q "kubeadm join --config /run/kubeadm/kubeadm-join-config.yaml" /var/lib/cloud/instance/user-data.txt && grep -q success /run/cluster-api/bootstrap-success.complete ; then
  /etc/eks/generate-kube-vip-manifest.sh $VIP
fi

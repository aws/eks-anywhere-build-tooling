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

DNI=$(ip -br link | grep -Ev 'lo|ens3|docker0' | awk '{print $1}')
log::info "Generating kube-vip manifest"
log::info "Using DNI: $DNI"
log::info "Using kube-vip: $VIP"
log::info "Using kube-vip image: $KUBE_VIP_IMAGE"

cat<<EOF >/etc/kubernetes/manifests/kube-vip.yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  name: kube-vip
  namespace: kube-system
spec:
  containers:
  - args:
    - manager
    env:
    - name: vip_arp
      value: "true"
    - name: port
      value: "6443"
    - name: vip_cidr
      value: "32"
    - name: cp_enable
      value: "true"
    - name: cp_namespace
      value: kube-system
    - name: vip_ddns
      value: "false"
    - name: vip_leaderelection
      value: "true"
    - name: vip_leaseduration
      value: "15"
    - name: vip_renewdeadline
      value: "10"
    - name: vip_retryperiod
      value: "2"
    - name: address
      value: $VIP
    image: $KUBE_VIP_IMAGE
    imagePullPolicy: IfNotPresent
    name: kube-vip
    resources: {}
    securityContext:
      capabilities:
        add:
        - NET_ADMIN
        - NET_RAW
    volumeMounts:
    - mountPath: /etc/kubernetes/admin.conf
      name: kubeconfig
  hostNetwork: true
  volumes:
  - hostPath:
      path: /etc/kubernetes/admin.conf
    name: kubeconfig
status: {}
EOF

log::info "Generated kube vip manifest successfully at /etc/kubernetes/manifests/kube-vip.yaml"
log::info "Printing content of kube vip manifest"
log::info "$(cat /etc/kubernetes/manifests/kube-vip.yaml)"
# restore stdout and stderr
exec 1>&3 2>&4

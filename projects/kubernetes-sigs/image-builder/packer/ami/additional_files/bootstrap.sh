#!/bin/bash
source /etc/eks/logging.sh

set -uxo pipefail

SCRIPT_LOG=/var/log/eks-bootstrap.log
touch $SCRIPT_LOG

# save stdout and stderr to file descriptors 3 and 4,
# then redirect them to "$SCRIPT_LOG"
# restore stdout and stderr at bottom of script
exec 3>&1 4>&2 >>$SCRIPT_LOG 2>&1

# Check number of variables, given worker nodes don't need kube-vip
if [ $# -eq 2 ]
then
    KUBE_VIP_IMAGE=$1
    VIP=$2
fi

# static ip configuration needs 5 parameters for cp nodes
if [ $# -eq 5 ]
then
    KUBE_VIP_IMAGE=$1
    VIP=$2
    STATIC_IP=$3
    GATEWAY=$4
    CIDR=$5
fi

# static ip configuration needs 3 parameters for worker nodes
if [ $# -eq 3 ]
then
    STATIC_IP=$1
    GATEWAY=$2
    CIDR=$3
fi

# Using instance id as a unique hostname before we implement hostname in capas
INSTANCE_ID=$(curl 169.254.169.254/latest/meta-data/instance-id | sed -r 's/[.]+/-/g')
hostnamectl set-hostname "$INSTANCE_ID"
log::info "Using hostname: $INSTANCE_ID"
log::info "preserve_hostname: true" > /etc/cloud/cloud.cfg.d/99_hostname.cfg

cat <<EOF > /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
127.0.0.1   $INSTANCE_ID
EOF

# configure DNI
DNI=$(ip -br link | grep -Ev 'lo|ens3|docker0' | awk '{print $1}')
while [ -z "$DNI" ]
do
  # creating DNI is a separate api call, which has some delays
  log::info "DNI is not ready, retry after 5 s"
  sleep 5
  DNI=$(ip -br link | grep -Ev 'lo|ens3|docker0' | awk '{print $1}')
done
log::info "Using DNI: $DNI"

MAC=$(ip -br link | grep -Ev 'lo|ens3|docker0' |  awk '{print $3}')
DEFAULT_GATEWAY=$(ip r | grep default | awk '{print $3}')

log::info "Using MAC: $MAC"
log::info "Using default gateway: $DEFAULT_GATEWAY"

# dhcp configuration needs 2 parameters for cp nodes, 0 parameters for worker nodes
if [ $# -le 2 ]
then
log::info "Configuring DHCP"
cat<<EOF >/etc/netplan/config.yaml
network:
    version: 2
    renderer: networkd
    ethernets:
        $DNI:
            set-name: $DNI
            dhcp4: true
            dhcp-identifier: mac
            dhcp4-overrides:
                route-metric: 50
                send-hostname: true
                hostname: $INSTANCE_ID
            match:
                macaddress: $MAC
        ens3:
            routes:
                - to: 169.254.169.254
                  via: $DEFAULT_GATEWAY
EOF
fi

# dhcp configuration needs 5 parameters for cp nodes, nodes 3 parameters for worker nodes
if [ $# -ge 3 ]
then
log::info "Configuring static ip: $STATIC_IP/$CIDR"
cat<<EOF >/etc/netplan/config.yaml
network:
    version: 2
    renderer: networkd
    ethernets:
        $DNI:
            set-name: $DNI
            addresses:
                - $STATIC_IP/$CIDR
            routes:
                - to: default
                  via: $GATEWAY
                  metric: 50
            match:
                macaddress: $MAC
        ens3:
            routes:
                - to: 169.254.169.254
                  via: $DEFAULT_GATEWAY
EOF
fi

netplan --debug apply

# wait ip leasing via DHCP for 10 s
sleep 10
MY_IP=$(ip route show | grep default | grep "$DNI" | awk '{print $9}')
while [ -z "$MY_IP" ]
do
    # if ip is not ready in 10 s, need to retry netplan apply
    log::info "IP is not ready, retrying"
    netplan --debug apply
    sleep 10
    MY_IP=$(ip route show | grep default | grep "$DNI" | awk '{print $9}')
done
log::info "IP leased from DHCP: $MY_IP"
log::info "netplan applied successfully"

# other network configuration
cat <<EOF >/etc/sysctl.d/kubernetes.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF

sysctl --system

swapoff -a

log::info "network configuration finished"

# if vip is not provided, it's a worker node, we don't need kube-vip manifest
# if `kubeadm init` command doesn't exist in the user-data, it's not the first control plane node, we should generate the kube-vip manifest after the `kubeadm join` command finishes
if [ $# -eq 2 ] || [ $# -eq 5 ]
then
    if zgrep -q "kubeadm init --config /run/kubeadm/kubeadm.yaml" /var/lib/cloud/instance/user-data.txt
  then
    log::info "This is first control plane node, generating kube-vip manifest"
    /etc/eks/generate-kube-vip-manifest.sh "$KUBE_VIP_IMAGE" "$VIP"
  fi
else
  log::info "No VIP provided, this is worker node"
fi

# restore stdout and stderr
exec 1>&3 2>&4

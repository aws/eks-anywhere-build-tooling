#!/bin/bash

set -o errexit
set -o nounset
set -x

backup_file() {
  file_path=$1
  backup_folder=$2

  backedup_file="$backup_folder/$(basename "$file_path").bk"

  if test -f "$backedup_file"; then
    return
  fi

  cp "$file_path" "$backedup_file"
}

backup_and_replace() {
  old_file=$1
  backup_folder=$2
  new_file=$3

  backup_file "$old_file" "$backup_folder" && cp "$new_file" "$old_file"
}

script_dir() {
  echo $(dirname "$(realpath "$0")")
}

upgrade_components_dir() {
   echo "$(dirname "$(script_dir)")" 
}

upgrade_components_bin_dir() {
  echo "$(upgrade_components_dir)/binaries"
}

upgrade_components_kubernetes_bin_dir() {
  echo "$(upgrade_components_bin_dir)/kubernetes/usr/bin"
}

kubeadm_in_first_cp(){
  kube_version=$1
  etcd_version="${2:-NO_UPDATE}"

  components_dir=$(upgrade_components_kubernetes_bin_dir)

  backup_and_replace /usr/bin/kubeadm "$components_dir" "$components_dir/kubeadm"


  kubeadm_config_backup="${components_dir}/kubeadm-config.backup.yaml"
  new_kubeadm_config="${components_dir}/kubeadm-config.yaml"
  kubectl get cm -n kube-system kubeadm-config -ojsonpath='{.data.ClusterConfiguration}' --kubeconfig /etc/kubernetes/admin.conf > "$kubeadm_config_backup"

  if [ "$etcd_version" != "NO_UPDATE" ]; then
    sed -zE "s/(imageRepository: public.ecr.aws\/eks-distro\/etcd-io\n\s+imageTag: )[^\n]*/\1${etcd_version}/" "$kubeadm_config_backup" > "$new_kubeadm_config"
  fi

  # the kubelet config appears to lose values, in the case of a kind cluster the failSwapOn:false
  echo "---" >> "$new_kubeadm_config"
  kubectl get cm -n kube-system kubelet-config -ojsonpath='{.data.kubelet}' --kubeconfig /etc/kubernetes/admin.conf >> "$new_kubeadm_config"

  # Backup and delete coredns configmap. If the CM doesn't exist, kubeadm will skip its upgrade.
  # This is desirable for 2 reasons:
  # - CAPI already takes care of coredns upgrades
  # - kubeadm will fail when verifying the current version of coredns bc the image tag created by
  #   eks-s is not recognised by the migration verification logic https://github.com/coredns/corefile-migration/blob/master/migration/versions.go
  # Ideally we will instruct kubeadm to just skip coredns upgrade during this phase, but
  # it doesn't seem like there is an option.
  # TODO: consider using --skip-phases to skip addons/coredns once the feature flag is supported in kubeadm upgrade command
  backup_and_delete_coredns_config "$components_dir"

  kubeadm version
  kubeadm upgrade plan --ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration --config "$new_kubeadm_config"
  kubeadm upgrade apply "$kube_version" --config "$new_kubeadm_config" --ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration --allow-experimental-upgrades --yes

  restore_coredns_config "$components_dir"
}

kubeadm_in_rest_cp(){
  components_dir=$(upgrade_components_kubernetes_bin_dir)

  backup_and_replace /usr/bin/kubeadm "$components_dir" "$components_dir/kubeadm"

  # Backup and delete coredns configmap. If the CM doesn't exist, kubeadm will skip its upgrade.
  # This is desirable for 2 reasons:
  # - CAPI already takes care of coredns upgrades
  # - kubeadm will fail when verifying the current version of coredns bc the image tag created by
  #   eks-s is not recognised by the migration verification logic https://github.com/coredns/corefile-migration/blob/master/migration/versions.go
  # Ideally we will instruct kubeadm to just skip coredns upgrade during this phase, but
  # it doesn't seem like there is an option.
  # TODO: consider using --skip-phases to skip addons/coredns once the feature flag is supported in kubeadm upgrade command
  backup_and_delete_coredns_config "$components_dir"

  kubeadm version
  kubeadm upgrade node --ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration
  restore_coredns_config "$components_dir"
}

backup_and_delete_coredns_config(){
  components_dir=$1
  # Backup and delete coredns configmap. If the CM doesn't exist, kubeadm will skip its upgrade.
  # This is desirable for 2 reasons:
  # - CAPI already takes care of coredns upgrades
  # - kubeadm will fail when verifying the current version of coredns bc the image tag created by
  #   eks-s is not recognised by the migration verification logic https://github.com/coredns/corefile-migration/blob/master/migration/versions.go
  # Ideally we will instruct kubeadm to just skip coredns upgrade during this phase, but
  # it doesn't seem like there is an option.
  # TODO: consider using --skip-phases to skip addons/coredns once the feature flag is supported in kubeadm upgrade command
  coredns_backup="${components_dir}/coredns.yaml"
  coredns=$(kubectl get cm -n kube-system coredns -oyaml --kubeconfig /etc/kubernetes/admin.conf --ignore-not-found=true)
  if [ -n "$coredns" ]; then
    echo "$coredns" >"$coredns_backup"
  fi
  kubectl delete cm -n kube-system coredns --kubeconfig /etc/kubernetes/admin.conf --ignore-not-found=true
}

restore_coredns_config(){
  components_dir=$1
  coredns_backup="${components_dir}/coredns.yaml"
  # Restore coredns config from backup
  kubectl create -f "$coredns_backup" --kubeconfig /etc/kubernetes/admin.conf
}

kubeadm_in_worker() {
  components_dir=$(upgrade_components_kubernetes_bin_dir)

  backup_and_replace /usr/bin/kubeadm "$components_dir" "$components_dir/kubeadm"

  kubeadm version
  kubeadm upgrade node 
}

kubelet_and_kubectl() {
  kube_version=$(kubeadm version -oshort)

  components_dir=$(upgrade_components_kubernetes_bin_dir)

  backup_and_replace /usr/bin/kubectl "$components_dir" "$components_dir/kubectl"

  systemctl stop kubelet
  backup_and_replace /usr/bin/kubelet "$components_dir" "$components_dir/kubelet"

  # KubeletCredentialProviders support became GA in k8s v1.26, and the feature gate was removed in k8s v1.28.
  # For in-place upgrades, we should remove this feature gate if it exists on nodes running k8s v1.26 and above.
  if [[ "$kube_version" != v1.25* ]]; then
    update_kubelet_extra_args
  fi

  systemctl daemon-reload
  systemctl restart kubelet
}

update_kubelet_extra_args() {
  kubelet_extra_args=$(cat /etc/sysconfig/kubelet)
  feature_gate=" --feature-gates=KubeletCredentialProviders=true"
  if [[ $kubelet_extra_args == *$feature_gate* ]]; then
    kubelet_extra_args=${kubelet_extra_args//"$feature_gate"/}
    mkdir "$components_dir/extraargs"
    backup_file /etc/sysconfig/kubelet "$components_dir/extraargs"
    echo "$kubelet_extra_args" > /etc/sysconfig/kubelet
  fi
}

upgrade_containerd() {
  components_dir=$(upgrade_components_bin_dir)

  containerd --version

  # before nsenter, copy /eksa-upgrades to /tmp/eksa-upgrades
  # IDEA: similar to the kubeadm flow, we could loop the folder structure
  # and for each binary/config file backup the existing file before copying 
  # so there is a rollback path
  cp -rf $components_dir/containerd/* /

  containerd --version

  systemctl daemon-reload
  # systemctl stop containerd

  systemctl restart containerd

}

cni_plugins() {  
  components_dir=$(upgrade_components_bin_dir)
  
  /opt/cni/bin/loopback --version

  # before nsenter, copy /eksa-upgrades to /tmp/eksa-upgrades
  # IDEA: similar to the kubeadm flow, we could loop the folder structure
  # and for each binary/config file backup the existing file before copying 
  # so there is a rollback path
  cp -rf $components_dir/cni-plugins/* /

  /opt/cni/bin/loopback --version

  # rm -rf /foo/eksa-upgrades
}

print_status() {
  systemctl status containerd
  systemctl status kubelet
  kubeadm version
}

print_status_and_cleanup() {
  print_status

  components_dir=$(upgrade_components_dir)
  echo "Deleting all leftover upgrade components at ${components_dir}"
  rm -rf "$components_dir"
}

$@
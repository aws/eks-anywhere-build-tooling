package upgrade

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

const (
	kubeAdmBinDir            = "/usr/bin/kubeadm"
	etcdImageRepo            = "public.ecr.aws/eks-distro/etcd-io"
	noEtcdUpdate             = "NO_UPDATE"
	yamlSeparatorWithNewLine = "---\n"
	staticKubeVipPath        = "/etc/kubernetes/manifests/kube-vip.yaml"
	kubeConfigPath           = "/etc/kubernetes/admin.conf"
)

func (u *Upgrader) KubeAdmInFirstCP(ctx context.Context) error {
	componentsDir, err := u.upgradeComponentsKubernetesBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade componenets kubernetes binary directory: %v", err)
	}

	if err = u.BackUpAndReplace(kubeAdmBinDir, componentsDir, fmt.Sprintf("%s/kubeadm", componentsDir)); err != nil {
		return fmt.Errorf("backing up and replacing kubeadm binary: %v", err)
	}
	logger.Info("Backed up and replaced kubeadm binary sucessfully")

	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", componentsDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", componentsDir)

	out, err := u.ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath)
	if err != nil {
		return fmt.Errorf("executing command 'kubectl get cm -n kube-system kubeadm-config -ojsonpath='{.data.ClusterConfiguration}' --kubeconfig %s: %s", kubeConfigPath, string(out))
	}

	trimmedOut := strings.Trim(string(out), "`'")
	err = u.WriteFile(kubeAdmConfigBackUp, []byte(trimmedOut), 0o640)
	if err != nil {
		return fmt.Errorf("writing kubeadm config to backup file: %v", err)
	}

	if u.etcdVersion != noEtcdUpdate {
		if err = u.updateEtcdVersion(kubeAdmConfigBackUp, newKubeAdmConfig, u.etcdVersion); err != nil {
			return fmt.Errorf("updating etcd version to %s: %v", u.etcdVersion, err)
		}
	}

	if err = u.appendKubeletConfig(ctx, newKubeAdmConfig); err != nil {
		return fmt.Errorf("appending kubelet config: %v", err)
	}

	if err = u.backUpAndDeleteCoreDNSConfig(ctx, componentsDir); err != nil {
		return fmt.Errorf("backing up and deleting coreDNS config: %v", err)
	}

	kubeAdmVersion, err := u.ExecCommand(ctx, "kubeadm", "version")
	if err != nil {
		return fmt.Errorf("executing command 'kubeadm version': %v", kubeAdmVersion)
	}
	logger.Info("current version of kubeadm", "cmd", "kubeadm version", "output", string(kubeAdmVersion))

	kubeAdmUpgPlan, err := u.ExecCommand(ctx, "kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig)
	if err != nil {
		return fmt.Errorf("executing command 'kubeadm upgrade plan --ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration --config %s': %v", newKubeAdmConfig, string(kubeAdmUpgPlan))
	}
	logger.Info("components to be upgraded with kubeadm", "output", string(kubeAdmUpgPlan))

	kubeAdmUpg, err := u.ExecCommand(ctx, "kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes")
	if err != nil {
		return fmt.Errorf("executing command 'kubeadm upgrade apply --ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration --config %s --allow-experimental-upgrades --yes': %v", newKubeAdmConfig, string(kubeAdmUpg))
	}
	logger.Info("verbose output for kubeadm upgrade", "output", string(kubeAdmUpg))

	upgCmpDir, err := u.upgradeComponentsDir()
	if err != nil {
		return fmt.Errorf("getting upgrade componenets directory: %v", err)
	}

	newKubeVipConfigPath := fmt.Sprintf("%s/kube-vip.yaml", upgCmpDir)
	if err := u.copy(staticKubeVipPath, fmt.Sprintf("%s/kube-vip.backup.yaml", upgCmpDir)); err != nil {
		return copyError(staticKubeVipPath, fmt.Sprintf("%s/kube-vip.backup.yaml", upgCmpDir), err)
	}

	if err := u.copy(newKubeVipConfigPath, staticKubeVipPath); err != nil {
		return copyError(newKubeVipConfigPath, staticKubeVipPath, err)
	}

	if err = u.restoreCoreDNSConfig(ctx, componentsDir); err != nil {
		return fmt.Errorf("restoring coreDNS config: %v", err)
	}
	logger.Info("Upgraded kubeadm in First Control Plane successfully!", "version", u.kubernetesVersion)

	return nil
}

func (u *Upgrader) KubeAdmInRestCP(ctx context.Context) error {
	componentsDir, err := u.upgradeComponentsKubernetesBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade componenets kubernetes binary directory: %v", err)
	}

	if err = u.BackUpAndReplace(kubeAdmBinDir, componentsDir, fmt.Sprintf("%s/kubeadm", componentsDir)); err != nil {
		return fmt.Errorf("backing up and replacing kubeadm binary: %v", err)
	}
	logger.Info("Backed up and replaced kubeadm binary sucessfully")

	if err = u.backUpAndDeleteCoreDNSConfig(ctx, componentsDir); err != nil {
		return fmt.Errorf("backing up and deleting coreDNS config: %v", err)
	}

	kubeAdmVersion, err := u.ExecCommand(ctx, "kubeadm", "version")
	if err != nil {
		return fmt.Errorf("executing command 'kubeadm version': %v", kubeAdmVersion)
	}
	logger.Info("current version of kubeadm", "cmd", "kubeadm version", "output", string(kubeAdmVersion))

	kubeAdmUpg, err := u.ExecCommand(ctx, "kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration")
	if err != nil {
		return fmt.Errorf("executing command 'kubeadm upgrade node --ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration': %v", string(kubeAdmUpg))
	}
	logger.Info("verbose output for kubeadm upgrade", "output", string(kubeAdmUpg))

	upgCmpDir, err := u.upgradeComponentsDir()
	if err != nil {
		return fmt.Errorf("getting upgrade componenets directory: %v", err)
	}

	newKubeVipConfigPath := fmt.Sprintf("%s/kube-vip.yaml", upgCmpDir)
	if err := u.copy(staticKubeVipPath, fmt.Sprintf("%s/kube-vip.backup.yaml", upgCmpDir)); err != nil {
		return copyError(staticKubeVipPath, fmt.Sprintf("%s/kube-vip.backup.yaml", upgCmpDir), err)
	}

	if err := u.copy(newKubeVipConfigPath, staticKubeVipPath); err != nil {
		return copyError(newKubeVipConfigPath, staticKubeVipPath, err)
	}

	if err = u.restoreCoreDNSConfig(ctx, componentsDir); err != nil {
		return fmt.Errorf("restoring coreDNS config: %v", err)
	}
	logger.Info("Kubeadm in First Control Plane upgraded successfully", "version", u.kubernetesVersion)

	return nil
}

func (u *Upgrader) KubeAdmInWorker(ctx context.Context) error {
	componentsDir, err := u.upgradeComponentsKubernetesBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade componenets kubernetes binary directory: %v", err)
	}

	if err = u.BackUpAndReplace(kubeAdmBinDir, componentsDir, fmt.Sprintf("%s/kubeadm", componentsDir)); err != nil {
		return fmt.Errorf("backing up and replacing kubeadm binary: %v", err)
	}
	logger.Info("Backed up and replaced kubeadm binary sucessfully")

	kubeAdmVersion, err := u.ExecCommand(ctx, "kubeadm", "version")
	if err != nil {
		return fmt.Errorf("executing command 'kubeadm version': %v", kubeAdmVersion)
	}
	logger.Info("current version of kubeadm", "cmd", "kubeadm version", "output", string(kubeAdmVersion))

	kubeAdmUpg, err := u.ExecCommand(ctx, "kubeadm", "upgrade", "node")
	if err != nil {
		return fmt.Errorf("executing command 'kubeadm upgrade node': %v", string(kubeAdmUpg))
	}
	logger.Info("verbose output for kubeadm upgrade", "output", string(kubeAdmUpg))
	logger.Info("Kubeadm in Worker Node upgraded successfully")

	return nil
}

func (u *Upgrader) updateEtcdVersion(oldKubeAdmConf, newKubeAdmConf, version string) error {
	conf, err := u.ReadFile(oldKubeAdmConf)
	if err != nil {
		return err
	}
	lines := strings.Split(string(conf), "\n")
	for i, line := range lines {
		if strings.Contains(line, etcdImageRepo) {
			imageTag := strings.Split(lines[i+1], ":")
			// the space in the below string is for yaml formatting and should not be removed
			imageTag[1] = fmt.Sprintf(" %s", version)
			lines[i+1] = strings.Join(imageTag, ":")
		}
	}
	updatedConf := strings.Join(lines, "\n")
	err = u.WriteFile(newKubeAdmConf, []byte(updatedConf), 0o640)
	if err != nil {
		return err
	}
	return nil
}

func (u *Upgrader) appendKubeletConfig(ctx context.Context, kubeAdmConf string) error {
	conf, err := u.ReadFile(kubeAdmConf)
	if err != nil {
		return err
	}
	conf = append(conf, []byte(yamlSeparatorWithNewLine)...)
	out, err := u.ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath)
	if err != nil {
		return fmt.Errorf("executing command 'kubectl get cm -n kube-system kubelet-config -ojsonpath='{.data.kubelet}' --kubeconfig %s': %v", kubeConfigPath, string(out))
	}
	logger.Info("kubelet config as string", "out", string(out))

	trimmedOut := strings.Trim(string(out), "`'")
	conf = append(conf, []byte(trimmedOut)...)
	err = u.WriteFile(kubeAdmConf, []byte(conf), 0o640)
	if err != nil {
		return fmt.Errorf("writing kubelet config to kubeadm file: %v", err)
	}

	logger.Info("Appended Kubelet Config to Kubeadm config file")
	return nil
}

// Backup and delete coredns configmap. If the CM doesn't exist, kubeadm will skip its upgrade.
// This is desirable for 2 reasons:
//  1. CAPI already takes care of coredns upgrades
//  2. kubeadm will fail when verifying the current version of coredns bc the image tag created by  eks-a
//     is not recognised by the migration verification logic https://github.com/coredns/corefile-migration/blob/master/migration/versions.go
//
// Ideally we will instruct kubeadm to just skip coredns upgrade during this phase, but
// it doesn't seem like there is an option.
// TODO: consider using --skip-phases to skip addons/coredns once the feature flag is supported in kubeadm upgrade command
func (u *Upgrader) backUpAndDeleteCoreDNSConfig(ctx context.Context, cmpDir string) error {
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", cmpDir)
	coreDNSConf, err := u.ExecCommand(ctx, "kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true")
	if err != nil {
		return fmt.Errorf("executing command 'kubectl get cm -n kube-system coredns -oyaml --kubeconfig %s': %v", kubeConfigPath, string(coreDNSConf))
	}
	if len(coreDNSConf) > 0 {
		logger.Info("coreDNS config as string", "out", string(coreDNSConf))
		err = u.WriteFile(coreDNSBackup, coreDNSConf, 0o644)
		if err != nil {
			return err
		}
	}
	out, err := u.ExecCommand(ctx, "kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true")
	if err != nil {
		return fmt.Errorf("kubectl delete cm -n kube-system coredns -oyaml --kubeconfig %s --ignore-not-found=true': %v", kubeConfigPath, string(out))
	}

	logger.Info("Backed up and deleted CoreDNS config")
	return nil
}

func (u *Upgrader) restoreCoreDNSConfig(ctx context.Context, cmpDir string) error {
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", cmpDir)
	out, err := u.ExecCommand(ctx, "kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath)
	if err != nil {
		return fmt.Errorf("executing command 'kubectl create -f %s --kubeconfig %s': %v", coreDNSBackup, kubeConfigPath, string(out))
	}

	logger.Info("Restored CoreDNS config successfully!")
	return nil
}

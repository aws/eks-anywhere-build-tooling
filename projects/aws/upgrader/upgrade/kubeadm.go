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

// KubeAdmInFirstCP upgrades the first control plane node
//
// As part of upgrade:
//  1. backs up the existing kubeadm binary, replace with the new binary and backs up existing kubeadm cluster-config.
//  2. updates the cluster config with latest etcd version and appends the kubelet config.
//  3. back up and delete coreDNS config as capi handles coreDNS upgrade.
//  4. run kubadm upgrade commands and copy over new kube-vip config.
//  5. Restore coreDNS config back once the kubeadm upgrade commands are complete.
func (u *InPlaceUpgrader) KubeAdmInFirstCP(ctx context.Context) error {
	componentsDir, err := u.upgradeComponentsKubernetesBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade components kubernetes binary directory: %v", err)
	}

	if err = u.BackUpAndReplace(kubeAdmBinDir, componentsDir, fmt.Sprintf("%s/kubeadm", componentsDir)); err != nil {
		return fmt.Errorf("backing up and replacing kubeadm binary: %v", err)
	}
	logger.Info("backed up and replaced kubeadm binary successfully")

	kubeAdmConfigBackUp := fmt.Sprintf("%s/kubeadm-config.backup.yaml", componentsDir)
	newKubeAdmConfig := fmt.Sprintf("%s/kubeadm-config.yaml", componentsDir)

	getClusterConfigCmd := []string{"kubectl", "get", "cm", "-n", "kube-system", "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath}
	out, err := u.ExecCommand(ctx, getClusterConfigCmd[0], getClusterConfigCmd[1:]...)
	if err != nil {
		return execError(getClusterConfigCmd, string(out))
	}

	trimmedOut := strings.Trim(string(out), "`'")
	err = u.WriteFile(kubeAdmConfigBackUp, []byte(trimmedOut), fileMode416)
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

	kubeAdmVersionCmd := []string{"kubeadm", "version"}
	version, err := u.ExecCommand(ctx, kubeAdmVersionCmd[0], kubeAdmVersionCmd[1:]...)
	if err != nil {
		return execError(kubeAdmVersionCmd, string(version))
	}
	logger.Info("current version of kubeadm", "cmd", "kubeadm version", "output", string(version))

	kubeAdmUpgPlanCmd := []string{"kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig}
	kubeAdmUpgPlan, err := u.ExecCommand(ctx, kubeAdmUpgPlanCmd[0], kubeAdmUpgPlanCmd[1:]...)
	if err != nil {
		return execError(kubeAdmUpgPlanCmd, string(kubeAdmUpgPlan))
	}
	logger.Info("components to be upgraded with kubeadm", "output", string(kubeAdmUpgPlan))

	kubeAdmUpgCmd := []string{"kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes"}
	kubeAdmUpg, err := u.ExecCommand(ctx, kubeAdmUpgCmd[0], kubeAdmUpgCmd[1:]...)
	if err != nil {
		return execError(kubeAdmUpgCmd, string(kubeAdmUpg))
	}
	logger.Info("verbose output for kubeadm upgrade", "output", string(kubeAdmUpg))

	upgCmpDir, err := u.upgradeComponentsDir()
	if err != nil {
		return fmt.Errorf("getting upgrade components directory: %v", err)
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
	logger.Info("kubeadm upgrade in first control plane successful!", "version", u.kubernetesVersion)

	return nil
}

// KubeAdmInRestCP upgrades the rest of control plane nodes
//
// As part of upgrade:
//  1. backs up the existing kubeadm binary and replace with the new binary.
//  2. back up and delete coreDNS config as capi handles coreDNS upgrade.
//  3. run kubadm upgrade commands and copy over new kube-vip config.
//  4. Restore coreDNS config back once the kubeadm upgrade commands are complete.
func (u *InPlaceUpgrader) KubeAdmInRestCP(ctx context.Context) error {
	componentsDir, err := u.upgradeComponentsKubernetesBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade components kubernetes binary directory: %v", err)
	}

	if err = u.BackUpAndReplace(kubeAdmBinDir, componentsDir, fmt.Sprintf("%s/kubeadm", componentsDir)); err != nil {
		return fmt.Errorf("backing up and replacing kubeadm binary: %v", err)
	}
	logger.Info("Backed up and replaced kubeadm binary successfully")

	if err = u.backUpAndDeleteCoreDNSConfig(ctx, componentsDir); err != nil {
		return fmt.Errorf("backing up and deleting coreDNS config: %v", err)
	}

	kubeAdmVersionCmd := []string{"kubeadm", "version"}
	version, err := u.ExecCommand(ctx, kubeAdmVersionCmd[0], kubeAdmVersionCmd[1:]...)
	if err != nil {
		return execError(kubeAdmVersionCmd, string(version))
	}
	logger.Info("current version of kubeadm", "cmd", "kubeadm version", "output", string(version))

	kubeAdmUpgNodeCmd := []string{"kubeadm", "upgrade", "node", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration"}
	kubeAdmUpg, err := u.ExecCommand(ctx, kubeAdmUpgNodeCmd[0], kubeAdmUpgNodeCmd[1:]...)
	if err != nil {
		return execError(kubeAdmUpgNodeCmd, string(kubeAdmUpg))
	}
	logger.Info("verbose output for kubeadm upgrade", "output", string(kubeAdmUpg))

	upgCmpDir, err := u.upgradeComponentsDir()
	if err != nil {
		return fmt.Errorf("getting upgrade components directory: %v", err)
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
	logger.Info("kubeadm upgrade in control plane successful!")

	return nil
}

// KubeAdmInWorker upgrades the worker nodes
//
// As part of upgrade:
//  1. backs up the existing kubeadm binary and replace with the new binary.
//  2. run kubadm upgrade commands and copy over new kube-vip config.
func (u *InPlaceUpgrader) KubeAdmInWorker(ctx context.Context) error {
	componentsDir, err := u.upgradeComponentsKubernetesBinDir()
	if err != nil {
		return fmt.Errorf("getting upgrade components kubernetes binary directory: %v", err)
	}

	if err = u.BackUpAndReplace(kubeAdmBinDir, componentsDir, fmt.Sprintf("%s/kubeadm", componentsDir)); err != nil {
		return fmt.Errorf("backing up and replacing kubeadm binary: %v", err)
	}
	logger.Info("Backed up and replaced kubeadm binary successfully")

	kubeAdmVersionCmd := []string{"kubeadm", "version"}
	version, err := u.ExecCommand(ctx, kubeAdmVersionCmd[0], kubeAdmVersionCmd[1:]...)
	if err != nil {
		return execError(kubeAdmVersionCmd, string(version))
	}
	logger.Info("current version of kubeadm", "cmd", "kubeadm version", "output", string(version))

	kubeAdmUpgNodeCmd := []string{"kubeadm", "upgrade", "node"}
	kubeAdmUpg, err := u.ExecCommand(ctx, kubeAdmUpgNodeCmd[0], kubeAdmUpgNodeCmd[1:]...)
	if err != nil {
		return execError(kubeAdmUpgNodeCmd, string(kubeAdmUpg))
	}
	logger.Info("verbose output for kubeadm upgrade", "output", string(kubeAdmUpg))
	logger.Info("kubeadm upgrade in worker node successful!")

	return nil
}

// updateEtcdVersion updates the etcd image to the latest tag for that kubernetes version in the kubeadm config file
func (u *InPlaceUpgrader) updateEtcdVersion(oldKubeAdmConf, newKubeAdmConf, version string) error {
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
	err = u.WriteFile(newKubeAdmConf, []byte(updatedConf), fileMode416)
	if err != nil {
		return err
	}
	return nil
}

// appendKubeletConfig retreives kubelet-config and appends it to the existing kubeadm-config file
func (u *InPlaceUpgrader) appendKubeletConfig(ctx context.Context, kubeAdmConf string) error {
	conf, err := u.ReadFile(kubeAdmConf)
	if err != nil {
		return err
	}
	conf = append(conf, []byte(yamlSeparatorWithNewLine)...)
	getKubeletConfCmd := []string{"kubectl", "get", "cm", "-n", "kube-system", "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath}
	out, err := u.ExecCommand(ctx, getKubeletConfCmd[0], getKubeletConfCmd[1:]...)
	if err != nil {
		return execError(getKubeletConfCmd, string(out))
	}
	logger.Info("kubelet config as string", "out", string(out))

	trimmedOut := strings.Trim(string(out), "`'")
	conf = append(conf, []byte(trimmedOut)...)
	err = u.WriteFile(kubeAdmConf, []byte(conf), fileMode416)
	if err != nil {
		return fmt.Errorf("writing kubelet config to kubeadm file: %v", err)
	}

	logger.Info("appended kubelet config to kubeadm config file")
	return nil
}

// backUpAndDeleteCoreDNSConfig backs up and deletes the coreDNS config map.
// Backup and delete coredns configmap. If the CM doesn't exist, kubeadm will skip its upgrade.
// This is desirable for 2 reasons:
//  1. CAPI already takes care of coredns upgrades
//  2. kubeadm will fail when verifying the current version of coredns bc the image tag created by  eks-a
//     is not recognised by the migration verification logic https://github.com/coredns/corefile-migration/blob/master/migration/versions.go
//
// Ideally we will instruct kubeadm to just skip coredns upgrade during this phase, but
// it doesn't seem like there is an option.
// TODO: consider using --skip-phases to skip addons/coredns once the feature flag is supported in kubeadm upgrade command
func (u *InPlaceUpgrader) backUpAndDeleteCoreDNSConfig(ctx context.Context, cmpDir string) error {
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", cmpDir)
	getCoreDNSConfCmd := []string{"kubectl", "get", "cm", "-n", "kube-system", "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true"}
	coreDNSConf, err := u.ExecCommand(ctx, getCoreDNSConfCmd[0], getCoreDNSConfCmd[1:]...)
	if err != nil {
		return execError(getCoreDNSConfCmd, string(coreDNSConf))
	}
	if len(coreDNSConf) > 0 {
		logger.Info("coreDNS config as string", "out", string(coreDNSConf))
		err = u.WriteFile(coreDNSBackup, coreDNSConf, fileMode416)
		if err != nil {
			return err
		}
	}
	deleteCoreDNSConfig := []string{"kubectl", "delete", "cm", "-n", "kube-system", "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true"}
	out, err := u.ExecCommand(ctx, deleteCoreDNSConfig[0], deleteCoreDNSConfig[1:]...)
	if err != nil {
		return execError(deleteCoreDNSConfig, string(out))
	}

	logger.Info("backed up and deleted coreDNS config")
	return nil
}

func (u *InPlaceUpgrader) restoreCoreDNSConfig(ctx context.Context, cmpDir string) error {
	coreDNSBackup := fmt.Sprintf("%s/coredns.yaml", cmpDir)
	createCoreDNSConfCmd := []string{"kubectl", "create", "-f", coreDNSBackup, "--kubeconfig", kubeConfigPath}
	out, err := u.ExecCommand(ctx, createCoreDNSConfCmd[0], createCoreDNSConfCmd[1:]...)
	if err != nil {
		return execError(createCoreDNSConfCmd, string(out))
	}

	logger.Info("restored coreDNS config successfully!")
	return nil
}

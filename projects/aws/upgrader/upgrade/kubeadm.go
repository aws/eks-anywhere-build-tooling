package upgrade

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
	"golang.org/x/mod/semver"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"
)

const (
	kubeAdmBinDir            = "/usr/bin/kubeadm"
	etcdImageRepo            = "public.ecr.aws/eks-distro/etcd-io"
	noEtcdUpdate             = "NO_UPDATE"
	yamlSeparatorWithNewLine = "---\n"
	staticKubeVipPath        = "/etc/kubernetes/manifests/kube-vip.yaml"
	kubeConfigPath           = "/etc/kubernetes/admin.conf"
	upgradeConfigurationKind = "UpgradeConfiguration"
	kubeadmAPIv1beta4        = "kubeadm.k8s.io/v1beta4"
	kubeadmCMName            = "kubeadm-config"
	kubeSystemNS             = "kube-system"
	kubeVersion130           = "v1.30"
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

	getClusterConfigCmd := []string{"kubectl", "get", "cm", "-n", kubeSystemNS, "kubeadm-config", "-ojsonpath='{.data.ClusterConfiguration}'", "--kubeconfig", kubeConfigPath}
	out, err := u.ExecCommand(ctx, getClusterConfigCmd[0], getClusterConfigCmd[1:]...)
	if err != nil {
		return execError(getClusterConfigCmd, string(out))
	}

	trimmedOut := strings.Trim(string(out), "`'")
	err = u.WriteFile(kubeAdmConfigBackUp, []byte(trimmedOut), fileMode640)
	if err != nil {
		return fmt.Errorf("writing kubeadm config to backup file: %v", err)
	}

	if u.etcdVersion != noEtcdUpdate {
		if err = u.updateEtcdVersion(kubeAdmConfigBackUp, newKubeAdmConfig, u.etcdVersion); err != nil {
			return fmt.Errorf("updating etcd version to %s: %v", u.etcdVersion, err)
		}
	}

	if err = u.backUpAndDeleteCoreDNSConfig(ctx, componentsDir); err != nil {
		return fmt.Errorf("backing up and deleting coreDNS config: %v", err)
	}

	kubeAdmVersionCmd := []string{"kubeadm", "version", "-oshort"}
	version, err := u.ExecCommand(ctx, kubeAdmVersionCmd[0], kubeAdmVersionCmd[1:]...)
	if err != nil {
		return execError(kubeAdmVersionCmd, string(version))
	}
	logger.Info("current version of kubeadm", "cmd", "kubeadm version -oshort", "output", string(version))

	// K8s version passed to the upgrader object is of the form vMajor.Minor.Patch-eksd-tag
	// so it's safe to parse the version
	kubeVersion := semver.MajorMinor(u.kubernetesVersion)

	// From version 1.30 and above kubeadm upgrade needs special handling from legacy flow
	// as --config flag starts supporting and actual kubeadm upgradeConfiguration type and not cluster configuration type
	// Ref: https://github.com/kubernetes/kubernetes/pull/123068
	// Issue: https://github.com/kubernetes/kubeadm/issues/3054
	if semver.Compare(kubeVersion, kubeVersion130) >= 0 {
		updatedClusterConfig, err := u.ReadFile(newKubeAdmConfig)
		if err != nil {
			return err
		}
		err = u.kubeAdmUpgradeVersion130AndAbove(ctx, componentsDir, string(updatedClusterConfig))
		if err != nil {
			return err
		}
	} else {
		if err = u.appendKubeletConfig(ctx, newKubeAdmConfig); err != nil {
			return fmt.Errorf("appending kubelet config: %v", err)
		}
		kubeAdmUpgPlanCmd := []string{"kubeadm", "upgrade", "plan", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig}
		kubeAdmUpgPlan, err := u.ExecCommand(ctx, kubeAdmUpgPlanCmd[0], kubeAdmUpgPlanCmd[1:]...)
		if err != nil {
			return execError(kubeAdmUpgPlanCmd, string(kubeAdmUpgPlan))
		}
		logger.Info("components to be upgraded with kubeadm", "output", string(kubeAdmUpgPlan))

		kubeAdmUpgCmd := []string{"kubeadm", "upgrade", "apply", "--ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", "--config", newKubeAdmConfig, "--allow-experimental-upgrades", "--yes", "--force"}
		kubeAdmUpg, err := u.ExecCommand(ctx, kubeAdmUpgCmd[0], kubeAdmUpgCmd[1:]...)
		if err != nil {
			return execError(kubeAdmUpgCmd, string(kubeAdmUpg))
		}
		logger.Info("verbose output for kubeadm upgrade", "output", string(kubeAdmUpg))
	}
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

// kubeAdmUpgradeVersion130AndAbove upgrades first CP node for K8s version 1.30 & above
//
// As part of the upgrades:
//  1. Update the kubeadm clusterConfig config map with latest etcd version
//  2. Uses a kubeadm UpgradeConfiguration type for the upgrade command with apply and plan configurations
func (u *InPlaceUpgrader) kubeAdmUpgradeVersion130AndAbove(ctx context.Context, cmpDir, clusterConfig string) error {
	kubeAdmConfCM := generateKubeAdmConfCM(clusterConfig)
	kubeAdmConfCMData, err := yaml.Marshal(kubeAdmConfCM)
	if err != nil {
		return fmt.Errorf("marshaling kubeadm-config config map: %v", err)
	}

	KubeAdmConfCMYaml := fmt.Sprintf("%s/kubeadm-config-cm.yaml", cmpDir)
	err = u.WriteFile(KubeAdmConfCMYaml, kubeAdmConfCMData, fileMode640)
	if err != nil {
		return fmt.Errorf("writing kubeadm upgrade config: %v", err)
	}

	applyKubeAdmConfCMCmd := []string{"kubectl", "apply", "-f", KubeAdmConfCMYaml, "--kubeconfig", kubeConfigPath}
	out, err := u.ExecCommand(ctx, applyKubeAdmConfCMCmd[0], applyKubeAdmConfCMCmd[1:]...)
	if err != nil {
		return execError(applyKubeAdmConfCMCmd, string(out))
	}
	logger.Info("updated config map kubeadm-config on cluster")

	upgradeConfig := generateKubeAdmUpgradeConfig(u.kubernetesVersion)

	upgradeConfigData, err := yaml.Marshal(&upgradeConfig)
	if err != nil {
		return fmt.Errorf("marshaling kubeadm upgrade config: %v", err)
	}

	kubeAdmUpgradeConfigYaml := fmt.Sprintf("%s/kubeadm-upgrade-config.yaml", cmpDir)
	err = u.WriteFile(kubeAdmUpgradeConfigYaml, upgradeConfigData, fileMode640)
	if err != nil {
		return fmt.Errorf("writing kubeadm upgrade config: %v", err)
	}
	logger.Info("generated kubeadm upgrade config for k8s version >= 1.30", "fileLocation", kubeAdmUpgradeConfigYaml, "k8sVersion", u.kubernetesVersion)

	kubeAdmUpgPlanCmd := []string{"kubeadm", "upgrade", "plan", "--config", kubeAdmUpgradeConfigYaml}
	kubeAdmUpgPlan, err := u.ExecCommand(ctx, kubeAdmUpgPlanCmd[0], kubeAdmUpgPlanCmd[1:]...)
	if err != nil {
		return execError(kubeAdmUpgPlanCmd, string(kubeAdmUpgPlan))
	}
	logger.Info("components to be upgraded with kubeadm", "output", string(kubeAdmUpgPlan))

	kubeAdmUpgCmd := []string{"kubeadm", "upgrade", "apply", "--config", kubeAdmUpgradeConfigYaml}
	kubeAdmUpg, err := u.ExecCommand(ctx, kubeAdmUpgCmd[0], kubeAdmUpgCmd[1:]...)
	if err != nil {
		return execError(kubeAdmUpgCmd, string(kubeAdmUpg))
	}
	logger.Info("verbose output for kubeadm upgrade", "output", string(kubeAdmUpg))

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

	kubeAdmVersionCmd := []string{"kubeadm", "version", "-oshort"}
	version, err := u.ExecCommand(ctx, kubeAdmVersionCmd[0], kubeAdmVersionCmd[1:]...)
	if err != nil {
		return execError(kubeAdmVersionCmd, string(version))
	}
	logger.Info("current version of kubeadm", "cmd", "kubeadm version -oshort", "output", string(version))

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

	kubeAdmVersionCmd := []string{"kubeadm", "version", "-oshort"}
	version, err := u.ExecCommand(ctx, kubeAdmVersionCmd[0], kubeAdmVersionCmd[1:]...)
	if err != nil {
		return execError(kubeAdmVersionCmd, string(version))
	}
	logger.Info("current version of kubeadm", "cmd", "kubeadm version -oshort", "output", string(version))

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
	err = u.WriteFile(newKubeAdmConf, []byte(updatedConf), fileMode640)
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
	getKubeletConfCmd := []string{"kubectl", "get", "cm", "-n", kubeSystemNS, "kubelet-config", "-ojsonpath='{.data.kubelet}'", "--kubeconfig", kubeConfigPath}
	out, err := u.ExecCommand(ctx, getKubeletConfCmd[0], getKubeletConfCmd[1:]...)
	if err != nil {
		return execError(getKubeletConfCmd, string(out))
	}
	logger.Info("kubelet config as string", "out", string(out))

	trimmedOut := strings.Trim(string(out), "`'")
	conf = append(conf, []byte(trimmedOut)...)
	err = u.WriteFile(kubeAdmConf, []byte(conf), fileMode640)
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
	getCoreDNSConfCmd := []string{"kubectl", "get", "cm", "-n", kubeSystemNS, "coredns", "-oyaml", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true"}
	coreDNSConf, err := u.ExecCommand(ctx, getCoreDNSConfCmd[0], getCoreDNSConfCmd[1:]...)
	if err != nil {
		return execError(getCoreDNSConfCmd, string(coreDNSConf))
	}
	if len(coreDNSConf) > 0 {
		logger.Info("coreDNS config as string", "out", string(coreDNSConf))
		err = u.WriteFile(coreDNSBackup, coreDNSConf, fileMode640)
		if err != nil {
			return err
		}
	}
	deleteCoreDNSConfig := []string{"kubectl", "delete", "cm", "-n", kubeSystemNS, "coredns", "--kubeconfig", kubeConfigPath, "--ignore-not-found=true"}
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

func generateKubeAdmConfCM(clusterConfig string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeadmCMName,
			Namespace: kubeSystemNS,
		},
		Data: map[string]string{"ClusterConfiguration": clusterConfig},
	}
}

func generateKubeAdmUpgradeConfig(version string) UpgradeConfiguration {
	preflightErrorsList := []string{"CoreDNSUnsupportedPlugins", "CoreDNSMigration"}
	return UpgradeConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       upgradeConfigurationKind,
			APIVersion: kubeadmAPIv1beta4,
		},
		Apply: UpgradeApplyConfiguration{
			KubernetesVersion:         version,
			AllowExperimentalUpgrades: ptr.To(true),
			ForceUpgrade:              ptr.To(true),
			EtcdUpgrade:               ptr.To(true),
			IgnorePreflightErrors:     preflightErrorsList,
		},
		Plan: UpgradePlanConfiguration{
			KubernetesVersion:         version,
			AllowExperimentalUpgrades: ptr.To(true),
			IgnorePreflightErrors:     preflightErrorsList,
			PrintConfig:               ptr.To(true),
		},
	}
}

package upgrade

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// UpgradeConfiguration contains a list of options that are specific to "kubeadm upgrade" subcommands.
type UpgradeConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// Apply holds a list of options that are specific to the "kubeadm upgrade apply" command.
	// +optional
	Apply UpgradeApplyConfiguration `json:"apply,omitempty"`

	// Plan holds a list of options that are specific to the "kubeadm upgrade plan" command.
	// +optional
	Plan UpgradePlanConfiguration `json:"plan,omitempty"`
}

// UpgradeApplyConfiguration contains a list of configurable options which are specific to the  "kubeadm upgrade apply" command.
type UpgradeApplyConfiguration struct {
	// KubernetesVersion is the target version of the control plane.
	// +optional
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// AllowExperimentalUpgrades instructs kubeadm to show unstable versions of Kubernetes as an upgrade
	// alternative and allows upgrading to an alpha/beta/release candidate version of Kubernetes.
	// Default: false
	// +optional
	AllowExperimentalUpgrades *bool `json:"allowExperimentalUpgrades,omitempty"`

	// Enable AllowRCUpgrades will show release candidate versions of Kubernetes as an upgrade alternative and
	// allows upgrading to a release candidate version of Kubernetes.
	// +optional
	AllowRCUpgrades *bool `json:"allowRCUpgrades,omitempty"`

	// CertificateRenewal instructs kubeadm to execute certificate renewal during upgrades.
	// Defaults to true.
	// +optional
	CertificateRenewal *bool `json:"certificateRenewal,omitempty"`

	// DryRun tells if the dry run mode is enabled, don't apply any change if it is and just output what would be done.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`

	// EtcdUpgrade instructs kubeadm to execute etcd upgrade during upgrades.
	// Defaults to true.
	// +optional
	EtcdUpgrade *bool `json:"etcdUpgrade,omitempty"`

	// ForceUpgrade flag instructs kubeadm to upgrade the cluster without prompting for confirmation.
	// +optional
	ForceUpgrade *bool `json:"forceUpgrade,omitempty"`

	// IgnorePreflightErrors provides a slice of pre-flight errors to be ignored during the upgrade process, e.g. 'IsPrivilegedUser,Swap'.
	// Value 'all' ignores errors from all checks.
	// +optional
	IgnorePreflightErrors []string `json:"ignorePreflightErrors,omitempty"`

	// PrintConfig specifies whether the configuration file that will be used in the upgrade should be printed or not.
	// +optional
	PrintConfig *bool `json:"printConfig,omitempty"`

	// SkipPhases is a list of phases to skip during command execution.
	// NOTE: This field is currently ignored for "kubeadm upgrade apply", but in the future it will be supported.
	SkipPhases []string

	// ImagePullSerial specifies if image pulling performed by kubeadm must be done serially or in parallel.
	// Default: true
	// +optional
	ImagePullSerial *bool `json:"imagePullSerial,omitempty"`
}

// UpgradePlanConfiguration contains a list of configurable options which are specific to the "kubeadm upgrade plan" command.
type UpgradePlanConfiguration struct {
	// KubernetesVersion is the target version of the control plane.
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// AllowExperimentalUpgrades instructs kubeadm to show unstable versions of Kubernetes as an upgrade
	// alternative and allows upgrading to an alpha/beta/release candidate version of Kubernetes.
	// Default: false
	// +optional
	AllowExperimentalUpgrades *bool `json:"allowExperimentalUpgrades,omitempty"`

	// Enable AllowRCUpgrades will show release candidate versions of Kubernetes as an upgrade alternative and
	// allows upgrading to a release candidate version of Kubernetes.
	// +optional
	AllowRCUpgrades *bool `json:"allowRCUpgrades,omitempty"`

	// DryRun tells if the dry run mode is enabled, don't apply any change if it is and just output what would be done.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`

	// IgnorePreflightErrors provides a slice of pre-flight errors to be ignored during the upgrade process, e.g. 'IsPrivilegedUser,Swap'.
	// Value 'all' ignores errors from all checks.
	// +optional
	IgnorePreflightErrors []string `json:"ignorePreflightErrors,omitempty"`

	// PrintConfig specifies whether the configuration file that will be used in the upgrade should be printed or not.
	// +optional
	PrintConfig *bool `json:"printConfig,omitempty"`
}

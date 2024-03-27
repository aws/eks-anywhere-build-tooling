package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/upgrader/upgrade"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var upgradeNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Upgrade node",
	Long:  "Use InPlace Upgrader upgrade node to upgrade kubeadm on the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		nodeType := viper.GetString("type")
		k8sVersion := viper.GetString("k8sVersion")
		if k8sVersion == "" {
			return errors.New("k8sVersion flag has to be set for upgrade node command")
		}

		etcdVersion := viper.GetString("etcdVersion")
		if etcdVersion == "" {
			etcdVersion = "NO_UPDATE"
		}

		err = upgradeNode(cmd.Context(), nodeType, k8sVersion, etcdVersion)
		if err != nil {
			log.Fatalf("upgrade node failed: %v", err)
		}
		return nil
	},
}

func init() {
	var err error

	upgradeCmd.AddCommand(upgradeNodeCmd)
	upgradeNodeCmd.Flags().String("type", "", "Node type flag")
	upgradeNodeCmd.Flags().String("k8sVersion", "", "kubernetes version flag")
	upgradeNodeCmd.Flags().String("etcdVersion", "", "etcd version flag")
	err = viper.BindPFlags(upgradeNodeCmd.Flags())
	if err != nil {
		log.Fatalf("Error initializing flags: %v", err)
	}
}

func upgradeNode(ctx context.Context, nodeType, k8sVersion, etcdVersion string) error {
	upg := upgrade.NewUpgrader(upgrade.WithKubernetesVersion(k8sVersion), upgrade.WithEtcdVersion(etcdVersion))

	switch nodeType {
	case "FirstCP":
		if err := upg.KubeAdmInFirstCP(ctx); err != nil {
			return fmt.Errorf("upgrading kubeadm in first controlplane node: %v", err)
		}
	case "RestCP":
		if err := upg.KubeAdmInRestCP(ctx); err != nil {
			return fmt.Errorf("upgrading kubeadm in controlplane node: %v", err)
		}
	case "Worker":
		if err := upg.KubeAdmInWorker(ctx); err != nil {
			return fmt.Errorf("upgrading kubeadm in worker node: %v", err)
		}
	default:
		return fmt.Errorf("invalid node type, please specify one of the three types: FirstCP, RestCP or Worker")
	}

	return nil
}

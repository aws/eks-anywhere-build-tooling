package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

var rootCmd = &cobra.Command{
	Use:              "upgrader",
	Short:            "EKS Anywhere InPlace upgrader",
	Long:             `Use EKS Anywhere InPlace upgrader to upgrade your nodes InPlace`,
	PersistentPreRun: rootPersistentPreRun,
}

func init() {
	rootCmd.PersistentFlags().IntP("verbosity", "v", 0, "Set the log level verbosity")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatalf("failed to bind flags for root: %v", err)
	}
}

func rootPersistentPreRun(cmd *cobra.Command, args []string) {
	if err := initLogger(); err != nil {
		log.Fatal(err)
	}
}

func initLogger() error {
	if err := logger.Init(viper.GetInt("verbosity")); err != nil {
		return fmt.Errorf("failed init zap logger in root command: %v", err)
	}

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

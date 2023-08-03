package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// rootCmd is the top-level version-tracker command used to track projects and their versions.
var rootCmd = &cobra.Command{
	Use:              "version-tracker",
	Short:            "Amazon EKS Anywhere Build-tooling Version Tracker",
	Long:             "Use version-tracker to check and update the Git tag and Go version used to build projects in aws/eks-anywhere-build-tooling",
	PersistentPreRun: rootPersistentPreRun,
}

func init() {
	rootCmd.PersistentFlags().IntP("verbosity", "v", 0, "Set the logging verbosity level")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatalf("failed to bind flags to root command: %v", err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func rootPersistentPreRun(cmd *cobra.Command, args []string) {
	if err := initLogger(); err != nil {
		log.Fatal(err)
	}
}

func initLogger() error {
	if err := logger.Init(viper.GetInt("verbosity")); err != nil {
		return fmt.Errorf("failed to init Zap logger in root command: %v", err)
	}

	return nil
}

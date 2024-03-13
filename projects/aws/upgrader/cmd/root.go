package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

var rootCmd = &cobra.Command{
	Use:              "upgrader",
	Short:            "EKS Anywhere InPlace upgrader",
	Long:             `Use EKS Anywhere InPlace upgrader to upgrade your nodes in place`,
	PersistentPreRun: rootPersistentPreRun,
}

func init() {
	rootCmd.PersistentFlags().IntP("verbosity", "v", 0, "Set the log level verbosity")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatalf("failed to bind flags for root: %v", err)
	}
}

func prerunCmdBindFlags(cmd *cobra.Command, args []string) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		err := viper.BindPFlag(flag.Name, flag)
		if err != nil {
			log.Fatalf("Error initializing flags: %v", err)
		}
	})
}

func rootPersistentPreRun(cmd *cobra.Command, args []string) {
	if err := initLogger(); err != nil {
		log.Fatalf("failed to init root command: %v", err)
	}
}

func initLogger() error {
	if err := logger.Init(viper.GetInt("verbosity")); err != nil {
		return fmt.Errorf("init zap logger: %v", err)
	}

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

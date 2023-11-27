package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "image-builder",
	Short: "Amazon EKS Anywhere Image Builder",
	Long:  `Use image-builder to build your own EKS Anywhere node image`,
}

func init() {
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatalf("failed to bind flags for root: %v", err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

package cmd

import (
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade command",
	Long:  "Use InPlace Upgrader upgrade to run different upgrade commands on the node",
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

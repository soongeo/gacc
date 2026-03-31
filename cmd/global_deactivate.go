package cmd

import (
	"fmt"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/spf13/cobra"
)

var globalDeactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Clear globally active gacc account settings.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌍 Clearing global gacc account settings...")

		_ = git.UnsetGlobalUserConfig()
		_ = git.UnsetGlobalSSHCommand()

		fmt.Println("✅ Global Git identity and SSH command cleared.")
	},
}

func init() {
	globalCmd.AddCommand(globalDeactivateCmd)
}

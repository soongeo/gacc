package cmd

import "github.com/spf13/cobra"

var globalCmd = &cobra.Command{
	Use:   "global",
	Short: "Manage the globally active gacc account.",
}

func init() {
	rootCmd.AddCommand(globalCmd)
}

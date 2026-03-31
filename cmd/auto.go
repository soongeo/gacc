package cmd

import "github.com/spf13/cobra"

var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Manage path-based automatic account switching via Git includeIf.",
}

func init() {
	rootCmd.AddCommand(autoCmd)
}

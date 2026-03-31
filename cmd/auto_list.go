package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var autoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured automatic includeIf rules.",
	Run: func(cmd *cobra.Command, args []string) {
		rules := listStoredAutoRules()
		if len(rules) == 0 {
			fmt.Println("No automatic account rules configured.")
			return
		}

		fmt.Println("Automatic account rules:")
		for _, rule := range rules {
			fmt.Printf("  - %s -> %s\n", rule.Directory, rule.Account)
		}
	},
}

func init() {
	autoCmd.AddCommand(autoListCmd)
}

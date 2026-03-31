package cmd

import (
	"fmt"
	"os"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/spf13/cobra"
)

var autoRemoveCmd = &cobra.Command{
	Use:   "remove [name] [directory]",
	Short: "Remove an automatic account rule for a directory.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := args[0]
		directory := args[1]

		rule, found, err := findStoredAutoRule(account, directory)
		if err != nil {
			fmt.Printf("❌ Failed to resolve rule: %v\n", err)
			os.Exit(1)
		}
		if !found {
			fmt.Printf("❌ Auto rule for '%s' under '%s' not found.\n", account, directory)
			os.Exit(1)
		}

		if err := git.RemoveGlobalIncludeIf(rule.Condition, rule.Include); err != nil {
			fmt.Printf("❌ Failed to remove includeIf rule: %v\n", err)
			os.Exit(1)
		}

		if rule.Include != "" {
			if err := os.Remove(rule.Include); err != nil && !os.IsNotExist(err) {
				fmt.Printf("❌ Failed to remove include file: %v\n", err)
				os.Exit(1)
			}
		}

		if err := deleteAutoRule(account, rule.Condition); err != nil {
			fmt.Printf("❌ Failed to update stored auto rules: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Auto rule removed: '%s' will no longer be auto-selected under %s\n", account, rule.Condition)
	},
}

func init() {
	autoCmd.AddCommand(autoRemoveCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
)

var autoAddCmd = &cobra.Command{
	Use:               "add [name] [directory]",
	Short:             "Automatically use an account for repositories under a directory.",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: accountNameCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {
		account := args[0]
		directory := args[1]

		accounts, err := ssh.ListAccounts()
		if err != nil {
			fmt.Printf("❌ Failed to load accounts: %v\n", err)
			os.Exit(1)
		}

		found := false
		for _, existing := range accounts {
			if existing == account {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("❌ Account '%s' not found.\n", account)
			os.Exit(1)
		}

		condition, err := git.NormalizeGitDirCondition(directory)
		if err != nil {
			fmt.Printf("❌ Failed to normalize directory: %v\n", err)
			os.Exit(1)
		}

		includePath, err := writeIncludeFile(account, condition)
		if err != nil {
			fmt.Printf("❌ Failed to write include file: %v\n", err)
			os.Exit(1)
		}

		if err := git.AddGlobalIncludeIf(condition, includePath); err != nil {
			fmt.Printf("❌ Failed to register includeIf rule: %v\n", err)
			os.Exit(1)
		}

		if err := saveAutoRule(autoRule{
			Account:   account,
			Directory: directory,
			Condition: condition,
			Include:   includePath,
		}); err != nil {
			fmt.Printf("❌ Failed to save auto rule metadata: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Auto rule added: '%s' will be used under %s\n", account, condition)
	},
}

func init() {
	autoCmd.AddCommand(autoAddCmd)
}

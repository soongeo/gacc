package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered Git accounts.",
	Run: func(cmd *cobra.Command, args []string) {
		accounts, err := ssh.ListAccounts()
		if err != nil {
			fmt.Printf("❌ Error reading SSH configuration: %v\n", err)
			os.Exit(1)
		}

		if len(accounts) == 0 {
			fmt.Println("🤷 No registered accounts found.")
			fmt.Println("👉 Use 'gacc add [name]' to register a new account.")
			return
		}

		status, err := collectResolvedStatus()
		if err != nil {
			fmt.Printf("❌ Failed to resolve current status: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("📋 Registered Git accounts:")
		for _, account := range accounts {
			labels := accountLabelsForList(account, status)
			if len(labels) == 0 {
				fmt.Printf("  - %s\n", account)
				continue
			}
			fmt.Printf("  - %s [%s]\n", account, strings.Join(labels, ", "))
		}

		fmt.Println()
		if status.InGitRepo {
			fmt.Printf("Current repo: %s\n", status.RepoRoot)
		}
		fmt.Printf("Auto match: %s\n", valueOrDefault(func() string {
			if status.AutoRule == nil {
				return ""
			}
			return status.AutoRule.Account
		}(), "(none)"))
		fmt.Printf("Global default: %s\n", valueOrDefault(globalAccountFromSSHCommand(status.GlobalSSHCmd), "(none)"))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

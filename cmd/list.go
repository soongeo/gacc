package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/soongeo/gacc/pkg/git"
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

		var activeAccount string
		if git.IsInsideWorkTree() {
			if remoteUrl, err := git.GetRemoteURL("origin"); err == nil {
				for _, account := range accounts {
					targetHost := fmt.Sprintf("github.com-%s", account)
					if strings.Contains(remoteUrl, targetHost+":") || strings.Contains(remoteUrl, targetHost+"/") {
						activeAccount = account
						break
					}
				}
			}
		}

		fmt.Println("📋 Registered Git accounts:")
		for _, account := range accounts {
			if account == activeAccount {
				fmt.Printf("  - %s 🌟 (active)\n", account)
			} else {
				fmt.Printf("  - %s\n", account)
			}
		}
		
		fmt.Println()
		if activeAccount == "" {
			fmt.Println("💡 To use an account, run: gacc activate [name]")
		} else {
			fmt.Println("💡 To deactivate the current account, run: gacc deactivate")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

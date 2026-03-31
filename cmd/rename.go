package cmd

import (
	"fmt"
	"os"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var renameCmd = &cobra.Command{
	Use:   "rename [old-name] [new-name]",
	Short: "Rename a registered gacc account alias.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName := args[0]
		newName := args[1]

		accounts, err := ssh.ListAccounts()
		if err != nil {
			fmt.Printf("❌ Failed to load accounts: %v\n", err)
			os.Exit(1)
		}

		exists := false
		for _, account := range accounts {
			if account == oldName {
				exists = true
			}
			if account == newName {
				fmt.Printf("❌ Account '%s' already exists.\n", newName)
				os.Exit(1)
			}
		}
		if !exists {
			fmt.Printf("❌ Account '%s' does not exist.\n", oldName)
			os.Exit(1)
		}

		if err := ssh.RenameAccount(oldName, newName); err != nil {
			fmt.Printf("❌ Failed to rename account: %v\n", err)
			os.Exit(1)
		}

		accountsData := viper.GetStringMap("accounts")
		if oldData, ok := accountsData[oldName]; ok {
			accountsData[newName] = oldData
			delete(accountsData, oldName)
			viper.Set("accounts", accountsData)
			if err := viper.WriteConfig(); err != nil {
				fmt.Printf("❌ Failed to update config: %v\n", err)
				os.Exit(1)
			}
		}

		if git.IsInsideWorkTree() {
			remoteURL, err := git.GetRemoteURL("origin")
			if err == nil {
				newURL, parseErr := git.ParseAndReplaceRemoteHost(remoteURL, newName)
				if parseErr == nil && newURL != remoteURL {
					_ = git.SetRemoteURL("origin", newURL)
				}
			}
		}

		fmt.Printf("✅ Account '%s' renamed to '%s'.\n", oldName, newName)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

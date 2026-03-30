package cmd

import (
	"fmt"
	"os"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var activateCmd = &cobra.Command{
	Use:   "activate [name]",
	Short: "Activate a gacc account for the current Git repository.",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: accountNameCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {
		accountName := args[0]
		
		fmt.Printf("🚀 Activating account '%s' for current project...\n", accountName)

		// 1. Verify we are in a Git repo
		if !git.IsInsideWorkTree() {
			fmt.Println("❌ Error: Current directory is not a Git repository.")
			os.Exit(1)
		}

		// 2. Verify account exists in SSH config
		accounts, err := ssh.ListAccounts()
		if err != nil {
			fmt.Printf("❌ Failed to read SSH accounts: %v\n", err)
			os.Exit(1)
		}
		
		accountExists := false
		for _, acc := range accounts {
			if acc == accountName {
				accountExists = true
				break
			}
		}
		if !accountExists {
			fmt.Printf("❌ Account '%s' not found. Please add it first using 'gacc add %s'.\n", accountName, accountName)
			os.Exit(1)
		}

		// 3. Update Git Remote URL
		fmt.Println("\n[1/2] Updating origin URL for SSH alias routing...")
		remoteUrl, err := git.GetRemoteURL("origin")
		if err != nil {
			fmt.Printf("⚠️ Warning: Could not get 'origin' remote URL. Attempting to proceed without updating remote. error: %v\n", err)
		} else {
			newUrl, err := git.ParseAndReplaceRemoteHost(remoteUrl, accountName)
			if err != nil {
				fmt.Printf("❌ Failed to parse remote URL: %v\n", err)
				os.Exit(1)
			}
			if newUrl != remoteUrl {
				if err := git.SetRemoteURL("origin", newUrl); err != nil {
					fmt.Printf("❌ Failed to update origin URL: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("✅ Remote 'origin' URL updated to: %s\n", newUrl)
			} else {
				fmt.Println("✅ Remote 'origin' URL is already configured correctly.")
			}
		}

		// 4. Update Git User Config
		fmt.Println("\n[2/2] Overriding local Git user configurations...")
		name := viper.GetString("accounts." + accountName + ".name")
		email := viper.GetString("accounts." + accountName + ".email")

		if name != "" && email != "" {
			if err := git.SetLocalUserConfig(name, email); err != nil {
				fmt.Printf("❌ Failed to set local Git commit configs: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Local Git config set to %s <%s>\n", name, email)
		} else {
			fmt.Printf("⚠️ No profile info (name/email) stored for '%s'. Skipping git config mutation.\n", accountName)
			fmt.Println("   Note: If you run 'add' again, this information will be seamlessly fetched & saved.")
		}

		fmt.Printf("\n🎉 Success! The current project is now using '%s'.\n", accountName)
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}

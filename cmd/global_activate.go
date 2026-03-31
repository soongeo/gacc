package cmd

import (
	"fmt"
	"os"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var globalActivateCmd = &cobra.Command{
	Use:               "activate [name]",
	Short:             "Activate a gacc account globally for Git operations.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: accountNameCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {
		accountName, err := resolveAccountForOptionalArg(args)
		if err != nil {
			fmt.Printf("❌ Failed to resolve account: %v\n", err)
			os.Exit(1)
		}

		accounts, err := ssh.ListAccounts()
		if err != nil {
			fmt.Printf("❌ Failed to read SSH accounts: %v\n", err)
			os.Exit(1)
		}

		accountExists := false
		for _, account := range accounts {
			if account == accountName {
				accountExists = true
				break
			}
		}
		if !accountExists {
			fmt.Printf("❌ Account '%s' not found. Please add it first using 'gacc add %s'.\n", accountName, accountName)
			os.Exit(1)
		}

		name := viper.GetString("accounts." + accountName + ".name")
		email := viper.GetString("accounts." + accountName + ".email")
		privateKeyPath, err := ssh.PrivateKeyPath(accountName)
		if err != nil {
			fmt.Printf("❌ Failed to resolve SSH key path: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🌍 Activating account '%s' globally...\n", accountName)

		if name != "" || email != "" {
			if err := git.SetGlobalUserConfig(name, email); err != nil {
				fmt.Printf("❌ Failed to set global Git identity: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Global Git identity set to %s <%s>\n", name, email)
		} else {
			fmt.Printf("⚠️ No profile info (name/email) stored for '%s'. Skipping global user config mutation.\n", accountName)
		}

		sshCommand := fmt.Sprintf("ssh -i \"%s\" -o IdentitiesOnly=yes", privateKeyPath)
		if err := git.SetGlobalSSHCommand(sshCommand); err != nil {
			fmt.Printf("❌ Failed to set global SSH command: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Global SSH command set for '%s'\n", accountName)

		fmt.Printf("\n🎉 Success! Global Git operations now default to '%s'.\n", accountName)
	},
}

func init() {
	globalCmd.AddCommand(globalActivateCmd)
}

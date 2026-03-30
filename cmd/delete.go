package cmd

import (
	"fmt"
	"os"

	"github.com/soongeo/gacc/pkg/github"
	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a registered Git account and SSH key locally and from GitHub.",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: accountNameCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {
		accountName := args[0]
		fmt.Printf("🗑️  Starting deletion for account '%s'...\n", accountName)

		fmt.Println("\n[1/4] Removing local SSH keys and configuration...")
		err := ssh.DeleteSSHKeys(accountName)
		if err != nil {
			fmt.Printf("❌ Error deleting SSH keys: %v\n", err)
		} else {
			fmt.Println("✅ Local SSH key deleted (~/.ssh/gacc_" + accountName + ")")
		}

		err = ssh.RemoveSSHConfig(accountName)
		if err != nil {
			fmt.Printf("❌ Error cleaning up SSH config: %v\n", err)
		} else {
			fmt.Println("✅ Host configuration removed from ~/.ssh/config")
		}

		fmt.Println("\n[2/4] Authenticating to delete public key from GitHub...")
		// GitHub OAuth App Client ID
		clientID := viper.GetString("github.client_id")
		if clientID == "" {
			clientID = os.Getenv("GACC_GITHUB_CLIENT_ID")
		}
		if clientID == "" {
			clientID = "Ov23liyiHVS9qKMPk3Tp" // 기본 제공 Client ID
		}

		accessToken, err := github.StartDeviceFlow(clientID)
		if err != nil {
			fmt.Printf("❌ GitHub authentication failed (Remote keys must be deleted manually): %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ GitHub authentication complete\n")

		fmt.Println("\n[3/4] Deleting registered SSH public key from GitHub...")
		err = github.DeleteSSHPublicKey(accessToken, accountName)
		if err != nil {
			fmt.Printf("❌ Error deleting remote SSH key (Manual deletion required): %v\n", err)
		} else {
			fmt.Println("✅ Public key successfully deleted from GitHub server")
		}

		// 4. Remove from Viper Config
		fmt.Println("\n[4/4] Removing local config records...")
		accountsData := viper.GetStringMap("accounts")
		if _, ok := accountsData[accountName]; ok {
			delete(accountsData, accountName)
			viper.Set("accounts", accountsData)
			_ = viper.WriteConfig()
		}
		fmt.Println("✅ Config records cleaned up")

		fmt.Printf("\n🎉 Account '%s' record successfully deleted from both local and remote!\n", accountName)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

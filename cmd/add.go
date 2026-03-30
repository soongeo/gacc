package cmd

import (
	"fmt"
	"os"

	"github.com/soongeo/gacc/pkg/github"
	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new Git account and SSH key.",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: accountNameCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {
		accountName := args[0]
		fmt.Printf("🛠️  Starting setup for account '%s'...\n", accountName)

		fmt.Println("\n[1/4] Generating SSH key...")
		pubKey, err := ssh.GenerateAndSaveEd25519(accountName)
		if err != nil {
			fmt.Printf("❌ Failed to generate SSH key: %v\n", err)
			os.Exit(1)
		}
		
		err = ssh.UpdateSSHConfig(accountName)
		if err != nil {
			fmt.Printf("❌ Failed to update SSH config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ SSH key generated and ~/.ssh/config updated\n")

		fmt.Println("\n[2/4] Starting GitHub Authentication (Device Flow)...")
		
		// GitHub OAuth App Client ID
		clientID := viper.GetString("github.client_id")
		if clientID == "" {
			clientID = os.Getenv("GACC_GITHUB_CLIENT_ID")
		}
		if clientID == "" {
			clientID = "Ov23liyiHVS9qKMPk3Tp" // Default Client ID
		}

		accessToken, err := github.StartDeviceFlow(clientID)
		if err != nil {
			fmt.Printf("❌ GitHub authentication failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ GitHub authentication complete\n")

		fmt.Println("\n[3/4] Registering SSH public key to GitHub...")
		err = github.AddSSHPublicKey(accessToken, accountName, pubKey)
		if err != nil {
			fmt.Printf("❌ Failed to upload SSH key: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n[4/4] Fetching GitHub Profile Info...")
		userInfo, err := github.GetUserInfo(accessToken)
		if err == nil {
			viper.Set("accounts."+accountName+".name", userInfo.Name)
			viper.Set("accounts."+accountName+".email", userInfo.Email)
			_ = viper.WriteConfig()
			fmt.Printf("✅ Saved profile: %s (%s)\n", userInfo.Name, userInfo.Email)
		} else {
			fmt.Printf("⚠️ Warning: Could not fetch user profile: %v\n", err)
		}
		
		fmt.Printf("\n🎉 Account '%s' setup completed successfully!\n", accountName)
		fmt.Printf("You can now test the connection with: 'ssh -T github.com-%s'\n", accountName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

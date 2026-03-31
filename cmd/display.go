package cmd

import (
	"fmt"
	"os"

	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
)

var displayCmd = &cobra.Command{
	Use:               "display [name]",
	Short:             "Display the public SSH key for an account.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: accountNameCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {
		accountName, err := resolveAccountForOptionalArg(args)
		if err != nil {
			fmt.Printf("❌ Failed to resolve account: %v\n", err)
			os.Exit(1)
		}

		publicKey, err := ssh.ReadPublicKey(accountName)
		if err != nil {
			fmt.Printf("❌ Failed to read public key for '%s': %v\n", accountName, err)
			os.Exit(1)
		}

		fmt.Printf("🔑 Public key for '%s':\n\n%s\n", accountName, publicKey)
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)
}

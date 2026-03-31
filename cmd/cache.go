package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:               "cache [name]",
	Short:             "Add an account SSH key to ssh-agent.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: accountNameCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {
		accountName, err := resolveAccountForOptionalArg(args)
		if err != nil {
			fmt.Printf("❌ Failed to resolve account: %v\n", err)
			os.Exit(1)
		}

		privateKeyPath, err := ssh.PrivateKeyPath(accountName)
		if err != nil {
			fmt.Printf("❌ Failed to resolve private key path: %v\n", err)
			os.Exit(1)
		}

		execCmd := exec.Command("ssh-add", privateKeyPath)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Stdin = os.Stdin

		if err := execCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to cache SSH key for '%s': %v\n", accountName, err)
			os.Exit(1)
		}

		fmt.Printf("✅ SSH key for '%s' added to ssh-agent.\n", accountName)
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show local, auto, and global gacc status for the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		status, err := collectResolvedStatus()
		if err != nil {
			fmt.Printf("❌ Failed to resolve status: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Current directory: %s\n", status.CWD)
		if status.InGitRepo {
			fmt.Printf("Git repository: %s\n", status.RepoRoot)
		} else {
			fmt.Println("Git repository: (not inside a git repository)")
		}

		fmt.Println()
		fmt.Println("Local manual status:")
		if status.ManualAccount == "" && status.LocalUserName == "" && status.LocalUserEmail == "" && status.LocalSSHCommand == "" {
			fmt.Println("  - account: (not set)")
		} else {
			fmt.Printf("  - account: %s\n", valueOrDefault(status.ManualAccount, "(not set)"))
			fmt.Printf("  - identity: %s\n", formatIdentity(status.LocalUserName, status.LocalUserEmail))
			fmt.Printf("  - sshCommand: %s\n", valueOrDefault(status.LocalSSHCommand, "(not set)"))
			fmt.Printf("  - remote: %s\n", valueOrDefault(status.ManualRemoteURL, "(not set)"))
		}

		fmt.Println()
		fmt.Println("Auto status:")
		if status.AutoRule == nil {
			fmt.Println("  - rule: (not matched)")
		} else {
			fmt.Printf("  - account: %s\n", status.AutoRule.Account)
			fmt.Printf("  - directory: %s\n", status.AutoRule.Directory)
			fmt.Printf("  - condition: %s\n", status.AutoRule.Condition)
		}

		fmt.Println()
		fmt.Println("Global status:")
		fmt.Printf("  - identity: %s\n", formatIdentity(status.GlobalUserName, status.GlobalUserEmail))
		fmt.Printf("  - sshCommand: %s\n", valueOrDefault(status.GlobalSSHCmd, "(not set)"))
		fmt.Printf("  - account: %s\n", valueOrDefault(globalAccountFromSSHCommand(status.GlobalSSHCmd), "(not inferred)"))
	},
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

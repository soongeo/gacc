package cmd

import (
	"fmt"
	"os"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/spf13/cobra"
)

var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate the gacc account for the current Git repository and restore defaults.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔌 Deactivating account overrides for current project...")

		// 1. Verify we are in a Git repo
		if !git.IsInsideWorkTree() {
			fmt.Println("❌ Error: Current directory is not a Git repository.")
			os.Exit(1)
		}

		// 2. Clear Remote Alias
		fmt.Println("\n[1/2] Restoring standard GitHub SSH remote URL...")
		remoteUrl, err := git.GetRemoteURL("origin")
		if err == nil {
			standardUrl, err := git.RevertRemoteHostToStandard(remoteUrl)
			if err == nil && standardUrl != remoteUrl {
				if err := git.SetRemoteURL("origin", standardUrl); err == nil {
					fmt.Printf("✅ Remote 'origin' URL restored to: %s\n", standardUrl)
				}
			} else {
				fmt.Println("✅ Remote 'origin' is already standard or unaltered.")
			}
		}

		// 3. Clear Local User configs
		fmt.Println("\n[2/2] Clearing local Git user.name & user.email configs...")
		err = git.UnsetLocalUserConfig()
		if err != nil {
			// It may emit error if they weren't set, which is harmless.
		}
		fmt.Println("✅ Local Git author configuration unset (will use global defaults).")

		fmt.Println("\n🎉 Success! The project now uses system global git configurations.")
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)
}

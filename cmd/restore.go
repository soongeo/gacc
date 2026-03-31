package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore [archive]",
	Short: "Restore gacc config and managed SSH keys from a backup archive.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := gaccConfigFilePath()
		if err != nil {
			fmt.Printf("❌ Failed to resolve gacc config path: %v\n", err)
			os.Exit(1)
		}

		sshDir, err := ssh.SSHDir()
		if err != nil {
			fmt.Printf("❌ Failed to resolve SSH directory: %v\n", err)
			os.Exit(1)
		}

		archiveFile, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("❌ Failed to open backup archive: %v\n", err)
			os.Exit(1)
		}
		defer archiveFile.Close()

		gzipReader, err := gzip.NewReader(archiveFile)
		if err != nil {
			fmt.Printf("❌ Failed to read gzip archive: %v\n", err)
			os.Exit(1)
		}
		defer gzipReader.Close()

		tarReader := tar.NewReader(gzipReader)
		var restoredAccounts []string

		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("❌ Failed while reading archive: %v\n", err)
				os.Exit(1)
			}

			switch header.Name {
			case "config/config.yaml":
				if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
					fmt.Printf("❌ Failed to prepare config directory: %v\n", err)
					os.Exit(1)
				}
				if err := writeFileFromArchive(configPath, tarReader, os.FileMode(header.Mode)); err != nil {
					fmt.Printf("❌ Failed to restore gacc config: %v\n", err)
					os.Exit(1)
				}
			case "manifest/accounts.txt":
				content, err := io.ReadAll(tarReader)
				if err != nil {
					fmt.Printf("❌ Failed to read account manifest: %v\n", err)
					os.Exit(1)
				}
				for _, line := range strings.Split(string(content), "\n") {
					line = strings.TrimSpace(line)
					if line != "" {
						restoredAccounts = append(restoredAccounts, line)
					}
				}
			default:
				if strings.HasPrefix(header.Name, "ssh/") {
					targetPath := filepath.Join(sshDir, filepath.Base(header.Name))
					if err := os.MkdirAll(filepath.Dir(targetPath), 0700); err != nil {
						fmt.Printf("❌ Failed to prepare SSH directory: %v\n", err)
						os.Exit(1)
					}
					if err := writeFileFromArchive(targetPath, tarReader, os.FileMode(header.Mode)); err != nil {
						fmt.Printf("❌ Failed to restore SSH file '%s': %v\n", header.Name, err)
						os.Exit(1)
					}
				}
			}
		}

		if len(restoredAccounts) == 0 {
			accounts, err := ssh.ListAccounts()
			if err == nil {
				restoredAccounts = accounts
			}
		}

		for _, account := range restoredAccounts {
			if err := ssh.UpdateSSHConfig(account); err != nil {
				fmt.Printf("❌ Failed to restore SSH config entry for '%s': %v\n", account, err)
				os.Exit(1)
			}
		}

		fmt.Printf("✅ Restore completed from %s\n", args[0])
	},
}

func writeFileFromArchive(targetPath string, reader io.Reader, mode os.FileMode) error {
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return os.WriteFile(targetPath, content, mode)
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}

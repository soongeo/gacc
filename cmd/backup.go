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

var backupCmd = &cobra.Command{
	Use:   "backup [archive]",
	Short: "Backup gacc config and managed SSH keys to a tar.gz archive.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputPath := defaultBackupArchiveName()
		if len(args) == 1 && strings.TrimSpace(args[0]) != "" {
			outputPath = args[0]
		}

		configPath, err := gaccConfigFilePath()
		if err != nil {
			fmt.Printf("❌ Failed to resolve gacc config path: %v\n", err)
			os.Exit(1)
		}

		accounts, err := ssh.ListAccounts()
		if err != nil {
			fmt.Printf("❌ Failed to load accounts: %v\n", err)
			os.Exit(1)
		}

		file, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("❌ Failed to create backup archive: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		gzipWriter := gzip.NewWriter(file)
		defer gzipWriter.Close()

		tarWriter := tar.NewWriter(gzipWriter)
		defer tarWriter.Close()

		if err := addFileToArchive(tarWriter, configPath, "config/config.yaml"); err != nil && !os.IsNotExist(err) {
			fmt.Printf("❌ Failed to add config to backup: %v\n", err)
			os.Exit(1)
		}

		manifest := strings.Join(accounts, "\n")
		if err := writeArchiveContent(tarWriter, "manifest/accounts.txt", []byte(manifest), 0644); err != nil {
			fmt.Printf("❌ Failed to add account manifest to backup: %v\n", err)
			os.Exit(1)
		}

		for _, account := range accounts {
			privateKeyPath, err := ssh.PrivateKeyPath(account)
			if err != nil {
				fmt.Printf("❌ Failed to resolve private key path for '%s': %v\n", account, err)
				os.Exit(1)
			}
			publicKeyPath := privateKeyPath + ".pub"

			if err := addFileToArchive(tarWriter, privateKeyPath, filepath.ToSlash(filepath.Join("ssh", filepath.Base(privateKeyPath)))); err != nil {
				fmt.Printf("❌ Failed to add private key for '%s': %v\n", account, err)
				os.Exit(1)
			}
			if err := addFileToArchive(tarWriter, publicKeyPath, filepath.ToSlash(filepath.Join("ssh", filepath.Base(publicKeyPath)))); err != nil {
				fmt.Printf("❌ Failed to add public key for '%s': %v\n", account, err)
				os.Exit(1)
			}
		}

		fmt.Printf("✅ Backup created at %s\n", outputPath)
	},
}

func addFileToArchive(tarWriter *tar.Writer, sourcePath, archivePath string) error {
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}
	header.Name = archivePath

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}

func writeArchiveContent(tarWriter *tar.Writer, archivePath string, content []byte, mode int64) error {
	header := &tar.Header{
		Name: archivePath,
		Mode: mode,
		Size: int64(len(content)),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	_, err := tarWriter.Write(content)
	return err
}

func init() {
	rootCmd.AddCommand(backupCmd)
}

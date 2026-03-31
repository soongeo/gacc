package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/viper"
)

func gaccConfigFilePath() (string, error) {
	if used := viper.ConfigFileUsed(); used != "" {
		return used, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "gacc", "config.yaml"), nil
}

func detectActiveAccount() string {
	accounts, err := ssh.ListAccounts()
	if err != nil || !git.IsInsideWorkTree() {
		return ""
	}

	remoteURL, err := git.GetRemoteURL("origin")
	if err != nil {
		return ""
	}

	for _, account := range accounts {
		targetHost := fmt.Sprintf("github.com-%s", account)
		if strings.Contains(remoteURL, targetHost+":") || strings.Contains(remoteURL, targetHost+"/") {
			return account
		}
	}

	return ""
}

func chooseAccountInteractively(accounts []string) (string, error) {
	if len(accounts) == 0 {
		return "", fmt.Errorf("no accounts available")
	}

	fmt.Println("Select an account:")
	for i, account := range accounts {
		fmt.Printf("  %d. %s\n", i+1, account)
	}
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(accounts) {
		return "", fmt.Errorf("invalid selection")
	}

	return accounts[index-1], nil
}

func resolveAccountForOptionalArg(args []string) (string, error) {
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		return args[0], nil
	}

	activeAccount := detectActiveAccount()
	if activeAccount != "" {
		return activeAccount, nil
	}

	accounts, err := ssh.ListAccounts()
	if err != nil {
		return "", err
	}

	return chooseAccountInteractively(accounts)
}

func defaultBackupArchiveName() string {
	return fmt.Sprintf("gacc-backup-%s.tar.gz", time.Now().Format("20060102150405"))
}

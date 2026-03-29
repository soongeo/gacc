package ssh

import (
	"os"
	"path/filepath"
	"strings"
)

// ListAccounts returns a list of configured gacc account names present in ~/.ssh/config.
func ListAccounts() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(homeDir, ".ssh", "config")

	content, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var accounts []string
	lines := strings.Split(string(content), "\n")
	prefix := "Host github.com-"

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			accountName := strings.TrimPrefix(trimmed, prefix)
			if accountName != "" {
				accounts = append(accounts, accountName)
			}
		}
	}

	return accounts, nil
}

package ssh

import (
	"os"
	"strings"
)

// ListAccounts returns a list of configured gacc account names present in ~/.ssh/config.
func ListAccounts() ([]string, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

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

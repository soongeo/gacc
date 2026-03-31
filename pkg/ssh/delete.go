package ssh

import (
	"fmt"
	"os"
	"strings"
)

// DeleteSSHKeys 로컬의 SSH 키 파일들을 삭제합니다.
func DeleteSSHKeys(accountName string) error {
	privPath, err := PrivateKeyPath(accountName)
	if err != nil {
		return err
	}
	pubPath := privPath + ".pub"

	// Delete private key
	if err := os.Remove(privPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete private key: %w", err)
	}

	// Delete public key
	if err := os.Remove(pubPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete public key: %w", err)
	}

	return nil
}

// RemoveSSHConfig ~/.ssh/config 에서 해당 계정의 Host 블록을 제거합니다.
func RemoveSSHConfig(accountName string) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	skipMode := false
	hostAlias := fmt.Sprintf("Host github.com-%s", accountName)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if skipMode {
			if strings.HasPrefix(trimmed, "Host ") || strings.HasPrefix(trimmed, "Match ") {
				skipMode = false
			} else {
				continue
			}
		}

		if trimmed == hostAlias {
			skipMode = true
			continue
		}

		if !skipMode {
			newLines = append(newLines, line)
		}
	}

	// 연속된 빈 줄 정리 (선택 사항)
	var cleaned []string
	for i, line := range newLines {
		if strings.TrimSpace(line) == "" {
			// 연속된 빈 줄이거나 첫 줄이 빈 경우면 제외
			if i == 0 || (i > 0 && strings.TrimSpace(newLines[i-1]) == "") {
				continue
			}
		}
		cleaned = append(cleaned, line)
	}

	newContent := strings.Join(cleaned, "\n")
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return err
	}

	return nil
}

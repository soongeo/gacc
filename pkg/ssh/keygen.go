package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

func GenerateAndSaveEd25519(accountName string) (string, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}

	// PEM 인코딩 (Private)
	privateKeyBlock, err := ssh.MarshalPrivateKey(priv, "")
	if err != nil {
		return "", fmt.Errorf("failed to encode private key: %w", err)
	}
	privatePEM := pem.EncodeToMemory(privateKeyBlock)

	// OpenSSH 인코딩 (Public)
	sshPubKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("failed to convert public key: %w", err)
	}
	publicBytes := ssh.MarshalAuthorizedKey(sshPubKey)
	publicStr := strings.TrimSpace(string(publicBytes)) + " gacc-" + accountName

	// ~/.ssh 디렉토리 경로
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	sshDir := filepath.Join(homeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return "", err
	}

	privPath := filepath.Join(sshDir, fmt.Sprintf("gacc_%s", accountName))
	pubPath := privPath + ".pub"

	// 파일 쓰기
	if err := os.WriteFile(privPath, privatePEM, 0600); err != nil {
		return "", fmt.Errorf("failed to save private key: %w", err)
	}
	if err := os.WriteFile(pubPath, []byte(publicStr+"\n"), 0644); err != nil {
		return "", fmt.Errorf("failed to save public key: %w", err)
	}

	return publicStr, nil
}

// UpdateSSHConfig ~/.ssh/config 파일에 GitHub 호스트를 추가합니다.
func UpdateSSHConfig(accountName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	sshDir := filepath.Join(homeDir, ".ssh")
	configPath := filepath.Join(sshDir, "config")

	hostAlias := fmt.Sprintf("github.com-%s", accountName)
	privKeyPath := filepath.Join(sshDir, fmt.Sprintf("gacc_%s", accountName))

	configBlock := fmt.Sprintf(`
Host %s
    HostName github.com
    User git
    IdentityFile %s
    IdentitiesOnly yes
`, hostAlias, privKeyPath)

	// 이미 존재하는지 확인
	content, err := os.ReadFile(configPath)
	if err == nil {
		if strings.Contains(string(content), hostAlias) {
			// 이미 존재함
			return nil
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// append로 쓰기
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(configBlock); err != nil {
		return err
	}

	return nil
}

package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SSHDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".ssh"), nil
}

func ConfigPath() (string, error) {
	sshDir, err := SSHDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(sshDir, "config"), nil
}

func PrivateKeyPath(accountName string) (string, error) {
	sshDir, err := SSHDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(sshDir, fmt.Sprintf("gacc_%s", accountName)), nil
}

func PublicKeyPath(accountName string) (string, error) {
	privateKeyPath, err := PrivateKeyPath(accountName)
	if err != nil {
		return "", err
	}
	return privateKeyPath + ".pub", nil
}

func ReadPublicKey(accountName string) (string, error) {
	publicKeyPath, err := PublicKeyPath(accountName)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func RenameAccount(oldName, newName string) error {
	oldPrivatePath, err := PrivateKeyPath(oldName)
	if err != nil {
		return err
	}
	oldPublicPath, err := PublicKeyPath(oldName)
	if err != nil {
		return err
	}
	newPrivatePath, err := PrivateKeyPath(newName)
	if err != nil {
		return err
	}
	newPublicPath, err := PublicKeyPath(newName)
	if err != nil {
		return err
	}

	if _, err := os.Stat(newPrivatePath); err == nil {
		return fmt.Errorf("account '%s' already exists", newName)
	}

	if err := os.Rename(oldPrivatePath, newPrivatePath); err != nil {
		return fmt.Errorf("failed to rename private key: %w", err)
	}
	if err := os.Rename(oldPublicPath, newPublicPath); err != nil {
		_ = os.Rename(newPrivatePath, oldPrivatePath)
		return fmt.Errorf("failed to rename public key: %w", err)
	}

	publicKey, err := os.ReadFile(newPublicPath)
	if err == nil {
		updatedKey := rewritePublicKeyComment(string(publicKey), "gacc-"+newName)
		if writeErr := os.WriteFile(newPublicPath, []byte(updatedKey), 0644); writeErr != nil {
			return fmt.Errorf("failed to update public key comment: %w", writeErr)
		}
	}

	if err := RemoveSSHConfig(oldName); err != nil {
		return err
	}
	if err := UpdateSSHConfig(newName); err != nil {
		return err
	}

	return nil
}

func rewritePublicKeyComment(publicKey, comment string) string {
	fields := strings.Fields(strings.TrimSpace(publicKey))
	if len(fields) < 2 {
		return strings.TrimSpace(publicKey) + "\n"
	}
	return fmt.Sprintf("%s %s %s\n", fields[0], fields[1], comment)
}

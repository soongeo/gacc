package git

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsInsideWorkTree checks if the current directory is inside a Git repository.
func IsInsideWorkTree() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// GetRemoteURL returns the URL of the specified remote branch (e.g. "origin").
func GetRemoteURL(remote string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", remote)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("could not get remote URL for '%s': %w", remote, err)
	}
	return strings.TrimSpace(out.String()), nil
}

func GetConfig(scope, key string) (string, error) {
	args := []string{"config"}
	if scope != "" {
		args = append(args, scope)
	}
	args = append(args, "--get", key)

	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", fmt.Errorf("failed to get git config %s: %w", key, err)
	}
	return strings.TrimSpace(out.String()), nil
}

func CurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Clean(dir), nil
}

func WorkTreeRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to detect work tree root: %w", err)
	}
	return filepath.Clean(strings.TrimSpace(out.String())), nil
}

// SetRemoteURL sets a new URL for the specified remote branch.
func SetRemoteURL(remote, url string) error {
	cmd := exec.Command("git", "remote", "set-url", remote, url)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not set remote URL: %w", err)
	}
	return nil
}

// SetLocalUserConfig configures the user.name and user.email for the local repository.
func SetLocalUserConfig(name, email string) error {
	if name != "" {
		if err := exec.Command("git", "config", "--local", "user.name", name).Run(); err != nil {
			return fmt.Errorf("failed to set local user.name: %w", err)
		}
	}
	if email != "" {
		if err := exec.Command("git", "config", "--local", "user.email", email).Run(); err != nil {
			return fmt.Errorf("failed to set local user.email: %w", err)
		}
	}
	return nil
}

// SetGlobalUserConfig configures the user.name and user.email globally.
func SetGlobalUserConfig(name, email string) error {
	if name != "" {
		if err := exec.Command("git", "config", "--global", "user.name", name).Run(); err != nil {
			return fmt.Errorf("failed to set global user.name: %w", err)
		}
	}
	if email != "" {
		if err := exec.Command("git", "config", "--global", "user.email", email).Run(); err != nil {
			return fmt.Errorf("failed to set global user.email: %w", err)
		}
	}
	return nil
}

// UnsetLocalUserConfig unsets the local user.name and user.email overrides.
func UnsetLocalUserConfig() error {
	_ = exec.Command("git", "config", "--local", "--unset", "user.name").Run()
	_ = exec.Command("git", "config", "--local", "--unset", "user.email").Run()
	return nil
}

// UnsetGlobalUserConfig unsets the global user.name and user.email values.
func UnsetGlobalUserConfig() error {
	_ = exec.Command("git", "config", "--global", "--unset", "user.name").Run()
	_ = exec.Command("git", "config", "--global", "--unset", "user.email").Run()
	return nil
}

// SetGlobalSSHCommand configures the global core.sshCommand value.
func SetGlobalSSHCommand(command string) error {
	if err := exec.Command("git", "config", "--global", "core.sshCommand", command).Run(); err != nil {
		return fmt.Errorf("failed to set global core.sshCommand: %w", err)
	}
	return nil
}

// UnsetGlobalSSHCommand unsets the global core.sshCommand value.
func UnsetGlobalSSHCommand() error {
	_ = exec.Command("git", "config", "--global", "--unset", "core.sshCommand").Run()
	return nil
}

// SetLocalSSHCommand configures the local (per-repo) core.sshCommand value.
func SetLocalSSHCommand(command string) error {
	if err := exec.Command("git", "config", "--local", "core.sshCommand", command).Run(); err != nil {
		return fmt.Errorf("failed to set local core.sshCommand: %w", err)
	}
	return nil
}

// UnsetLocalSSHCommand unsets the local (per-repo) core.sshCommand value.
func UnsetLocalSSHCommand() error {
	_ = exec.Command("git", "config", "--local", "--unset", "core.sshCommand").Run()
	return nil
}

// ParseAndReplaceRemoteHost takes an existing remote (e.g., git@github.com:usr/repo.git or https://github.com/usr/repo.git)
// and returns the modified string that uses the target gacc host alias (e.g., git@github.com-alias:usr/repo.git).
func ParseAndReplaceRemoteHost(url string, accountAlias string) (string, error) {
	targetHost := fmt.Sprintf("github.com-%s", accountAlias)

	// If it is already targeting the exact host, do nothing.
	if strings.Contains(url, targetHost+":") || strings.Contains(url, targetHost+"/") {
		return url, nil
	}

	// Case 1: Standard SSH git@github.com:user/repo.git
	if strings.HasPrefix(url, "git@github.com:") {
		return strings.Replace(url, "git@github.com:", "git@"+targetHost+":", 1), nil
	}

	// Case 2: HTTPS https://github.com/user/repo.git -> Convert to SSH alias
	if strings.HasPrefix(url, "https://github.com/") {
		replaced := strings.TrimPrefix(url, "https://github.com/")
		return "git@" + targetHost + ":" + replaced, nil
	}

	// Case 3: It was using a different gacc alias, e.g., git@github.com-other:user/repo.git
	if strings.HasPrefix(url, "git@github.com-") {
		parts := strings.SplitN(url, ":", 2)
		if len(parts) == 2 {
			return "git@" + targetHost + ":" + parts[1], nil
		}
	}

	return "", errors.New("unsupported remote URL format (only standard GitHub URLs are supported)")
}

// RevertRemoteHostToStandard takes an aliased remote (e.g., git@github.com-alias:usr/repo.git)
// and returns standard SSH remote (e.g., git@github.com:usr/repo.git)
func RevertRemoteHostToStandard(url string) (string, error) {
	if strings.HasPrefix(url, "git@github.com-") {
		parts := strings.SplitN(url, ":", 2)
		if len(parts) == 2 {
			return "git@github.com:" + parts[1], nil
		}
	}
	return url, nil
}

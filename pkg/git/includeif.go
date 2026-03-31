package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func NormalizeGitDirCondition(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	cleaned := filepath.Clean(absDir)
	if !strings.HasSuffix(cleaned, string(filepath.Separator)) {
		cleaned += string(filepath.Separator)
	}

	return filepath.ToSlash(cleaned), nil
}

func AddGlobalIncludeIf(gitDirCondition, includePath string) error {
	key := fmt.Sprintf("includeIf.gitdir:%s.path", gitDirCondition)
	cmd := exec.Command("git", "config", "--global", "--replace-all", key, includePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add includeIf rule: %w", err)
	}
	return nil
}

func RemoveGlobalIncludeIf(gitDirCondition, includePath string) error {
	key := fmt.Sprintf("includeIf.gitdir:%s.path", gitDirCondition)
	cmd := exec.Command("git", "config", "--global", "--unset-all", key, includePath)
	if err := cmd.Run(); err != nil {
		unsetCmd := exec.Command("git", "config", "--global", "--unset-all", key)
		if unsetErr := unsetCmd.Run(); unsetErr != nil {
			return fmt.Errorf("failed to remove includeIf rule: %w", err)
		}
	}
	return nil
}

type GlobalIncludeRule struct {
	Condition string
	Path      string
}

func ListGlobalIncludeIf() ([]GlobalIncludeRule, error) {
	cmd := exec.Command("git", "config", "--global", "--get-regexp", "^includeIf\\..*\\.path$")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return []GlobalIncludeRule{}, nil
		}
		return nil, fmt.Errorf("failed to list includeIf rules: %w", err)
	}

	pattern := regexp.MustCompile(`^includeIf\.gitdir:(.+)\.path\s+(.+)$`)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	rules := make([]GlobalIncludeRule, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		match := pattern.FindStringSubmatch(line)
		if len(match) != 3 {
			continue
		}

		rules = append(rules, GlobalIncludeRule{
			Condition: match[1],
			Path:      match[2],
		})
	}

	return rules, nil
}

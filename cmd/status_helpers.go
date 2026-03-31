package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/soongeo/gacc/pkg/ssh"
)

type resolvedStatus struct {
	CWD             string
	InGitRepo       bool
	RepoRoot        string
	ManualAccount   string
	ManualRemoteURL string
	LocalUserName   string
	LocalUserEmail  string
	LocalSSHCommand string
	AutoRule        *autoRule
	GlobalUserName  string
	GlobalUserEmail string
	GlobalSSHCmd    string
}

func collectResolvedStatus() (resolvedStatus, error) {
	cwd, err := git.CurrentDir()
	if err != nil {
		return resolvedStatus{}, err
	}

	status := resolvedStatus{CWD: cwd}

	status.GlobalUserName, _ = git.GetConfig("--global", "user.name")
	status.GlobalUserEmail, _ = git.GetConfig("--global", "user.email")
	status.GlobalSSHCmd, _ = git.GetConfig("--global", "core.sshCommand")

	autoRule := detectAutoRuleForDir(cwd)
	if autoRule != nil {
		status.AutoRule = autoRule
	}

	if !git.IsInsideWorkTree() {
		return status, nil
	}

	status.InGitRepo = true
	status.RepoRoot, _ = git.WorkTreeRoot()
	status.LocalUserName, _ = git.GetConfig("--local", "user.name")
	status.LocalUserEmail, _ = git.GetConfig("--local", "user.email")
	status.LocalSSHCommand, _ = git.GetConfig("--local", "core.sshCommand")
	status.ManualAccount = detectActiveAccount()
	if status.ManualAccount != "" {
		status.ManualRemoteURL, _ = git.GetRemoteURL("origin")
	}

	return status, nil
}

func detectAutoRuleForDir(dir string) *autoRule {
	normalizedDir := filepath.ToSlash(filepath.Clean(dir)) + "/"
	rules := listStoredAutoRules()
	if len(rules) == 0 {
		return nil
	}

	sort.Slice(rules, func(i, j int) bool {
		return len(rules[i].Condition) > len(rules[j].Condition)
	})

	for _, rule := range rules {
		condition := rule.Condition
		if condition == "" {
			continue
		}
		if strings.HasPrefix(normalizedDir, condition) {
			matched := rule
			return &matched
		}
	}

	return nil
}

func accountLabelsForList(account string, status resolvedStatus) []string {
	var labels []string

	if account == status.ManualAccount {
		labels = append(labels, "local")
	}
	if status.AutoRule != nil && account == status.AutoRule.Account {
		labels = append(labels, "auto")
	}
	if globalAccountFromSSHCommand(status.GlobalSSHCmd) == account {
		labels = append(labels, "global")
	}

	return labels
}

func globalAccountFromSSHCommand(command string) string {
	for _, account := range mustListAccounts() {
		privateKeyPath, err := ssh.PrivateKeyPath(account)
		if err != nil {
			continue
		}
		if strings.Contains(command, privateKeyPath) {
			return account
		}
	}
	return ""
}

func mustListAccounts() []string {
	accounts, err := ssh.ListAccounts()
	if err != nil {
		return nil
	}
	return accounts
}

func formatIdentity(name, email string) string {
	if name == "" && email == "" {
		return "(not set)"
	}
	if name == "" {
		return email
	}
	if email == "" {
		return name
	}
	return fmt.Sprintf("%s <%s>", name, email)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/soongeo/gacc/pkg/git"
	"github.com/soongeo/gacc/pkg/ssh"
	"github.com/spf13/viper"
)

type autoRule struct {
	Account   string
	Directory string
	Condition string
	Include   string
}

func autoRuleKey(account, directory string) string {
	slug := strings.ToLower(account + "_" + directory)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "_")
	slug = strings.Trim(slug, "_")
	if slug == "" {
		return "default"
	}
	return slug
}

func includeDirPath() (string, error) {
	configPath, err := gaccConfigFilePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(configPath), "includeif"), nil
}

func includeFilePath(account, condition string) (string, error) {
	dir, err := includeDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, autoRuleKey(account, condition)+".gitconfig"), nil
}

func loadAutoRules() map[string]any {
	rules := viper.GetStringMap("auto_rules")
	if rules == nil {
		return map[string]any{}
	}
	return rules
}

func saveAutoRule(rule autoRule) error {
	rules := loadAutoRules()
	key := autoRuleKey(rule.Account, rule.Condition)
	rules[key] = map[string]string{
		"account":   rule.Account,
		"directory": rule.Directory,
		"condition": rule.Condition,
		"include":   rule.Include,
	}
	viper.Set("auto_rules", rules)
	return viper.WriteConfig()
}

func deleteAutoRule(account, condition string) error {
	rules := loadAutoRules()
	delete(rules, autoRuleKey(account, condition))
	viper.Set("auto_rules", rules)
	return viper.WriteConfig()
}

func listStoredAutoRules() []autoRule {
	rawRules := loadAutoRules()
	rules := make([]autoRule, 0, len(rawRules))

	for _, raw := range rawRules {
		entry, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		rules = append(rules, autoRule{
			Account:   fmt.Sprint(entry["account"]),
			Directory: fmt.Sprint(entry["directory"]),
			Condition: fmt.Sprint(entry["condition"]),
			Include:   fmt.Sprint(entry["include"]),
		})
	}

	return rules
}

func findStoredAutoRule(account, directory string) (autoRule, bool, error) {
	condition, err := git.NormalizeGitDirCondition(directory)
	if err != nil {
		return autoRule{}, false, err
	}

	for _, rule := range listStoredAutoRules() {
		if rule.Account == account && rule.Condition == condition {
			return rule, true, nil
		}
	}

	return autoRule{}, false, nil
}

func writeIncludeFile(account, condition string) (string, error) {
	name := strings.TrimSpace(viper.GetString("accounts." + account + ".name"))
	email := strings.TrimSpace(viper.GetString("accounts." + account + ".email"))
	privateKeyPath, err := ssh.PrivateKeyPath(account)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if name != "" || email != "" {
		builder.WriteString("[user]\n")
		if name != "" {
			builder.WriteString(fmt.Sprintf("\tname = %s\n", name))
		}
		if email != "" {
			builder.WriteString(fmt.Sprintf("\temail = %s\n", email))
		}
	}
	builder.WriteString("[core]\n")
	builder.WriteString(fmt.Sprintf("\tsshCommand = ssh -i \"%s\" -o IdentitiesOnly=yes\n", privateKeyPath))

	includePath, err := includeFilePath(account, condition)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(includePath), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(includePath, []byte(builder.String()), 0644); err != nil {
		return "", err
	}

	return includePath, nil
}

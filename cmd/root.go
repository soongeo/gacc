package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gacc",
	Short: "gacc is a Git Account & SSH Key Manager",
	// Removed CompletionOptions.HiddenDefaultCmd to expose the completion command
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Welcome to gacc! Try 'gacc --help'.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// 사용하는 OS의 기본 설정 디렉토리를 가져옵니다. (ex: ~/.config)
		configDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		// gacc 전용 설정 폴더 (ex: ~/.config/gacc)
		gaccConfigPath := filepath.Join(configDir, "gacc")
		if _, err := os.Stat(gaccConfigPath); os.IsNotExist(err) {
			_ = os.MkdirAll(gaccConfigPath, 0755)
		}

		viper.AddConfigPath(gaccConfigPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		// 기본 파일이 없으면 생성
		defaultConfig := filepath.Join(gaccConfigPath, "config.yaml")
		if _, err := os.Stat(defaultConfig); os.IsNotExist(err) {
			_ = os.WriteFile(defaultConfig, []byte(""), 0644)
		}
	}

	viper.AutomaticEnv() // 환경 변수 자동 로드

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("설정 파일 로드 완료:", viper.ConfigFileUsed())
	}
}

func accountNameCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	initConfig() // ensure config is loaded
	accounts := viper.GetStringMap("accounts")
	var names []string
	for name := range accounts {
		names = append(names, name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

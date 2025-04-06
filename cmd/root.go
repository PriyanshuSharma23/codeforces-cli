/*
Copyright © 2025 Priyanshu Sharma inbox.priyanshu@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var logger = log.New(os.Stdout, "", log.LstdFlags)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codeforces-cli",
	Short: "A CLI tool to streamline competitive programming on Codeforces",
	Long: `codeforces-cli is a powerful command-line application designed to
enhance your competitive programming experience on Codeforces.

This tool helps you efficiently manage problems by integrating with the Competitive
Companion browser extension, automatically preparing your working environment with
boilerplate code, and setting up test cases. Moreover, it allows easy compilation
and execution of your solutions against sample test cases, providing quick feedback
on their correctness.

Key Features:

- Automatic problem directory creation with predefined structure
- Integration with Competitive Companion for seamless problem import
- Utilization of customizable templates to kickstart your solutions
- Execution and testing of solutions with straightforward commands
- Supports flexible configuration for different languages and editors

With codeforces-cli, you can focus more on solving problems rather than setting up
your development environment.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.codeforces-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the --config flag
		viper.SetConfigFile(cfgFile)
		fmt.Println("[config] Using config file from flag:", cfgFile)
	} else {
		// Get OS-specific config dir
		configDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		fmt.Println("[config] Config dir:", configDir)

		// Default: ~/.config/codeforces-cli/config.yaml (or %AppData%\codeforces-cli\config.yaml)
		defaultPath := filepath.Join(configDir, "codeforces-cli", "config.yaml")

		if _, err := os.Stat(defaultPath); err == nil {
			viper.SetConfigFile(defaultPath)
			fmt.Println("[config] Found primary config at:", defaultPath)
		} else {
			// Fallback: ~/codeforces-cli.yaml
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			fallback := filepath.Join(home, "codeforces-cli.yaml")
			fmt.Println("[config] Trying fallback config at:", fallback)

			viper.SetConfigFile(fallback)
		}
	}

	// Specify the format and allow env overrides
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Try to read the selected config
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "✅ Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "⚠️ No config file found or failed to read:", err)
	}

	// Set default values
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	defaultProblemPath := filepath.Join(home, "codeforces", "problems")
	defaultTemplatePath := filepath.Join(home, "codeforces", "templates", "main.cpp")

	viper.SetDefault("root", defaultProblemPath)
	viper.SetDefault("language", "py")
	viper.SetDefault("programFile", "main")
	viper.SetDefault("buildCommand", "")
	viper.SetDefault("executeCommand", "python3 {{.Path}}")
	viper.SetDefault("testCaseInputPrefix", "input")
	viper.SetDefault("testCaseOutputPrefix", "output")
	viper.SetDefault("port", 10045)
	viper.SetDefault("editorCommand", "nvim {{.Path}}")
	viper.SetDefault("templatePath", defaultTemplatePath)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/git-mask/git-mask/core"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-mask",
	Short: "Manage and apply identity profiles for Git",
	Long: `Define and use different authoring names, emails, GPG signing keys,
SSH authorization keys, faking time zones as git-mask profiles. Matching Git
repo URL with custom rules to apply different profiles automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := core.LookPath("git")
		cobra.CheckErr(err)
		fmt.Print(paths)
		repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
			DetectDotGit:          true,
			EnableDotGitCommonDir: true,
		})
		cobra.CheckErr(err)
		remote, err := repo.Remote("origin")
		cobra.CheckErr(err)
		fmt.Print(remote.Config().URLs)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-mask.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".git-mask" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".git-mask")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

type GitMaskMode int

const (
	// git-mask is invoked by real git
	PluginMode GitMaskMode = iota
	// git-mask is replacing the real git entrypoint
	WrapperMode
)

func exeName(path string) string {
	name := filepath.Base(path)
	if ext := filepath.Ext(name); ext != "" {
		name = name[:len(name)-len(ext)]
	}
	return strings.ToLower(name)
}

func DetectGitMaskMode() (mode GitMaskMode, err error) {
	self, err := os.Executable()
	if err != nil {
		return
	}
	self = exeName(self)
	if self == "git" {
		return WrapperMode, nil
	}
	return PluginMode, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	mode, err := DetectGitMaskMode()
	cobra.CheckErr(err)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

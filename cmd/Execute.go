package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// Execute executes the root command once Initialize has been called
func Execute() error {
	defaultCfgFile = fmt.Sprintf("%s.%s", PROGRAM, cfgExt)

	rootCmd.Use = PROGRAM
	rootCmd.Version = fmt.Sprintf("%s, version %s-(%s-%s)", PROGRAM, VERSION, runtime.GOARCH, runtime.GOOS)

	// ? should not require RAW mode just go get the dimensions...
	if w, _, err := terminal.GetSize(0); err == nil {
		termWidth = w
	}

	// Load Configuration on Initialize
	defaultConfig()
	cobra.OnInitialize(loadConfig)

	// Attach Flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s)", defaultCfgFile))
	rootCmd.PersistentFlags().BoolVar(&textMode, "text", false, "force text (console) user interface")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable debug logging for troubleshooting sake")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Add Commands
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}

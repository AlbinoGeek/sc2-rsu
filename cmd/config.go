package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/kataras/golog"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgExt         = "yaml"
	cfgFile        string
	defaultCfgFile string
)

func loadConfig() {
	if cfgFile != "" {
		// load configuration from the path specified via Flags
		viper.SetConfigFile(cfgFile)
	} else {
		// otherwise, search in the following paths:
		// 1) User's Home Directory
		// 2) Shell's Working Directory
		// 3) System Configuration Directory (on Linux)
		// 4) Executable's Parent Directory
		if home, err := homedir.Dir(); err == nil {
			viper.AddConfigPath(home)
		}

		if wd, err := os.Getwd(); err == nil {
			viper.AddConfigPath(wd)
		}

		if runtime.GOOS == "linux" {
			viper.AddConfigPath("/etc")
		}

		if ed, err := os.Executable(); err == nil {
			viper.AddConfigPath(filepath.Dir(ed))
		}

		viper.SetConfigName(defaultCfgFile)
	}

	viper.SetConfigType("yaml")

	// Also read configuration from environment variables
	viper.AutomaticEnv()

	// Finally, load the config
	if err := viper.ReadInConfig(); err != nil {
		if _, existErr := err.(viper.ConfigFileNotFoundError); !existErr {
			golog.Warnf("failed loading configuration: %v", err)
		}
	} else {
		golog.Debugf("using configuration: %v", viper.ConfigFileUsed())
	}

	if viper.GetBool("verbose") {
		golog.SetLevel("debug")
	}
}

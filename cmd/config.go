package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kataras/golog"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/utils"
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

func saveConfig() error {
	if cfgFile == "" {
		cfgFile = viper.ConfigFileUsed()
	}
	if cfgFile == "" {
		cfgFile = defaultCfgFile
	}

	if err := viper.WriteConfigAs(cfgFile); err != nil {
		return fmt.Errorf("unable to save configuration: %v", err)
	}

	golog.Debugf("Wrote Configuration: %v", cfgFile)
	return nil
}

func setAPIkey(key string) error {
	if !utils.ValidAPIKey(key) {
		return errors.New("invalid API key format")
	}

	if viper.GetString("apikey") == key {
		golog.Info("API key already in configuration! Doing nothing.")
		return nil
	}

	viper.Set("apikey", key)
	if err := saveConfig(); err != nil {
		return err
	}

	golog.Info("API Key set in configuration!")
	return nil
}

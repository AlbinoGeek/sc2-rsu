package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
)

var (
	cfgExt         = "yaml"
	cfgFile        string
	defaultCfgFile string
	defaults       = map[string]interface{}{
		"theme.iconInlineSize":     20, // 20
		"theme.padding":            4,
		"theme.scrollBarSize":      12, // 16
		"theme.scrollBarSmallSize": 3,
		"theme.textSize":           16,
		"update.automatic.enabled": false,
		"update.check.enabled":     true,
		"update.check.period":      time.Duration(minimumUpdatePeriod).String(),
	}
)

func defaultConfig() {
	for key, val := range defaults {
		viper.SetDefault(key, val)
	}
}

func loadConfig() {
	if cfgFile != "" {
		// load configuration from the path specified via Flags
		viper.SetConfigFile(cfgFile)
	} else {
		// otherwise, search in the following paths:
		// 1) User's Home Directory
		viper.AddConfigPath("$HOME/")

		// 2) Shell's Working Directory
		viper.AddConfigPath(".")

		// 3) System Configuration Directory
		if runtime.GOOS == "linux" {
			viper.AddConfigPath("/etc")
		}

		// 4) Executable's Parent Directory
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
	if !sc2replaystats.ValidAPIKey(key) {
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

func getToonEnabled(toon string) bool {
	golog.Debugf("getToonEnabled(%s)", toon)

	confToons := viper.GetStringSlice("toons")

	// ! if there are none in config, enable all by default
	if len(confToons) == 0 {
		return true
	}

	for _, t := range confToons {
		if t == toon {
			return true
		}
	}

	return false
}

func setToons(toons []string) error {
	in := func(needle string, haystack []string) bool {
		for _, n := range haystack {
			if n == needle {
				return true
			}
		}

		return false
	}

	var changed bool
	confToons := viper.GetStringSlice("toons")
	for _, t := range toons {
		if !in(t, confToons) {
			golog.Debug("setToons: added ", t)
			confToons = append(confToons, t)
			changed = true
		}
	}

	for i := len(confToons) - 1; i >= 0; i-- {
		if !in(confToons[i], toons) {
			golog.Debug("setToons: removed ", confToons[i])
			confToons = append(confToons[:i], confToons[i+1:]...)
			changed = true
		}
	}

	// don't save configuration if there were no changes made
	if !changed {
		return nil
	}

	viper.Set("toons", confToons)

	return saveConfig()
}

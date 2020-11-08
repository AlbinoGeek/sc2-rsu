package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/bgentry/speakeasy"
	"github.com/kataras/golog"
	"github.com/mxschmitt/playwright-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use: "login <apikey or email>",
	Args: func(cmd *cobra.Command, args []string) error {
		if l := len(args); l != 1 {
			return fmt.Errorf("wrong argument count: expected 1, got %d", l)
		}

		// is it an API key?
		if strings.Count(args[0], ";") == 2 {
			return nil
		}

		// is it an email address?
		if err := checkmail.ValidateFormat(args[0]); err != nil {
			return fmt.Errorf("email address: %v", err)
		}

		return nil
	},
	Short: "Add an sc2replaystats account to the config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		// is it an API key?
		if strings.Count(args[0], ";") == 2 {
			return addAPIKey(args[0])
		}

		// is it an email address?
		t := time.Now()
		if err := login(args[0]); err != nil {
			return err
		}
		golog.Debugf("Login completed in %s", time.Since(t))

		return nil
	},
}

var (
	sc2rsBase    = "https://sc2replaystats.com"
	sc2rsAPIBase = "https://api.sc2replaystats.com"
)

func login(email string) error {
	fmt.Printf(`
	============================================================
	We are about to login to sc2replaystats for you to obtain or
	generate your API key. We will have to ask you for your pass
	word, which we WILL NOT SAVE -- once the login form has been
	loaded. If you want to avoid providing your account password
	please call this command with your API key instead. Example:
	%s login <apikey>
	============================================================

`, os.Args[0])

	golog.Debugf("Setting up browser...")
	pw, err := playwright.Run()
	if err != nil {
		golog.Fatalf("failed to setup login browser: %v", err)
	}
	defer pw.Stop()

	golog.Debugf("Launching browser...")
	browser, err := pw.Chromium.Launch()
	if err != nil {
		golog.Fatalf("failed to launch login browser: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return err
	}
	defer page.Close()

	golog.Infof("Navigating to login page...")
	if _, err = page.Goto(fmt.Sprintf("%s/Account/signin", sc2rsBase)); err != nil {
		return err
	}

	golog.Debugf("Filling login form...")
	if input, err := page.QuerySelector("css=input[name='email']"); err == nil && input != nil {
		if err = input.Fill(email); err != nil {
			return err
		}
	}

	tt := time.Now()
	password, err := speakeasy.Ask(fmt.Sprintf("Password for sc2ReplayStats account %s: ", email))
	golog.Debugf("User input took %s", time.Since(tt))
	if err != nil {
		return err
	}
	if input, err := page.QuerySelector("css=input[name='password']"); err == nil && input != nil {
		if err = input.Fill(password); err != nil {
			return err
		}
	}

	golog.Debugf("Submitting login form...")
	if btn, err := page.QuerySelector("css=input[value='Sign In']"); err == nil && btn != nil {
		if err = btn.Click(); err != nil {
			return err
		}
	}

	url := page.URL()
	if !strings.Contains(url, "display") {
		return fmt.Errorf("unexpected URL, login probably failed: %v", url)
	}

	parts := strings.Split(url, "/")
	accid := strings.Split(parts[len(parts)-1], "#")[0]
	golog.Infof("Success! Logged in to account #%v", accid)

	if _, err = page.Goto(fmt.Sprintf("%s/account/settings/%v", sc2rsBase, accid)); err != nil {
		return err
	}

	golog.Debugf("Waiting for settings page to load...")
	e, err := page.WaitForSelector("*css=a[data-toggle='tab'] >> text=API Access")
	if err != nil {
		return err
	}
	golog.Debugf("Clicking 'API Access'...")
	if err = e.Click(); err != nil {
		return err
	}
	golog.Debugf("Finding API key...")
	e, err = page.QuerySelector("*css=.form-group >> text=Authorization Key")
	if e == nil || err != nil {
		golog.Infof("Generating new API key...")
		e, err = page.QuerySelector("text=Generate New API Key")
		if err != nil {
			return err
		}
		if err = e.Click(); err != nil {
			return err
		}
		_, err = page.WaitForSelector("*css=a[data-toggle='tab'] >> text=API Access")
		if err != nil {
			return err
		}
		e, err = page.QuerySelector("*css=.form-group >> text=Authorization Key")
	}
	if err != nil || e == nil {
		return fmt.Errorf("unable to find API key: %v", err)
	}

	t, err := e.InnerText()
	if err != nil {
		return err
	}

	key := strings.Trim(strings.Split(t, ": ")[1], " \r\n\t")
	return addAPIKey(key)
}

func addAPIKey(key string) error {
	if len(key) < 80 || len(key) > 100 {
		return fmt.Errorf("API key format invalid")
	}

	keys := viper.GetStringSlice("apikeys")
	for _, k := range keys {
		if k == key {
			golog.Infof("API key already in configuration! Doing nothing.")
			return nil
		}
	}

	keys = append(keys, key)
	viper.Set("apikeys", keys)
	if err := saveConfig(); err != nil {
		return err
	}

	golog.Info("API Key added to configuration!")
	return nil
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
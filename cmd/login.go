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

	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
)

var (
	loginWarning = `
============================================================
We are about to login to sc2replaystats for you to obtain or
generate your API key. We will have to ask you for your pass
word, which we WILL NOT SAVE -- once the login form has been
loaded. If you want to avoid providing your account password
please call this command with your API key instead. Example:
%s login <apikey>
============================================================

`

	loginCmd = &cobra.Command{
		Use: "login <apikey or email>",
		Args: func(cmd *cobra.Command, args []string) error {
			if l := len(args); l != 1 {
				return fmt.Errorf("wrong argument count: expected 1, got %d", l)
			}

			// is it an API key?
			if sc2replaystats.ValidAPIKey(args[0]) {
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
			if sc2replaystats.ValidAPIKey(args[0]) {
				return setAPIkey(args[0])
			}

			// is it an email address?
			t := time.Now()
			if err := login(args[0]); err != nil {
				return fmt.Errorf("email login error: %v", err)
			}
			golog.Debugf("Login completed in %s", time.Since(t))

			return nil
		},
	}
)

func login(email string) error {
	fmt.Printf(loginWarning, os.Args[0])

	golog.Debugf("Setting up browser...")
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("failed to setup signin browser: %v", err)
	}
	defer pw.Stop()

	golog.Debugf("Launching browser...")
	browser, err := pw.Chromium.Launch()
	if err != nil {
		return fmt.Errorf("failed to initialize signin browser 1: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("failed to initialize signin browser 2: %v", err)
	}
	defer page.Close()

	golog.Infof("Navigating to login page...")
	if _, err = page.Goto(fmt.Sprintf("%s/Account/signin", sc2replaystats.WebRoot)); err != nil {
		return fmt.Errorf("failed to navigate to signin page: %v", err)
	}

	golog.Debugf("Filling login form...")
	input, err := page.QuerySelector("css=input[name='email']")
	if err != nil || input == nil {
		return fmt.Errorf("[signin] failed to locate email field: %v", err)
	}
	if err = input.Fill(email); err != nil {
		return fmt.Errorf("[signin] failed to fill email field: %v", err)
	}

	tt := time.Now()
	password, err := speakeasy.Ask(fmt.Sprintf("Password for sc2ReplayStats account %s: ", email))
	golog.Debugf("User input took %s", time.Since(tt))
	if err != nil {
		return fmt.Errorf("failed to prompt user for password: %v", err)
	}

	if input, err = page.QuerySelector("css=input[name='password']"); err != nil || input == nil {
		return fmt.Errorf("[signin] failed to locate password field: %v", err)
	}
	if err = input.Fill(password); err != nil {
		return fmt.Errorf("[signin] failed to fill password field: %v", err)
	}

	golog.Debugf("Submitting login form...")
	if input, err = page.QuerySelector("css=input[value='Sign In']"); err != nil || input == nil {
		return fmt.Errorf("[signin] failed to locate submit button: %v", err)
	}
	if err = input.Click(); err != nil {
		return fmt.Errorf("[signin] failed to click submit button: %v", err)
	}

	url := page.URL()
	if !strings.Contains(url, "display") {
		if alert, err := page.QuerySelector("css=.alert-danger"); err == nil && alert != nil {
			if text, err := alert.InnerText(); err == nil {
				return fmt.Errorf("[signin] login failed, sc2replaystats says: %v", text)
			}
		}

		return fmt.Errorf("[signin] unexpected redirect URL, login probably failed: %v", url)
	}

	parts := strings.Split(url, "/")
	accid := strings.Split(parts[len(parts)-1], "#")[0]
	golog.Infof("Success! Logged in to account #%v", accid)

	if _, err = page.Goto(fmt.Sprintf("%s/account/settings/%v", sc2replaystats.WebRoot, accid)); err != nil {
		return fmt.Errorf("failed to navigate to settings page: %v", err)
	}

	golog.Debugf("Waiting for settings page to load...")
	e, err := page.WaitForSelector("*css=a[data-toggle='tab'] >> text=API Access")
	if err != nil {
		return fmt.Errorf("[settings] failed to locate API Access section: %v", err)
	}
	golog.Debugf("Clicking 'API Access'...")
	if err = e.Click(); err != nil {
		return fmt.Errorf("[settings] failed to click API Access: %v", err)
	}
	golog.Debugf("Finding API key...")
	e, err = page.QuerySelector("*css=.form-group >> text=Authorization Key")
	if e == nil || err != nil {
		golog.Infof("Generating new API key...")
		e, err = page.QuerySelector("text=Generate New API Key")
		if err != nil {
			return fmt.Errorf("[settings] failed to locate Generate New API Key button: %v", err)
		}
		if err = e.Click(); err != nil {
			return fmt.Errorf("[settings] failed to click Generate New API Key button: %v", err)
		}
		e, err = page.WaitForSelector("*css=.form-group >> text=Authorization Key")
	}
	if err != nil || e == nil {
		return fmt.Errorf("[settings] failed to locate \"Authorization Key\" (API Key): %v", err)
	}

	t, err := e.InnerText()
	if err != nil {
		return fmt.Errorf("[settings] failed to resolve \"Authorization Key\" (API Key) Text: %v", err)
	}

	return setAPIkey(strings.Trim(strings.Split(t, ": ")[1], " \r\n\t"))
}

package cmd

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
	"github.com/AlbinoGeek/sc2-rsu/sc2utils"
)

type windowSettings struct {
	*windowBase

	// do we have unsaved changes in the form?
	unsaved bool

	// widgets
	apiKey       *widget.Entry
	autoDownload *widget.Check
	checkUpdates *widget.Check
	replaysRoot  *widget.Entry
	updatePeriod *widget.Entry
}

// TODO: candidate for refactor
func (w *windowSettings) Init() {
	w.windowBase.Window = w.windowBase.app.NewWindow("Settings")

	w.apiKey = widget.NewEntry()
	w.apiKey.SetText(viper.GetString("apiKey"))
	w.apiKey.Validator = func(key string) (err error) {
		if !sc2replaystats.ValidAPIKey(key) {
			err = errors.New("invalid API key format")
		}
		return
	}
	w.apiKey.OnChanged = func(string) {
		w.unsaved = true
	}

	w.autoDownload = widget.NewCheck("Automatically Download Updates?", func(checked bool) {
		w.unsaved = true
	})
	w.autoDownload.SetChecked(viper.GetBool("update.automatic.enabled"))

	w.updatePeriod = widget.NewEntry()
	w.updatePeriod.SetText(getUpdateDuration().String())
	w.updatePeriod.Validator = func(period string) (err error) {
		_, err = time.ParseDuration(period)
		return
	}
	w.updatePeriod.OnChanged = func(string) {
		w.unsaved = true
	}

	w.checkUpdates = widget.NewCheck("Check for Updates Periodically?", func(checked bool) {
		w.unsaved = true
		if checked {
			w.autoDownload.Enable()
			w.updatePeriod.Enable()
		} else {
			w.autoDownload.Disable()
			w.updatePeriod.Disable()
		}
	})
	w.checkUpdates.SetChecked(viper.GetBool("update.check.enabled"))
	if !w.checkUpdates.Checked {
		w.autoDownload.Disable()
		w.updatePeriod.Disable()
	}
	w.unsaved = false // otherwise set by the above line

	w.replaysRoot = widget.NewEntry()
	w.replaysRoot.SetText(viper.GetString("replaysRoot"))
	w.replaysRoot.OnChanged = func(string) {
		w.unsaved = true
	}

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(5, 5))
	w.SetContent(widget.NewVBox(
		widget.NewCard(fmt.Sprintf("%s Settings", PROGRAM), "", widget.NewVBox(
			w.checkUpdates,
			w.autoDownload,
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("Check Every"),
				w.updatePeriod,
			),
		)),
		spacer,
		widget.NewCard("sc2ReplayStats Account", "", widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("API Key"),
				widget.NewHScrollContainer(w.apiKey),
			),
			widget.NewButtonWithIcon("Login and Generate it for me...", theme.ComputerIcon(), w.openLogin),
		)),
		spacer,
		widget.NewCard("StarCraft II", "", widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("Replays Root"),
				widget.NewHScrollContainer(w.replaysRoot),
			),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2),
				widget.NewButtonWithIcon("Find it for me...", theme.SearchIcon(), func() { go w.findReplaysRoot() }),
				widget.NewButtonWithIcon("Browse...", theme.FolderOpenIcon(), func() {
					dlg := dialog.NewFolderOpen(w.browseReplaysRoot, w)
					dlg.Resize(w.Canvas().Size().Subtract(fyne.NewSize(20, 20))) // ! can't be larger than the settings window
					dlg.Show()
				}),
			),
		)),
		spacer,
		layout.NewSpacer(),
		widget.NewSeparator(),
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(2),
			widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
				w.onClose()
			}),
			widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), w.save),
		),
	))

	w.SetCloseIntercept(w.onClose)
	w.SetOnClosed(func() {
		w.SetWindow(nil)
	})

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			w.onClose()
		}
	})

	w.SetPadded(false)
	w.SetFixedSize(true)

	w.Resize(fyne.NewSize(600, 600))
	w.CenterOnScreen()
	w.Show()
}

// TODO: candidate for refactor
func (w *windowSettings) browseReplaysRoot(uri fyne.ListableURI, err error) {
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	if uri == nil {
		return // cancelled
	}

	root := uri.String()
	if strings.HasPrefix(root, "file://") {
		root = root[7:] // ? is this reliable
	}

	// TODO: record the newly found accounts if confirmed
	w.confirmValidReplaysRoot(root, func() {
		w.unsaved = true
		w.replaysRoot.SetText(root)
	})
}

// confirmValidReplaysRoot checks whether there are any accounts found at a
// given root, and if not, asks the user if they would like to use this root
// regardless. If they confirm, or accounts were found, callback is called.
func (w *windowSettings) confirmValidReplaysRoot(root string, callback func()) {
	if accs, err := sc2utils.EnumerateAccounts(root); err == nil && len(accs) > 0 {
		callback()
		return
	}

	dialog.ShowConfirm("Invalid Directory!",
		fmt.Sprintf("We could not find any accounts in that directory.\nAre you sure you want to use it anyways?\n\n%s", root),
		func(ok bool) {
			if ok {
				callback()
			}
		}, w)
}

// TODO: candidate for refactor
func (w *windowSettings) findReplaysRoot() {
	scanRoot := "/"
	if home, err := os.UserHomeDir(); err == nil {
		scanRoot = home
	}

	dlg := dialog.NewProgressInfinite("Searching for Replays Root...",
		"Please wait while we search for a valid Replays folder.\nThis could take several minutes.", w)
	dlg.Show()
	roots, err := sc2utils.FindReplaysRoot(scanRoot)
	dlg.Hide()

	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	if len(roots) == 0 {
		dialog.ShowError(errors.New("no replay directories found"), w)
		return
	}

	if len(roots) == 1 {
		w.confirmValidReplaysRoot(roots[0], func() {
			w.unsaved = true
			w.replaysRoot.SetText(roots[0])
			dialog.ShowInformation("Replays Root Found!", "We found your replays directory!", w)
		})
		return
	}

	selected := -1
	longest := ""
	for _, s := range roots {
		if l := len(s); l > len(longest) {
			longest = s
		}
	}

	listWidget := widget.NewList(func() int {
		return len(roots)
	}, func() fyne.CanvasObject {
		return widget.NewLabel(longest)
	}, func(id int, obj fyne.CanvasObject) {
		obj.(*widget.Label).SetText(roots[id])
	})
	listWidget.OnSelected = func(id int) {
		selected = id
	}
	dlg2 := dialog.NewCustomConfirm("Multiple Possible Roots Found",
		"Select", "Cancel", widget.NewHScrollContainer(listWidget), func(ok bool) {
			if !ok {
				return
			}

			w.confirmValidReplaysRoot(roots[selected], func() {
				w.unsaved = true
				w.replaysRoot.SetText(roots[selected])
			})
		}, w)

	size := fyne.MeasureText(longest, theme.TextSize(), fyne.TextStyle{})
	size.Height *= len(roots)

	dlg2.Resize(fyne.NewSize(60, 144).Add(size))
	dlg2.Show()
}

func (w *windowSettings) onClose() {
	if !w.unsaved {
		w.Close()
		return
	}

	dialog.ShowConfirm("Unsaved Changes",
		"You have not saved your settings.\nAre you sure you want to discard amy changes?",
		func(ok bool) {
			if ok {
				w.Close()
			}
		}, w)
}

func (w *windowSettings) openLogin() {
	user := widget.NewEntry()
	pass := widget.NewPasswordEntry()

	// TODO: actually write a different warning for the gui instead of recycling the cli warning
	warning := widget.NewLabel(loginWarning[:strings.LastIndexByte(loginWarning[:len(loginWarning)-1], '.')+1])
	warning.Wrapping = fyne.TextWrapWord
	vbox := widget.NewVBox(
		warning,
		layout.NewSpacer(),
		widget.NewForm(
			widget.NewFormItem("Email", user),
			widget.NewFormItem("Password", pass),
		),
		layout.NewSpacer(),
	)

	dlg := dialog.NewCustomConfirm("Login to sc2replaystats", "Login", "Cancel", vbox, func(ok bool) {
		if ok {
			dlg2 := dialog.NewProgressInfinite("1) Login", "Setting up our login browser...", w)
			dlg2.Show()
			pw, browser, page, err := newBrowser()
			if pw != nil {
				defer pw.Stop()
			}
			if browser != nil {
				defer browser.Close()
			}
			if page != nil {
				defer page.Close()
			}
			dlg2.Hide()

			if err != nil {
				dialog.ShowError(fmt.Errorf("failed setting up browser: %v", err), w)
				return
			}

			dlg2 = dialog.NewProgressInfinite("2) Login", "Please wait while we try to login to sc2replaystats...", w)
			dlg2.Show()
			accid, err := login(page, user.Text, pass.Text)
			dlg2.Hide()
			if err != nil {
				dialog.ShowError(fmt.Errorf("login error: %v", err), w)
				return
			}
			dlg2 = dialog.NewProgressInfinite("3) Login", "Finding or Generating API Key...", w)
			dlg2.Show()
			key, err := extractAPIKey(page, accid)
			dlg2.Hide()
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to get API key: %v", err), w)
				return
			}

			w.apiKey.SetText(key)
			w.apiKey.Validate()
		}
	}, w)

	vbox.Resize(fyne.NewSize(999, 280))
	dlg.Resize(fyne.NewSize(420, 280))
	dlg.Show()
}

func (w *windowSettings) save() {
	if err := w.validate(); err != nil {
		dialog.ShowError(err, w)
		return
	}

	viper.Set("replaysRoot", w.replaysRoot.Text)
	viper.Set("update.automatic.enabled", w.autoDownload.Checked)
	viper.Set("update.check.enabled", w.checkUpdates.Checked)

	if err := saveConfig(); err != nil {
		dialog.ShowError(err, w)
		return
	}

	dialog.ShowInformation("Saved!", "Your settings have been saved.", w.ui.windows[WindowMain].GetWindow())
	w.unsaved = false
	w.Close()
}

func (w *windowSettings) validate() error {
	if err := w.apiKey.Validate(); err != nil {
		return fmt.Errorf("invalid value for \"API Key\": %v", err)
	}
	if err := w.replaysRoot.Validate(); err != nil {
		return fmt.Errorf("invalid value for \"Replays Root\": %v", err)
	}
	if err := w.updatePeriod.Validate(); err != nil {
		return fmt.Errorf("invalid value for \"Check Every\": %v", err)
	}

	return nil
}

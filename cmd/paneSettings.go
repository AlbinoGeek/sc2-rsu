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
	"fyne.io/fyne/container"
	"fyne.io/fyne/data/binding"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
	"github.com/AlbinoGeek/sc2-rsu/fynex"
	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
	"github.com/AlbinoGeek/sc2-rsu/sc2utils"
)

type paneSettings struct {
	fynex.Pane

	// do we have unsaved changes in the form?
	unsaved         binding.Bool
	validationError binding.String

	apiKey       binding.String
	autoDownload binding.Bool
	checkUpdates binding.Bool
	replaysRoot  binding.String

	// widgets
	updatePeriod *widget.Entry
}

func makePaneSettings(w gui.Window) fynex.Pane {
	p := &paneSettings{
		Pane:            fynex.NewPaneWithIcon("Settings", theme.SettingsIcon(), w),
		unsaved:         binding.NewBool(),
		validationError: binding.NewString(),
	}

	p.Init()
	return p
}

func (settings *paneSettings) boundCheck(confKey, label string) (binding.Bool, *widget.Check) {
	b := binding.NewBool()
	b.Set(viper.GetBool(confKey))
	b.AddListener(binding.NewDataListener(func() {
		settings.unsaved.Set(true)
	}))

	return b, widget.NewCheckWithData(label, b)
}

func (settings *paneSettings) boundEntry(confKey, placeHolder string, fn func(string) error) (binding.String, *widget.Entry) {
	b := validatedString(fn)
	b.Set(viper.GetString(confKey))
	b.AddListener(binding.NewDataListener(func() {
		settings.unsaved.Set(true)
	}))

	e := widget.NewEntryWithData(b)
	e.SetPlaceHolder(placeHolder)
	e.SetOnValidationChanged(settings.validationChanged)

	return b, e
}

func (settings *paneSettings) validationChanged(err error) {
	if err != nil {
		settings.validationError.Set(err.Error())
	} else {
		settings.validationError.Set("")
	}
}

// TODO: candidate for refactor
func (settings *paneSettings) Init() {
	var (
		apiKeyEntry, replaysRootEntry        *widget.Entry
		autoDownloadCheck, checkUpdatesCheck *widget.Check
	)

	settings.apiKey, apiKeyEntry = settings.boundEntry("apiKey", "API Key", func(s string) (err error) {
		if s != "" && !sc2replaystats.ValidAPIKey(s) {
			err = errors.New("invalid API key format")
		}

		return
	})

	settings.autoDownload, autoDownloadCheck = settings.boundCheck("update.automatic.enabled", "Automatically Download Updates?")

	settings.updatePeriod = widget.NewEntry()
	settings.updatePeriod.SetText(getUpdateDuration().String())
	settings.updatePeriod.Validator = func(period string) (err error) {
		_, err = time.ParseDuration(period)
		return
	}
	settings.updatePeriod.OnChanged = func(string) {
		settings.unsaved.Set(true)
	}

	settings.checkUpdates, checkUpdatesCheck = settings.boundCheck("update.check.enabled", "Check for Updates Periodically?")
	settings.checkUpdates.AddListener(binding.NewDataListener(func() {
		if checked, _ := settings.checkUpdates.Get(); checked {
			autoDownloadCheck.Enable()
			settings.updatePeriod.Enable()
		} else {
			autoDownloadCheck.Disable()
			settings.updatePeriod.Disable()
		}
	}))

	settings.replaysRoot, replaysRootEntry = settings.boundEntry("replaysRoot", "Replays Root", nil)

	btnSave := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), settings.save)
	btnSave.Importance = widget.HighImportance

	pad := theme.Padding()
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(pad, pad))

	w := settings.GetWindow().GetWindow()
	settings.SetContent(container.NewBorder(
		nil,
		container.NewVBox(
			widget.NewSeparator(),
			btnSave,
		),
		nil,
		nil,
		container.NewVScroll(widget.NewVBox(
			fynex.NewTextWithStyle("StarCraft II", fyne.TextAlignLeading, fynex.StyleHeading5()),
			container.NewHScroll(replaysRootEntry),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2),
				widget.NewButtonWithIcon("Find it for me...", theme.SearchIcon(), func() { go settings.findReplaysRoot() }),
				widget.NewButtonWithIcon("Select folder...", theme.FolderOpenIcon(), func() {
					dlg := dialog.NewFolderOpen(settings.browseReplaysRoot, w)
					dlg.Show()
					dlg.Resize(w.Canvas().Size().Subtract(fyne.NewSize(20, 20))) // ! can't be larger than the settings window
				}),
			),
			spacer,
			fynex.NewTextWithStyle("sc2ReplayStats", fyne.TextAlignLeading, fynex.StyleHeading5()),
			container.NewHScroll(apiKeyEntry),
			widget.NewButtonWithIcon("Login and Generate it for me...", theme.ComputerIcon(), settings.openLogin),
			spacer,
			fynex.NewTextWithStyle("Updates", fyne.TextAlignLeading, fynex.StyleHeading5()),
			checkUpdatesCheck,
			autoDownloadCheck,
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("Check Every"),
				settings.updatePeriod,
			),
		)),
	))

	settings.unsaved.AddListener(binding.NewDataListener(func() {
		if b, _ := settings.unsaved.Get(); b {
			btnSave.Enable()
		} else {
			btnSave.Disable()
		}
	}))
	settings.unsaved.Set(false)
}

// TODO: candidate for refactor
func (settings *paneSettings) browseReplaysRoot(uri fyne.ListableURI, err error) {
	if err != nil {
		dialog.ShowError(err, settings.GetWindow().GetWindow())
		return
	}

	if uri == nil {
		return // cancelled
	}

	root := strings.TrimPrefix(uri.String(), "file://")

	// TODO: record the newly found accounts if confirmed
	settings.confirmValidReplaysRoot(root, func() {
		settings.unsaved.Set(true)
		settings.replaysRoot.Set(root)
	})
}

// confirmValidReplaysRoot checks whether there are any accounts found at a
// given root, and if not, asks the user if they would like to use this root
// regardless. If they confirm, or accounts were found, callback is called.
func (settings *paneSettings) confirmValidReplaysRoot(root string, callback func()) {
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
		}, settings.GetWindow().GetWindow())
}

// TODO: candidate for refactor
func (settings *paneSettings) findReplaysRoot() {
	w := settings.GetWindow().GetWindow()
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
		settings.confirmValidReplaysRoot(roots[0], func() {
			settings.unsaved.Set(true)
			settings.replaysRoot.Set(roots[0])
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
		"Select", "Cancel", container.NewHScroll(listWidget), func(ok bool) {
			if !ok {
				return
			}

			settings.confirmValidReplaysRoot(roots[selected], func() {
				settings.unsaved.Set(true)
				settings.replaysRoot.Set(roots[selected])
			})
		}, w)

	size := fyne.MeasureText(longest, theme.TextSize(), fyne.TextStyle{})
	size.Height *= float32(len(roots))

	dlg2.Show()
	dlg2.Resize(fyne.NewSize(60, 144).Add(size))
}

func (settings *paneSettings) onClose() {
	w := settings.GetWindow().GetWindow()

	if unsaved, _ := settings.unsaved.Get(); !unsaved {
		return
	}

	dialog.ShowConfirm("Unsaved Changes",
		"You have not saved your settings.\nAre you sure you want to discard amy changes?",
		func(ok bool) {
			settings.unsaved.Set(false) // ignore unsaved changes
		}, w)
}

func (settings *paneSettings) openLogin() {
	w := settings.GetWindow().GetWindow()
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

			if err := settings.apiKey.Set(key); err != nil {
				dialog.ShowError(err, w)
				return
			}
		}
	}, w)

	dlg.Show()
	vbox.Resize(fyne.NewSize(999, 280))
	dlg.Resize(fyne.NewSize(420, 280))
}

func (settings *paneSettings) save() {
	w := settings.GetWindow().GetWindow()

	if err, _ := settings.validationError.Get(); err != "" {
		dialog.ShowError(errors.New(err), w)
		return
	}

	main := settings.GetWindow().(*windowMain)

	s, _ := settings.replaysRoot.Get()
	if main.gettingStarted == 2 && s != "" {
		main.nav.Select(3) // ! ID BASED IS ERROR PRONE
		// main.openGettingStart/ed3()
	}

	s, _ = settings.apiKey.Get()
	if main.gettingStarted == 3 && s != "" {
		main.nav.Select(3) // ! ID BASED IS ERROR PRONE
		// main.openGettingStarted4()
	}

	var changes bool

	if oldKey := viper.Get("apikey"); oldKey != s {
		viper.Set("apikey", s)

		changes = true

		// Use the new apiKey immediately
		sc2api = sc2replaystats.New(s)
	}

	s, _ = settings.replaysRoot.Get()
	if oldRoot := viper.Get("replaysRoot"); oldRoot != s {
		viper.Set("replaysRoot", s)

		changes = true
	}

	if changes {
		main.accounts.Refresh()
		main.setupUploader()
	}

	b, _ := settings.autoDownload.Get()
	viper.Set("update.automatic.enabled", b)

	b, _ = settings.checkUpdates.Get()
	viper.Set("update.check.enabled", b)

	viper.Set("version", VERSION)

	if err := saveConfig(); err != nil {
		dialog.ShowError(err, w)

		return
	}

	dialog.ShowInformation("Saved!", "Your settings have been saved.", w)

	settings.unsaved.Set(false)
}

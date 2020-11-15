package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
	"github.com/AlbinoGeek/sc2-rsu/sc2utils"
)

// TODO: GUI could be a struct reducing global variables

var (
	fyneApp     fyne.App
	mainWindow  fyne.Window
	replaysRoot *widget.Entry
	settings    fyne.Window
	unsaved     bool
)

func gui() error {
	fyneApp = app.New()
	fyneApp.Settings().SetTheme(theme.DarkTheme())

	mainWindow = fyneApp.NewWindow(fmt.Sprintf("SC2ReplayStats Uploader (%s)", PROGRAM))
	mainWindow.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("Menu",
			fyne.NewMenuItem("Check for Updates", guiCheckUpdate),
			fyne.NewMenuItem("Settings", guiSettings),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Browse Source", guiOpenGithub("")),
			fyne.NewMenuItem("Report Bug", guiOpenGithub("issues/new?assignees=AlbinoGeek&labels=bug&template=bug-report.md&title=%5BBUG%5D")),
			fyne.NewMenuItem("Request Feature", guiOpenGithub("issues/new?assignees=AlbinoGeek&labels=enhancement&template=feature-request.md&title=%5BFEATURE+REQUEST%5D")),
		),
	))

	mainWindow.Resize(fyne.NewSize(420, 360))
	mainWindow.CenterOnScreen()
	mainWindow.Show()

	// accs, err := findAccounts(replaysRoot)
	// if err != nil {
	// 	return err
	// }

	// paths := make([]string, 0)
	// for _, a := range accs {
	// 	p := filepath.Join(replaysRoot, a, "Replays", "Multiplayer")
	// 	if f, err := os.Stat(p); err == nil && f.IsDir() {
	// 		paths = append(paths, p)
	// 	}
	// }

	// golog.Info("Starting Automatic Replay Uploader...")
	// sc2api = sc2replaystats.New(key)
	// go automaticUpload(paths)

	hello := widget.NewLabel("Hello Fyne!")
	mainWindow.SetContent(widget.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	fyneApp.Run()
	return nil
}

func guiCheckUpdate() {
	dlg := dialog.NewProgressInfinite("Check for Updates",
		"Checking for new releases...",
		mainWindow)
	dlg.Show()

	rel := checkUpdate()
	if rel == nil {
		dlg.Hide()
		dialog.ShowInformation("Check for Updates",
			fmt.Sprintf("You are running version %s.\nNo updates are available at this time.", VERSION),
			mainWindow)
		return
	}

	dlg.Hide()
	dialog.ShowConfirm("Update Available!",
		fmt.Sprintf("You are running version %s.\nAn update is avaialble: %s\nWould you like us to download it now?", VERSION, rel.GetTagName()),
		guiDoUpdate(rel),
		mainWindow)
}

func guiDoUpdate(rel *github.RepositoryRelease) func(bool) {
	return func(ok bool) {
		if !ok {
			return
		}

		// TODO: display download progress, filename and size
		dlg := dialog.NewProgressInfinite("Downloading Update",
			fmt.Sprintf("Downloading version %s now...", rel.GetTagName()),
			mainWindow)
		dlg.Show()

		err := downloadUpdate(rel)
		dlg.Hide()

		if err != nil {
			dialog.ShowError(err, mainWindow)
		} else {
			dialog.ShowInformation("Update Complete!", "Please close the program and start the new binary.", mainWindow)
		}
	}
}

func guiSettings() {
	if settings == nil {
		guiSettingsInit()
	}

	settings.Show()
	settings.CenterOnScreen()
}

func guiSettingsBrowseReplaysRoot(uri fyne.ListableURI, err error) {
	if err != nil {
		dialog.ShowError(err, settings)
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
	if accs, err := sc2utils.EnumerateAccounts(root); err != nil || len(accs) == 0 {
		dialog.ShowConfirm("Invalid Directory!",
			fmt.Sprintf("We could not find any accounts in that directory.\nAre you sure you want to use it anyways?\n\n%s", root),
			func(ok bool) {
				if ok {
					unsaved = true
					replaysRoot.SetText(root)
				}
			},
			settings)
		return
	}

	unsaved = true
	replaysRoot.SetText(root)
}

// TODO: refactor
func guiSettingsInit() {
	settings = fyneApp.NewWindow("Settings")

	updatePeriod := widget.NewEntry()
	updatePeriod.SetText(getUpdateDuration().String())
	updatePeriod.Validator = func(period string) (err error) {
		_, err = time.ParseDuration(period)
		return
	}

	autoDownload := widget.NewCheck("Automatically Download Updates?", func(checked bool) {
		unsaved = true
		viper.Set("update.automatic.enabled", checked)
	})
	autoDownload.SetChecked(viper.GetBool("update.automatic.enabled"))

	checkUpdates := widget.NewCheck("Check for Updates Periodically?", func(checked bool) {
		unsaved = true
		if checked {
			autoDownload.Enable()
			updatePeriod.Enable()
		} else {
			autoDownload.Disable()
			updatePeriod.Disable()
		}
		viper.Set("update.check.enabled", checked)
	})
	checkUpdates.SetChecked(viper.GetBool("update.check.enabled"))

	if !checkUpdates.Checked {
		autoDownload.Disable()
		updatePeriod.Disable()
	}

	apiKey := widget.NewEntry()
	apiKey.SetText(viper.GetString("apiKey"))
	apiKey.Validator = func(key string) (err error) {
		if !sc2replaystats.ValidAPIKey(key) {
			err = errors.New("invalid API key format")
		}
		return
	}

	replaysRoot = widget.NewEntry()
	replaysRoot.SetText(viper.GetString("replaysRoot"))

	settings.SetContent(widget.NewVBox(
		widget.NewCard(fmt.Sprintf("%s Settings", PROGRAM), "", widget.NewVBox(
			checkUpdates,
			autoDownload,
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("Check Every"),
				updatePeriod,
			),
		)),
		widget.NewCard("sc2ReplayStats Account", "", widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("API Key"),
				widget.NewHScrollContainer(apiKey),
			),
			widget.NewButtonWithIcon("Find it for me...", theme.ComputerIcon(), func() {
				// ! IMPLEMENT LOGIN FORM
			}),
		)),
		widget.NewCard("StarCraft II", "", widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("Replays Root"),
				widget.NewHScrollContainer(replaysRoot),
			),
			widget.NewButtonWithIcon("Browse...", theme.FolderOpenIcon(), func() {
				dlg := dialog.NewFolderOpen(guiSettingsBrowseReplaysRoot, settings)
				dlg.Resize(fyne.NewSize(1000, 1000)) // ! can't be larger than the settings window
				dlg.Show()
			}),
			// ! requires refactor of findReplaysRoot
			// widget.NewButtonWithIcon("Search", theme.SearchIcon(), func() {
			// }),
			// ),
		)),
		layout.NewSpacer(),
		widget.NewSeparator(),
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(2),
			widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
				unsaved = false
				loadConfig()
				settings.Close()
			}),
			widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
				if err := saveConfig(); err != nil {
					dialog.ShowError(err, settings)
				} else {
					unsaved = false
					dialog.ShowInformation("Saved!", "Your settings have been saved.", settings)
				}
			}),
		),
	))

	settings.Resize(fyne.NewSize(420, 420))
	settings.SetFixedSize(true)
	settings.SetPadded(false)
	settings.SetOnClosed(func() {
		// ? for some reason we can't just re-use the window
		settings = nil
	})

	shouldClose := func() {
		if !unsaved {
			settings.Close()
			return
		}

		dialog.ShowConfirm("Unsaved Changes",
			"You have not saved your settings.\nDo you want to discard your changes?",
			func(ok bool) {
				if ok {
					settings.Close()
				}
			},
			settings)
	}

	settings.SetCloseIntercept(shouldClose)
	settings.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			shouldClose()
		}
	})
}

func guiOpenGithub(slug string) func() {
	u, _ := url.Parse(fmt.Sprintf("https://github.com/%s/%s/%s", ghOwner, ghRepo, slug))
	return func() {
		if err := fyneApp.OpenURL(u); err != nil {
			dialog.ShowError(err, mainWindow)
		}
	}
}

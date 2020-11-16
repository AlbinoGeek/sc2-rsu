package cmd

import (
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"os"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
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
	about       fyne.Window
	unsaved     bool
)

func gui() error {
	fyneApp = app.New()
	fyneApp.Settings().SetTheme(theme.DarkTheme())

	guiMainInit()

	fyneApp.Run()
	return nil
}

func guiMainInit() {
	mainWindow = fyneApp.NewWindow(fmt.Sprintf("SC2ReplayStats Uploader (%s)", PROGRAM))
	mainWindow.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Check for Updates", guiCheckUpdate),
			fyne.NewMenuItem("Settings", guiSettings),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("Report Bug", guiOpenGithub("issues/new?assignees=AlbinoGeek&labels=bug&template=bug-report.md&title=%5BBUG%5D")),
			fyne.NewMenuItem("Request Feature", guiOpenGithub("issues/new?assignees=AlbinoGeek&labels=enhancement&template=feature-request.md&title=%5BFEATURE+REQUEST%5D")),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("About", guiAbout),
		),
	))

	hello := widget.NewLabel("Hello Fyne!")
	mainWindow.SetContent(widget.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	mainWindow.Resize(fyne.NewSize(420, 360))
	mainWindow.CenterOnScreen()
	mainWindow.Show()

	// choice := ""
	// listWidget := widget.NewSelect(data, func(s string) {
	// 	choice = s
	// })
	// dlg2 := dialog.NewCustomConfirm("Multiple Possible Roots Found",
	// 	"Select", "Cancel", listWidget, func(ok bool) {
	// 		if !ok {
	// 			return
	// 		}
	// 		// viper.Set("replaysRoot", roots[selected])
	// 		// replaysRoot.SetText(roots[selected])
	// 		viper.Set("replaysRoot", choice)
	// 		replaysRoot.SetText(choice)
	// 	}, settings)

	// // ! need a way better way to figure out the size of the dialog
	// dlg2.Resize(fyne.NewSize(100+12*len(roots[0]), 140+30*len(roots)))
	// dlg2.Show()
	// return

	if viper.GetString("version") == "" {
		guiFirstRun()
	}
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

func guiFirstRun() {
	// modal := widget.NewModalPopUp(
	// 	widget.NewCard("Welcome!", "First-Time Setup",
	// 		widget.NewVBox(
	// 			widget.NewLabel("You are only two steps away from having your replays automatically uploaded to sc2replaystats!"),
	// 		),
	// 	), mainWindow.Canvas())
	// modal.Show()
}

func guiAbout() {
	if about == nil {
		guiAboutInit()
	}

	about.Show()
	about.CenterOnScreen()
}

func guiAboutInit() {
	about = fyneApp.NewWindow("About")

	u, _ := url.Parse(ghLink(""))
	about.SetContent(
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewVBox(
			widget.NewHBox(
				layout.NewSpacer(),
				widget.NewCard(PROGRAM, "", nil),
				layout.NewSpacer(),
			),
			widget.NewHBox(
				layout.NewSpacer(),
				widget.NewForm(
					widget.NewFormItem("Author", widget.NewLabel(ghOwner)),
					widget.NewFormItem("Version", widget.NewLabel(VERSION)),
				),
				layout.NewSpacer(),
			),
			widget.NewHBox(
				layout.NewSpacer(),
				widget.NewHyperlink("Browse Source", u),
				layout.NewSpacer(),
			),
		)),
	)

	about.Resize(fyne.NewSize(200, 160))
	about.SetFixedSize(true)
	about.SetPadded(false)
	about.SetOnClosed(func() {
		about = nil
	})

	about.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			about.Close()
		}
	})
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
					viper.Set("replaysRoot", root)
					replaysRoot.SetText(root)
				}
			},
			settings)
		return
	}

	unsaved = true
	viper.Set("replaysRoot", root)
	replaysRoot.SetText(root)
}

func guiSettingsFindReplaysRoot(entry *widget.Entry) func() {
	scanRoot := "/"
	if home, err := os.UserHomeDir(); err == nil {
		scanRoot = home
	}

	return func() {
		dlg := dialog.NewProgressInfinite("Searching for Replays Root...",
			"Please wait while we search for a valid Replays folder.\nThis could take several minutes.",
			settings)
		dlg.Show()
		roots, err := sc2utils.FindReplaysRoot(scanRoot)
		dlg.Hide()

		if err != nil {
			dialog.ShowError(err, settings)
			return
		}

		if len(roots) == 0 {
			dialog.ShowError(errors.New("no replay directories found"), settings)
			return
		}

		if len(roots) == 1 {
			accs, err := sc2utils.EnumerateAccounts(roots[0])
			if err != nil {
				dialog.ShowError(fmt.Errorf("error scanning for accounts: %v", err), settings)
				return
			}

			dialog.ShowInformation("Replays Root Found!",
				fmt.Sprintf("We found your replays directory!\nIt contains %d account/toons.\n%s", len(accs), roots[0]),
				settings)
			entry.SetText(roots[0])
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
		// choice := ""
		// listWidget := widget.NewSelect(roots, func(s string) {
		// 	choice = s
		// })
		// listWidget := widget.NewEntry()
		dlg2 := dialog.NewCustomConfirm("Multiple Possible Roots Found",
			"Select", "Cancel", widget.NewHScrollContainer(listWidget), func(ok bool) {
				if !ok {
					return
				}
				_ = selected
				// viper.Set("replaysRoot", roots[selected])
				// replaysRoot.SetText(roots[selected])
				// viper.Set("replaysRoot", choice)
				// replaysRoot.SetText(choice)
			}, settings)

		size := fyne.MeasureText(longest, theme.TextSize(), fyne.TextStyle{})
		size.Height *= len(roots)

		dlg2.Resize(fyne.NewSize(60, 144).Add(size))
		dlg2.Show()
	}
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

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(5, 5))
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
		spacer,
		widget.NewCard("sc2ReplayStats Account", "", widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("API Key"),
				widget.NewHScrollContainer(apiKey),
			),
			widget.NewButtonWithIcon("Login and Generate it for me...", theme.ComputerIcon(), func() {
				// ! IMPLEMENT LOGIN FORM
			}),
		)),
		spacer,
		widget.NewCard("StarCraft II", "", widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("Replays Root"),
				widget.NewHScrollContainer(replaysRoot),
			),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2),
				widget.NewButtonWithIcon("Find it for me...", theme.SearchIcon(), guiSettingsFindReplaysRoot(replaysRoot)),
				widget.NewButtonWithIcon("Browse...", theme.FolderOpenIcon(), func() {
					dlg := dialog.NewFolderOpen(guiSettingsBrowseReplaysRoot, settings)
					dlg.Resize(settings.Canvas().Size().Subtract(fyne.NewSize(20, 20))) // ! can't be larger than the settings window
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
				unsaved = false
				loadConfig()
				settings.Close()
			}),
			widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
				if err := saveConfig(); err != nil {
					dialog.ShowError(err, settings)
				} else {
					unsaved = false
					dialog.ShowInformation("Saved!", "Your settings have been saved.", mainWindow)
					settings.Close()
				}
			}),
		),
	))

	settings.Resize(fyne.NewSize(600, 600))
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

func ghLink(slug string) string {
	return fmt.Sprintf("https://github.com/%s/%s/%s", ghOwner, ghRepo, slug)
}

func guiOpenGithub(slug string) func() {
	u, _ := url.Parse(ghLink(slug))
	return func() {
		if err := fyneApp.OpenURL(u); err != nil {
			dialog.ShowError(err, mainWindow)
		}
	}
}

package cmd

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/viper"
)

type windowMain struct {
	*windowBase
	modal *widget.PopUp
}

func (w *windowMain) Init() {
	w.windowBase.Window = w.windowBase.app.NewWindow("SC2ReplayStats Uploader")

	w.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Check for Updates", func() { go w.checkUpdate() }),
			fyne.NewMenuItem("Settings", func() { w.ui.OpenWindow(WindowSettings) }),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("Report Bug", w.ui.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=bug&template=bug-report.md&title=%5BBUG%5D")),
			fyne.NewMenuItem("Request Feature", w.ui.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=enhancement&template=feature-request.md&title=%5BFEATURE+REQUEST%5D")),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("About", func() { w.ui.OpenWindow(WindowAbout) }),
		),
	))

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(widget.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	// closing the main window should quit the application
	w.SetCloseIntercept(func() {
		// Close "About" if it's open
		if win := w.ui.windows[WindowAbout].GetWindow(); win != nil {
			win.Close()
		}

		win := w.ui.windows[WindowSettings]
		if win.GetWindow() != nil {
			settings := win.(*windowSettings)
			if settings.unsaved {
				settings.onClose()
				return
			}
			settings.Close()
		}

		w.Close()
		w.app.Quit()
	})

	w.Resize(fyne.NewSize(420, 360))
	w.CenterOnScreen()
	w.Show()

	if viper.GetString("version") == "" || viper.GetString("apikey") == "" {
		w.openGettingStarted1()
	}
}

func labelWithWrapping(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

func (w *windowMain) checkUpdate() {
	dlg := dialog.NewProgressInfinite("Check for Updates", "Checking for new releases...", w)
	dlg.Show()
	rel := checkUpdate()
	dlg.Hide()

	if rel == nil {
		dialog.ShowInformation("Check for Updates",
			fmt.Sprintf("You are running version %s.\nNo updates are available at this time.", VERSION), w)
		return
	}

	dialog.ShowConfirm("Update Available!",
		fmt.Sprintf("You are running version %s.\nAn update is avaialble: %s\nWould you like us to download it now?", VERSION, rel.GetTagName()),
		w.doUpdate(rel), w)
}

func (w *windowMain) doUpdate(rel *github.RepositoryRelease) func(bool) {
	return func(ok bool) {
		if !ok {
			return
		}

		// otherwise we might block the fyne event queue...
		go func() {
			// TODO: display download progress, filename and size
			dlg := dialog.NewProgressInfinite("Downloading Update",
				fmt.Sprintf("Downloading version %s now...", rel.GetTagName()), w)
			dlg.Show()
			err := downloadUpdate(rel)
			dlg.Hide()

			if err != nil {
				dialog.ShowError(err, w)
			} else {
				dialog.ShowInformation("Update Complete!", "Please close the program and start the new binary.", w)
			}
		}()
	}
}

func (w *windowMain) WizardModal(skipText, nextText string, skipFn, nextFn func(), contents ...fyne.CanvasObject) {
	if skipFn == nil {
		skipFn = func() { w.modal.Hide() }
	}
	if nextFn == nil {
		nextFn = func() { w.modal.Hide() }
	}

	buttons := make([]fyne.CanvasObject, 0)
	if skipText != "" {
		btn := widget.NewButtonWithIcon(skipText, theme.CancelIcon(), skipFn)
		btn.Importance = widget.LowImportance
		buttons = append(buttons, btn)
	}
	if nextText != "" {
		btn := widget.NewButtonWithIcon(nextText, theme.NavigateNextIcon(), nextFn)
		btn.Importance = widget.HighImportance
		buttons = append(buttons, btn)
	}

	// Re-use existing Modal
	if w.modal != nil {
		box := w.modal.Content.(*widget.Box)
		box.Children[0].(*widget.Card).Content.(*fyne.Container).Objects = contents
		box.Children[len(box.Children)-1] = fyne.NewContainerWithLayout(
			layout.NewGridLayout(len(buttons)),
			buttons...,
		)

		w.modal.Show()
		// box.Resize(fyne.NewSize(320, 220)) // 3
		w.modal.Refresh()
		return
	}

	// Create Fresh Modal
	container := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		contents...,
	)

	box := widget.NewVBox(
		widget.NewCard("Welcome!", "First-Time Setup",
			container,
		),
		layout.NewSpacer(),
		widget.NewSeparator(),
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(len(buttons)),
			buttons...,
		),
	)

	w.modal = widget.NewModalPopUp(box, w.Canvas())
	w.modal.Show()

	w.modal.Resize(fyne.NewSize(320, 220)) // 1
	box.Resize(fyne.NewSize(320, 220))     // 2
}

func (w *windowMain) openGettingStarted1() {
	w.WizardModal("Skip", "Next", nil, func() {
		if viper.GetString("replaysroot") == "" {
			w.openGettingStarted2()
		} else {
			w.modal.Hide()
		}
	},
		labelWithWrapping("You are only two steps away from having your replays automatically uploaded!"),
		labelWithWrapping("1) We will find your Replays Directory"),
		labelWithWrapping("2) We will find your sc2replaystats API Key"),
	)
}

func (w *windowMain) openGettingStarted2() {
	btnSettings := widget.NewButtonWithIcon("Open Settings", theme.SettingsIcon(), func() {
		w.ui.OpenWindow(WindowSettings)
	})
	btnSettings.Importance = widget.HighImportance

	// TODO: Refactor this to actually have the settings UI, not just direct the user to settings
	w.WizardModal("", "Next", nil, func() {
		if viper.GetString("apikey") == "" {
			w.openGettingStarted3()
		} else {
			w.modal.Hide()
		}
	},
		labelWithWrapping("First thing's first. We need to find your StarCraft II Replays Directory!"),
		btnSettings,
	)
}

func (w *windowMain) openGettingStarted3() {
	btnSettings := widget.NewButtonWithIcon("Open Settings", theme.SettingsIcon(), func() {
		w.ui.OpenWindow(WindowSettings)
	})
	btnSettings.Importance = widget.HighImportance

	// TODO: Refactor this to actually have the settings UI, not just direct the user to settings
	w.WizardModal("", "Next", nil, func() {
		w.modal.Hide()
	},
		labelWithWrapping("Lastly, please set your sc2replaystats API key. If you do not know how to find this, use the \"Login and find it for me\" button to have us login to your account and generate one on your behalf."),
		btnSettings,
	)
}

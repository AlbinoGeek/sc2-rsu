package cmd

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"

	"github.com/google/go-github/v32/github"
)

type windowMain struct {
	*windowBase
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

	w.Resize(fyne.NewSize(420, 360))
	w.CenterOnScreen()
	w.Show()
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

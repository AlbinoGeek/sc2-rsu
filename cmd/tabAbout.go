package cmd

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
	"github.com/google/go-github/v32/github"
)

type tabAbout struct {
	*gui.TabBase
}

func makeTabAbout(w gui.Window) gui.Tab {
	tab := &tabAbout{
		TabBase: gui.NewTabWithIcon("", theme.InfoIcon(), w),
	}

	tab.Init()
	tab.Refresh()
	return tab
}

func (t *tabAbout) Init() {
	main := t.GetWindow().(*windowMain)
	sourceURL, _ := url.Parse(ghLink(""))

	t.SetContent(widget.NewVBox(
		// widget.NewHBox(
		// layout.NewSpacer(),
		// newHeader(PROGRAM),
		// layout.NewSpacer(),
		// layout.NewSpacer(),
		// ),
		widget.NewForm(
			widget.NewFormItem("Author", widget.NewLabel(ghOwner)),
			widget.NewFormItem("Version", widget.NewLabel(VERSION)),
		),
		widget.NewButtonWithIcon("Check for Updates", theme.ViewRefreshIcon(), func() { go t.checkUpdate() }),
		widget.NewButtonWithIcon("Request Feedback", feedbackIcon, main.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=enhancement&template=feature-request.md&title=%5BFEATURE+REQUEST%5D")),
		widget.NewButtonWithIcon("Report A Bug", reportBugIcon, main.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=bug&template=bug-report.md&title=%5BBUG%5D")),
		widget.NewHyperlink("Browse Source", sourceURL),
	))
}

func (t *tabAbout) checkUpdate() {
	w := t.GetWindow().GetWindow()

	dlg := dialog.NewProgressInfinite("Check for Updates", "Checking for new releases...", w)
	dlg.Show()
	rel := checkUpdate()
	dlg.Hide()

	if rel == nil {
		dialog.ShowInformation("Check for Updates", "No updates are available at this time.", w)
		return
	}

	dialog.ShowConfirm("Update Available!",
		fmt.Sprintf("You are running version %s.\nAn update is available: %s\nWould you like us to download it now?", VERSION, rel.GetTagName()),
		t.doUpdate(rel), w)
}

func (t *tabAbout) doUpdate(rel *github.RepositoryRelease) func(bool) {
	return func(ok bool) {
		if !ok {
			return
		}

		w := t.GetWindow().GetWindow()

		// otherwise we might block the fyne event queue...
		go func() {
			// TODO: display download progress, filename and size
			dlg := dialog.NewProgressInfinite("Downloading Update",
				fmt.Sprintf("Downloading version %s nomain...", rel.GetTagName()), w)
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

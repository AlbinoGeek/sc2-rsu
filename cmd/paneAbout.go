package cmd

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
	"github.com/AlbinoGeek/sc2-rsu/cmd/gui/fynewidget"
	"github.com/google/go-github/v32/github"
)

type paneAbout struct {
	fynewidget.Pane
}

func makePaneAbout(w gui.Window) fynewidget.Pane {
	p := &paneAbout{
		fynewidget.NewPaneWithIcon("Help & Feedback", feedbackIcon, w),
	}

	main := w.(*windowMain)
	sourceURL, _ := url.Parse(ghLink(""))

	p.SetContent(widget.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Author", widget.NewLabel(ghOwner)),
			widget.NewFormItem("Version", widget.NewLabel(VERSION)),
		),
		widget.NewButtonWithIcon("Check for Updates", theme.ViewRefreshIcon(), func() { go p.checkUpdate() }),
		widget.NewButtonWithIcon("Request Feedback", feedbackIcon, main.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=enhancement&template=feature-request.md&title=%5BFEATURE+REQUEST%5D")),
		widget.NewButtonWithIcon("Report A Bug", reportBugIcon, main.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=bug&template=bug-report.md&title=%5BBUG%5D")),
		widget.NewHyperlink("Browse Source", sourceURL),
	))

	return p
}

func (p *paneAbout) checkUpdate() {
	w := p.GetWindow().GetWindow()

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
		p.doUpdate(rel), w)
}

func (p *paneAbout) doUpdate(rel *github.RepositoryRelease) func(bool) {
	return func(ok bool) {
		if !ok {
			return
		}

		w := p.GetWindow().GetWindow()

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

package cmd

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
)

type windowAbout struct {
	*gui.WindowBase
}

func (about *windowAbout) Init() {
	w := about.App.NewWindow("About")
	about.SetWindow(w)

	u, _ := url.Parse(ghLink(""))

	w.SetContent(
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewVBox(
			widget.NewHBox(
				layout.NewSpacer(),
				newText(PROGRAM, 1.6, true),
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

	w.SetOnClosed(func() {
		about.SetWindow(nil)
	})

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			w.Close()
		}
	})

	w.SetPadded(false)
	w.SetFixedSize(true)

	w.Resize(fyne.NewSize(200, 160))
	w.CenterOnScreen()
	w.Show()
}

package cmd

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type windowAbout struct {
	*windowBase
}

func (w *windowAbout) Init() {
	w.windowBase.Window = w.windowBase.app.NewWindow("About")

	u, _ := url.Parse(ghLink(""))

	w.SetContent(
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

	w.SetOnClosed(func() {
		w.SetWindow(nil)
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

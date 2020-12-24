package cmd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
	"github.com/AlbinoGeek/sc2-rsu/fynemd"
)

type paneDeveloper struct {
	fynemd.Pane

	container *fyne.Container
}

func makePaneDeveloper(w gui.Window) fynemd.Pane {
	p := &paneDeveloper{
		Pane: fynemd.NewPaneWithIcon("Developer", codeIcon, w),
	}

	p.container = container.NewVBox()
	p.SetContent(container.NewVScroll(p.container))
	p.Init()
	return p
}

func (t *paneDeveloper) Init() {
	main := t.GetWindow().(*windowMain)

	t.container.Add(fynemd.NewTextWithStyle("AppBar Top", fyne.TextAlignLeading, fynemd.StyleHeading5()))

	var (
		denseBtn    = widget.NewButton("Toggle Dense", nil)
		extendedBtn = widget.NewButton("Toggle Extended", nil)

		appTitleEntry    = widget.NewEntry()
		navTitleEntry    = widget.NewEntry()
		navSubtitleEntry = widget.NewEntry()
	)

	denseBtn.OnTapped = func() {
		if main.topbar.Dense {
			denseBtn.Importance = widget.MediumImportance
		} else {
			denseBtn.Importance = widget.HighImportance
		}
		main.topbar.SetDense(!main.topbar.Dense)
	}

	extendedBtn.OnTapped = func() {
		if main.topbar.Extended {
			extendedBtn.Importance = widget.MediumImportance
		} else {
			extendedBtn.Importance = widget.HighImportance
		}
		main.topbar.SetExtended(!main.topbar.Extended)
	}

	t.container.Add(fyne.NewContainerWithLayout(
		layout.NewGridLayout(2),
		denseBtn,
		extendedBtn,
	))

	appTitleEntry.PlaceHolder = "Title Text"
	appTitleEntry.OnChanged = func(str string) {
		main.topbar.SetTitle(str)
	}
	t.container.Add(fyne.NewContainerWithLayout(
		layout.NewFormLayout(),
		widget.NewLabel("Title"),
		appTitleEntry,
	))

	// ---

	t.container.Add(fynemd.NewTextWithStyle("NavDrawer Left", fyne.TextAlignLeading, fynemd.StyleHeading5()))

	hideBtn := widget.NewButton("Toggle Visibility", nil)
	hideBtn.Importance = widget.HighImportance
	hideBtn.OnTapped = func() {
		if main.topbar.NavClosed {
			hideBtn.Importance = widget.HighImportance
			main.topbar.SetNavClosed(false)
		} else {
			hideBtn.Importance = widget.MediumImportance
			main.topbar.SetNavClosed(true)
		}
	}
	t.container.Add(hideBtn)

	navTitleEntry.PlaceHolder = "Title Text"
	navTitleEntry.OnChanged = func(str string) {
		main.nav.SetTitle(str)
	}
	t.container.Add(fyne.NewContainerWithLayout(
		layout.NewFormLayout(),
		widget.NewLabel("Title"),
		navTitleEntry,
	))

	navSubtitleEntry.PlaceHolder = "Subtitle Text"
	navSubtitleEntry.OnChanged = func(str string) {
		main.nav.SetSubtitle(str)
	}
	t.container.Add(fyne.NewContainerWithLayout(
		layout.NewFormLayout(),
		widget.NewLabel("Subtitle"),
		navSubtitleEntry,
	))

	t.container.Refresh()
}

package cmd

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
)

const (
	// WindowAbout is an index string used in GraphicalApplication.window
	WindowAbout = "About"

	// WindowMain is an index string used in GraphicalApplication.window
	WindowMain = "Main"

	// WindowSettings is an index string used in GraphicalApplication.window
	WindowSettings = "Settings"
)

func newUI() *graphicalInterface {
	ui := new(graphicalInterface)
	ui.app = app.New()

	ui.app.Settings().SetTheme(theme.DarkTheme())

	ui.windows = make(map[string]Window)
	ui.windows[WindowMain] = &windowMain{
		windowBase: &windowBase{app: ui.app, ui: ui}}
	ui.windows[WindowAbout] = &windowAbout{
		windowBase: &windowBase{app: ui.app, ui: ui}}
	ui.windows[WindowSettings] = &windowSettings{
		windowBase: &windowBase{app: ui.app, ui: ui}}
	ui.windows[WindowMain].Init()

	return ui
}

type graphicalInterface struct {
	app     fyne.App
	windows map[string]Window
}

// OpenGitHub launches the user's browser to a given GitHub URL relative to
// this project's repository root
func (ui *graphicalInterface) OpenGitHub(slug string) func() {
	u, _ := url.Parse(ghLink(slug))
	return func() {
		if err := ui.app.OpenURL(u); err != nil {
			dialog.ShowError(err, ui.windows[WindowMain].GetWindow())
		}
	}
}

// OpenWindow shows a given window, initializing it first if needed
func (ui *graphicalInterface) OpenWindow(windowName string) {
	if w, ok := ui.windows[windowName]; ok {
		if w.GetWindow() == nil {
			w.Init()
		} else {
			w.Show()
		}
	}
}

// Run starts the fyne.App and shows the main window
func (ui *graphicalInterface) Run() {
	ui.app.Run()
}

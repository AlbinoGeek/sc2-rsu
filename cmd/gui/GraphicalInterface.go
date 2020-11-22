package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
)

// New returns a wrapped fyne.App
func New() *GraphicalInterface {
	ui := new(GraphicalInterface)
	ui.App = app.New()

	ui.Theme = DarkerTheme()
	ui.App.Settings().SetTheme(ui.Theme)

	return ui
}

// GraphicalInterface represents a wrapped fyne.App
type GraphicalInterface struct {
	App     fyne.App
	Theme   fyne.Theme
	Primary Window
	Windows map[string]Window
}

// Init assigns the Windows and Primary, to keep track of managed fyne.Windows
func (ui *GraphicalInterface) Init(windows map[string]Window, primary string) {
	ui.Windows = windows
	ui.Primary = windows[primary]
	ui.Primary.Init()
}

// OpenWindow shows a given window, initializing it first if needed
func (ui *GraphicalInterface) OpenWindow(windowName string) {
	if w, ok := ui.Windows[windowName]; ok {
		if w.GetWindow() == nil {
			w.Init()
		} else {
			w.Show()
		}
	}
}

// Run starts the fyne.App and shows the main window
func (ui *GraphicalInterface) Run() {
	ui.App.Run()
}

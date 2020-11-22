package gui

import "fyne.io/fyne"

// WindowBase implements the Window common pieces interface
type WindowBase struct {
	App fyne.App
	UI  *GraphicalInterface
	w   fyne.Window
}

// GetWindow returns reference to the fyne.Window we are managing
func (b *WindowBase) GetWindow() fyne.Window {
	return b.w
}

// SetWindow immediately changes the fyne.Window we are managing
func (b *WindowBase) SetWindow(w fyne.Window) {
	b.w = w
}

func (b *WindowBase) Init() {
	b.SetWindow(b.App.NewWindow("Hello, World!"))
}

func (b *WindowBase) Show() {
	b.Show()
}

func (b *WindowBase) Hide() {
	b.Show()
}

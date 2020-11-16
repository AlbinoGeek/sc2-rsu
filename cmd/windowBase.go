package cmd

import "fyne.io/fyne"

type windowBase struct {
	fyne.Window
	app fyne.App
	ui  *graphicalInterface
}

func (b *windowBase) GetWindow() fyne.Window {
	return b.Window
}

func (b *windowBase) SetWindow(w fyne.Window) {
	b.Window = w
}

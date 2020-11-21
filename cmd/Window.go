package cmd

import (
	"fyne.io/fyne"
)

// Window represents a managed fyne.Window
type Window interface {
	GetWindow() fyne.Window
	SetWindow(fyne.Window)
	Init()
	Hide()
	Show()
}

package fynemd

import "fyne.io/fyne"

// NavItem ...
type NavItem interface {
	GetContent() fyne.CanvasObject
	GetLabel() fyne.CanvasObject
	// GetIcon() fyne.Resource
	GetTitle() string
}

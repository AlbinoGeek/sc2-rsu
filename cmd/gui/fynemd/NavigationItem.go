package fynemd

import "fyne.io/fyne"

// NavigationItem ...
type NavigationItem interface {
	GetContent() fyne.CanvasObject
	GetLabel() fyne.CanvasObject
	GetIcon() fyne.Resource
	GetTitle() string
}

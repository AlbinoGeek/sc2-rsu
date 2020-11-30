package fynewidget

import "fyne.io/fyne"

// NavigationItem ...
type NavigationItem interface {
	GetContent() fyne.CanvasObject
	GetLabel() fyne.CanvasObject
}

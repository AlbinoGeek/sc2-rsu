package fynemd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// NavigationSeparator ...
type NavigationSeparator struct {
	res fyne.CanvasObject
}

// NewNavigationSeparator ...
func NewNavigationSeparator() NavigationItem {
	return &NavigationSeparator{
		res: widget.NewSeparator(),
	}
}

// GetContent ...
func (*NavigationSeparator) GetContent() fyne.CanvasObject { return nil }

// GetLabel ...
func (l *NavigationSeparator) GetLabel() fyne.CanvasObject { return l.res }

// GetIcon ...
func (*NavigationSeparator) GetIcon() fyne.Resource { return nil }

// GetTitle ...
func (*NavigationSeparator) GetTitle() string { return "" }
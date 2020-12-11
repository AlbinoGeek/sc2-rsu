package fynemd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// NavigationSeparator represents a visual dividing line to be shown between
// other NavigationItem(s) within a NavigationDrawer -- it should not recieve
// focus or be a candidate for selection because it contains no content.
type NavigationSeparator struct {
	res fyne.CanvasObject
}

// NewNavigationSeparator returns a new NavigationDrawer divider
func NewNavigationSeparator() NavigationItem {
	return &NavigationSeparator{
		res: widget.NewSeparator(),
	}
}

// GetContent ...
func (*NavigationSeparator) GetContent() fyne.CanvasObject { return nil }

// GetIcon ...
func (*NavigationSeparator) GetIcon() fyne.Resource { return nil }

// GetLabel ...
func (l *NavigationSeparator) GetLabel() fyne.CanvasObject { return l.res }

// GetTitle ...
func (*NavigationSeparator) GetTitle() string { return "" }

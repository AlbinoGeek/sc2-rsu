package fynemd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// NavSeparator represents a visual dividing line to be shown between
// other NavItem(s) within a NavDrawer -- it should not recieve
// focus or be a candidate for selection because it contains no content.
type NavSeparator struct {
	res fyne.CanvasObject
}

// NewNavSeparator returns a new NavDrawer divider
func NewNavSeparator() NavItem {
	return &NavSeparator{
		res: widget.NewSeparator(),
	}
}

// GetContent ...
func (*NavSeparator) GetContent() fyne.CanvasObject { return nil }

// GetIcon ...
func (*NavSeparator) GetIcon() fyne.Resource { return nil }

// GetLabel ...
func (l *NavSeparator) GetLabel() fyne.CanvasObject { return l.res }

// GetTitle ...
func (*NavSeparator) GetTitle() string { return "" }

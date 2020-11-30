package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

// TabBase implements the Tab common pieces interface
type TabBase struct {
	ti *container.TabItem
	w  Window
}

// NewTab creates a managed container.TabItem
func NewTab(text string, w Window) *TabBase {
	return &TabBase{
		ti: container.NewTabItem(text, nil),
		w:  w,
	}
}

// NewTabWithIcon creates a managed container.TabItemWithIcon
func NewTabWithIcon(text string, icon fyne.Resource, w Window) *TabBase {
	return &TabBase{
		ti: container.NewTabItemWithIcon(text, icon, nil),
		w:  w,
	}
}

// FocusLocked if true, prevents the user from navigating away
func (b *TabBase) FocusLocked() bool {
	return false
}

// GetTab returns reference to the container.TabItem we manage
func (b *TabBase) GetTab() *container.TabItem {
	return b.ti
}

// GetWindow returns reference to the Window we belong to
func (b *TabBase) GetWindow() Window {
	return b.w
}

// SetContent replaces the contents of the managed container.TabItem
func (b *TabBase) SetContent(c fyne.CanvasObject) {
	b.ti.Content = c
}

// Init creates a new "Hello, World" tab (should be overridden!)
func (b *TabBase) Init() {
	b.SetContent(widget.NewLabel("Hello, World!"))
}

// Refresh re-creates the tab's contents if required (should be overriden!)
func (b *TabBase) Refresh() {}

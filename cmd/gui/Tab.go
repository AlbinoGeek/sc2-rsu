package gui

import (
	"fyne.io/fyne/container"
)

// Tab represents a managed container.TabItem
type Tab interface {
	FocusLocked() bool
	GetTab() *container.TabItem
	GetWindow() Window
	Init()
	Refresh()
}

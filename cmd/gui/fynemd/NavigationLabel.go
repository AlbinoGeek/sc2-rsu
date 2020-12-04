package fynemd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// NavigationLabel ...
type NavigationLabel struct {
	content fyne.CanvasObject
	icon    fyne.Resource
	text    string

	res fyne.CanvasObject
}

// NewNavigationLabel ...
func NewNavigationLabel(label string, content fyne.CanvasObject) NavigationItem {
	return &NavigationLabel{
		content: content,
		text:    label,
	}
}

// NewNavigationLabelWithIcon ...
func NewNavigationLabelWithIcon(label string, icon fyne.Resource, content fyne.CanvasObject) NavigationItem {
	return &NavigationLabel{
		content: content,
		icon:    theme.NewThemedResource(icon, nil),
		text:    label,
	}
}

// GetContent ...
func (l *NavigationLabel) GetContent() fyne.CanvasObject { return l.content }

// GetLabel ...
func (l *NavigationLabel) GetLabel() fyne.CanvasObject {
	if l.res != nil {
		if b, ok := l.res.(*widget.Button); ok {
			refresh := false
			if b.Icon != l.icon {
				b.Icon = l.icon
				refresh = true
			}
			if b.Text != l.text {
				b.Text = l.text
			}
			if refresh {
				b.Refresh()
			}
		}

		return l.res
	}

	b := widget.NewButtonWithIcon(l.text, l.icon, nil)
	b.Alignment = widget.ButtonAlignLeading
	b.HideShadow = true
	b.Importance = widget.LowImportance
	l.res = b
	return b
}

// GetIcon ...
func (l *NavigationLabel) GetIcon() fyne.Resource { return l.icon }

// GetTitle ...
func (l *NavigationLabel) GetTitle() string { return l.text }

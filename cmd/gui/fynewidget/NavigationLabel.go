package fynewidget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// NavigationLabel ...
type NavigationLabel struct {
	content fyne.CanvasObject
	icon    fyne.Resource
	label   string

	res fyne.CanvasObject
}

// NewNavigationLabel ...
func NewNavigationLabel(label string, content fyne.CanvasObject) NavigationItem {
	return &NavigationLabel{
		content: content,
		label:   label,
	}
}

// NewNavigationLabelWithIcon ...
func NewNavigationLabelWithIcon(label string, icon fyne.Resource, content fyne.CanvasObject) NavigationItem {
	return &NavigationLabel{
		content: content,
		icon:    theme.NewThemedResource(icon, nil),
		label:   label,
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
			if b.Text != l.label {
				b.Text = l.label
			}
			if refresh {
				b.Refresh()
			}
		}

		return l.res
	}

	b := widget.NewButtonWithIcon(l.label, l.icon, nil)
	b.Alignment = widget.ButtonAlignLeading
	b.HideShadow = true
	b.Importance = widget.LowImportance
	l.res = b
	return b
}

package fynemd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// NavLabel ...
type NavLabel struct {
	content fyne.CanvasObject
	label   fyne.CanvasObject
	icon    fyne.Resource
	text    string
}

// NewNavLabel ...
func NewNavLabel(label string, content fyne.CanvasObject) NavItem {
	return &NavLabel{
		content: content,
		text:    label,
	}
}

// NewNavLabelWithIcon ...
func NewNavLabelWithIcon(label string, icon fyne.Resource, content fyne.CanvasObject) NavItem {
	return &NavLabel{
		content: content,
		icon:    theme.NewThemedResource(icon, nil),
		text:    label,
	}
}

// GetContent ...
func (l *NavLabel) GetContent() fyne.CanvasObject { return l.content }

// GetLabel ...
func (l *NavLabel) GetLabel() fyne.CanvasObject {
	if l.label == nil {
		b := &widget.Button{
			Alignment:  widget.ButtonAlignLeading,
			Importance: widget.LowImportance,
			Text:       l.text,
			Icon:       l.icon,
		}

		b.ExtendBaseWidget(b)

		l.label = b

		return b
	}

	if b, ok := l.label.(*widget.Button); ok {
		refresh := false

		if b.Icon != l.icon {
			b.Icon = l.icon
			refresh = true
		}

		if b.Text != l.text {
			b.Text = l.text
			refresh = true
		}

		if refresh {
			b.Refresh()
		}
	}

	return l.label
}

// GetIcon ...
func (l *NavLabel) GetIcon() fyne.Resource { return l.icon }

// GetTitle ...
func (l *NavLabel) GetTitle() string { return l.text }

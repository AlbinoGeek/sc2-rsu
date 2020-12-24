package fynemd

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// NavDrawer ...
type NavDrawer struct {
	widget.BaseWidget

	OnDeselect func(NavItem) bool
	OnSelect   func(NavItem)

	items      []NavItem
	objects    []fyne.CanvasObject
	objectLock sync.RWMutex
	selected   int

	image     *widget.Icon      // dup: objects[0]
	separator *widget.Separator // dup: objects[3]
	subtitle  *canvas.Text      // dup: objects[2]
	title     *canvas.Text      // dup: objects[1]
}

// NewNavDrawer ...
func NewNavDrawer(title, subtitle string, items ...NavItem) *NavDrawer {
	sub := NewScaledText(TextSizeBody2, subtitle)
	sub.Color = theme.DisabledTextColor()

	ret := &NavDrawer{
		items:     items,
		image:     widget.NewIcon(theme.CancelIcon()),
		separator: widget.NewSeparator(),
		subtitle:  sub,
		title:     NewScaledText(TextSizeHeading5, title),
	}
	ret.objects = []fyne.CanvasObject{
		ret.title,
		ret.subtitle,
		ret.separator,
	}

	ret.image.Hide()
	ret.ExtendBaseWidget(ret)

	return ret
}

// CreateRenderer returns a new WidgetRenderer for this widget.
// This should not be called by regular code, it is used internally to render a widget.
//
// Implements: fyne.Widget
func (nav *NavDrawer) CreateRenderer() fyne.WidgetRenderer {
	return &navDrawerRenderer{nav: nav}
}

// Select ...
func (nav *NavDrawer) Select(id int) {
	if nav.OnDeselect != nil {
		// they can keepfocus (example: unsaved changes)
		if !nav.OnDeselect(nav.items[nav.selected]) {
			return
		}
	}

	nav.objectLock.RLock()

	// ! 4+ hard-coded
	// ! (*widget.Button) hard-coded
	if b, ok := nav.objects[4+nav.selected].(*widget.Button); ok {
		b.Style = widget.DefaultButton
		b.Refresh()
	}

	nav.selected = id

	// ! 4+ hard-coded
	// ! (*widget.Button) hard-coded
	if b, ok := nav.objects[4+nav.selected].(*widget.Button); ok {
		b.Style = widget.PrimaryButton
		b.Refresh()
	}

	if nav.OnSelect != nil {
		nav.OnSelect(nav.items[nav.selected])
	}

	nav.objectLock.RUnlock()
}

// SetImage ...
func (nav *NavDrawer) SetImage(image fyne.Resource) {
	nav.image.SetResource(image)
	nav.image.Hidden = image == nil
	nav.Refresh()
}

// SetSubtitle ...
func (nav *NavDrawer) SetSubtitle(subtitle string) {
	nav.subtitle.Hidden = subtitle == ""
	nav.subtitle.Text = subtitle
	nav.Refresh()
}

// SetTitle ...
func (nav *NavDrawer) SetTitle(title string) {
	nav.title.Hidden = title == ""
	nav.title.Text = title
	nav.Refresh()
}

package fynex

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// navDrawerRenderer defines the behaviour of a NavDraawer's implementation.
// This is returned from a widget's declarative object through the CreateRenderer()
// function and should be exactly one instance per widget in memory.
//
// Implements: fyne.WidgetRenderer
type navDrawerRenderer struct {
	nav *NavDrawer
}

// BackgroundColor returns the color that should be used to draw the background of this renderer’s widget.
//
// Deprecated: Widgets will no longer have a background to support hover and selection indication in collection widgets.
// If a widget requires a background color or image, this can be achieved by using a canvas.Rect or canvas.Image
// as the first child of a MaxLayout, followed by the rest of the widget components.
//
// Implements: fyne.WidgetRenderer
func (l *navDrawerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

// Destroy is for internal use.
//
// Implements: fyne.WidgetRenderer
func (l *navDrawerRenderer) Destroy() {}

// Layout is a hook that is called if the widget needs to be laid out.
// This should never call Refresh.
//
// Implements: fyne.WidgetRenderer
// TODO : ALIGN ELEMENTS ACCORDING TO MATERIAL DESIGN SPECS
func (l *navDrawerRenderer) Layout(space fyne.Size) {
	var (
		hasImage    = l.nav.image.Visible()
		hasTitle    = l.nav.title.Text != ""
		hasSubtitle = l.nav.subtitle.Text != ""
		hasSep      = hasSubtitle
	)

	pos := fyne.NewPos(Padding, Padding/2)

	if hasImage {
		l.nav.image.Resize(fyne.NewSize(40, 40))
		l.nav.image.Move(pos)
		pos.Y += l.nav.image.Size().Height + Padding/2
	}

	if hasTitle {
		l.nav.title.Resize(l.nav.title.MinSize())
		l.nav.title.Move(pos)
		pos.Y += l.nav.title.Size().Height + Padding/2
	}

	sepSize := fyne.NewSize(space.Width, 1)

	l.nav.subtitle.Resize(l.nav.subtitle.MinSize())
	l.nav.subtitle.Move(pos)

	l.nav.separator.Resize(sepSize)

	pos.X = 0

	l.nav.separator.Move(pos)

	qpad := Padding / 4

	if !hasSep {
		l.nav.subtitle.Hide()
		l.nav.separator.Hide()

		if !hasImage && !hasTitle {
			pos.Y = 0
		}
	} else {
		pos.Y += l.nav.subtitle.Size().Height + qpad
	}

	for i, o := range l.Objects()[4:] {
		// object returned by NavSeparator.GetLabel()
		if sep, ok := o.(*widget.Separator); ok {
			sep.Resize(sepSize)

			sep.Move(fyne.NewPos(0, pos.Y+qpad+2))
			pos.Y += Padding

			continue
		}

		// object returned by NavLabel.GetLabel()
		if b, ok := o.(*widget.Button); ok {
			if b.OnTapped == nil {
				b.OnTapped = func(j int) func() {
					return func() {
						l.nav.Select(j)
					}
				}(i)
			}
		}

		// resizing for a button-like object
		size := o.MinSize()
		size.Width = space.Width
		size.Height += Padding
		o.Resize(size)
		o.Move(pos)
		pos.Y += size.Height + qpad
	}
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
//
// Implements: fyne.WidgetRenderer
func (l *navDrawerRenderer) MinSize() fyne.Size {
	size := fyne.NewSize(Padding, Padding)

	sep := l.nav.separator.Position()
	size.Add(fyne.NewSize(sep.X, sep.Y))

	for _, o := range l.Objects()[4:] {
		if o == nil || !o.Visible() {
			continue
		}

		// at least as wide as the widest child
		childSize := o.MinSize()
		size = size.Max(childSize)

		// and the height of the child + padding
		// TODO: handle separators (which have more padding) vs buttons
		size.Height += childSize.Height + Padding
	}

	// hard minimum size: 128x128
	return size.Max(fyne.NewSize(128, 128))
}

// Objects returns all objects that should be drawn.
//
// Implements: fyne.WidgetRenderer
func (l *navDrawerRenderer) Objects() []fyne.CanvasObject {
	l.nav.objectLock.Lock()
	l.nav.objects = []fyne.CanvasObject{l.nav.image, l.nav.title, l.nav.subtitle, l.nav.separator}

	for _, o := range l.nav.items {
		if o != nil {
			l.nav.objects = append(l.nav.objects, o.GetLabel())
		}
	}

	l.nav.objectLock.Unlock()

	return l.nav.objects
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
//
// Implements: fyne.WidgetRenderer
func (l *navDrawerRenderer) Refresh() {
	for _, o := range l.Objects() {
		if o == nil || !o.Visible() {
			continue
		}

		o.Refresh()
	}
}

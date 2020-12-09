package fynemd

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type navigationDrawerRenderer struct {
	nav *NavigationDrawer
}

// NavigationDrawerRenderer ...
func NavigationDrawerRenderer(nav *NavigationDrawer) fyne.WidgetRenderer {
	return &navigationDrawerRenderer{
		nav: nav,
	}
}

// BackgroundColor
//
// Implements: fyne.WidgetRenderer
func (l *navigationDrawerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

// Destroy
//
// Implements: fyne.WidgetRenderer
func (l *navigationDrawerRenderer) Destroy() {}

// Layout
//
// Implements: fyne.WidgetRenderer
// TODO : ALIGN ELEMENTS ACCORDING TO MATERIAL DESIGN SPECS
func (l *navigationDrawerRenderer) Layout(space fyne.Size) {
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
		if sep, ok := o.(*widget.Separator); ok {
			sep.Resize(sepSize)
			sep.Move(fyne.NewPos(0, pos.Y+qpad-1))
			pos.Y += Padding
			continue
		}

		if b, ok := o.(*widget.Button); ok {
			if b.OnTapped == nil {
				b.OnTapped = func(j int) func() {
					return func() {
						l.nav.Select(j)
					}
				}(i)
			}
		}

		size := o.MinSize()
		size.Width = space.Width
		size.Height += Padding
		o.Resize(size)
		o.Move(pos)
		pos.Y += size.Height + Padding/2
	}
}

// MinSize
//
// Implements: fyne.WidgetRenderer
func (l *navigationDrawerRenderer) MinSize() fyne.Size {
	size := fyne.NewSize(Padding, Padding)
	for _, o := range l.Objects() {
		if o == nil || !o.Visible() {
			continue
		}

		childSize := o.MinSize()
		size = size.Max(childSize)
		size.Height += childSize.Height + Padding/2
	}

	return size.Max(fyne.NewSize(128, 128)).Add(fyne.NewSize(Padding, 0))
}

// Objects
//
// Implements: fyne.WidgetRenderer
func (l *navigationDrawerRenderer) Objects() []fyne.CanvasObject {
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

// Refresh
//
// Implements: fyne.WidgetRenderer
func (l *navigationDrawerRenderer) Refresh() {
	for _, o := range l.Objects() {
		if o == nil || !o.Visible() {
			continue
		}

		o.Refresh()
	}
}

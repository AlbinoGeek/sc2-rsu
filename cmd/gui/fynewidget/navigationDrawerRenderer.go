package fynewidget

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

func (l *navigationDrawerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (l *navigationDrawerRenderer) Destroy() {}

// Layout implements fyne.WidgetRenderer.Layout
func (l *navigationDrawerRenderer) Layout(space fyne.Size) {
	pad := theme.Padding()

	pos := fyne.NewPos(0, 0)
	if l.nav.image.Visible() {
		l.nav.image.Resize(fyne.NewSize(48, 48))
		l.nav.image.Move(pos)
		pos.Y += l.nav.image.Size().Height
	}

	l.nav.title.Resize(l.nav.title.MinSize())
	l.nav.title.Move(pos)
	pos.Y += l.nav.title.Size().Height

	sepSize := fyne.NewSize(space.Width, 1)
	l.nav.subtitle.Resize(l.nav.subtitle.MinSize())
	l.nav.subtitle.Move(pos)

	l.nav.separator.Resize(sepSize)
	l.nav.separator.Move(fyne.NewPos(0, pos.Y))

	if l.nav.subtitle.Text == "" {
		l.nav.subtitle.Hide()
		l.nav.separator.Hide()
	} else {
		pos.Y += l.nav.subtitle.Size().Height + pad
	}

	pos.Y += pad
	for i, o := range l.Objects()[4:] {
		if sep, ok := o.(*widget.Separator); ok {
			sep.Resize(sepSize)
			sep.Move(fyne.NewPos(0, pos.Y+pad-1))
			pos = pos.Add(fyne.NewPos(0, pad*2))
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
		size.Height += pad
		o.Resize(size)
		o.Move(pos)
		pos = pos.Add(fyne.NewPos(0, size.Height))
	}
}

func (l *navigationDrawerRenderer) Objects() []fyne.CanvasObject {
	l.nav.objects = []fyne.CanvasObject{l.nav.image, l.nav.title, l.nav.subtitle, l.nav.separator}
	for _, o := range l.nav.items {
		if o != nil {
			l.nav.objects = append(l.nav.objects, o.GetLabel())
		}
	}

	return l.nav.objects
}

func (l *navigationDrawerRenderer) Refresh() {
	for _, o := range l.Objects() {
		if o == nil || !o.Visible() {
			continue
		}

		o.Refresh()
	}
}

// MinSize implements fyne.WidgetRenderer.MinSize
func (l *navigationDrawerRenderer) MinSize() fyne.Size {
	pad := theme.Padding()

	size := fyne.NewSize(pad, pad)
	for _, o := range l.Objects() {
		if o == nil || !o.Visible() {
			continue
		}

		childSize := o.MinSize()
		size = size.Max(childSize)
		size.Height += childSize.Height + pad
	}

	return size
}

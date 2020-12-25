package fynex

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var (
	// IconSize used with (in display pixels) used with Material Design elements
	IconSize = 24

	// Padding (in display pixels) used with Material Design elements
	Padding = 16
)

// AppBar ...
type AppBar struct {
	widget.BaseWidget

	Dense     bool
	Extended  bool
	NavClosed bool
	navIcon   *widget.Button
	title     *canvas.Text

	actions []*widget.Button
	nav     *NavDrawer
	objects []fyne.CanvasObject
}

// NewAppBar ...
func NewAppBar(title string) *AppBar {
	bar := &AppBar{
		title: NewTextWithStyle(title, fyne.TextAlignLeading, StyleHeading5()),
	}

	bar.navIcon = widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		if bar.NavClosed {
			bar.SetNavClosed(false)
		} else {
			bar.SetNavClosed(true)
		}
	})
	bar.navIcon.Importance = widget.LowImportance

	bar.ExtendBaseWidget(bar)

	return bar
}

// Refresh ...
func (bar *AppBar) Refresh() {
	bar.ExtendBaseWidget(bar)
	bar.BaseWidget.Refresh()
}

// SetDense ...
func (bar *AppBar) SetDense(dense bool) {
	bar.Dense = dense
	bar.Refresh()
}

// SetExtended ...
func (bar *AppBar) SetExtended(extended bool) {
	bar.Extended = extended
	bar.Refresh()
}

// SetTitle ...
func (bar *AppBar) SetTitle(title string) {
	bar.title.Text = title
	bar.Refresh()
}

// SetNav ...
func (bar *AppBar) SetNav(nav *NavDrawer) {
	bar.nav = nav
	bar.Refresh()
}

// SetNavClosed ...
func (bar *AppBar) SetNavClosed(closed bool) {
	if bar.nav == nil {
		return
	}

	if bar.NavClosed = closed; bar.NavClosed {
		bar.nav.Hide()
	} else {
		bar.nav.Show()
	}

	bar.Refresh()
}

// CreateRenderer returns a new WidgetRenderer for this widget.
// This should not be called by regular code, it is used internally to render a widget.
//
// Implements: fyne.Widget
func (bar *AppBar) CreateRenderer() fyne.WidgetRenderer {
	return &appBarRenderer{
		bar: bar,
	}
}

// --

// appBarRenderer defines the behaviour of a AppBar's implementation.
// This is returned from a widget's declarative object through the CreateRenderer()
// function and should be exactly one instance per widget in memory.
//
// Implements: fyne.WidgetRenderer
type appBarRenderer struct {
	bar *AppBar
}

// BackgroundColor returns the color that should be used to draw the background of this renderer’s widget.
//
// Deprecated: Widgets will no longer have a background to support hover and selection indication in collection widgets.
// If a widget requires a background color or image, this can be achieved by using a canvas.Rect or canvas.Image
// as the first child of a MaxLayout, followed by the rest of the widget components.
//
// Implements: fyne.WidgetRenderer
func (*appBarRenderer) BackgroundColor() color.Color {
	return theme.PrimaryColor()
}

// Destroy is for internal use.
//
// Implements: fyne.WidgetRenderer
func (*appBarRenderer) Destroy() {}

// Layout is a hook that is called if the widget needs to be laid out.
// This should never call Refresh.
//
// Implements: fyne.WidgetRenderer
// ! should respond to theme values
func (br *appBarRenderer) Layout(space fyne.Size) {
	pos := fyne.NewPos(Padding, Padding)

	if br.bar.Dense {
		pos.Y = Padding / 2
	}

	br.bar.navIcon.Move(pos)
	br.bar.navIcon.Resize(fyne.NewSize(IconSize, IconSize))

	pos.Y -= Padding / 3

	if br.bar.nav == nil || !br.bar.NavClosed {
		br.bar.navIcon.Hide()
	} else {
		pos.X += br.bar.navIcon.Size().Width + Padding
		br.bar.navIcon.Show()
	}

	br.bar.title.Move(pos)
	br.bar.title.Resize(br.bar.title.MinSize())

	// TODO: Layout actions from right

	// TODO: if len(actions) > 3 { actionsMenu = actions[2:] } ...
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
//
// Implements: fyne.WidgetRenderer
func (br *appBarRenderer) MinSize() fyne.Size {
	// pad := theme.Padding()

	size := fyne.NewSize(360, 56) // material specs

	if br.bar.Dense {
		size.Height = 40
	} else if br.bar.Extended {
		size.Height = 128
	}

	// enough space for the text
	size = size.Max(fyne.MeasureText(br.bar.title.Text, theme.TextSize(), br.bar.title.TextStyle))

	// TODO: enough space for Action buttons when implemented

	return size
}

// Objects returns all objects that should be drawn.
//
// Implements: fyne.WidgetRenderer
func (br *appBarRenderer) Objects() []fyne.CanvasObject {
	br.bar.objects = []fyne.CanvasObject{
		br.bar.navIcon, br.bar.title,
	}

	// TODO: append Action buttons when implemented

	return br.bar.objects
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
//
// Implements: fyne.WidgetRenderer
func (br *appBarRenderer) Refresh() {
	// state mismatch -- navIcon visibility must change
	if br.bar.navIcon.Visible() != br.bar.NavClosed {
		br.Layout(br.bar.Size())
	}

	for _, o := range br.Objects() {
		if o == nil || !o.Visible() {
			continue
		}

		o.Refresh()
	}
}

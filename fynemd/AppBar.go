package fynemd

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
	Title     string

	actions []*widget.Button
	nav     *NavDrawer
	objects []fyne.CanvasObject
}

// NewAppBar ...
func NewAppBar(title string) *AppBar {
	bar := &AppBar{
		Title: title,
	}

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
	bar.Title = title
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
	rend := &appBarRenderer{
		bar: bar,
	}

	rend.Init()

	return rend
}

// --

// appBarRenderer defines the behaviour of a AppBar's implementation.
// This is returned from a widget's declarative object through the CreateRenderer()
// function and should be exactly one instance per widget in memory.
//
// Implements: fyne.WidgetRenderer
type appBarRenderer struct {
	bar *AppBar

	navIcon *widget.Button
	title   *canvas.Text
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

func (br *appBarRenderer) Init() {
	br.navIcon = widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		if br.bar.NavClosed {
			br.bar.SetNavClosed(false)
		} else {
			br.bar.SetNavClosed(true)
		}
	})
	br.title = NewScaledText(TextSizeHeading5, br.bar.Title)
	// br.title.TextStyle.Bold = true
	br.bar.objects = []fyne.CanvasObject{
		br.navIcon, br.title,
	}
}

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

	br.navIcon.Move(pos)
	br.navIcon.Resize(fyne.NewSize(IconSize, IconSize))

	pos.Y -= Padding / 3

	if br.bar.nav == nil || !br.bar.NavClosed {
		br.navIcon.Hide()
	} else {
		pos.X += br.navIcon.Size().Width + Padding
		br.navIcon.Show()
	}

	br.title.Move(pos)
	br.title.Resize(br.title.MinSize())

	if br.bar.Title == "" && br.title.Visible() {
		br.title.Hide()
	}

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

	// TODO: enough space for NavIcon if visible

	// enough space for the text
	size = size.Max(fyne.MeasureText(br.bar.Title, theme.TextSize(), br.bar.nav.title.TextStyle))

	// TODO: enough space for action buttons when implemented

	return size
}

// Objects returns all objects that should be drawn.
//
// Implements: fyne.WidgetRenderer
func (br *appBarRenderer) Objects() []fyne.CanvasObject {
	return br.bar.objects
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
//
// Implements: fyne.WidgetRenderer
func (br *appBarRenderer) Refresh() {
	// state mismatch -- navIcon visibility must change
	if br.navIcon.Visible() != br.bar.NavClosed {
		br.Layout(br.bar.Size())
	}

	for _, o := range br.Objects() {
		if o == nil || !o.Visible() {
			continue
		}

		o.Refresh()
	}
}
